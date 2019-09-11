package wrapped

import (
	"encoding/json"
	"errors"
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/BRBussy/bizzle/pkg/search/criterion/substring"
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

func (w *Wrapped) UnmarshalJSON(data []byte) error {
	// unmarshal into intermediate holder to avoid an infinite recursion loop
	var holder unmarshalHolder
	if err := json.Unmarshal(data, &holder); err != nil {
		return errors.New("unmarshalling failed: " + err.Error())
	}
	// update data on wrapped criterion
	w.Type = holder.Type
	w.Value = holder.Value

	// unwrap the wrapped criterion
	return w.Unwrap()
}

func (w Wrapped) Unwrap() error {
	switch w.Type {
	case searchCriterion.Substring:
		var unmarshalledCriterion substring.Criterion
		if err := json.Unmarshal(w.Value, &unmarshalledCriterion); err != nil {
			return errors.New("unmarshalling failed: " + err.Error())
		}
		w.Criterion = unmarshalledCriterion
	default:
		return errors.New("invalid/unsupported criterion type")
	}

	if w.Criterion == nil {
		return errors.New("criterion still nil")
	}

	if err := w.Criterion.IsValid(); err != nil {
		return errors.New("criterion not valid: " + err.Error())
	}

	return nil
}
