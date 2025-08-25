package selector

import (
	"testing"

	"go.baoshuo.dev/csslexer"

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
		expectedSelectors []*SimpleSelector
		expectedValid     bool
	}{
		{
			"valid compound selector with tag and class",
			"div.class",
			[]*SimpleSelector{
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
			true,
		},
		{
			"valid compound selector with id and attribute",
			"div#id[attr='value']",
			[]*SimpleSelector{
				{
					Match:    SelectorMatchTag,
					Value:    "div",
					Relation: SelectorRelationSubSelector,
				},
				{
					Match:    SelectorMatchId,
					Value:    "id",
					Relation: SelectorRelationSubSelector,
				},
				{
					Match:     SelectorMatchAttributeExact,
					Value:     "attr",
					Relation:  SelectorRelationSubSelector,
					AttrValue: "value",
					AttrMatch: SelectorAttributeMatchCaseSensitive,
				},
			},
			true,
		},
		{
			"valid compound selector with only class",
			".class",
			[]*SimpleSelector{
				{
					Match:    SelectorMatchClass,
					Value:    "class",
					Relation: SelectorRelationSubSelector,
				},
			},
			true,
		},
		{
			"compound selector with multiple classes",
			".btn.primary.large",
			[]*SimpleSelector{
				{
					Match:    SelectorMatchClass,
					Value:    "btn",
					Relation: SelectorRelationSubSelector,
				},
				{
					Match:    SelectorMatchClass,
					Value:    "primary",
					Relation: SelectorRelationSubSelector,
				},
				{
					Match:    SelectorMatchClass,
					Value:    "large",
					Relation: SelectorRelationSubSelector,
				},
			},
			true,
		},
		{
			"compound selector with multiple attributes",
			"input[type=text][required][disabled]",
			[]*SimpleSelector{
				{
					Match:    SelectorMatchTag,
					Value:    "input",
					Relation: SelectorRelationSubSelector,
				},
				{
					Match:     SelectorMatchAttributeExact,
					Value:     "type",
					Relation:  SelectorRelationSubSelector,
					AttrValue: "text",
					AttrMatch: SelectorAttributeMatchCaseSensitive,
				},
				{
					Match:    SelectorMatchAttributeSet,
					Value:    "required",
					Relation: SelectorRelationSubSelector,
				},
				{
					Match:    SelectorMatchAttributeSet,
					Value:    "disabled",
					Relation: SelectorRelationSubSelector,
				},
			},
			true,
		},
		{
			"universal selector",
			"*",
			[]*SimpleSelector{
				{
					Match:    SelectorMatchUniversalTag,
					Relation: SelectorRelationSubSelector,
				},
			},
			true,
		},
		{
			"universal selector with class",
			"*.warning",
			[]*SimpleSelector{
				{
					Match:    SelectorMatchClass,
					Value:    "warning",
					Relation: SelectorRelationSubSelector,
				},
			},
			true,
		},
		{
			"id only",
			"#main",
			[]*SimpleSelector{
				{
					Match:    SelectorMatchId,
					Value:    "main",
					Relation: SelectorRelationSubSelector,
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
				if !sel.Equal(expectedSel) {
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
		expectedSelector *Selector
		expectError      bool
	}{
		{
			"valid complex selector with tag and class",
			"div.class",
			&Selector{
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
			false,
		},
		{
			"valid complex selector with id and attribute",
			"div#id[attr='value']",
			&Selector{
				Flag: 0,
				Selectors: []*SimpleSelector{
					{
						Match:    SelectorMatchTag,
						Value:    "div",
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:    SelectorMatchId,
						Value:    "id",
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:     SelectorMatchAttributeExact,
						Value:     "attr",
						Relation:  SelectorRelationSubSelector,
						AttrValue: "value",
						AttrMatch: SelectorAttributeMatchCaseSensitive,
					},
				},
			},
			false,
		},
		{
			"valid complex selector with combinators",
			"div > .class + #id ~ [attr='value']",
			&Selector{
				Flag: SelectorFlagContainsComplexSelector,
				Selectors: []*SimpleSelector{
					{
						Match:    SelectorMatchTag,
						Value:    "div",
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:    SelectorMatchClass,
						Value:    "class",
						Relation: SelectorRelationChild,
					},
					{
						Match:    SelectorMatchId,
						Value:    "id",
						Relation: SelectorRelationDirectAdjacent,
					},
					{
						Match:     SelectorMatchAttributeExact,
						Value:     "attr",
						Relation:  SelectorRelationIndirectAdjacent,
						AttrValue: "value",
						AttrMatch: SelectorAttributeMatchCaseSensitive,
					},
				},
			},
			false,
		},
		{
			"descendant combinator",
			"nav ul",
			&Selector{
				Flag: SelectorFlagContainsComplexSelector,
				Selectors: []*SimpleSelector{
					{
						Match:    SelectorMatchTag,
						Value:    "nav",
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:    SelectorMatchTag,
						Value:    "ul",
						Relation: SelectorRelationDescendant,
					},
				},
			},
			false,
		},
		{
			"child combinator",
			"article > p",
			&Selector{
				Flag: SelectorFlagContainsComplexSelector,
				Selectors: []*SimpleSelector{
					{
						Match:    SelectorMatchTag,
						Value:    "article",
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:    SelectorMatchTag,
						Value:    "p",
						Relation: SelectorRelationChild,
					},
				},
			},
			false,
		},
		{
			"adjacent sibling combinator",
			"h1 + p",
			&Selector{
				Flag: SelectorFlagContainsComplexSelector,
				Selectors: []*SimpleSelector{
					{
						Match:    SelectorMatchTag,
						Value:    "h1",
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:    SelectorMatchTag,
						Value:    "p",
						Relation: SelectorRelationDirectAdjacent,
					},
				},
			},
			false,
		},
		{
			"general sibling combinator",
			"h1 ~ p",
			&Selector{
				Flag: SelectorFlagContainsComplexSelector,
				Selectors: []*SimpleSelector{
					{
						Match:    SelectorMatchTag,
						Value:    "h1",
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:    SelectorMatchTag,
						Value:    "p",
						Relation: SelectorRelationIndirectAdjacent,
					},
				},
			},
			false,
		},
		{
			"complex selector with multiple combinators",
			"main article > header h1.title",
			&Selector{
				Flag: SelectorFlagContainsComplexSelector,
				Selectors: []*SimpleSelector{
					{
						Match:    SelectorMatchTag,
						Value:    "main",
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:    SelectorMatchTag,
						Value:    "article",
						Relation: SelectorRelationDescendant,
					},
					{
						Match:    SelectorMatchTag,
						Value:    "header",
						Relation: SelectorRelationChild,
					},
					{
						Match:    SelectorMatchTag,
						Value:    "h1",
						Relation: SelectorRelationDescendant,
					},
					{
						Match:    SelectorMatchClass,
						Value:    "title",
						Relation: SelectorRelationSubSelector,
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

			if !tc.expectError && !selector.Equal(tc.expectedSelector) {
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
		expectedCombinator SelectorRelationType
	}{
		{
			"child combinator",
			">",
			SelectorRelationChild,
		},
		{
			"adjacent sibling combinator",
			"+",
			SelectorRelationDirectAdjacent,
		},
		{
			"general sibling combinator",
			"~",
			SelectorRelationIndirectAdjacent,
		},
		{
			"descendant combinator (whitespace)",
			"   ",
			SelectorRelationDescendant,
		},
		{
			"no combinator",
			"div",
			SelectorRelationSubSelector,
		},
		{
			"child combinator with whitespace",
			"  >  ",
			SelectorRelationChild,
		},
		{
			"adjacent with whitespace",
			"  +  ",
			SelectorRelationDirectAdjacent,
		},
		{
			"general sibling with whitespace",
			"  ~  ",
			SelectorRelationIndirectAdjacent,
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
