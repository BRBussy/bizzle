package exercise

import (
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/exercise"
	exerciseAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/admin"
	"github.com/rs/zerolog/log"
)

var initialExercises = []exercise.Exercise{
	exercise.ArmCurl{
		SomeField: "Field data",
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
