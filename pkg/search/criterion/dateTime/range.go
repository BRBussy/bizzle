package dateTime

import (
	"encoding/json"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	"time"
)

type Range struct {
	Field string     `json:"field"`
	Start RangeValue `json:"start"`
	End   RangeValue `json:"end"`
}

type RangeValue struct {
	Date      time.Time `json:"date"`
	Inclusive bool      `json:"inclusive"`
	Ignore    bool      `json:"ignore"`
}

func (r Range) IsValid() error {
	reasonsInvalid := make([]string, 0)

	if r.Field == "" {
		reasonsInvalid = append(reasonsInvalid, "field is blank")
	}

	if len(reasonsInvalid) > 0 {
		return criterion.ErrInvalid{Reasons: reasonsInvalid}
	}

	return nil
}

func (r Range) Type() criterion.Type {
	return criterion.DateTimeRangeCriterionType
}

func (r Range) ToFilter() map[string]interface{} {
	startOperator := "$gt"
	if !r.Start.Ignore {
		if r.Start.Inclusive {
			startOperator = "$gte"
		}
	}

	endOperator := "$lt"
	if !r.End.Ignore {
		if r.End.Inclusive {
			endOperator = "$lte"
		}
	}

	if !r.Start.Ignore && r.End.Ignore {
		// only consider start date
		return map[string]interface{}{r.Field: map[string]interface{}{startOperator: r.Start.Date}}
	} else if r.Start.Ignore && !r.End.Ignore {
		// only consider end date
		return map[string]interface{}{r.Field: map[string]interface{}{endOperator: r.End.Date}}
	} else if !(r.Start.Ignore || r.End.Ignore) {
		// consider both start and end dates
		return map[string]interface{}{r.Field: map[string]interface{}{
			startOperator: r.Start.Date,
			endOperator:   r.End.Date,
		}}
	}
	// consider neither
	return map[string]interface{}{r.Field: map[string]interface{}{}}
}

func (r Range) ToJSON() (string, json.RawMessage, error) {
	data, err := json.Marshal(struct {
		Type  string     `json:"type"`
		Start RangeValue `json:"start"`
		End   RangeValue `json:"end"`
	}{
		Type:  r.Type().String(),
		Start: r.Start,
		End:   r.End,
	})
	return r.Field, data, err
}
