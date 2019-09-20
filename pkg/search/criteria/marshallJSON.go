package criteria

import (
	"encoding/json"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	operationCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/operation"
)

func (s Serialized) MarshalJSON() ([]byte, error) {
	serializedCriteria := make(map[string]json.RawMessage)
	for _, criterion := range s.Criteria {
		field, jsonMessage, err := serialize(criterion)
		if err != nil {
			return nil, err
		}
		serializedCriteria[field] = jsonMessage
	}
	return json.Marshal(serializedCriteria)
}

func serialize(criterion searchCriterion.Criterion) (string, json.RawMessage, error) {
	switch typedCriterion := criterion.(type) {
	case operationCriterion.And:
		return "", nil, ErrInvalidSerializedCriteria{Reasons: []string{"and only allowed if contained within or"}}
	case operationCriterion.Or:
		orSerializedArray := make([]map[string]json.RawMessage, 0)
		for _, c := range typedCriterion.Criteria {
			switch typedOrElement := c.(type) {
			case operationCriterion.And:
				andSerializedObject := make(map[string]json.RawMessage)
				for _, c := range typedOrElement.Criteria {
					field, jsonMessage, err := serialize(c)
					if err != nil {
						return "", nil, ErrMarshal{Reasons: []string{
							"or element",
							err.Error(),
						}}
					}
					andSerializedObject[field] = jsonMessage
				}
				orSerializedArray = append(orSerializedArray, andSerializedObject)

			default:
				field, jsonMessage, err := serialize(c)
				if err != nil {
					return "", nil, ErrMarshal{Reasons: []string{
						"or element",
						err.Error(),
					}}
				}
				orSerializedArray = append(orSerializedArray, map[string]json.RawMessage{field: jsonMessage})
			}
		}
		orSerializedArrayData, err := json.Marshal(orSerializedArray)
		return searchCriterion.OROperator, orSerializedArrayData, err

	default:
		field, jsonMessageData, err := criterion.ToJSON()
		if err != nil {
			return "", nil, ErrMarshal{Reasons: []string{err.Error()}}
		}
		return field, jsonMessageData, nil
	}
}
