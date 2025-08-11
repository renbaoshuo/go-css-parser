package cssparser

type allowedRuleType int

const (
	qualifiedRuleTypeStyle allowedRuleType = 1 << iota
	qualifiedRuleTypeKeyframes
)

func (t allowedRuleType) Has(ruleType allowedRuleType) bool {
	return t&ruleType != 0
}

const topLevelAllowedRules = qualifiedRuleTypeStyle
