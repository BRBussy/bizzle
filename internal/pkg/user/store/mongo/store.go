package mongo

import (
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/scope"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	scopeAdmin scope.Admin
	validator  validationValidator.Validator
	collection *mongo.Collection
}

func New(
	validator validationValidator.Validator,
	scopeAdmin scope.Admin,
	database *mongo.Database,
) (userStore.Store, error) {
	// get user collection
	userCollection := database.Collection("user")

	// setup collection indices
	if err := userCollection.SetupIndices([]mongoDriver.IndexModel{
		mongo.NewUniqueIndex("id"),
		mongo.NewUniqueIndex("email"),
	}); err != nil {
		log.Error().Err(err).Msg("error setting up user collection indices")
		return nil, err
	}

	return &store{
		scopeAdmin: scopeAdmin,
		collection: database.Collection("user"),
		validator:  validator,
	}, nil
}

func (s *store) CreateOne(request userStore.CreateOneRequest) (*userStore.CreateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if err := s.collection.CreateOne(request.User); err != nil {
		log.Error().Err(err).Msg("error creating user")
		return nil, err
	}
	return &userStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request userStore.FindOneRequest) (*userStore.FindOneResponse, error) {
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

	var result user.User
	if err := s.collection.FindOne(&result, applyScopeToIdentifierResponse.ScopedIdentifier); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("finding one user")
			return nil, err
		}
	}
	return &userStore.FindOneResponse{User: result}, nil
}

func (s *store) FindMany(request userStore.FindManyRequest) (*userStore.FindManyResponse, error) {
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

	var records []user.User
	count, err := s.collection.FindMany(&records, applyScopeToCriteriaResponse.ScopedCriteria, request.Query)
	if err != nil {
		log.Error().Err(err).Msg("finding users")
		return nil, bizzleException.ErrUnexpected{}
	}
	if records == nil {
		records = make([]user.User, 0)
	}

	return &userStore.FindManyResponse{
		Records: records,
		Total:   count,
	}, nil
}

func (s *store) UpdateOne(request userStore.UpdateOneRequest) (*userStore.UpdateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	applyScopeToIdentifierResponse, err := s.scopeAdmin.ApplyScopeToIdentifier(
		scope.ApplyScopeToIdentifierRequest{
			Claims:            request.Claims,
			IdentifierToScope: request.User.ID,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("could not apply scope to identifier")
		return nil, bizzleException.ErrUnexpected{}
	}

	if err := s.collection.UpdateOne(request.User, applyScopeToIdentifierResponse.ScopedIdentifier); err != nil {
		log.Error().Err(err).Msg("updating user")
		return nil, err
	}
	return &userStore.UpdateOneResponse{}, nil
}
