package budget

type Entry struct {
	Date        int64    `json:"date" bson:"date"`
	Description string   `json:"description" bson:"description"`
	Amount      float64  `json:"amount" bson:"amount"`
	Category    Category `json:"category" bson:"category"`
	Identifier  string   `json:"identifier" bson:"identifier"`
}
