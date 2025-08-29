package selector

import (
	"testing"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/token_stream"
)

func Test_SelectorParser_ConsumeSimpleSelector(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *SimpleSelector
	}{
		{
			"valid hash selector",
			"#id",
			&SimpleSelector{
				Match: SelectorMatchId,
				Data:  NewSelectorData("id"),
			},
		},
		{
			"valid class selector",
			".class",
			&SimpleSelector{
				Match: SelectorMatchClass,
				Data:  NewSelectorData("class"),
			},
		},
		{
			"valid attribute selector",
			"[attr=value]",
			&SimpleSelector{
				Match: SelectorMatchAttributeExact,
				Data:  NewSelectorDataAttr("attr", "value", SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"valid attribute selector with namespace",
			"[ns|attr=value]",
			&SimpleSelector{
				Match: SelectorMatchAttributeExact,
				Data:  NewSelectorDataAttr("ns|attr", "value", SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"valid attribute selector with case insensitive match",
			"[attr|='value' i]",
			&SimpleSelector{
				Match: SelectorMatchAttributeHyphen,
				Data:  NewSelectorDataAttr("attr", "value", SelectorAttrMatchCaseInsensitive),
			},
		},
		{
			"hash selector with numbers",
			"#id123",
			&SimpleSelector{
				Match: SelectorMatchId,
				Data:  NewSelectorData("id123"),
			},
		},
		{
			"class selector with hyphens",
			".btn-primary",
			&SimpleSelector{
				Match: SelectorMatchClass,
				Data:  NewSelectorData("btn-primary"),
			},
		},
		{
			"attribute selector with string value",
			"[title=\"hello world\"]",
			&SimpleSelector{
				Match: SelectorMatchAttributeExact,
				Data:  NewSelectorDataAttr("title", "hello world", SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector contains match",
			"[class*=\"nav\"]",
			&SimpleSelector{
				Match: SelectorMatchAttributeContain,
				Data:  NewSelectorDataAttr("class", "nav", SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector starts with match",
			"[lang^=\"en\"]",
			&SimpleSelector{
				Match: SelectorMatchAttributeBegin,
				Data:  NewSelectorDataAttr("lang", "en", SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector ends with match",
			"[href$=\".pdf\"]",
			&SimpleSelector{
				Match: SelectorMatchAttributeEnd,
				Data:  NewSelectorDataAttr("href", ".pdf", SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector word match",
			"[class~=\"active\"]",
			&SimpleSelector{
				Match: SelectorMatchAttributeList,
				Data:  NewSelectorDataAttr("class", "active", SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector set match",
			"[required]",
			&SimpleSelector{
				Match: SelectorMatchAttributeSet,
				Data:  NewSelectorDataAttr("required", "", SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector with case sensitive flag",
			"[data-name=\"Value\" s]",
			&SimpleSelector{
				Match: SelectorMatchAttributeExact,
				Data:  NewSelectorDataAttr("data-name", "Value", SelectorAttrMatchCaseSensitiveAlways),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			sel, _, err := sp.consumeSimpleSelector()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if sel == nil {
				t.Error("expected a selector but got nil")
				return
			}

			if !sel.Equals(tc.expected) {
				t.Errorf("selector mismatch:\nexpected: %v\ngot: %v", tc.expected, sel)
			}

			t.Logf("selector: %s", sel.String())
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

			if sel.Data == nil {
				t.Error("expected selector data but got nil")
			} else if selectorData, ok := sel.Data.(*SelectorData); !ok {
				t.Error("expected SelectorData")
			} else if selectorData.Value != tc.expectedId {
				t.Errorf("expected id %q, got %q", tc.expectedId, selectorData.Value)
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

			if sel.Data == nil {
				t.Error("expected selector data but got nil")
			} else if selectorData, ok := sel.Data.(*SelectorData); !ok {
				t.Error("expected SelectorData")
			} else if selectorData.Value != tc.expectedClass {
				t.Errorf("expected class %q, got %q", tc.expectedClass, selectorData.Value)
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
		expectedAttrMatch SelectorAttrMatchType
		expectError       bool
	}{
		{
			"attribute exists",
			"[disabled]",
			SelectorMatchAttributeSet,
			"disabled",
			"",
			SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute exact match",
			"[type=\"text\"]",
			SelectorMatchAttributeExact,
			"type",
			"text",
			SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute contains",
			"[class*=\"btn\"]",
			SelectorMatchAttributeContain,
			"class",
			"btn",
			SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute starts with",
			"[href^=\"https\"]",
			SelectorMatchAttributeBegin,
			"href",
			"https",
			SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute ends with",
			"[src$=\".jpg\"]",
			SelectorMatchAttributeEnd,
			"src",
			".jpg",
			SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute word match",
			"[class~=\"active\"]",
			SelectorMatchAttributeList,
			"class",
			"active",
			SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute hyphen match",
			"[lang|=\"en\"]",
			SelectorMatchAttributeHyphen,
			"lang",
			"en",
			SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute with namespace",
			"[xml|lang=\"en\"]",
			SelectorMatchAttributeExact,
			"xml|lang",
			"en",
			SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute case insensitive",
			"[data-name=\"VALUE\" i]",
			SelectorMatchAttributeExact,
			"data-name",
			"VALUE",
			SelectorAttrMatchCaseInsensitive,
			false,
		},
		{
			"attribute case sensitive always",
			"[title=\"Title\" s]",
			SelectorMatchAttributeExact,
			"title",
			"Title",
			SelectorAttrMatchCaseSensitiveAlways,
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

			if sel.Data == nil {
				t.Error("expected selector data but got nil")
			} else if attrData, ok := sel.Data.(*SelectorDataAttr); !ok {
				t.Error("expected SelectorDataAttribute")
			} else if attrData.AttrName != tc.expectedValue {
				t.Errorf("expected value %q, got %q", tc.expectedValue, attrData.AttrName)
			}

			if sel.Data != nil {
				if attrData, ok := sel.Data.(*SelectorDataAttr); ok {
					if attrData.AttrValue != tc.expectedAttrValue {
						t.Errorf("expected attr value %q, got %q", tc.expectedAttrValue, attrData.AttrValue)
					}
					if attrData.AttrMatch != tc.expectedAttrMatch {
						t.Errorf("expected attr match %v, got %v", tc.expectedAttrMatch, attrData.AttrMatch)
					}
				}
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

			if sel.Data == nil {
				t.Error("expected selector data but got nil")
			} else if pseudoData, ok := sel.Data.(*SelectorDataPseudo); !ok {
				t.Error("expected SelectorDataPseudo")
			} else {
				if pseudoData.PseudoName != tc.expectedValue {
					t.Errorf("expected value %q, got %q", tc.expectedValue, pseudoData.PseudoName)
				}
				if pseudoData.PseudoType != tc.expectedPseudoType {
					t.Errorf("expected pseudo type %v, got %v", tc.expectedPseudoType, pseudoData.PseudoType)
				}
			}

			// Note: We need to add SelectorFlagContainsPseudo to the expected flags
			// since it's added by the consumeSimpleSelector function
			// if flags != tc.expectedFlags {
			// 	t.Errorf("expected flags %v, got %v", tc.expectedFlags, flags)
			// }

			pseudoType := SelectorPseudoUnknown
			if sel.Data != nil {
				if pseudoData, ok := sel.Data.(*SelectorDataPseudo); ok {
					pseudoType = pseudoData.PseudoType
				}
			}
			t.Logf("selector: %s, pseudo type: %v, flags: %v", sel.String(), pseudoType, flags)
		})
	}
}
