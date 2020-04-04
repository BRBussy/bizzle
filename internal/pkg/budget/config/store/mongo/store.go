package mongo

import (
	budgetConfig "github.com/BRBussy/bizzle/internal/pkg/budget/config"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/config/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/scope"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	scopeAdmin scope.Admin
	validator  validationValidator.Validator
	collection *mongo.Collection
}

// New creates a new mongo budget config store
func New(
	validator validationValidator.Validator,
	scopeAdmin scope.Admin,
	database *mongo.Database,
) (budgetEntryStore.Store, error) {
	// get budgetConfig collection
	budgetEntryCollection := database.Collection("budgetConfig")

	// setup collection indices
	if err := budgetEntryCollection.SetupIndices([]mongoDriver.IndexModel{
		mongo.NewUniqueIndex("id"),
		mongo.NewUniqueIndex("ownerID"),
	}); err != nil {
		log.Error().Err(err).Msg("error setting up budgetConfig collection indices")
		return nil, err
	}

	return &store{
		validator:  validator,
		collection: budgetEntryCollection,
		scopeAdmin: scopeAdmin,
	}, nil
}

func (s *store) CreateOne(request budgetEntryStore.CreateOneRequest) (*budgetEntryStore.CreateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	if err := s.collection.CreateOne(request.Config); err != nil {
		log.Error().Err(err).Msg("creating role")
		return nil, err
	}
	return &budgetEntryStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request budgetEntryStore.FindOneRequest) (*budgetEntryStore.FindOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	applyScopeToIdentifierResponse, err := s.scopeAdmin.ApplyScopeToIdentifier(
		scope.ApplyScopeToIdentifierRequest{
			Claims:            request.Claims,
			IdentifierToScope: request.Identifier,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("could not apply scope to identifier")
		return nil, bizzleException.ErrUnexpected{}
	}

	var result budgetConfig.Config
	if err := s.collection.FindOne(&result, applyScopeToIdentifierResponse.ScopedIdentifier); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("finding one budgetConfig")
			return nil, err
		}
	}

	return &budgetEntryStore.FindOneResponse{
		Config: result,
	}, nil
}

func (s *store) FindMany(request budgetEntryStore.FindManyRequest) (*budgetEntryStore.FindManyResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	applyScopeToCriteriaResponse, err := s.scopeAdmin.ApplyScopeToCriteria(
		scope.ApplyScopeToCriteriaRequest{
			Claims:          request.Claims,
			CriteriaToScope: request.Criteria,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("could not apply scope to criteria")
		return nil, bizzleException.ErrUnexpected{}
	}

	var records []budgetConfig.Config
	count, err := s.collection.FindMany(&records, applyScopeToCriteriaResponse.ScopedCriteria, request.Query)
	if err != nil {
		log.Error().Err(err).Msg("finding exercises")
		return nil, bizzleException.ErrUnexpected{}
	}
	if records == nil {
		records = make([]budgetConfig.Config, 0)
	}

	return &budgetEntryStore.FindManyResponse{
		Records: records,
		Total:   count,
	}, nil
}

func (s *store) UpdateOne(request budgetEntryStore.UpdateOneRequest) (*budgetEntryStore.UpdateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	applyScopeToIdentifierResponse, err := s.scopeAdmin.ApplyScopeToIdentifier(
		scope.ApplyScopeToIdentifierRequest{
			Claims:            request.Claims,
			IdentifierToScope: request.Config.ID,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("could not apply scope to identifier")
		return nil, bizzleException.ErrUnexpected{}
	}

	if err := s.collection.UpdateOne(request.Config, applyScopeToIdentifierResponse.ScopedIdentifier); err != nil {
		log.Error().Err(err).Msg("updating budgetConfig")
		return nil, err
	}

	return &budgetEntryStore.UpdateOneResponse{}, nil
}

func (s *store) DeleteOne(request budgetEntryStore.DeleteOneRequest) (*budgetEntryStore.DeleteOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	applyScopeToIdentifierResponse, err := s.scopeAdmin.ApplyScopeToIdentifier(
		scope.ApplyScopeToIdentifierRequest{
			Claims:            request.Claims,
			IdentifierToScope: request.Identifier,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("could not apply scope to identifier")
		return nil, bizzleException.ErrUnexpected{}
	}

	if err := s.collection.DeleteOne(applyScopeToIdentifierResponse.ScopedIdentifier); err != nil {
		log.Error().Err(err).Msg("updating budgetConfig")
		return nil, err
	}

	return &budgetEntryStore.DeleteOneResponse{}, nil
}
