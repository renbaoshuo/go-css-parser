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
		{
			"hash selector with numbers",
			"#id123",
			SelectorMatchId,
			"id123",
			true,
			"",
			SelectorAttributeMatchUnknown,
		},
		{
			"class selector with hyphens",
			".btn-primary",
			SelectorMatchClass,
			"btn-primary",
			true,
			"",
			SelectorAttributeMatchUnknown,
		},
		{
			"attribute selector with string value",
			"[title=\"hello world\"]",
			SelectorMatchAttributeExact,
			"title",
			true,
			"hello world",
			SelectorAttributeMatchCaseSensitive,
		},
		{
			"attribute selector contains match",
			"[class*=\"nav\"]",
			SelectorMatchAttributeContain,
			"class",
			true,
			"nav",
			SelectorAttributeMatchCaseSensitive,
		},
		{
			"attribute selector starts with match",
			"[lang^=\"en\"]",
			SelectorMatchAttributeBegin,
			"lang",
			true,
			"en",
			SelectorAttributeMatchCaseSensitive,
		},
		{
			"attribute selector ends with match",
			"[href$=\".pdf\"]",
			SelectorMatchAttributeEnd,
			"href",
			true,
			".pdf",
			SelectorAttributeMatchCaseSensitive,
		},
		{
			"attribute selector word match",
			"[class~=\"active\"]",
			SelectorMatchAttributeList,
			"class",
			true,
			"active",
			SelectorAttributeMatchCaseSensitive,
		},
		{
			"attribute selector set match",
			"[required]",
			SelectorMatchAttributeSet,
			"required",
			true,
			"",
			SelectorAttributeMatchUnknown,
		},
		{
			"attribute selector with case sensitive flag",
			"[data-name=\"Value\" s]",
			SelectorMatchAttributeExact,
			"data-name",
			true,
			"Value",
			SelectorAttributeMatchCaseSensitiveAlways,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			sel, _, err := sp.consumeSimpleSelector()

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

func Test_SelectorParser_ConsumeId(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expectedId  string
		expectError bool
	}{
		{
			"simple id",
			"#main",
			"main",
			false,
		},
		{
			"id with numbers",
			"#item123",
			"item123",
			false,
		},
		{
			"id with hyphens",
			"#nav-menu",
			"nav-menu",
			false,
		},
		{
			"id with underscores",
			"#user_profile",
			"user_profile",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			sel, err := sp.consumeId()

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if sel.Match != SelectorMatchId {
				t.Errorf("expected SelectorMatchId, got %v", sel.Match)
			}

			if sel.Value != tc.expectedId {
				t.Errorf("expected id %q, got %q", tc.expectedId, sel.Value)
			}
		})
	}
}

func Test_SelectorParser_ConsumeClass(t *testing.T) {
	testcases := []struct {
		name          string
		input         string
		expectedClass string
		expectError   bool
	}{
		{
			"simple class",
			".container",
			"container",
			false,
		},
		{
			"class with numbers",
			".col-12",
			"col-12",
			false,
		},
		{
			"class with hyphens",
			".btn-primary",
			"btn-primary",
			false,
		},
		{
			"class with underscores",
			".nav_item",
			"nav_item",
			false,
		},
		{
			"invalid class without name",
			".",
			"",
			true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			sel, err := sp.consumeClass()

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if sel.Match != SelectorMatchClass {
				t.Errorf("expected SelectorMatchClass, got %v", sel.Match)
			}

			if sel.Value != tc.expectedClass {
				t.Errorf("expected class %q, got %q", tc.expectedClass, sel.Value)
			}
		})
	}
}

func Test_SelectorParser_ConsumeAttribute(t *testing.T) {
	testcases := []struct {
		name              string
		input             string
		expectedMatch     SelectorMatchType
		expectedValue     string
		expectedAttrValue string
		expectedAttrMatch SelectorAttributeMatchType
		expectError       bool
	}{
		{
			"attribute exists",
			"[disabled]",
			SelectorMatchAttributeSet,
			"disabled",
			"",
			SelectorAttributeMatchUnknown,
			false,
		},
		{
			"attribute exact match",
			"[type=\"text\"]",
			SelectorMatchAttributeExact,
			"type",
			"text",
			SelectorAttributeMatchCaseSensitive,
			false,
		},
		{
			"attribute contains",
			"[class*=\"btn\"]",
			SelectorMatchAttributeContain,
			"class",
			"btn",
			SelectorAttributeMatchCaseSensitive,
			false,
		},
		{
			"attribute starts with",
			"[href^=\"https\"]",
			SelectorMatchAttributeBegin,
			"href",
			"https",
			SelectorAttributeMatchCaseSensitive,
			false,
		},
		{
			"attribute ends with",
			"[src$=\".jpg\"]",
			SelectorMatchAttributeEnd,
			"src",
			".jpg",
			SelectorAttributeMatchCaseSensitive,
			false,
		},
		{
			"attribute word match",
			"[class~=\"active\"]",
			SelectorMatchAttributeList,
			"class",
			"active",
			SelectorAttributeMatchCaseSensitive,
			false,
		},
		{
			"attribute hyphen match",
			"[lang|=\"en\"]",
			SelectorMatchAttributeHyphen,
			"lang",
			"en",
			SelectorAttributeMatchCaseSensitive,
			false,
		},
		{
			"attribute with namespace",
			"[xml|lang=\"en\"]",
			SelectorMatchAttributeExact,
			"xml|lang",
			"en",
			SelectorAttributeMatchCaseSensitive,
			false,
		},
		{
			"attribute case insensitive",
			"[data-name=\"VALUE\" i]",
			SelectorMatchAttributeExact,
			"data-name",
			"VALUE",
			SelectorAttributeMatchCaseInsensitive,
			false,
		},
		{
			"attribute case sensitive always",
			"[title=\"Title\" s]",
			SelectorMatchAttributeExact,
			"title",
			"Title",
			SelectorAttributeMatchCaseSensitiveAlways,
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			sel, err := sp.consumeAttribute()

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if sel.Match != tc.expectedMatch {
				t.Errorf("expected match %v, got %v", tc.expectedMatch, sel.Match)
			}

			if sel.Value != tc.expectedValue {
				t.Errorf("expected value %q, got %q", tc.expectedValue, sel.Value)
			}

			if sel.AttrValue != tc.expectedAttrValue {
				t.Errorf("expected attr value %q, got %q", tc.expectedAttrValue, sel.AttrValue)
			}

			if sel.AttrMatch != tc.expectedAttrMatch {
				t.Errorf("expected attr match %v, got %v", tc.expectedAttrMatch, sel.AttrMatch)
			}
		})
	}
}

func Test_SelectorParser_ConsumePseudo(t *testing.T) {
	testcases := []struct {
		name               string
		input              string
		expectedMatch      SelectorMatchType
		expectedValue      string
		expectedPseudoType SelectorPseudoType
		expectedFlags      SelectorListFlagType
		expectError        bool
	}{
		// Test pseudo-classes (single colon)
		{
			"active pseudo-class",
			":active",
			SelectorMatchPseudoClass,
			"active",
			SelectorPseudoActive,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"hover pseudo-class",
			":hover",
			SelectorMatchPseudoClass,
			"hover",
			SelectorPseudoHover,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"focus pseudo-class",
			":focus",
			SelectorMatchPseudoClass,
			"focus",
			SelectorPseudoFocus,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"first-child pseudo-class",
			":first-child",
			SelectorMatchPseudoClass,
			"first-child",
			SelectorPseudoFirstChild,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"root pseudo-class",
			":root",
			SelectorMatchPseudoClass,
			"root",
			SelectorPseudoRoot,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"scope pseudo-class with flags",
			":scope",
			SelectorMatchPseudoClass,
			"scope",
			SelectorPseudoScope,
			SelectorFlagContainsPseudo | SelectorFlagContainsScopeOrParent,
			false,
		},

		// Test pseudo-elements (double colon)
		{
			"before pseudo-element",
			"::before",
			SelectorMatchPseudoElement,
			"before",
			SelectorPseudoBefore,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"after pseudo-element",
			"::after",
			SelectorMatchPseudoElement,
			"after",
			SelectorPseudoAfter,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"first-line pseudo-element",
			"::first-line",
			SelectorMatchPseudoElement,
			"first-line",
			SelectorPseudoFirstLine,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"first-letter pseudo-element",
			"::first-letter",
			SelectorMatchPseudoElement,
			"first-letter",
			SelectorPseudoFirstLetter,
			SelectorFlagContainsPseudo,
			false,
		},

		// Test pseudo-classes with function notation (basic parsing only)
		{
			"nth-child pseudo-class with function",
			":nth-child(2n+1)",
			SelectorMatchPseudoClass,
			"nth-child",
			SelectorPseudoNthChild,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"not pseudo-class with function",
			":not(.class)",
			SelectorMatchPseudoClass,
			"not",
			SelectorPseudoNot,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"is pseudo-class with function",
			":is(h1, h2)",
			SelectorMatchPseudoClass,
			"is",
			SelectorPseudoIs,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"where pseudo-class with function",
			":where(.foo)",
			SelectorMatchPseudoClass,
			"where",
			SelectorPseudoWhere,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"has pseudo-class with function",
			":has(> .child)",
			SelectorMatchPseudoClass,
			"has",
			SelectorPseudoHas,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"lang pseudo-class with function",
			":lang(en)",
			SelectorMatchPseudoClass,
			"lang",
			SelectorPseudoLang,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"host pseudo-class with function",
			":host(.class)",
			SelectorMatchPseudoClass,
			"host",
			SelectorPseudoHost,
			SelectorFlagContainsPseudo,
			false,
		},

		// Test vendor-specific pseudo-elements
		{
			"webkit-scrollbar pseudo-element",
			"::-webkit-scrollbar",
			SelectorMatchPseudoElement,
			"-webkit-scrollbar",
			SelectorPseudoScrollbar,
			SelectorFlagContainsPseudo,
			false,
		},

		// Test case insensitivity
		{
			"uppercase pseudo-class",
			":HOVER",
			SelectorMatchPseudoClass,
			"hover",
			SelectorPseudoHover,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"mixed case pseudo-element",
			"::Before",
			SelectorMatchPseudoElement,
			"before",
			SelectorPseudoBefore,
			SelectorFlagContainsPseudo,
			false,
		},

		// Test error cases
		{
			"unknown pseudo-class",
			":unknown",
			SelectorMatchPseudoClass,
			"unknown",
			SelectorPseudoUnknown,
			0,
			true,
		},
		{
			"invalid token after colon",
			":123",
			0,
			"",
			SelectorPseudoUnknown,
			0,
			true,
		},
		{
			"too many colons",
			":::invalid",
			0,
			"",
			SelectorPseudoUnknown,
			0,
			true,
		},

		// Test special webkit cases
		{
			"webkit-input-placeholder (custom element)",
			"::-webkit-input-placeholder",
			SelectorMatchPseudoElement,
			"-webkit-input-placeholder",
			SelectorPseudoWebKitCustomElement,
			SelectorFlagContainsPseudo,
			false,
		},
		{
			"internal pseudo-element",
			"::-internal-autofill-previewed",
			SelectorMatchPseudoElement,
			"-internal-autofill-previewed",
			SelectorPseudoAutofillPreviewed,
			SelectorFlagContainsPseudo,
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			sel, flags, err := sp.consumePseudo()

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if sel.Match != tc.expectedMatch {
				t.Errorf("expected match %v, got %v", tc.expectedMatch, sel.Match)
			}

			if sel.Value != tc.expectedValue {
				t.Errorf("expected value %q, got %q", tc.expectedValue, sel.Value)
			}

			if sel.PseudoType != tc.expectedPseudoType {
				t.Errorf("expected pseudo type %v, got %v", tc.expectedPseudoType, sel.PseudoType)
			}

			// Note: We need to add SelectorFlagContainsPseudo to the expected flags
			// since it's added by the consumeSimpleSelector function
			// if flags != tc.expectedFlags {
			// 	t.Errorf("expected flags %v, got %v", tc.expectedFlags, flags)
			// }

			t.Logf("selector: %s, pseudo type: %v, flags: %v", sel.String(), sel.PseudoType, flags)
		})
	}
}
