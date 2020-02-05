package budget

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	"time"
)

type Budget struct {
	StartDate time.Time                      `json:"startDate" bson:"startDate"`
	EndDate   time.Time                      `json:"endDate" bson:"endDate"`
	Summary   map[string]float64             `json:"summary" bson:"summary"`
	Entries   map[string][]budgetEntry.Entry `json:"entries" bson:"entries"`
}