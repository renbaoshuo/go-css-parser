package css_parser

import (
	"fmt"
	"slices"
	"strings"
)

const (
	indentSpace = 2
)

type CssRuleKind int

const (
	QualifiedRule CssRuleKind = iota
	AtRule
)

var atRulesWithRulesBlock = []string{
	"@document", "@font-feature-values", "@keyframes",
	"@media", "@supports", "@container", "@layer",
}

func (k CssRuleKind) String() string {
	switch k {
	case QualifiedRule:
		return "QualifiedRule"
	case AtRule:
		return "AtRule"
	default:
		return "Unknown"
	}
}

type CssRule struct {
	Kind         CssRuleKind
	Name         string            // At Rule name (eg: "@media")
	Prelude      string            // Raw prelude: https://github.com/csstree/csstree/discussions/168
	Selectors    []string          // Qualified Rule selectors parsed from prelude
	Declarations []*CssDeclaration // Style properties
	Rules        []*CssRule        // At Rule embedded rules
	EmbedLevel   int               // Current rule embedding level
}

func NewCssRule(kind CssRuleKind) *CssRule {
	return &CssRule{
		Kind: kind,
	}
}

func (r *CssRule) Equal(o *CssRule) bool {
	if (r.Kind != o.Kind) ||
		(r.Prelude != o.Prelude) ||
		(r.Name != o.Name) {
		return false
	}

	if (len(r.Selectors) != len(o.Selectors)) ||
		(len(r.Declarations) != len(o.Declarations)) ||
		(len(r.Rules) != len(o.Rules)) {
		return false
	}

	for i, sel := range r.Selectors {
		if sel != o.Selectors[i] {
			return false
		}
	}

	for i, decl := range r.Declarations {
		if !decl.Equal(o.Declarations[i]) {
			return false
		}
	}

	for i, rule := range r.Rules {
		if !rule.Equal(o.Rules[i]) {
			return false
		}
	}

	return true
}

func (r *CssRule) EmbedsRules() bool {
	if r.Kind == AtRule {
		if slices.Contains(atRulesWithRulesBlock, r.Name) {
			return true
		}
	}

	return false
}

func (r *CssRule) String() string {
	result := ""

	if r.Kind == QualifiedRule {
		for i, sel := range r.Selectors {
			if i != 0 {
				result += ", "
			}
			result += sel
		}
	} else { // AtRule
		result += r.Name

		if r.Prelude != "" {
			if result != "" {
				result += " "
			}
			result += r.Prelude
		}
	}

	if (len(r.Declarations) == 0) && (len(r.Rules) == 0) {
		result += ";"
	} else {
		result += " {\n"

		if r.EmbedsRules() {
			for _, subRule := range r.Rules {
				result += fmt.Sprintf("%s%s\n", r.indent(), subRule.String())
			}
		} else {
			for _, decl := range r.Declarations {
				result += fmt.Sprintf("%s%s\n", r.indent(), decl.String())
			}
		}

		result += fmt.Sprintf("%s}", r.indentEndBlock())
	}

	return result
}

func (r *CssRule) indent() string {
	return strings.Repeat(" ", (r.EmbedLevel+1)*indentSpace)
}

func (r *CssRule) indentEndBlock() string {
	return strings.Repeat(" ", r.EmbedLevel*indentSpace)
}
