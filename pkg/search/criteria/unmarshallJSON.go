package criteria

import (
	"encoding/json"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	numberCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/number"
	operationCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/operation"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/text"
	"github.com/rs/zerolog/log"
)

func (s *Serialized) UnmarshalJSON(data []byte) error {
	// confirm that given data is not nil
	if data == nil {
		err := ErrInvalidSerializedCriteria{Reasons: []string{"json criterion data is nil"}}
		log.Error().Err(err)
		return err
	}

	// unmarshal into Serialized section of Serialized
	if err := json.Unmarshal(data, &s.Serialized); err != nil {
		err = ErrUnmarshal{Reasons: []string{"json unmarshal", err.Error()}}
		log.Error().Err(err)
		return err
	}

	// parse each key value pair of the serialized criterion to a valid criterion
	s.Criteria = make([]searchCriterion.Criterion, 0)
	for operationOrField, value := range s.Serialized {
		parsedCriterion, err := parse(operationOrField, value)
		if err != nil {
			log.Error().Err(err).Msg("parsing criterion")
			return err
		}
		s.Criteria = append(s.Criteria, parsedCriterion)
	}

	return nil
}

type typeHolder struct {
	Type searchCriterion.Type `json:"type"`
}

func parse(operationOrField string, value json.RawMessage) (searchCriterion.Criterion, error) {
	var parsedCriterion searchCriterion.Criterion

	switch operationOrField {
	case searchCriterion.OROperator:
		var orSerializedArray []map[string]json.RawMessage
		if err := json.Unmarshal(value, &orSerializedArray); err != nil {
			err = ErrUnmarshal{Reasons: []string{
				"or array unmarshal",
				err.Error(),
			}}
			log.Error().Err(err)
			return nil, err
		}
		var orCriterion operationCriterion.Or
		orCriterion.Criteria = make([]searchCriterion.Criterion, 0)
		for _, serializedCriterion := range orSerializedArray {
			if len(serializedCriterion) > 1 {
				var andCriterion = operationCriterion.And{Criteria: make([]searchCriterion.Criterion, 0)}
				for operationOrField, value := range serializedCriterion {
					crit, err := parse(operationOrField, value)
					if err != nil {
						err = ErrUnmarshal{Reasons: []string{
							"element in or",
							err.Error(),
						}}
						return nil, err
					}
					andCriterion.Criteria = append(andCriterion.Criteria, crit)
				}
				orCriterion.Criteria = append(orCriterion.Criteria, andCriterion)
			} else {
				for operationOrField, value := range serializedCriterion {
					crit, err := parse(operationOrField, value)
					if err != nil {
						err = ErrUnmarshal{Reasons: []string{
							"element in or",
							err.Error(),
						}}
						return nil, err
					}
					orCriterion.Criteria = append(orCriterion.Criteria, crit)
				}
			}
		}
		parsedCriterion = orCriterion

	default:
		th := new(typeHolder)
		if err := json.Unmarshal(value, th); err != nil {
			err = ErrUnmarshal{Reasons: []string{
				"criterion object unmarshal",
				err.Error(),
			}}
			log.Error().Err(err)
			return nil, err
		}
		switch th.Type {
		case searchCriterion.TextSubstringCriterionType:
			var typedCriterion stringCriterion.Substring
			if err := json.Unmarshal(value, &typedCriterion); err != nil {
				err = ErrUnmarshal{Reasons: []string{
					"string substring",
					err.Error(),
				}}
				return nil, err
			}
			typedCriterion.Field = operationOrField
			parsedCriterion = typedCriterion

		case searchCriterion.TextExactCriterionType:
			var typedCriterion stringCriterion.Exact
			if err := json.Unmarshal(value, &typedCriterion); err != nil {
				err = ErrUnmarshal{Reasons: []string{
					"string exact",
					err.Error(),
				}}
				return nil, err
			}
			typedCriterion.Field = operationOrField
			parsedCriterion = typedCriterion

		case searchCriterion.NumberRangeCriterionType:
			var typedCriterion numberCriterion.Range
			if err := json.Unmarshal(value, &typedCriterion); err != nil {
				err = ErrUnmarshal{Reasons: []string{
					"number range",
					err.Error(),
				}}
				return nil, err
			}
			typedCriterion.Field = operationOrField
			parsedCriterion = typedCriterion

		case searchCriterion.NumberExactCriterionType:
			var typedCriterion numberCriterion.Exact
			if err := json.Unmarshal(value, &typedCriterion); err != nil {
				err = ErrUnmarshal{Reasons: []string{
					"number exact",
					err.Error(),
				}}
				return nil, err
			}
			typedCriterion.Field = operationOrField
			parsedCriterion = typedCriterion

		default:
			return nil, ErrInvalidSerializedCriteria{Reasons: []string{
				"invalid field criterion type",
				th.Type.String(),
			}}
		}
	}

	// check that parsed criterion is valid
	if err := parsedCriterion.IsValid(); err != nil {
		return nil, ErrInvalidSerializedCriteria{Reasons: []string{
			err.Error(),
		}}
	}

	return parsedCriterion, nil
}
