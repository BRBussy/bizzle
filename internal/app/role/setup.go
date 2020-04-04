package role

import (
	budgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin"
	budgetEntryAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/admin"
	budgetCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store"
	exerciseAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/admin"
	sessionAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/session/admin"
	sessionStore "github.com/BRBussy/bizzle/internal/pkg/exercise/session/store"
	exerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	securityPermission "github.com/BRBussy/bizzle/internal/pkg/security/permission"
	securityRole "github.com/BRBussy/bizzle/internal/pkg/security/role"
	roleAdmin "github.com/BRBussy/bizzle/internal/pkg/security/role/admin"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
)

var initialRoles = []securityRole.Role{
	{
		Name: "user",
		Permissions: []securityPermission.Permission{
			exerciseStore.FindOneService,
			exerciseStore.FindManyService,
			sessionStore.FindManyService,
			sessionStore.FindOneService,
			sessionAdmin.CreateOneService,
			budgetEntryAdmin.XLSXStandardBankStatementToBudgetCompositeEntriesService,
			budgetEntryAdmin.DuplicateCheckService,
			budgetEntryAdmin.CreateOneService,
			budgetEntryAdmin.CreateManyService,
			budgetEntryAdmin.UpdateOneService,
			budgetEntryAdmin.UpdateManyService,
			budgetEntryAdmin.DeleteOneService,
			budgetEntryStore.FindManyService,
			budgetAdmin.GetBudgetForDateRangeService,
			budgetCategoryRuleStore.FindManyService,
		},
	},
	{
		Name:        "system",
		Permissions: []securityPermission.Permission{},
	},
}

var rootOnlyPermissions = []securityPermission.Permission{
	roleStore.FindOneService,
	exerciseAdmin.CreateOneService,
}

func Setup(
	admin roleAdmin.Admin,
	store roleStore.Store,
) error {
	// retrieve root role
	var rootRole securityRole.Role
	var rootRoleCopy securityRole.Role
	findOneResponse, err := store.FindOne(roleStore.FindOneRequest{Identifier: identifier.Name("root")})
	switch err.(type) {
	case mongo.ErrNotFound:
		// root role not found, it should be created
		log.Info().Msg("creating root role")
		createOneResponse, err := admin.CreateOne(&roleAdmin.CreateOneRequest{
			Role: securityRole.Role{
				Name: "root",
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("creating root role")
			return err
		}
		// set root role
		rootRole = createOneResponse.Role
		rootRoleCopy = rootRole

	default:
		// there was some error retrieving the role
		log.Error().Err(err).Msg("finding root role")
		return err

	case nil:
		// root role found, set it
		rootRole = findOneResponse.Role
		rootRoleCopy = rootRole
	}
	// set root role permissions to root only permissions
	rootRole.Permissions = rootOnlyPermissions

	// for every initial role to create
	for i := range initialRoles {
		// add all permissions associated with role to root role
		for pI := range initialRoles[i].Permissions {
			rootRole.AddUniquePermission(initialRoles[i].Permissions[pI])
		}

		// try and retrieve the role
		findOneResponse, err := store.FindOne(roleStore.FindOneRequest{Identifier: identifier.Name(initialRoles[i].Name)})
		if err != nil {
			switch err.(type) {
			case mongo.ErrNotFound:
				// role was not found, create it and move on to next role
				createResponse, err := admin.CreateOne(&roleAdmin.CreateOneRequest{Role: initialRoles[i]})
				if err != nil {
					log.Error().Err(err).Msg("creating role")
					return err
				}
				findOneResponse = roleStore.FindOneResponse{Role: createResponse.Role}
				continue

			default:
				// there was some error retrieving the role
				log.Error().Err(err).Msg("finding role")
				return err
			}
		}
		// set id on initial permission to prevent incorrect compare result
		initialRoles[i].ID = findOneResponse.Role.ID

		// compare them to see if an update is required
		if !securityRole.CompareRoles(initialRoles[i], findOneResponse.Role) {
			// update as required
			log.Info().Msg("updating role " + initialRoles[i].Name)
			if _, err := admin.UpdateOne(&roleAdmin.UpdateOneRequest{Role: initialRoles[i]}); err != nil {
				log.Error().Err(err).Msg("updating role")
				return err
			}
		}
	}

	// check if root role should be updated
	if !securityRole.CompareRoles(rootRole, rootRoleCopy) {
		if _, err := admin.UpdateOne(&roleAdmin.UpdateOneRequest{Role: rootRole}); err != nil {
			log.Error().Err(err).Msg("updating root role")
			return err
		}
	}

	return nil
}
