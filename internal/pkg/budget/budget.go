package budget

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	budgetEntryCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	"time"
)

type Budget struct {
	StartDate time.Time                      `json:"startDate" bson:"startDate"`
	EndDate   time.Time                      `json:"endDate" bson:"endDate"`
	Summary   map[string]CategoryTotal       `json:"summary" bson:"summary"`
	Entries   map[string][]budgetEntry.Entry `json:"entries" bson:"entries"`
	TotalIn   CompareTotal                   `json:"totalIn" bson:"totalIn"`
	TotalOut  CompareTotal                   `json:"totalOut" bson:"totalOut"`
	Net       float64                        `json:"net" bson:"net"`
}

type CategoryTotal struct {
	BudgetEntryCategoryRule budgetEntryCategoryRule.CategoryRule `json:"budgetEntryCategoryRule" bson:"budgetEntryCategoryRule"`
	Amount                  float64                              `json:"amount" bson:"amount"`
}

type CompareTotal struct {
	Expected float64 `json:"expected" bson:"expected"`
	Actual   float64 `json:"actual" bson:"actual"`
}
