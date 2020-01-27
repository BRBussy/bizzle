package budget

type Budget struct {
	Month   string               `json:"month" bson:"month"`
	Year    int                  `json:"year" bson:"year"`
	Summary map[Category]float64 `json:"summary" bson:"summary"`
	Entries map[Category][]Entry `json:"entries" bson:"entries"`
}

const DateFormat = "2006-01-02"

type Entry struct {
	Date        string   `json:"date" bson:"date"`
	Description string   `json:"description" bson:"description"`
	Amount      float64  `json:"amount" bson:"amount"`
	Category    Category `json:"category" bson:"category"`
}
