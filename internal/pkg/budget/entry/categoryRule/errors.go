package categoryRule

type ErrCouldNotClassify struct {
	Reason string
}

func (e ErrCouldNotClassify) Error() string {
	return "could not classify: " + e.Reason
}

type ErrDefaultOtherBudgetEntryCategoryRuleNotSet struct {
}

func (e ErrDefaultOtherBudgetEntryCategoryRuleNotSet) Error() string {
	return "default other budget category rule not set"
}
