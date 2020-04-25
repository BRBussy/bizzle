package identifier

type Identifier interface {
	IsValid() error                   // Returns the validity of the UserIdentifier
	Type() Type                       // Returns the Type of the UserIdentifier
	ToFilter() map[string]interface{} // Returns a map filter to use to query the databases
	ToJSON() ([]byte, error)          // Returns json marshalled version of identifier
}
