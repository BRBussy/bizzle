package exercise

import (
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/exercise"
	exerciseAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/admin"
	"github.com/rs/zerolog/log"
)

var initialExercises = []exercise.Exercise{
	{
		Name:        "Arm Curl",
		MuscleGroup: exercise.BicepsMuscleGroup,
		Variant:     "Dumbbell",
		Description: "Stand and Curl",
	},
}

func Setup(exerciseAdminImp exerciseAdmin.Admin) error {

	for i := range initialExercises {
		log.Info().Msg(fmt.Sprintf(
			"creating exercise %d/%d",
			i+1,
			len(initialExercises),
		))
		if _, err := exerciseAdminImp.CreateOne(&exerciseAdmin.CreateOneRequest{
			Exercise: initialExercises[i],
		}); err != nil {
			log.Error().Err(err).Msg("creating exercise")
			return err
		}
	}

	return nil
}
