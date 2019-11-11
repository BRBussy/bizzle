package exercise

type Exercise interface {
	Type() Type              // Returns the Type of the exercise
	ToJSON() ([]byte, error) // Returns json marshalled version of exercise
}
