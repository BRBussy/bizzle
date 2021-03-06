package mongo

import (
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/role"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	validator  validationValidator.Validator
	collection *mongo.Collection
}

func New(
	validator validationValidator.Validator,
	database *mongo.Database,
) (roleStore.Store, error) {
	// get role collection
	roleCollection := database.Collection("role")

	// setup collection indices
	if err := roleCollection.SetupIndices([]mongoDriver.IndexModel{
		mongo.NewUniqueIndex("id"),
		mongo.NewUniqueIndex("name"),
	}); err != nil {
		log.Error().Err(err).Msg("error setting up role collection indices")
		return nil, err
	}

	return &store{
		collection: database.Collection("role"),
		validator:  validator,
	}, nil
}

func (s *store) CreateOne(request roleStore.CreateOneRequest) (*roleStore.CreateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if err := s.collection.CreateOne(request.Role); err != nil {
		log.Error().Err(err).Msg("creating role")
		return nil, err
	}
	return &roleStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request roleStore.FindOneRequest) (*roleStore.FindOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var result role.Role
	if err := s.collection.FindOne(&result, request.Identifier); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("finding one role")
			return nil, err
		}
	}
	return &roleStore.FindOneResponse{Role: result}, nil
}

func (s *store) FindMany(request roleStore.FindManyRequest) (*roleStore.FindManyResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var records []role.Role
	count, err := s.collection.FindMany(&records, request.Criteria, request.Query)
	if err != nil {
		log.Error().Err(err).Msg("finding roles")
		return nil, bizzleException.ErrUnexpected{}
	}
	if records == nil {
		records = make([]role.Role, 0)
	}

	return &roleStore.FindManyResponse{
		Records: records,
		Total:   count,
	}, nil
}

func (s *store) UpdateOne(request roleStore.UpdateOneRequest) (*roleStore.UpdateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if err := s.collection.UpdateOne(request.Role, request.Role.ID); err != nil {
		log.Error().Err(err).Msg("updating role")
		return nil, err
	}
	return &roleStore.UpdateOneResponse{}, nil
}
