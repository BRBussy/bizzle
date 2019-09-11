package criterion

// Criterion
type Criterion interface {
	IsValid() error                   // Returns the validity of the Criterion
	Type() Type                       // Returns the Type of the Criterion
	ToFilter() map[string]interface{} // Returns a map filter to use to query the databases
}

// Criteria is a slice of criterion that can be converted into a filter
type Criteria []Criterion

// ToFilter returns a filter combining all of the criterion in criteria
func (c Criteria) ToFilter() map[string]interface{} {
	filters := make([]map[string]interface{}, 0)
	for _, crit := range c {
		filters = append(filters, crit.ToFilter())
	}
	return map[string]interface{}{"$and": filters}
}
