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
		expectedData      string
		expectedValid     bool
		expectedAttrValue string
		expectedAttrMatch SelectorAttributeMatchType
	}{
		{
			"valid hash selector",
			"#id",
			SelectorMatchId,
			"id",
			true,
			"",
			SelectorAttributeMatchUnknown,
		},
		{
			"valid class selector",
			".class",
			SelectorMatchClass,
			"class",
			true,
			"",
			SelectorAttributeMatchUnknown,
		},
		{
			"valid attribute selector",
			"[attr=value]",
			SelectorMatchAttributeExact,
			"attr",
			true,
			"value",
			SelectorAttributeMatchCaseSensitive,
		},
		{
			"valid attribute selector with namespace",
			"[ns|attr=value]",
			SelectorMatchAttributeExact,
			"ns|attr",
			true,
			"value",
			SelectorAttributeMatchCaseSensitive,
		},
		{
			"valid attribute selector with case insensitive match",
			"[attr|='value' i]",
			SelectorMatchAttributeHyphen,
			"attr",
			true,
			"value",
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
				if sel.Value != tc.expectedData {
					t.Errorf("expected data %q, got %q", tc.expectedData, sel.Value)
				}
				if sel.AttrValue != tc.expectedAttrValue {
					t.Errorf("expected attr value %q, got %q", tc.expectedAttrValue, sel.AttrValue)
				}
				if sel.AttrMatch != tc.expectedAttrMatch {
					t.Errorf("expected attr match %q, got %q", tc.expectedAttrMatch, sel.AttrMatch)
				}

				t.Logf("selector: %s", sel.String())
			} else if tc.expectedValid {
				t.Error("expected a selector but got nil")
			}
		})
	}
}
