package criteria

import (
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
	operationCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/operation"
)

// Compare can be used to compare to criteria arrays of criterion
func Compare(a, b []searchCriterion.Criterion) bool {
	// check lengths
	if len(a) != len(b) {
		return false
	}

	// for every element in a
nextA:
	for ia := range a {
		// look through b for a match
		for ib := range b {
			// if a match is found go to next element ia
			switch typedA := a[ia].(type) {
			case operationCriterion.And:
				if CompareANDCriterion(typedA, b[ib]) {
					continue nextA
				}
			case operationCriterion.Or:
				if CompareORCriterion(typedA, b[ib]) {
					continue nextA
				}
			default:
				if a[ia] == b[ib] {
					continue nextA
				}
			}
		}
		// if execution reaches here ia was not found in b
		return false
	}
	// if execution reaches here every ia was found in b
	return true
}

func CompareANDCriterion(a operationCriterion.And, b searchCriterion.Criterion) bool {
	// check that a and b are both and criterion
	typedB, ok := b.(operationCriterion.And)
	if !ok {
		return false
	}
	// check lengths of a and b are the same
	if len(a.Criteria) != len(typedB.Criteria) {
		return false
	}
	return Compare(a.Criteria, typedB.Criteria)
}

func CompareORCriterion(a operationCriterion.Or, b searchCriterion.Criterion) bool {
	// check that a and b are both or criterion
	typedB, ok := b.(operationCriterion.Or)
	if !ok {
		return false
	}
	// check lengths of a and b are the same
	if len(a.Criteria) != len(typedB.Criteria) {
		return false
	}
	return Compare(a.Criteria, typedB.Criteria)
}
