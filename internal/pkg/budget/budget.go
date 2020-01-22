package budget

import "time"

type Budget struct {
	Month   time.Month           `json:"month" bson:"month"`
	Year    string               `json:"int" bson:"year"`
	Summary map[Category]float64 `json:"summary" bson:"summary"`
	Entries map[Category][]Entry `json:"entries" bson:"entries"`
}

type Entry struct {
	Date        int64    `json:"date" bson:"date"`
	Description string   `json:"description" bson:"description"`
	Amount      float64  `json:"amount" bson:"amount"`
	Category    Category `json:"category" bson:"category"`
	Identifier  string   `json:"identifier" bson:"identifier"`
}
