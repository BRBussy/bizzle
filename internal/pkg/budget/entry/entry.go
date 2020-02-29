package entry

import (
	"time"

	"github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

// Entry is a budget entry
type Entry struct {
	ID             identifier.ID `json:"id" bson:"id"`
	OwnerID        identifier.ID `validate:"required" json:"ownerID" bson:"ownerID"`
	Date           time.Time     `validate:"required" json:"date" bson:"date"`
	Description    string        `validate:"required" json:"description" bson:"description"`
	Amount         float64       `validate:"required" json:"amount" bson:"amount"`
	CategoryRuleID identifier.ID `json:"categoryRuleID" bson:"categoryRuleID"`
}

// UpperDuplicateMargin is the amount an entry's value can be above anothers and still be considered a dupilicate
const UpperDuplicateMargin float64 = 2

// LowerDuplicateMargin is the amount an entry's value can be below anothers and still be considered a dupilicate
const LowerDuplicateMargin float64 = 2

// ExactDuplicate is used to compare a budget entry with another to determine if they are exact duplicates
func (e Entry) ExactDuplicate(be Entry) bool {
	if !(e.Date.Day() == be.Date.Day() &&
		e.Date.Year() == be.Date.Year() &&
		e.Date.Month() == be.Date.Month()) {
		// if day, month and year are not the same these are not exact duplicates
		return false
	}

	if e.Description != be.Description {
		// if descriptions are not the same then these are not exact duplicates
		return false
	}

	// check if if amount falls within limit
	upperLimit := e.Amount + UpperDuplicateMargin
	lowerLimit := e.Amount - LowerDuplicateMargin
	if be.Amount < lowerLimit || be.Amount > upperLimit {
		return false
	}

	return true
}

// SuspectedDuplicate is used to compare a budget entry with another to determine if they are suspected duplicates
func (e Entry) SuspectedDuplicate(be Entry) bool {
	if !(e.Date.Day() == be.Date.Day() &&
		e.Date.Year() == be.Date.Year() &&
		e.Date.Month() == be.Date.Month()) {
		// if day, month and year are not the same these are not exact duplicates
		return false
	}

	// check if if amount falls within limit
	upperLimit := e.Amount + UpperDuplicateMargin
	lowerLimit := e.Amount - LowerDuplicateMargin
	if be.Amount < lowerLimit || be.Amount > upperLimit {
		return false
	}

	return true
}

// CompositeEntry is used to populate budget entry with corresponding category rule
type CompositeEntry struct {
	Entry        `bson:"inline"`
	CategoryRule categoryRule.CategoryRule `json:"categoryRule"`
}
