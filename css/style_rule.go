package css

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
	Type         StyleRuleType  // Type of the rule (AtRule or QualifiedRule)
	Selectors    []*Selector    // Selectors for the style rule
	Declarations []*Declaration // CSS declarations
	Rules        []*GenericRule // Child rules
}

// Equals compares two StyleRule instances
func (sr *StyleRule) Equals(other *StyleRule) bool {
	if other == nil {
		return false
	}

	if sr.Type != other.Type ||
		len(sr.Selectors) != len(other.Selectors) ||
		len(sr.Declarations) != len(other.Declarations) ||
		len(sr.Rules) != len(other.Rules) {
		return false
	}

	// Compare selectors
	for i, sel := range sr.Selectors {
		if !sel.Equals(other.Selectors[i]) {
			return false
		}
	}

	// Compare declarations
	for i, decl := range sr.Declarations {
		if !decl.Equals(other.Declarations[i]) {
			return false
		}
	}

	// Compare child rules
	for i, rule := range sr.Rules {
		if !rule.Equals(other.Rules[i]) {
			return false
		}
	}

	return true
}

// GenericRule represents a generic CSS rule
type GenericRule struct {
	// Placeholder for now - will be expanded as needed
}

// Equals compares two GenericRule instances
func (gr *GenericRule) Equals(other *GenericRule) bool {
	return other != nil
}
