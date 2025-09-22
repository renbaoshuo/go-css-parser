package css

import (
	"testing"
)

func TestStyleRuleType(t *testing.T) {
	tests := []struct {
		name     string
		ruleType StyleRuleType
		expected string
	}{
		{"unknown rule", StyleRuleTypeUnknown, "Unknown"},
		{"at rule", StyleRuleTypeAtRule, "AtRule"},
		{"qualified rule", StyleRuleTypeQualifiedRule, "QualifiedRule"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ruleType.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestStyleRuleEquals(t *testing.T) {
	// Create test data
	decl1 := &Declaration{Property: "color", Value: "red", Important: false}
	decl2 := &Declaration{Property: "background", Value: "blue", Important: true}

	selector1 := &Selector{
		Flag: SelectorFlagContainsPseudo,
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("test"),
			},
		},
	}

	selector2 := &Selector{
		Flag: SelectorFlagContainsComplexSelector,
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchId,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("main"),
			},
		},
	}

	rule1 := &StyleRule{
		Type:         StyleRuleTypeQualifiedRule,
		Selectors:    []*Selector{selector1},
		Declarations: []*Declaration{decl1},
		Rules:        []*GenericRule{},
	}

	rule2 := &StyleRule{
		Type:         StyleRuleTypeQualifiedRule,
		Selectors:    []*Selector{selector1},
		Declarations: []*Declaration{decl1},
		Rules:        []*GenericRule{},
	}

	rule3 := &StyleRule{
		Type:         StyleRuleTypeAtRule,
		Selectors:    []*Selector{selector1},
		Declarations: []*Declaration{decl1},
		Rules:        []*GenericRule{},
	}

	rule4 := &StyleRule{
		Type:         StyleRuleTypeQualifiedRule,
		Selectors:    []*Selector{selector2},
		Declarations: []*Declaration{decl1},
		Rules:        []*GenericRule{},
	}

	rule5 := &StyleRule{
		Type:         StyleRuleTypeQualifiedRule,
		Selectors:    []*Selector{selector1},
		Declarations: []*Declaration{decl2},
		Rules:        []*GenericRule{},
	}

	tests := []struct {
		name     string
		rule1    *StyleRule
		rule2    *StyleRule
		expected bool
	}{
		{"identical rules", rule1, rule2, true},
		{"same object", rule1, rule1, true},
		{"different types", rule1, rule3, false},
		{"different selectors", rule1, rule4, false},
		{"different declarations", rule1, rule5, false},
		{"nil comparison", rule1, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rule1.Equals(tt.rule2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStyleRuleEqualsWithDifferentCounts(t *testing.T) {
	decl1 := &Declaration{Property: "color", Value: "red", Important: false}
	decl2 := &Declaration{Property: "background", Value: "blue", Important: true}

	selector1 := &Selector{
		Flag: SelectorFlagContainsPseudo,
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("test"),
			},
		},
	}

	selector2 := &Selector{
		Flag: SelectorFlagContainsComplexSelector,
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchId,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("main"),
			},
		},
	}

	genericRule1 := &GenericRule{}
	genericRule2 := &GenericRule{}

	rule1 := &StyleRule{
		Type:         StyleRuleTypeQualifiedRule,
		Selectors:    []*Selector{selector1},
		Declarations: []*Declaration{decl1},
		Rules:        []*GenericRule{genericRule1},
	}

	rule2 := &StyleRule{
		Type:         StyleRuleTypeQualifiedRule,
		Selectors:    []*Selector{selector1, selector2},
		Declarations: []*Declaration{decl1},
		Rules:        []*GenericRule{genericRule1},
	}

	rule3 := &StyleRule{
		Type:         StyleRuleTypeQualifiedRule,
		Selectors:    []*Selector{selector1},
		Declarations: []*Declaration{decl1, decl2},
		Rules:        []*GenericRule{genericRule1},
	}

	rule4 := &StyleRule{
		Type:         StyleRuleTypeQualifiedRule,
		Selectors:    []*Selector{selector1},
		Declarations: []*Declaration{decl1},
		Rules:        []*GenericRule{genericRule1, genericRule2},
	}

	tests := []struct {
		name     string
		rule1    *StyleRule
		rule2    *StyleRule
		expected bool
	}{
		{"different selector count", rule1, rule2, false},
		{"different declaration count", rule1, rule3, false},
		{"different rule count", rule1, rule4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rule1.Equals(tt.rule2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGenericRuleEquals(t *testing.T) {
	rule1 := &GenericRule{}
	rule2 := &GenericRule{}

	tests := []struct {
		name     string
		rule1    *GenericRule
		rule2    *GenericRule
		expected bool
	}{
		{"both non-nil", rule1, rule2, true},
		{"same object", rule1, rule1, true},
		{"one nil", rule1, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rule1.Equals(tt.rule2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
