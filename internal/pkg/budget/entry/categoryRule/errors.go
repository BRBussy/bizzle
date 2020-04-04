package categoryRule

type ErrDefaultOtherBudgetEntryCategoryRuleNotSet struct {
}

func (e ErrDefaultOtherBudgetEntryCategoryRuleNotSet) Error() string {
	return "default other budget category rule not set"
}
