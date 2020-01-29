package budget

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
)

type Budget struct {
	Month   string                                       `json:"month" bson:"month"`
	Year    int                                          `json:"year" bson:"year"`
	Summary map[budgetEntry.Category]float64             `json:"summary" bson:"summary"`
	Entries map[budgetEntry.Category][]budgetEntry.Entry `json:"entries" bson:"entries"`
}

const DateFormat = "2006-01-02"
