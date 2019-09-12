package criteria

import (
	"encoding/json"
	"errors"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	numberCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/number"
	operationCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/operation"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/string"
	"github.com/rs/zerolog/log"
)

type Serialized struct {
	Serialized map[string]json.RawMessage
	Criteria   []searchCriterion.Criterion
}

func (s *Serialized) UnmarshalJSON(data []byte) error {
	// unmarshal into serialized section of Serialized
	if err := json.Unmarshal(data, &s.Serialized); err != nil {
		log.Error().Err(err).Msg("unmarshalling serialized criterion")
		return errors.New("unmarshalling failed: " + err.Error())
	}

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

type jsonObject []map[string]json.RawMessage

func parse(operationOrField string, value json.RawMessage) (searchCriterion.Criterion, error) {
	var parsedCriterion searchCriterion.Criterion

	switch operationOrField {
	case searchCriterion.OROperator:
		var oh jsonObject
		if err := json.Unmarshal(value, &oh); err != nil {
			log.Error().Err(err).Msg("unmarshalling wrapped criterion")
			return nil, errors.New("unmarshalling failed: " + err.Error())
		}
		var orCriterion operationCriterion.Or
		orCriterion.Criteria = make([]searchCriterion.Criterion, 0)
		for _, serializedCriterion := range oh {
			if len(serializedCriterion) > 1 {
				var andCriterion = operationCriterion.And{Criteria: make([]searchCriterion.Criterion, 0)}
				for operationOrField, value := range serializedCriterion {
					crit, err := parse(operationOrField, value)
					if err != nil {
						log.Error().Err(err).Msg("parsing an or criterion")
						return nil, errors.New("parsing an or criterion: " + err.Error())
					}
					andCriterion.Criteria = append(andCriterion.Criteria, crit)
				}
				orCriterion.Criteria = append(orCriterion.Criteria, andCriterion)
			} else {
				for operationOrField, value := range serializedCriterion {
					crit, err := parse(operationOrField, value)
					if err != nil {
						log.Error().Err(err).Msg("parsing an or criterion")
						return nil, errors.New("parsing an or criterion: " + err.Error())
					}
					orCriterion.Criteria = append(orCriterion.Criteria, crit)
				}
			}
		}
		parsedCriterion = orCriterion

	default:
		th := new(typeHolder)
		if err := json.Unmarshal(value, th); err != nil {
			log.Error().Err(err).Msg("unmarshalling wrapped criterion")
			return nil, errors.New("unmarshalling failed: " + err.Error())
		}
		switch th.Type {
		case searchCriterion.StringSubstringCriterionType:
			var typedCriterion stringCriterion.Substring
			if err := json.Unmarshal(value, &typedCriterion); err != nil {
				return nil, errors.New("unmarshalling failed: " + err.Error())
			}
			typedCriterion.Field = operationOrField
			parsedCriterion = typedCriterion

		case searchCriterion.StringExactCriterionType:
			var typedCriterion stringCriterion.Exact
			if err := json.Unmarshal(value, &typedCriterion); err != nil {
				return nil, errors.New("unmarshalling failed: " + err.Error())
			}
			typedCriterion.Field = operationOrField
			parsedCriterion = typedCriterion

		case searchCriterion.NumberRangeCriterionType:
			var typedCriterion numberCriterion.Range
			if err := json.Unmarshal(value, &typedCriterion); err != nil {
				return nil, errors.New("unmarshalling failed: " + err.Error())
			}
			typedCriterion.Field = operationOrField
			parsedCriterion = typedCriterion

		case searchCriterion.NumberExactCriterionType:
			var typedCriterion numberCriterion.Exact
			if err := json.Unmarshal(value, &typedCriterion); err != nil {
				return nil, errors.New("unmarshalling failed: " + err.Error())
			}
			typedCriterion.Field = operationOrField
			parsedCriterion = typedCriterion

		default:
			return nil, errors.New("invalid")
		}
	}

	// check that parsed criterion is valid
	if parsedCriterion == nil {
		return nil, errors.New("criterion still nil")
	}
	if err := parsedCriterion.IsValid(); err != nil {
		return nil, errors.New("criterion not valid: " + err.Error())
	}

	return parsedCriterion, nil
}
