package wrapped

import (
	"encoding/json"
	"errors"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/BRBussy/bizzle/pkg/search/criterion/operation/or"
	"github.com/BRBussy/bizzle/pkg/search/criterion/substring"
	"github.com/rs/zerolog/log"
)

type Wrapped struct {
	Type      searchCriterion.Type      `json:"type"`
	Value     json.RawMessage           `json:"value"`
	Criterion searchCriterion.Criterion `json:"-"`
}

type Slice []Wrapped

func (s Slice) ToCriteria() []searchCriterion.Criterion {
	criteria := make([]searchCriterion.Criterion, 0)
	for _, wrappedCriterion := range s {
		criteria = append(criteria, wrappedCriterion.Criterion)
	}
	return criteria
}

type unmarshalHolder struct {
	Type  searchCriterion.Type `json:"type"`
	Value json.RawMessage      `json:"value"`
}

type orWrapped struct {
	Criteria []Wrapped `json:"criteria"`
}

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
	}
	return nil
}

func (w *Wrapped) Unwrap() error {
	switch w.Type {
	case searchCriterion.Substring:
		var unmarshalledCriterion substring.Criterion
		if err := json.Unmarshal(w.Value, &unmarshalledCriterion); err != nil {
			return errors.New("unmarshalling failed: " + err.Error())
		}
		w.Criterion = unmarshalledCriterion

	case searchCriterion.OrCriterionType:
		var unmarshalledWrappedOr orWrapped
		var unmarshalledCriterion or.Criterion
		if err := json.Unmarshal(w.Value, &unmarshalledWrappedOr); err != nil {
			return errors.New("unmarshalling failed: " + err.Error())
		}
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

	if w.Criterion == nil {
		return errors.New("criterion still nil")
	}

	if err := w.Criterion.IsValid(); err != nil {
		return errors.New("criterion not valid: " + err.Error())
	}

	return nil
}

func Wrap(criterion searchCriterion.Criterion) (*Wrapped, error) {
	switch typedCriterion := criterion.(type) {
	case or.Criterion:
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
