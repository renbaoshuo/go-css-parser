package selector

import (
	"testing"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/css"
	"go.baoshuo.dev/cssparser/nesting"
	"go.baoshuo.dev/cssparser/token_stream"
)

func Test_SelectorParser_ConsumeName(t *testing.T) {
	testcases := []struct {
		name              string
		input             string
		expectedName      []rune
		expectedNamespace []rune
		expectedSuccess   bool
		expectedNextToken csslexer.Token
	}{
		{
			"valid name without namespace",
			"div",
			[]rune("div"),
			nil,
			true,
			csslexer.Token{Type: csslexer.EOFToken, Value: "", Raw: nil},
		},
		{
			"valid name with namespace",
			"ns|div",
			[]rune("div"),
			[]rune("ns"),
			true,
			csslexer.Token{Type: csslexer.EOFToken, Value: "", Raw: nil},
		},
		{
			"universal selector",
			"*",
			nil,
			nil,
			true,
			csslexer.Token{Type: csslexer.EOFToken, Value: "", Raw: nil},
		},
		{
			"name with id",
			"div#id",
			[]rune("div"),
			nil,
			true,
			csslexer.Token{Type: csslexer.HashToken, Value: "id", Raw: []rune("#id")},
		},
		{
			"invalid name with delimiter",
			"div|#id",
			nil,
			nil,
			false,
			csslexer.Token{Type: csslexer.DelimiterToken, Value: "|", Raw: []rune("|")},
		},
		{
			"universal selector with namespace",
			"ns|*",
			nil,
			[]rune("ns"),
			true,
			csslexer.Token{Type: csslexer.EOFToken, Value: "", Raw: nil},
		},
		{
			"empty namespace",
			"|div",
			[]rune("div"),
			[]rune("*"),
			true,
			csslexer.Token{Type: csslexer.EOFToken, Value: "", Raw: nil},
		},
		{
			"name with hyphens",
			"custom-element",
			[]rune("custom-element"),
			nil,
			true,
			csslexer.Token{Type: csslexer.EOFToken, Value: "", Raw: nil},
		},
		{
			"name with underscores",
			"my_element",
			[]rune("my_element"),
			nil,
			true,
			csslexer.Token{Type: csslexer.EOFToken, Value: "", Raw: nil},
		},
		{
			"name with numbers",
			"h1",
			[]rune("h1"),
			nil,
			true,
			csslexer.Token{Type: csslexer.EOFToken, Value: "", Raw: nil},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			name, namespace, success := sp.consumeName()

			if string(name) != string(tc.expectedName) {
				t.Errorf("expected name %q, got %q", tc.expectedName, name)
			}
			if string(namespace) != string(tc.expectedNamespace) {
				t.Errorf("expected namespace %q, got %q", tc.expectedNamespace, namespace)
			}
			if success != tc.expectedSuccess {
				t.Errorf("expected success %v, got %v", tc.expectedSuccess, success)
			}

			nextToken := ts.Peek()
			if nextToken.Type != tc.expectedNextToken.Type ||
				nextToken.Value != tc.expectedNextToken.Value ||
				string(nextToken.Raw) != string(tc.expectedNextToken.Raw) {
				t.Errorf("expected next token %v, got %v", tc.expectedNextToken, nextToken)
			}
		})
	}
}

func Test_SelectorParser_ConsumeCompoundSelector(t *testing.T) {
	testcases := []struct {
		name              string
		input             string
		expectedSelectors []*css.SimpleSelector
		expectedValid     bool
	}{
		{
			"valid compound selector with tag and class",
			"div.class",
			[]*css.SimpleSelector{
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
			true,
		},
		{
			"valid compound selector with id and attribute",
			"div#id[attr='value']",
			[]*css.SimpleSelector{
				{
					Match:    css.SelectorMatchTag,
					Data:     css.NewSelectorDataTag("", "div"),
					Relation: css.SelectorRelationSubSelector,
				},
				{
					Match:    css.SelectorMatchId,
					Data:     css.NewSelectorData("id"),
					Relation: css.SelectorRelationSubSelector,
				},
				{
					Match:    css.SelectorMatchAttributeExact,
					Relation: css.SelectorRelationSubSelector,
					Data:     css.NewSelectorDataAttr("attr", "value", css.SelectorAttrMatchCaseSensitive),
				},
			},
			true,
		},
		{
			"valid compound selector with only class",
			".class",
			[]*css.SimpleSelector{
				{
					Match:    css.SelectorMatchClass,
					Data:     css.NewSelectorData("class"),
					Relation: css.SelectorRelationSubSelector,
				},
			},
			true,
		},
		{
			"compound selector with multiple classes",
			".btn.primary.large",
			[]*css.SimpleSelector{
				{
					Match:    css.SelectorMatchClass,
					Data:     css.NewSelectorData("btn"),
					Relation: css.SelectorRelationSubSelector,
				},
				{
					Match:    css.SelectorMatchClass,
					Data:     css.NewSelectorData("primary"),
					Relation: css.SelectorRelationSubSelector,
				},
				{
					Match:    css.SelectorMatchClass,
					Data:     css.NewSelectorData("large"),
					Relation: css.SelectorRelationSubSelector,
				},
			},
			true,
		},
		{
			"compound selector with multiple attributes",
			"input[type=text][required][disabled]",
			[]*css.SimpleSelector{
				{
					Match:    css.SelectorMatchTag,
					Data:     css.NewSelectorDataTag("", "input"),
					Relation: css.SelectorRelationSubSelector,
				},
				{
					Match:    css.SelectorMatchAttributeExact,
					Relation: css.SelectorRelationSubSelector,
					Data:     css.NewSelectorDataAttr("type", "text", css.SelectorAttrMatchCaseSensitive),
				},
				{
					Match:    css.SelectorMatchAttributeSet,
					Relation: css.SelectorRelationSubSelector,
					Data:     css.NewSelectorDataAttr("required", "", css.SelectorAttrMatchCaseSensitive),
				},
				{
					Match:    css.SelectorMatchAttributeSet,
					Relation: css.SelectorRelationSubSelector,
					Data:     css.NewSelectorDataAttr("disabled", "", css.SelectorAttrMatchCaseSensitive),
				},
			},
			true,
		},
		{
			"universal selector",
			"*",
			[]*css.SimpleSelector{
				{
					Match:    css.SelectorMatchUniversalTag,
					Relation: css.SelectorRelationSubSelector,
					Data:     css.NewSelectorDataTag("", ""),
				},
			},
			true,
		},
		{
			"universal selector with class",
			"*.warning",
			[]*css.SimpleSelector{
				{
					Match:    css.SelectorMatchClass,
					Data:     css.NewSelectorData("warning"),
					Relation: css.SelectorRelationSubSelector,
				},
			},
			true,
		},
		{
			"universal selector with namespace",
			"ns|*",
			[]*css.SimpleSelector{
				{
					Match:    css.SelectorMatchUniversalTag,
					Data:     css.NewSelectorDataTag("ns", ""),
					Relation: css.SelectorRelationSubSelector,
				},
			},
			true,
		},
		{
			"id only",
			"#main",
			[]*css.SimpleSelector{
				{
					Match:    css.SelectorMatchId,
					Data:     css.NewSelectorData("main"),
					Relation: css.SelectorRelationSubSelector,
				},
			},
			true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			selectors, _ := sp.consumeCompoundSelector(nesting.NestingTypeNone)

			if len(selectors) != len(tc.expectedSelectors) {
				t.Errorf("expected %d selectors, got %d", len(tc.expectedSelectors), len(selectors))
				return
			}

			for i, sel := range selectors {
				expectedSel := tc.expectedSelectors[i]
				if !sel.Equals(expectedSel) {
					t.Errorf("selector %d mismatch: expected %q, got %q", i, expectedSel, sel)
				}
			}
		})
	}
}

func Test_SelectorParser_ConsumeComplexSelector(t *testing.T) {
	testcases := []struct {
		name             string
		input            string
		expectedSelector *css.Selector
		expectError      bool
	}{
		{
			"valid complex selector with tag and class",
			"div.class",
			&css.Selector{
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
			false,
		},
		{
			"valid complex selector with id and attribute",
			"div#id[attr='value']",
			&css.Selector{
				Flag: 0,
				Selectors: []*css.SimpleSelector{
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "div"),
						Relation: css.SelectorRelationSubSelector,
					},
					{
						Match:    css.SelectorMatchId,
						Data:     css.NewSelectorData("id"),
						Relation: css.SelectorRelationSubSelector,
					},
					{
						Match:    css.SelectorMatchAttributeExact,
						Relation: css.SelectorRelationSubSelector,
						Data:     css.NewSelectorDataAttr("attr", "value", css.SelectorAttrMatchCaseSensitive),
					},
				},
			},
			false,
		},
		{
			"valid complex selector with combinators",
			"div > .class + #id ~ [attr='value']",
			&css.Selector{
				Flag: css.SelectorFlagContainsComplexSelector,
				Selectors: []*css.SimpleSelector{
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "div"),
						Relation: css.SelectorRelationSubSelector,
					},
					{
						Match:    css.SelectorMatchClass,
						Data:     css.NewSelectorData("class"),
						Relation: css.SelectorRelationChild,
					},
					{
						Match:    css.SelectorMatchId,
						Data:     css.NewSelectorData("id"),
						Relation: css.SelectorRelationDirectAdjacent,
					},
					{
						Match:    css.SelectorMatchAttributeExact,
						Relation: css.SelectorRelationIndirectAdjacent,
						Data:     css.NewSelectorDataAttr("attr", "value", css.SelectorAttrMatchCaseSensitive),
					},
				},
			},
			false,
		},
		{
			"descendant combinator",
			"nav ul",
			&css.Selector{
				Flag: css.SelectorFlagContainsComplexSelector,
				Selectors: []*css.SimpleSelector{
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "nav"),
						Relation: css.SelectorRelationSubSelector,
					},
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "ul"),
						Relation: css.SelectorRelationDescendant,
					},
				},
			},
			false,
		},
		{
			"child combinator",
			"article > p",
			&css.Selector{
				Flag: css.SelectorFlagContainsComplexSelector,
				Selectors: []*css.SimpleSelector{
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "article"),
						Relation: css.SelectorRelationSubSelector,
					},
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "p"),
						Relation: css.SelectorRelationChild,
					},
				},
			},
			false,
		},
		{
			"adjacent sibling combinator",
			"h1 + p",
			&css.Selector{
				Flag: css.SelectorFlagContainsComplexSelector,
				Selectors: []*css.SimpleSelector{
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "h1"),
						Relation: css.SelectorRelationSubSelector,
					},
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "p"),
						Relation: css.SelectorRelationDirectAdjacent,
					},
				},
			},
			false,
		},
		{
			"general sibling combinator",
			"h1 ~ p",
			&css.Selector{
				Flag: css.SelectorFlagContainsComplexSelector,
				Selectors: []*css.SimpleSelector{
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "h1"),
						Relation: css.SelectorRelationSubSelector,
					},
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "p"),
						Relation: css.SelectorRelationIndirectAdjacent,
					},
				},
			},
			false,
		},
		{
			"complex selector with multiple combinators",
			"main article > header h1.title",
			&css.Selector{
				Flag: css.SelectorFlagContainsComplexSelector,
				Selectors: []*css.SimpleSelector{
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "main"),
						Relation: css.SelectorRelationSubSelector,
					},
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "article"),
						Relation: css.SelectorRelationDescendant,
					},
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "header"),
						Relation: css.SelectorRelationChild,
					},
					{
						Match:    css.SelectorMatchTag,
						Data:     css.NewSelectorDataTag("", "h1"),
						Relation: css.SelectorRelationDescendant,
					},
					{
						Match:    css.SelectorMatchClass,
						Data:     css.NewSelectorData("title"),
						Relation: css.SelectorRelationSubSelector,
					},
				},
			},
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			selector, err := sp.consumeComplexSelector(nesting.NestingTypeNone, true)

			if (err != nil) != tc.expectError {
				t.Errorf("expected error: %v, got: %v", tc.expectError, err)
				return
			}

			if !tc.expectError && !selector.Equals(tc.expectedSelector) {
				t.Errorf("expected selector %v, got %v", tc.expectedSelector, selector)
			}

			t.Logf("Parsed selector: %v", selector.String())
		})
	}
}

func Test_SelectorParser_ConsumeCombinator(t *testing.T) {
	testcases := []struct {
		name               string
		input              string
		expectedCombinator css.SelectorRelationType
	}{
		{
			"child combinator",
			">",
			css.SelectorRelationChild,
		},
		{
			"adjacent sibling combinator",
			"+",
			css.SelectorRelationDirectAdjacent,
		},
		{
			"general sibling combinator",
			"~",
			css.SelectorRelationIndirectAdjacent,
		},
		{
			"descendant combinator (whitespace)",
			"   ",
			css.SelectorRelationDescendant,
		},
		{
			"no combinator",
			"div",
			css.SelectorRelationSubSelector,
		},
		{
			"child combinator with whitespace",
			"  >  ",
			css.SelectorRelationChild,
		},
		{
			"adjacent with whitespace",
			"  +  ",
			css.SelectorRelationDirectAdjacent,
		},
		{
			"general sibling with whitespace",
			"  ~  ",
			css.SelectorRelationIndirectAdjacent,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			combinator := sp.consumeCombinator()

			if combinator != tc.expectedCombinator {
				t.Errorf("expected combinator %v, got %v", tc.expectedCombinator, combinator)
			}
		})
	}
}
