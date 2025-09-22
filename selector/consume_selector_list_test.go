package selector

import (
	"testing"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/css"
	"go.baoshuo.dev/cssparser/nesting"
	"go.baoshuo.dev/cssparser/token_stream"
)

func Test_SelectorParser_ConsumeComplexSelectorList(t *testing.T) {
	testcases := []struct {
		name              string
		input             string
		nestingType       nesting.NestingTypeType
		expectedCount     int
		expectedError     bool
		expectedSelectors []*css.Selector
	}{
		{
			name:          "single selector",
			input:         "div",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 1,
			expectedError: false,
			expectedSelectors: []*css.Selector{
				{
					Flag: 0,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchTag,
							Data:     css.NewSelectorDataTag("", "div"),
							Relation: css.SelectorRelationSubSelector,
						},
					},
				},
			},
		},
		{
			name:          "multiple selectors",
			input:         "div, .class, #id",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 3,
			expectedError: false,
			expectedSelectors: []*css.Selector{
				{
					Flag: 0,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchTag,
							Data:     css.NewSelectorDataTag("", "div"),
							Relation: css.SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: 0,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchClass,
							Data:     css.NewSelectorData("class"),
							Relation: css.SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: 0,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchId,
							Data:     css.NewSelectorData("id"),
							Relation: css.SelectorRelationSubSelector,
						},
					},
				},
			},
		},
		{
			name:          "complex selectors",
			input:         "div.class, p > span, .item + .next",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 3,
			expectedError: false,
			expectedSelectors: []*css.Selector{
				{
					Flag: 0,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchTag,
							Data:     css.NewSelectorDataTag("", "div"),
							Relation: css.SelectorRelationSubSelector,
						},
						{
							Match:    css.SelectorMatchClass,
							Data:     css.NewSelectorData("class"),
							Relation: css.SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: css.SelectorFlagContainsComplexSelector,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchTag,
							Data:     css.NewSelectorDataTag("", "p"),
							Relation: css.SelectorRelationSubSelector,
						},
						{
							Match:    css.SelectorMatchTag,
							Data:     css.NewSelectorDataTag("", "span"),
							Relation: css.SelectorRelationChild,
						},
					},
				},
				{
					Flag: css.SelectorFlagContainsComplexSelector,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchClass,
							Data:     css.NewSelectorData("item"),
							Relation: css.SelectorRelationSubSelector,
						},
						{
							Match:    css.SelectorMatchClass,
							Data:     css.NewSelectorData("next"),
							Relation: css.SelectorRelationDirectAdjacent,
						},
					},
				},
			},
		},
		{
			name:          "selector with trailing comma",
			input:         "div, .class",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 2,
			expectedError: false,
			expectedSelectors: []*css.Selector{
				{
					Flag: 0,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchTag,
							Data:     css.NewSelectorDataTag("", "div"),
							Relation: css.SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: 0,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchClass,
							Data:     css.NewSelectorData("class"),
							Relation: css.SelectorRelationSubSelector,
						},
					},
				},
			},
		},
		{
			name:          "selector with whitespace test",
			input:         "div , .class",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 2,
			expectedError: false,
			expectedSelectors: []*css.Selector{
				{
					Flag: css.SelectorFlagContainsComplexSelector,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchTag,
							Data:     css.NewSelectorDataTag("", "div"),
							Relation: css.SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: 0,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchClass,
							Data:     css.NewSelectorData("class"),
							Relation: css.SelectorRelationSubSelector,
						},
					},
				},
			},
		},
		{
			name:          "empty selector",
			input:         "",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 0,
			expectedError: true,
		},
		{
			name:          "selector ending with left brace",
			input:         "div {",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 1,
			expectedError: false,
			expectedSelectors: []*css.Selector{
				{
					Flag: css.SelectorFlagContainsComplexSelector,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchTag,
							Data:     css.NewSelectorDataTag("", "div"),
							Relation: css.SelectorRelationSubSelector,
						},
					},
				},
			},
		},
		{
			name:          "multiple selectors ending with left brace",
			input:         "div, .class {",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 2,
			expectedError: false,
			expectedSelectors: []*css.Selector{
				{
					Flag: 0,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchTag,
							Data:     css.NewSelectorDataTag("", "div"),
							Relation: css.SelectorRelationSubSelector,
						},
					},
				},
				{
					Flag: css.SelectorFlagContainsComplexSelector,
					Selectors: []*css.SimpleSelector{
						{
							Match:    css.SelectorMatchClass,
							Data:     css.NewSelectorData("class"),
							Relation: css.SelectorRelationSubSelector,
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
					if !sel.Equals(expectedSel) {
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
		nestingType   nesting.NestingTypeType
		expectedCount int
		expectedError bool
	}{
		{
			name:          "single tag selector",
			input:         "div",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:          "multiple selectors",
			input:         "h1, h2, h3",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 3,
			expectedError: false,
		},
		{
			name:          "complex selector with combinators",
			input:         "nav > ul li.active",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:          "attribute selectors",
			input:         "[type=text], [required]",
			nestingType:   nesting.NestingTypeNone,
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:          "class and id selectors",
			input:         ".btn, #main, .nav-item",
			nestingType:   nesting.NestingTypeNone,
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
