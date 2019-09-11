package wrapped

import (
	"encoding/json"
	"errors"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/BRBussy/bizzle/pkg/search/criterion/substring"
)

type Wrapped struct {
	Type      criterion.Type      `json:"type"`
	Value     json.RawMessage     `json:"value"`
	Criterion criterion.Criterion `json:"-"`
}

func (w *Wrapped) UnmarshalJSON(data []byte) error {
	unwrapped, err := w.Unwrap()
	if err != nil {
		return err
	}
	w.Criterion = unwrapped
	return nil
}

func (w Wrapped) Unwrap() (criterion.Criterion, error) {
	var result criterion.Criterion = nil

	switch w.Type {
	case criterion.Substring:
		var unmarshalledCriterion substring.Criterion
		if err := json.Unmarshal(w.Value, &unmarshalledCriterion); err != nil {
			return nil, errors.New("unmarshalling failed: " + err.Error())
		}
		result = unmarshalledCriterion
	default:
		return nil, errors.New("invalid/unsupported criterion type")
	}

	if result == nil {
		return nil, errors.New("criterion still nil")
	}

	if err := result.IsValid(); err != nil {
		return nil, errors.New("criterion not valid: " + err.Error())
	}

	return result, nil
}
