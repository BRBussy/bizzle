package reasonInvalid

import (
	"fmt"
	"strings"
)

// ReasonInvalid is a reason that an entity is invalid
type ReasonInvalid struct {
	Field string      `json:"field"`
	Type  Type        `json:"type"`
	Help  string      `json:"help"`
	Data  interface{} `json:"data"`
}

func (r ReasonInvalid) String() string {
	return fmt.Sprintf("%s - %s", r.Field, r.Type)
}

// ReasonsInvalid is list of reasons why an entity is invalid
type ReasonsInvalid []ReasonInvalid

func (r ReasonsInvalid) String() string {
	reasonsInvalid := make([]string, 0)
	for i := range r {
		reasonsInvalid = append(reasonsInvalid, r[i].String())
	}
	return strings.Join(reasonsInvalid, ", ")
}
