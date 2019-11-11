package exercise

type Type string

func (t Type) String() string {
	return string(t)
}

// arm exercises
const ArmCurlExerciseType Type = "ArmCurl"
