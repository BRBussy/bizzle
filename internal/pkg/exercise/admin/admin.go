package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/exercise"
)

type Admin interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
	UpdateOne(*UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "Exercise-Admin"

const CreateOneService = ServiceProvider + ".CreateOne"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	Exercise exercise.Exercise `validate:"required"`
}

type CreateOneResponse struct {
	Exercise exercise.Exercise
}

type UpdateOneRequest struct {
	Exercise exercise.Exercise `validate:"required"`
}

type UpdateOneResponse struct {
	Exercise exercise.Exercise
}
