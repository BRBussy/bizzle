package user

import (
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/text"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var rootUserToCreate = user.User{
	Name:  "root",
	Email: "root@bizzle.com",
	RoleIDs: []identifier.ID{
		"root",
	},
}

func Setup(
	userStoreImp userStore.Store,
	roleStoreImp roleStore.Store,
	rootPassword string,
) error {
	// find roles which need to be assigned to root user
	if len(rootUserToCreate.RoleIDs) > 0 {
		roleFindCriteria := make([]criterion.Criterion, 0)
		for i := range rootUserToCreate.RoleIDs {
			roleFindCriteria = append(
				roleFindCriteria,
				stringCriterion.Exact{
					Field: "name",
					Text:  rootUserToCreate.RoleIDs[i].String(),
				},
			)
		}
		roleFindResponse, err := roleStoreImp.FindMany(
			roleStore.FindManyRequest{
				Criteria: roleFindCriteria,
			},
		)
		if err != nil {
			log.Error().Err(err).Msg("error finding root user's roles")
			return bizzleException.ErrUnexpected{}
		}
		rootUserToCreate.RoleIDs = make([]identifier.ID, 0)
		for i := range roleFindResponse.Records {
			rootUserToCreate.RoleIDs = append(rootUserToCreate.RoleIDs, roleFindResponse.Records[i].ID)
		}
	}

	// populate root user ID and ownerID
	rootUserToCreate.ID = identifier.ID(uuid.NewV4().String())
	rootUserToCreate.OwnerID = rootUserToCreate.ID

	// hash and populate root user's password
	pwdHash, err := bcrypt.GenerateFromPassword(
		[]byte(rootPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Error().Err(err).Msg("error hashing root user's password")
		return bizzleException.ErrUnexpected{}
	}
	rootUserToCreate.Password = pwdHash

	// create the root user
	if _, err := userStoreImp.CreateOne(
		userStore.CreateOneRequest{
			User: rootUserToCreate,
		},
	); err != nil {
		log.Error().Err(err).Msg("creating bizzle root user")
		return bizzleException.ErrUnexpected{}
	}

	return nil
}
