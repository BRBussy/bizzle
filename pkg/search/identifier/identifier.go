package identifier

// Identifier
type Identifier interface {
	IsValid() error                   // Returns the validity of the identifier
	Type() Type                       // Returns the Type of the identifier
	ToFilter() map[string]interface{} // Returns a map filter to use to query the database
}
