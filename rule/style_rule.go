package rule

type StyleRuleType int

const (
	StyleRuleTypeUnknown StyleRuleType = iota

	StyleRuleTypeAtRule
	StyleRuleTypeQualifiedRule
)

func (srt StyleRuleType) String() string {
	switch srt {
	case StyleRuleTypeAtRule:
		return "AtRule"
	case StyleRuleTypeQualifiedRule:
		return "QualifiedRule"
	default:
		return "Unknown"
	}
}

// ------

type StyleRule struct {
	Type  StyleRuleType // Type of the rule (AtRule or QualifiedRule)
	Rules []*Rule
}
