package mongo

import (
	budgetEntryCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	budgetEntryCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/scope"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	validator  validationValidator.Validator
	collection *mongo.Collection
	scopeAdmin scope.Admin
}

func New(
	validator validationValidator.Validator,
	scopeAdmin scope.Admin,
	database *mongo.Database,
) (budgetEntryCategoryRuleStore.Store, error) {
	// get budgetEntryCategoryRule collection
	budgetEntryCollection := database.Collection("budgetEntryCategoryRule")

	// setup collection indices
	if err := budgetEntryCollection.SetupIndices([]mongoDriver.IndexModel{
		mongo.NewUniqueIndex("id"),
		mongo.NewUniqueIndex("id", "name"),
	}); err != nil {
		log.Error().Err(err).Msg("error setting up budgetEntryCategoryRule collection indices")
		return nil, err
	}

	return &store{
		validator:  validator,
		collection: budgetEntryCollection,
		scopeAdmin: scopeAdmin,
	}, nil
}

func (s *store) CreateOne(request budgetEntryCategoryRuleStore.CreateOneRequest) (*budgetEntryCategoryRuleStore.CreateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	if err := s.collection.CreateOne(request.CategoryRule); err != nil {
		log.Error().Err(err).Msg("creating category")
		return nil, err
	}
	return &budgetEntryCategoryRuleStore.CreateOneResponse{}, nil
}

func (s *store) CreateMany(request budgetEntryCategoryRuleStore.CreateManyRequest) (*budgetEntryCategoryRuleStore.CreateManyResponse, error) {
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

	return &budgetEntryCategoryRuleStore.CreateManyResponse{}, nil
}

func (s *store) FindOne(request budgetEntryCategoryRuleStore.FindOneRequest) (*budgetEntryCategoryRuleStore.FindOneResponse, error) {
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

	var result budgetEntryCategoryRule.CategoryRule
	if err := s.collection.FindOne(&result, applyScopeToIdentifierResponse.ScopedIdentifier); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("finding one budgetEntryCategoryRule")
			return nil, err
		}
	}

	return &budgetEntryCategoryRuleStore.FindOneResponse{
		CategoryRule: result,
	}, nil
}

func (s *store) FindMany(request budgetEntryCategoryRuleStore.FindManyRequest) (*budgetEntryCategoryRuleStore.FindManyResponse, error) {
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

	var records []budgetEntryCategoryRule.CategoryRule
	count, err := s.collection.FindMany(&records, applyScopeToCriteriaResponse.ScopedCriteria, request.Query)
	if err != nil {
		log.Error().Err(err).Msg("finding exercises")
		return nil, bizzleException.ErrUnexpected{}
	}
	if records == nil {
		records = make([]budgetEntryCategoryRule.CategoryRule, 0)
	}

	return &budgetEntryCategoryRuleStore.FindManyResponse{
		Records: records,
		Total:   count,
	}, nil
}

func (s *store) UpdateOne(request budgetEntryCategoryRuleStore.UpdateOneRequest) (*budgetEntryCategoryRuleStore.UpdateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if err := s.collection.UpdateOne(request.CategoryRule, request.CategoryRule.ID); err != nil {
		log.Error().Err(err).Msg("updating budgetEntryCategoryRule")
		return nil, err
	}

	return &budgetEntryCategoryRuleStore.UpdateOneResponse{}, nil
}
