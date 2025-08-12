package selector

import (
	"testing"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser"
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
			csslexer.Token{Type: csslexer.EOFToken, Data: nil},
		},
		{
			"valid name with namespace",
			"ns|div",
			[]rune("div"),
			[]rune("ns"),
			true,
			csslexer.Token{Type: csslexer.EOFToken, Data: nil},
		},
		{
			"universal selector",
			"*",
			nil,
			nil,
			true,
			csslexer.Token{Type: csslexer.EOFToken, Data: nil},
		},
		{
			"name with id",
			"div#id",
			[]rune("div"),
			nil,
			true,
			csslexer.Token{Type: csslexer.HashToken, Data: []rune("#id")},
		},
		{
			"invalid name with delimiter",
			"div|#id",
			nil,
			nil,
			false,
			csslexer.Token{Type: csslexer.DelimiterToken, Data: []rune("|")},
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
			if nextToken.Type != tc.expectedNextToken.Type || string(nextToken.Data) != string(tc.expectedNextToken.Data) {
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
					Data:     []rune("div"),
					Relation: SelectorRelationSubSelector,
				},
				{
					Match:    SelectorMatchClass,
					Data:     []rune("class"),
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
					Data:     []rune("div"),
					Relation: SelectorRelationSubSelector,
				},
				{
					Match:    SelectorMatchId,
					Data:     []rune("id"),
					Relation: SelectorRelationSubSelector,
				},
				{
					Match:     SelectorMatchAttributeExact,
					Data:      []rune("attr"),
					Relation:  SelectorRelationSubSelector,
					AttrValue: []rune("value"),
					AttrMatch: SelectorAttributeMatchCaseSensitive,
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

			selectors := sp.consumeCompoundSelector(cssparser.NestingTypeNone)

			if len(selectors) != len(tc.expectedSelectors) {
				t.Errorf("expected %d selectors, got %d", len(tc.expectedSelectors), len(selectors))
				return
			}

			for i, sel := range selectors {
				expectedSel := tc.expectedSelectors[i]
				if !sel.Equal(expectedSel) {
					t.Errorf("selector %d mismatch: expected %v, got %v", i, expectedSel, sel)
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
						Data:     []rune("div"),
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:    SelectorMatchClass,
						Data:     []rune("class"),
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
						Data:     []rune("div"),
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:    SelectorMatchId,
						Data:     []rune("id"),
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:     SelectorMatchAttributeExact,
						Data:      []rune("attr"),
						Relation:  SelectorRelationSubSelector,
						AttrValue: []rune("value"),
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
						Data:     []rune("div"),
						Relation: SelectorRelationSubSelector,
					},
					{
						Match:    SelectorMatchClass,
						Data:     []rune("class"),
						Relation: SelectorRelationChild,
					},
					{
						Match:    SelectorMatchId,
						Data:     []rune("id"),
						Relation: SelectorRelationDirectAdjacent,
					},
					{
						Match:     SelectorMatchAttributeExact,
						Data:      []rune("attr"),
						Relation:  SelectorRelationIndirectAdjacent,
						AttrValue: []rune("value"),
						AttrMatch: SelectorAttributeMatchCaseSensitive,
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

			selector, err := sp.consumeComplexSelector(cssparser.NestingTypeNone)

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
