package css

import (
	"testing"
)

func TestSelectorListFlagType(t *testing.T) {
	var flag SelectorListFlagType

	// Test initial state
	if flag.Has(SelectorFlagContainsPseudo) {
		t.Error("expected flag to not have SelectorFlagContainsPseudo initially")
	}

	// Test setting flags
	flag.Set(SelectorFlagContainsPseudo)
	if !flag.Has(SelectorFlagContainsPseudo) {
		t.Error("expected flag to have SelectorFlagContainsPseudo after setting")
	}

	// Test setting multiple flags
	flag.Set(SelectorFlagContainsComplexSelector)
	if !flag.Has(SelectorFlagContainsPseudo) {
		t.Error("expected flag to still have SelectorFlagContainsPseudo")
	}
	if !flag.Has(SelectorFlagContainsComplexSelector) {
		t.Error("expected flag to have SelectorFlagContainsComplexSelector")
	}

	// Test flag that wasn't set
	if flag.Has(SelectorFlagContainsScopeOrParent) {
		t.Error("expected flag to not have SelectorFlagContainsScopeOrParent")
	}
}

func TestSelectorAppend(t *testing.T) {
	sel := &Selector{}

	// Test initial state
	if len(sel.Selectors) != 0 {
		t.Errorf("expected 0 selectors, got %d", len(sel.Selectors))
	}

	// Test appending single selector
	simple1 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     NewSelectorData("test"),
	}
	sel.Append(simple1)

	if len(sel.Selectors) != 1 {
		t.Errorf("expected 1 selector, got %d", len(sel.Selectors))
	}
	if sel.Selectors[0] != simple1 {
		t.Error("expected appended selector to match")
	}

	// Test appending multiple selectors
	simple2 := &SimpleSelector{
		Match:    SelectorMatchId,
		Relation: SelectorRelationDescendant,
		Data:     NewSelectorData("main"),
	}
	simple3 := &SimpleSelector{
		Match:    SelectorMatchTag,
		Relation: SelectorRelationChild,
		Data:     NewSelectorDataTag("", "div"),
	}
	sel.Append(simple2, simple3)

	if len(sel.Selectors) != 3 {
		t.Errorf("expected 3 selectors, got %d", len(sel.Selectors))
	}

	// Test appending empty slice
	sel.Append()
	if len(sel.Selectors) != 3 {
		t.Errorf("expected 3 selectors after empty append, got %d", len(sel.Selectors))
	}
}

func TestSelectorPrepend(t *testing.T) {
	sel := &Selector{}

	simple1 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     NewSelectorData("first"),
	}
	simple2 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     NewSelectorData("second"),
	}

	sel.Append(simple1)
	sel.Prepend(simple2)

	if len(sel.Selectors) != 2 {
		t.Errorf("expected 2 selectors, got %d", len(sel.Selectors))
	}

	// Check order: simple2 should be first
	if sel.Selectors[0] != simple2 {
		t.Error("expected prepended selector to be first")
	}
	if sel.Selectors[1] != simple1 {
		t.Error("expected original selector to be second")
	}

	// Test prepending nil
	sel.Prepend(nil)
	if len(sel.Selectors) != 2 {
		t.Errorf("expected 2 selectors after prepending nil, got %d", len(sel.Selectors))
	}
}

func TestSelectorInsertBefore(t *testing.T) {
	sel := &Selector{}

	simple1 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     NewSelectorData("first"),
	}
	simple2 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     NewSelectorData("second"),
	}
	simple3 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     NewSelectorData("third"),
	}

	sel.Append(simple1, simple2)

	// Insert at beginning
	sel.InsertBefore(0, simple3)
	if len(sel.Selectors) != 3 {
		t.Errorf("expected 3 selectors, got %d", len(sel.Selectors))
	}
	if sel.Selectors[0] != simple3 {
		t.Error("expected inserted selector at beginning")
	}

	// Insert in middle
	simple4 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     NewSelectorData("fourth"),
	}
	sel.InsertBefore(2, simple4)
	if len(sel.Selectors) != 4 {
		t.Errorf("expected 4 selectors, got %d", len(sel.Selectors))
	}
	if sel.Selectors[2] != simple4 {
		t.Error("expected inserted selector at index 2")
	}

	// Test invalid index
	initialLen := len(sel.Selectors)
	sel.InsertBefore(-1, simple1)
	if len(sel.Selectors) != initialLen {
		t.Error("expected no change for negative index")
	}

	sel.InsertBefore(100, simple1)
	if len(sel.Selectors) != initialLen {
		t.Error("expected no change for out-of-bounds index")
	}

	// Test inserting nil
	sel.InsertBefore(1, nil)
	if len(sel.Selectors) != initialLen {
		t.Error("expected no change when inserting nil")
	}
}

func TestSelectorString(t *testing.T) {
	tests := []struct {
		name     string
		selector *Selector
		expected string
	}{
		{
			name: "empty selector",
			selector: &Selector{
				Selectors: []*SimpleSelector{},
			},
			expected: "",
		},
		{
			name: "single class selector",
			selector: &Selector{
				Selectors: []*SimpleSelector{
					{
						Match:    SelectorMatchClass,
						Relation: SelectorRelationSubSelector,
						Data:     NewSelectorData("test"),
					},
				},
			},
			expected: ".test",
		},
		{
			name: "compound selector",
			selector: &Selector{
				Selectors: []*SimpleSelector{
					{
						Match:    SelectorMatchTag,
						Relation: SelectorRelationSubSelector,
						Data:     NewSelectorDataTag("", "div"),
					},
					{
						Match:    SelectorMatchClass,
						Relation: SelectorRelationSubSelector,
						Data:     NewSelectorData("container"),
					},
					{
						Match:    SelectorMatchId,
						Relation: SelectorRelationSubSelector,
						Data:     NewSelectorData("main"),
					},
				},
			},
			expected: "div.container#main",
		},
		{
			name: "complex selector with descendant",
			selector: &Selector{
				Selectors: []*SimpleSelector{
					{
						Match:    SelectorMatchClass,
						Relation: SelectorRelationSubSelector,
						Data:     NewSelectorData("parent"),
					},
					{
						Match:    SelectorMatchClass,
						Relation: SelectorRelationDescendant,
						Data:     NewSelectorData("child"),
					},
				},
			},
			expected: ".parent .child",
		},
		{
			name: "complex selector with child combinator",
			selector: &Selector{
				Selectors: []*SimpleSelector{
					{
						Match:    SelectorMatchTag,
						Relation: SelectorRelationSubSelector,
						Data:     NewSelectorDataTag("", "ul"),
					},
					{
						Match:    SelectorMatchTag,
						Relation: SelectorRelationChild,
						Data:     NewSelectorDataTag("", "li"),
					},
				},
			},
			expected: "ul > li",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.selector.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSelectorEquals(t *testing.T) {
	sel1 := &Selector{
		Flag: SelectorFlagContainsPseudo,
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("test"),
			},
		},
	}

	sel2 := &Selector{
		Flag: SelectorFlagContainsPseudo,
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("test"),
			},
		},
	}

	sel3 := &Selector{
		Flag: SelectorFlagContainsComplexSelector,
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("test"),
			},
		},
	}

	sel4 := &Selector{
		Flag: SelectorFlagContainsPseudo,
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchId,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("main"),
			},
		},
	}

	sel5 := &Selector{
		Flag: SelectorFlagContainsPseudo,
		Selectors: []*SimpleSelector{
			{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("test"),
			},
			{
				Match:    SelectorMatchId,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("extra"),
			},
		},
	}

	tests := []struct {
		name     string
		sel1     *Selector
		sel2     *Selector
		expected bool
	}{
		{"identical selectors", sel1, sel2, true},
		{"same object", sel1, sel1, true},
		{"different flags", sel1, sel3, false},
		{"different simple selectors", sel1, sel4, false},
		{"different selector count", sel1, sel5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sel1.Equals(tt.sel2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSimpleSelectorString(t *testing.T) {
	tests := []struct {
		name     string
		selector *SimpleSelector
		expected string
	}{
		{
			name: "class selector",
			selector: &SimpleSelector{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationSubSelector,
				Data:     NewSelectorData("test"),
			},
			expected: ".test",
		},
		{
			name: "descendant class selector",
			selector: &SimpleSelector{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationDescendant,
				Data:     NewSelectorData("child"),
			},
			expected: " .child",
		},
		{
			name: "child combinator",
			selector: &SimpleSelector{
				Match:    SelectorMatchTag,
				Relation: SelectorRelationChild,
				Data:     NewSelectorDataTag("", "li"),
			},
			expected: " > li",
		},
		{
			name: "adjacent sibling",
			selector: &SimpleSelector{
				Match:    SelectorMatchTag,
				Relation: SelectorRelationDirectAdjacent,
				Data:     NewSelectorDataTag("", "p"),
			},
			expected: " + p",
		},
		{
			name: "general sibling",
			selector: &SimpleSelector{
				Match:    SelectorMatchTag,
				Relation: SelectorRelationIndirectAdjacent,
				Data:     NewSelectorDataTag("", "div"),
			},
			expected: " ~ div",
		},
		{
			name: "nil data",
			selector: &SimpleSelector{
				Match:    SelectorMatchClass,
				Relation: SelectorRelationSubSelector,
				Data:     nil,
			},
			expected: "[UnknownSelector]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.selector.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSimpleSelectorEquals(t *testing.T) {
	sel1 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     NewSelectorData("test"),
	}

	sel2 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     NewSelectorData("test"),
	}

	sel3 := &SimpleSelector{
		Match:    SelectorMatchId,
		Relation: SelectorRelationSubSelector,
		Data:     NewSelectorData("test"),
	}

	sel4 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationDescendant,
		Data:     NewSelectorData("test"),
	}

	sel5 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     NewSelectorData("different"),
	}

	selNilData1 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     nil,
	}

	selNilData2 := &SimpleSelector{
		Match:    SelectorMatchClass,
		Relation: SelectorRelationSubSelector,
		Data:     nil,
	}

	tests := []struct {
		name     string
		sel1     *SimpleSelector
		sel2     *SimpleSelector
		expected bool
	}{
		{"identical selectors", sel1, sel2, true},
		{"same object", sel1, sel1, true},
		{"different match type", sel1, sel3, false},
		{"different relation", sel1, sel4, false},
		{"different data", sel1, sel5, false},
		{"nil comparison", sel1, nil, false},
		{"both nil data", selNilData1, selNilData2, true},
		{"one nil data", sel1, selNilData1, false},
		{"other nil data", selNilData1, sel1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sel1.Equals(tt.sel2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSelectorRelationTypeString(t *testing.T) {
	tests := []struct {
		name     string
		relation SelectorRelationType
		expected string
	}{
		{"sub selector", SelectorRelationSubSelector, ""},
		{"descendant", SelectorRelationDescendant, " "},
		{"child", SelectorRelationChild, " > "},
		{"direct adjacent", SelectorRelationDirectAdjacent, " + "},
		{"indirect adjacent", SelectorRelationIndirectAdjacent, " ~ "},
		{"relative descendant", SelectorRelationRelativeDescendant, " "},
		{"relative child", SelectorRelationRelativeChild, " > "},
		{"relative direct adjacent", SelectorRelationRelativeDirectAdjacent, " + "},
		{"relative indirect adjacent", SelectorRelationRelativeIndirectAdjacent, " ~ "},
		{"unknown relation", SelectorRelationType(999), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.relation.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
