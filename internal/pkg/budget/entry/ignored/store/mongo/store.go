package mongo

import (
	budgetEntryIgnored "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored"
	budgetEntryIgnoredStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored/store"
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

// New creates a new mongo budget entry store
func New(
	validator validationValidator.Validator,
	scopeAdmin scope.Admin,
	database *mongo.Database,
) (budgetEntryIgnoredStore.Store, error) {
	// get budgetEntryIgnored collection
	budgetEntryCollection := database.Collection("budgetEntryIgnored")

	// setup collection indices
	if err := budgetEntryCollection.SetupIndices([]mongoDriver.IndexModel{
		mongo.NewUniqueIndex("id"),
		mongo.NewUniqueIndex("description"),
	}); err != nil {
		log.Error().Err(err).Msg("error setting up budgetEntryIgnored collection indices")
		return nil, err
	}

	return &store{
		validator:  validator,
		collection: budgetEntryCollection,
		scopeAdmin: scopeAdmin,
	}, nil
}

func (s *store) CreateOne(request budgetEntryIgnoredStore.CreateOneRequest) (*budgetEntryIgnoredStore.CreateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	if err := s.collection.CreateOne(request.Ignored); err != nil {
		log.Error().Err(err).Msg("creating ignored")
		return nil, err
	}
	return &budgetEntryIgnoredStore.CreateOneResponse{}, nil
}

func (s *store) CreateMany(request budgetEntryIgnoredStore.CreateManyRequest) (*budgetEntryIgnoredStore.CreateManyResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	documentsToCreate := make([]interface{}, 0)
	for _, entry := range request.Entries {
		documentsToCreate = append(documentsToCreate, entry)
	}
	if err := s.collection.CreateMany(documentsToCreate); err != nil {
		log.Error().Err(err).Msg("creating budget entries")
		return nil, err
	}

	return &budgetEntryIgnoredStore.CreateManyResponse{}, nil
}

func (s *store) FindOne(request budgetEntryIgnoredStore.FindOneRequest) (*budgetEntryIgnoredStore.FindOneResponse, error) {
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

	var result budgetEntryIgnored.Ignored
	if err := s.collection.FindOne(&result, applyScopeToIdentifierResponse.ScopedIdentifier); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("finding one budgetEntryIgnored")
			return nil, err
		}
	}

	return &budgetEntryIgnoredStore.FindOneResponse{
		Ignored: result,
	}, nil
}

func (s *store) FindMany(request budgetEntryIgnoredStore.FindManyRequest) (*budgetEntryIgnoredStore.FindManyResponse, error) {
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

	var records []budgetEntryIgnored.Ignored
	count, err := s.collection.FindMany(&records, applyScopeToCriteriaResponse.ScopedCriteria, request.Query)
	if err != nil {
		log.Error().Err(err).Msg("finding budget entries")
		return nil, bizzleException.ErrUnexpected{}
	}
	if records == nil {
		records = make([]budgetEntryIgnored.Ignored, 0)
	}

	return &budgetEntryIgnoredStore.FindManyResponse{
		Records: records,
		Total:   count,
	}, nil
}

func (s *store) UpdateOne(request budgetEntryIgnoredStore.UpdateOneRequest) (*budgetEntryIgnoredStore.UpdateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	applyScopeToIdentifierResponse, err := s.scopeAdmin.ApplyScopeToIdentifier(
		scope.ApplyScopeToIdentifierRequest{
			Claims:            request.Claims,
			IdentifierToScope: request.Ignored.ID,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("could not apply scope to identifier")
		return nil, bizzleException.ErrUnexpected{}
	}

	if err := s.collection.UpdateOne(request.Ignored, applyScopeToIdentifierResponse.ScopedIdentifier); err != nil {
		log.Error().Err(err).Msg("updating budgetEntryIgnored")
		return nil, err
	}

	return &budgetEntryIgnoredStore.UpdateOneResponse{}, nil
}

func (s *store) DeleteOne(request budgetEntryIgnoredStore.DeleteOneRequest) (*budgetEntryIgnoredStore.DeleteOneResponse, error) {
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
		log.Error().Err(err).Msg("updating budgetEntryIgnored")
		return nil, err
	}

	return &budgetEntryIgnoredStore.DeleteOneResponse{}, nil
}
