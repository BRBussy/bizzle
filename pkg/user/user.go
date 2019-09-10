package user

type User struct {
	Id string `json:"id" bson:"id"`

	// Personal Details
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`

	// System Details
	EmailAddress string   `json:"emailAddress" bson:"emailAddress"`
	Password     []byte   `json:"password" bson:"password"`
	Roles        []string `json:"roles" bson:"roles"`

	Registered bool `json:"registered" bson:"registered"`
}
