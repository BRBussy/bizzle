package wrapped

import (
	"encoding/json"
	"errors"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	numberCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/number"
	operationCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/operation"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/string"
	"github.com/rs/zerolog/log"
)

// Wrapped represents a wrapped criterion. Used in json marshalling and unmarshalling.
type Wrapped struct {
	Type      searchCriterion.Type      `json:"type"`
	Value     json.RawMessage           `json:"value"`
	Criterion searchCriterion.Criterion `json:"-"`
}

// Slice is a slice of Wrapped criterion
type Slice []Wrapped

// ToCriteria can be called on a slice of wrapped criterion to convert it into a slice of criteria
func (s Slice) ToCriteria() []searchCriterion.Criterion {
	criteria := make([]searchCriterion.Criterion, 0)
	for _, wrappedCriterion := range s {
		criteria = append(criteria, wrappedCriterion.Criterion)
	}
	return criteria
}

// unmarshalHolder is used during initial unmarshalling of a Wrapped criterion to avoid infinite recursion
type unmarshalHolder struct {
	Type  searchCriterion.Type `json:"type"`
	Value json.RawMessage      `json:"value"`
}

// orWrapped is used in an intermediate step when unmarshalling a wrapped or criterion
type orWrapped struct {
	Criteria []Wrapped `json:"criteria"`
}

// UnmarshallJSON is written for Wrapped so that it implements the Unmarshaller interface. This is called during unmarshalling.
func (w *Wrapped) UnmarshalJSON(data []byte) error {
	// unmarshal into intermediate holder to avoid an infinite recursion loop
	var holder unmarshalHolder
	if err := json.Unmarshal(data, &holder); err != nil {
		log.Error().Err(err).Msg("unmarshalling wrapped criterion")
		return errors.New("unmarshalling failed: " + err.Error())
	}

	// update data on wrapped criterion
	w.Type = holder.Type
	w.Value = holder.Value

	// unwrap the wrapped criterion
	if err := w.Unwrap(); err != nil {
		log.Error().Err(err).Msg("unwrapping wrapped criterion")
		return err
	}
	return nil
}

// UnWrap is called on a Wrapped criterion to populate it's Criterion interface field
func (w *Wrapped) Unwrap() error {
	switch w.Type {
	case searchCriterion.NumberRangeCriterionType:
		var unmarshalledCriterion numberCriterion.Range
		if err := json.Unmarshal(w.Value, &unmarshalledCriterion); err != nil {
			return errors.New("unmarshalling failed: " + err.Error())
		}
		w.Criterion = unmarshalledCriterion

	case searchCriterion.StringSubstringCriterionType:
		var unmarshalledCriterion stringCriterion.Substring
		if err := json.Unmarshal(w.Value, &unmarshalledCriterion); err != nil {
			return errors.New("unmarshalling failed: " + err.Error())
		}
		w.Criterion = unmarshalledCriterion

	case searchCriterion.StringExactCriterionType:
		var unmarshalledCriterion stringCriterion.Exact
		if err := json.Unmarshal(w.Value, &unmarshalledCriterion); err != nil {
			return errors.New("unmarshalling failed: " + err.Error())
		}
		w.Criterion = unmarshalledCriterion

	case searchCriterion.OperationOrCriterionType:
		var unmarshalledCriterion operationCriterion.Or
		// first unmarshal into a wrapped or
		var unmarshalledWrappedOr orWrapped
		if err := json.Unmarshal(w.Value, &unmarshalledWrappedOr); err != nil {
			return errors.New("unmarshalling failed: " + err.Error())
		}
		// then unwrap each wrapped criterion in the wrapped or
		unmarshalledCriterion.Criteria = make([]searchCriterion.Criterion, 0)
		for _, wrappedCrit := range unmarshalledWrappedOr.Criteria {
			if err := wrappedCrit.Unwrap(); err != nil {
				return errors.New("unwrappig an or crit failed: " + err.Error())
			}
			unmarshalledCriterion.Criteria = append(unmarshalledCriterion.Criteria, wrappedCrit.Criterion)
		}
		w.Criterion = unmarshalledCriterion

	default:
		return errors.New("invalid/unsupported criterion type: " + w.Type.String())
	}

	// check that unmarshalled criterion is valid
	if w.Criterion == nil {
		return errors.New("criterion still nil")
	}
	if err := w.Criterion.IsValid(); err != nil {
		return errors.New("criterion not valid: " + err.Error())
	}

	return nil
}

// Wrap is used to wrap a Criterion before passing into an untyped environment
func Wrap(criterion searchCriterion.Criterion) (*Wrapped, error) {
	switch typedCriterion := criterion.(type) {
	case operationCriterion.Or:
		wrappedOr := orWrapped{Criteria: make([]Wrapped, 0)}
		for _, critToWrap := range typedCriterion.Criteria {
			wrappedCrit, err := Wrap(critToWrap)
			if err != nil {
				return nil, errors.New("wrapping an or crit: " + err.Error())
			}
			wrappedOr.Criteria = append(wrappedOr.Criteria, *wrappedCrit)
		}
		value, err := json.Marshal(wrappedOr)
		if err != nil {
			return nil, errors.New("json marshalling: " + err.Error())
		}
		return &Wrapped{
			Type:  criterion.Type(),
			Value: value,
		}, nil

	default:
		value, err := json.Marshal(criterion)
		if err != nil {
			return nil, errors.New("json marshalling: " + err.Error())
		}
		return &Wrapped{
			Type:  criterion.Type(),
			Value: value,
		}, nil
	}
}
