package exercise

import (
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/exercise"
	exerciseAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/admin"
	exerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
)

var initialExercises = []exercise.Exercise{
	{
		Name:        "Arm Curl",
		MuscleGroup: exercise.BicepsMuscleGroup,
		Variant:     "Dumbbell",
		Description: "Perform curl with dumbbells",
	},
	{
		Name:        "Push Ups",
		MuscleGroup: exercise.PectoralMuscleGroup,
		Variant:     "Outside Shoulder",
		Description: "Perform push ups outside shoulder",
	},
	{
		Name:        "Lateral Raise",
		MuscleGroup: exercise.ShouldersMuscleGroup,
		Variant:     "Dumbbell",
		Description: "Raise dumbbells to the sides",
	},
	{
		Name:        "Lateral Raise",
		MuscleGroup: exercise.ShouldersMuscleGroup,
		Variant:     "Barbell",
		Description: "Raise barbell vertically to the sides",
	},
	{
		Name:        "Deltoid Raise",
		MuscleGroup: exercise.ShouldersMuscleGroup,
		Variant:     "Dumbbell",
		Description: "Raise dumbbells to the front",
	},
	{
		Name:        "Deltoid Raise",
		MuscleGroup: exercise.ShouldersMuscleGroup,
		Variant:     "Barbell",
		Description: "Raise barbell in an arc",
	},
	{
		Name:        "Plank",
		MuscleGroup: exercise.CoreMuscleGroup,
		Variant:     "Standard",
		Description: "Normal Plank Position",
	},
	{
		Name:        "Plank",
		MuscleGroup: exercise.CoreMuscleGroup,
		Variant:     "Side",
		Description: "Plank with one hand down, the other pointing up",
	},
}

func Setup(
	exerciseAdminImp exerciseAdmin.Admin,
	exerciseStoreImp exerciseStore.Store,
) error {

	// for every exercise to be created
	for i := range initialExercises {
		log.Info().Msg(fmt.Sprintf(
			"__________ exercise %d/%d - %s %s __________",
			i+1,
			len(initialExercises),
			initialExercises[i].Name,
			initialExercises[i].Variant,
		))

		// try and retrieve the exercise
		findOneResponse, err := exerciseStoreImp.FindOne(
			&exerciseStore.FindOneRequest{
				Identifier: identifier.NameVariant{
					Name:    initialExercises[i].Name,
					Variant: initialExercises[i].Variant,
				},
			},
		)
		switch err.(type) {
		case mongo.ErrNotFound:
			// not found, create
			log.Info().Msg("--> create")
			if _, err := exerciseAdminImp.CreateOne(
				&exerciseAdmin.CreateOneRequest{
					Exercise: initialExercises[i],
				},
			); err != nil {
				log.Error().Err(err).Msg("creating exercise")
				return err
			}

		case nil:
			// found, compare and update if necessary
			initialExercises[i].ID = findOneResponse.Exercise.ID
			if findOneResponse.Exercise != initialExercises[i] {
				log.Info().Msg("--> update")
				if _, err := exerciseAdminImp.UpdateOne(
					&exerciseAdmin.UpdateOneRequest{
						Exercise: initialExercises[i],
					},
				); err != nil {
					log.Error().Err(err).Msg("updating exercise")
					return err
				}
				continue
			}
			log.Info().Msg("--> do nothing")

		default:
			log.Error().Err(err).Msg("retrieving exercise")
		}
	}

	return nil
}
