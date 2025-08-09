package cssparser

type RuleType int

const (
	RuleTypeUnknown RuleType = iota

	RuleTypeAtRule
	RuleTypeQualifiedRule
)

func (rt RuleType) String() string {
	switch rt {
	case RuleTypeAtRule:
		return "AtRule"
	case RuleTypeQualifiedRule:
		return "QualifiedRule"
	default:
		return "Unknown"
	}
}

// ------

type Rule struct {
	Type RuleType // Type of the rule (AtRule or QualifiedRule)
}
