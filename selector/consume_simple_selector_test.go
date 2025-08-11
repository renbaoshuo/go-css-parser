package selector

import (
	"testing"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/token_stream"
)

func Test_SelectorParser_ConsumeSimpleSelector(t *testing.T) {
	testcases := []struct {
		name              string
		input             string
		expectedMatch     SelectorMatchType
		expectedData      []rune
		expectedValid     bool
		expectedAttrValue []rune
		expectedAttrMatch SelectorAttributeMatchType
	}{
		{
			"valid hash selector",
			"#id",
			SelectorMatchId,
			[]rune("id"),
			true,
			nil,
			SelectorAttributeMatchUnknown,
		},
		{
			"valid class selector",
			".class",
			SelectorMatchClass,
			[]rune("class"),
			true,
			nil,
			SelectorAttributeMatchUnknown,
		},
		{
			"valid attribute selector",
			"[attr=value]",
			SelectorMatchAttributeExact,
			[]rune("attr"),
			true,
			[]rune("value"),
			SelectorAttributeMatchCaseSensitive,
		},
		{
			"valid attribute selector with namespace",
			"[ns|attr=value]",
			SelectorMatchAttributeExact,
			[]rune("ns|attr"),
			true,
			[]rune("value"),
			SelectorAttributeMatchCaseSensitive,
		},
		{
			"valid attribute selector with case insensitive match",
			"[attr|='value' i]",
			SelectorMatchAttributeHyphen,
			[]rune("attr"),
			true,
			[]rune("value"),
			SelectorAttributeMatchCaseInsensitive,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			sel, err := sp.consumeSimpleSelector()

			if tc.expectedValid && err != nil {
				t.Errorf("expected valid, got %v", err)
				return
			}

			if sel != nil {
				if sel.Match != tc.expectedMatch {
					t.Errorf("expected type %q, got %q", tc.expectedMatch, sel.Match)
				}
				if string(sel.Data) != string(tc.expectedData) {
					t.Errorf("expected data %q, got %q", tc.expectedData, sel.Data)
				}
				if string(sel.AttrValue) != string(tc.expectedAttrValue) {
					t.Errorf("expected attr value %q, got %q", tc.expectedAttrValue, sel.AttrValue)
				}
				if sel.AttrMatch != tc.expectedAttrMatch {
					t.Errorf("expected attr match %q, got %q", tc.expectedAttrMatch, sel.AttrMatch)
				}
			} else if tc.expectedValid {
				t.Error("expected a selector but got nil")
			}
		})
	}
}
