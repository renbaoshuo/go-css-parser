package selector

import (
	"testing"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser"
	"go.baoshuo.dev/cssparser/token_stream"
)

func Test_SelectorParser_ConsumeComplexSelectorList(t *testing.T) {
	testcases := []struct {
		name              string
		input             string
		nestingType       cssparser.NestingTypeType
		expectedCount     int
		expectedError     bool
		expectedSelectors []*Selector
	}{
		{
			name:          "single selector",
			input:         "div",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 1,
			expectedError: false,
			expectedSelectors: []*Selector{
				{
					Flag: 0,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchTag,
							Value:    "div",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
			},
		},
		{
			name:          "multiple selectors",
			input:         "div, .class, #id",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 3,
			expectedError: false,
			expectedSelectors: []*Selector{
				{
					Flag: 0,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchTag,
							Value:    "div",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: 0,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchClass,
							Value:    "class",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: 0,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchId,
							Value:    "id",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
			},
		},
		{
			name:          "complex selectors",
			input:         "div.class, p > span, .item + .next",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 3,
			expectedError: false,
			expectedSelectors: []*Selector{
				{
					Flag: 0,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchTag,
							Value:    "div",
							Relation: SelectorRelationSubSelector,
						},
						{
							Match:    SelectorMatchClass,
							Value:    "class",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: SelectorFlagContainsComplexSelector,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchTag,
							Value:    "p",
							Relation: SelectorRelationSubSelector,
						},
						{
							Match:    SelectorMatchTag,
							Value:    "span",
							Relation: SelectorRelationChild,
						},
					},
				},
				{
					Flag: SelectorFlagContainsComplexSelector,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchClass,
							Value:    "item",
							Relation: SelectorRelationSubSelector,
						},
						{
							Match:    SelectorMatchClass,
							Value:    "next",
							Relation: SelectorRelationDirectAdjacent,
						},
					},
				},
			},
		},
		{
			name:          "selector with trailing comma",
			input:         "div, .class",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 2,
			expectedError: false,
			expectedSelectors: []*Selector{
				{
					Flag: 0,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchTag,
							Value:    "div",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: 0,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchClass,
							Value:    "class",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
			},
		},
		{
			name:          "selector with whitespace test",
			input:         "div , .class",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 2,
			expectedError: false,
			expectedSelectors: []*Selector{
				{
					Flag: SelectorFlagContainsComplexSelector,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchTag,
							Value:    "div",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: 0,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchClass,
							Value:    "class",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
			},
		},
		{
			name:          "empty selector",
			input:         "",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 0,
			expectedError: true,
		},
		{
			name:          "selector ending with left brace",
			input:         "div {",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 1,
			expectedError: false,
			expectedSelectors: []*Selector{
				{
					Flag: SelectorFlagContainsComplexSelector,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchTag,
							Value:    "div",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
			},
		},
		{
			name:          "multiple selectors ending with left brace",
			input:         "div, .class {",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 2,
			expectedError: false,
			expectedSelectors: []*Selector{
				{
					Flag: 0,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchTag,
							Value:    "div",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: SelectorFlagContainsComplexSelector,
					Selectors: []*SimpleSelector{
						{
							Match:    SelectorMatchClass,
							Value:    "class",
							Relation: SelectorRelationSubSelector,
						},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			selectors, err := sp.consumeComplexSelectorList(tc.nestingType)

			if tc.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(selectors) != tc.expectedCount {
				t.Errorf("expected %d selectors, got %d", tc.expectedCount, len(selectors))
				return
			}

			if tc.expectedSelectors != nil {
				for i, sel := range selectors {
					if i >= len(tc.expectedSelectors) {
						t.Errorf("unexpected selector at index %d", i)
						continue
					}
					expectedSel := tc.expectedSelectors[i]
					if !sel.Equal(expectedSel) {
						t.Errorf("selector %d mismatch:\nexpected: %v\ngot: %v", i, expectedSel, sel)
					}
				}
			}

			// Log the parsed selectors for debugging
			for i, sel := range selectors {
				t.Logf("Selector %d: %s (Flag: %d)", i, sel.String(), sel.Flag)
			}
		})
	}
}

func Test_ConsumeSelector(t *testing.T) {
	testcases := []struct {
		name          string
		input         string
		nestingType   cssparser.NestingTypeType
		expectedCount int
		expectedError bool
	}{
		{
			name:          "single tag selector",
			input:         "div",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:          "multiple selectors",
			input:         "h1, h2, h3",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 3,
			expectedError: false,
		},
		{
			name:          "complex selector with combinators",
			input:         "nav > ul li.active",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:          "attribute selectors",
			input:         "[type=text], [required]",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:          "class and id selectors",
			input:         ".btn, #main, .nav-item",
			nestingType:   cssparser.NestingTypeNone,
			expectedCount: 3,
			expectedError: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)

			selectors, err := ConsumeSelector(ts, tc.nestingType, nil)

			if tc.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(selectors) != tc.expectedCount {
				t.Errorf("expected %d selectors, got %d", tc.expectedCount, len(selectors))
				return
			}

			// Log the parsed selectors for debugging
			for i, sel := range selectors {
				t.Logf("Selector %d: %s", i, sel.String())
			}
		})
	}
}
