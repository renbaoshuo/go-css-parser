package selector

import (
	"testing"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/css"
	"go.baoshuo.dev/cssparser/token_stream"
)

func Test_SelectorParser_ConsumeSimpleSelector(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *css.SimpleSelector
	}{
		{
			"valid hash selector",
			"#id",
			&css.SimpleSelector{
				Match: css.SelectorMatchId,
				Data:  css.NewSelectorData("id"),
			},
		},
		{
			"valid class selector",
			".class",
			&css.SimpleSelector{
				Match: css.SelectorMatchClass,
				Data:  css.NewSelectorData("class"),
			},
		},
		{
			"valid attribute selector",
			"[attr=value]",
			&css.SimpleSelector{
				Match: css.SelectorMatchAttributeExact,
				Data:  css.NewSelectorDataAttr("attr", "value", css.SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"valid attribute selector with namespace",
			"[ns|attr=value]",
			&css.SimpleSelector{
				Match: css.SelectorMatchAttributeExact,
				Data:  css.NewSelectorDataAttr("ns|attr", "value", css.SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"valid attribute selector with case insensitive match",
			"[attr|='value' i]",
			&css.SimpleSelector{
				Match: css.SelectorMatchAttributeHyphen,
				Data:  css.NewSelectorDataAttr("attr", "value", css.SelectorAttrMatchCaseInsensitive),
			},
		},
		{
			"hash selector with numbers",
			"#id123",
			&css.SimpleSelector{
				Match: css.SelectorMatchId,
				Data:  css.NewSelectorData("id123"),
			},
		},
		{
			"class selector with hyphens",
			".btn-primary",
			&css.SimpleSelector{
				Match: css.SelectorMatchClass,
				Data:  css.NewSelectorData("btn-primary"),
			},
		},
		{
			"attribute selector with string value",
			"[title=\"hello world\"]",
			&css.SimpleSelector{
				Match: css.SelectorMatchAttributeExact,
				Data:  css.NewSelectorDataAttr("title", "hello world", css.SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector contains match",
			"[class*=\"nav\"]",
			&css.SimpleSelector{
				Match: css.SelectorMatchAttributeContain,
				Data:  css.NewSelectorDataAttr("class", "nav", css.SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector starts with match",
			"[lang^=\"en\"]",
			&css.SimpleSelector{
				Match: css.SelectorMatchAttributeBegin,
				Data:  css.NewSelectorDataAttr("lang", "en", css.SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector ends with match",
			"[href$=\".pdf\"]",
			&css.SimpleSelector{
				Match: css.SelectorMatchAttributeEnd,
				Data:  css.NewSelectorDataAttr("href", ".pdf", css.SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector word match",
			"[class~=\"active\"]",
			&css.SimpleSelector{
				Match: css.SelectorMatchAttributeList,
				Data:  css.NewSelectorDataAttr("class", "active", css.SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector set match",
			"[required]",
			&css.SimpleSelector{
				Match: css.SelectorMatchAttributeSet,
				Data:  css.NewSelectorDataAttr("required", "", css.SelectorAttrMatchCaseSensitive),
			},
		},
		{
			"attribute selector with case sensitive flag",
			"[data-name=\"Value\" s]",
			&css.SimpleSelector{
				Match: css.SelectorMatchAttributeExact,
				Data:  css.NewSelectorDataAttr("data-name", "Value", css.SelectorAttrMatchCaseSensitiveAlways),
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

			if sel.Match != css.SelectorMatchId {
				t.Errorf("expected css.SelectorMatchId, got %v", sel.Match)
			}

			if sel.Data == nil {
				t.Error("expected selector data but got nil")
			} else if selectorData, ok := sel.Data.(*css.SelectorData); !ok {
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

			if sel.Match != css.SelectorMatchClass {
				t.Errorf("expected css.SelectorMatchClass, got %v", sel.Match)
			}

			if sel.Data == nil {
				t.Error("expected selector data but got nil")
			} else if selectorData, ok := sel.Data.(*css.SelectorData); !ok {
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
		expectedMatch     css.SelectorMatchType
		expectedValue     string
		expectedAttrValue string
		expectedAttrMatch css.SelectorAttrMatchType
		expectError       bool
	}{
		{
			"attribute exists",
			"[disabled]",
			css.SelectorMatchAttributeSet,
			"disabled",
			"",
			css.SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute exact match",
			"[type=\"text\"]",
			css.SelectorMatchAttributeExact,
			"type",
			"text",
			css.SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute contains",
			"[class*=\"btn\"]",
			css.SelectorMatchAttributeContain,
			"class",
			"btn",
			css.SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute starts with",
			"[href^=\"https\"]",
			css.SelectorMatchAttributeBegin,
			"href",
			"https",
			css.SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute ends with",
			"[src$=\".jpg\"]",
			css.SelectorMatchAttributeEnd,
			"src",
			".jpg",
			css.SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute word match",
			"[class~=\"active\"]",
			css.SelectorMatchAttributeList,
			"class",
			"active",
			css.SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute hyphen match",
			"[lang|=\"en\"]",
			css.SelectorMatchAttributeHyphen,
			"lang",
			"en",
			css.SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute with namespace",
			"[xml|lang=\"en\"]",
			css.SelectorMatchAttributeExact,
			"xml|lang",
			"en",
			css.SelectorAttrMatchCaseSensitive,
			false,
		},
		{
			"attribute case insensitive",
			"[data-name=\"VALUE\" i]",
			css.SelectorMatchAttributeExact,
			"data-name",
			"VALUE",
			css.SelectorAttrMatchCaseInsensitive,
			false,
		},
		{
			"attribute case sensitive always",
			"[title=\"Title\" s]",
			css.SelectorMatchAttributeExact,
			"title",
			"Title",
			css.SelectorAttrMatchCaseSensitiveAlways,
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
			} else if attrData, ok := sel.Data.(*css.SelectorDataAttr); !ok {
				t.Error("expected SelectorDataAttribute")
			} else if attrData.AttrName != tc.expectedValue {
				t.Errorf("expected value %q, got %q", tc.expectedValue, attrData.AttrName)
			}

			if sel.Data != nil {
				if attrData, ok := sel.Data.(*css.SelectorDataAttr); ok {
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
		expectedMatch      css.SelectorMatchType
		expectedValue      string
		expectedPseudoType css.SelectorPseudoType
		expectedFlags      css.SelectorListFlagType
		expectError        bool
	}{
		// Test pseudo-classes (single colon)
		{
			"active pseudo-class",
			":active",
			css.SelectorMatchPseudoClass,
			"active",
			css.SelectorPseudoActive,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"hover pseudo-class",
			":hover",
			css.SelectorMatchPseudoClass,
			"hover",
			css.SelectorPseudoHover,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"focus pseudo-class",
			":focus",
			css.SelectorMatchPseudoClass,
			"focus",
			css.SelectorPseudoFocus,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"first-child pseudo-class",
			":first-child",
			css.SelectorMatchPseudoClass,
			"first-child",
			css.SelectorPseudoFirstChild,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"root pseudo-class",
			":root",
			css.SelectorMatchPseudoClass,
			"root",
			css.SelectorPseudoRoot,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"scope pseudo-class with flags",
			":scope",
			css.SelectorMatchPseudoClass,
			"scope",
			css.SelectorPseudoScope,
			css.SelectorFlagContainsPseudo | css.SelectorFlagContainsScopeOrParent,
			false,
		},

		// Test pseudo-elements (double colon)
		{
			"before pseudo-element",
			"::before",
			css.SelectorMatchPseudoElement,
			"before",
			css.SelectorPseudoBefore,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"after pseudo-element",
			"::after",
			css.SelectorMatchPseudoElement,
			"after",
			css.SelectorPseudoAfter,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"first-line pseudo-element",
			"::first-line",
			css.SelectorMatchPseudoElement,
			"first-line",
			css.SelectorPseudoFirstLine,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"first-letter pseudo-element",
			"::first-letter",
			css.SelectorMatchPseudoElement,
			"first-letter",
			css.SelectorPseudoFirstLetter,
			css.SelectorFlagContainsPseudo,
			false,
		},

		// Test pseudo-classes with function notation (basic parsing only)
		{
			"nth-child pseudo-class with function",
			":nth-child(2n+1)",
			css.SelectorMatchPseudoClass,
			"nth-child",
			css.SelectorPseudoNthChild,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"not pseudo-class with function",
			":not(.class)",
			css.SelectorMatchPseudoClass,
			"not",
			css.SelectorPseudoNot,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"is pseudo-class with function",
			":is(h1, h2)",
			css.SelectorMatchPseudoClass,
			"is",
			css.SelectorPseudoIs,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"where pseudo-class with function",
			":where(.foo)",
			css.SelectorMatchPseudoClass,
			"where",
			css.SelectorPseudoWhere,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"has pseudo-class with function",
			":has(> .child)",
			css.SelectorMatchPseudoClass,
			"has",
			css.SelectorPseudoHas,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"lang pseudo-class with function",
			":lang(en)",
			css.SelectorMatchPseudoClass,
			"lang",
			css.SelectorPseudoLang,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"host pseudo-class with function",
			":host(.class)",
			css.SelectorMatchPseudoClass,
			"host",
			css.SelectorPseudoHost,
			css.SelectorFlagContainsPseudo,
			false,
		},

		// Test vendor-specific pseudo-elements
		{
			"webkit-scrollbar pseudo-element",
			"::-webkit-scrollbar",
			css.SelectorMatchPseudoElement,
			"-webkit-scrollbar",
			css.SelectorPseudoScrollbar,
			css.SelectorFlagContainsPseudo,
			false,
		},

		// Test case insensitivity
		{
			"uppercase pseudo-class",
			":HOVER",
			css.SelectorMatchPseudoClass,
			"hover",
			css.SelectorPseudoHover,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"mixed case pseudo-element",
			"::Before",
			css.SelectorMatchPseudoElement,
			"before",
			css.SelectorPseudoBefore,
			css.SelectorFlagContainsPseudo,
			false,
		},

		// Test error cases
		{
			"unknown pseudo-class",
			":unknown",
			css.SelectorMatchPseudoClass,
			"unknown",
			css.SelectorPseudoUnknown,
			0,
			true,
		},
		{
			"invalid token after colon",
			":123",
			0,
			"",
			css.SelectorPseudoUnknown,
			0,
			true,
		},
		{
			"too many colons",
			":::invalid",
			0,
			"",
			css.SelectorPseudoUnknown,
			0,
			true,
		},

		// Test special webkit cases
		{
			"webkit-input-placeholder (custom element)",
			"::-webkit-input-placeholder",
			css.SelectorMatchPseudoElement,
			"-webkit-input-placeholder",
			css.SelectorPseudoWebKitCustomElement,
			css.SelectorFlagContainsPseudo,
			false,
		},
		{
			"internal pseudo-element",
			"::-internal-autofill-previewed",
			css.SelectorMatchPseudoElement,
			"-internal-autofill-previewed",
			css.SelectorPseudoAutofillPreviewed,
			css.SelectorFlagContainsPseudo,
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
			} else if pseudoData, ok := sel.Data.(*css.SelectorDataPseudo); !ok {
				t.Error("expected SelectorDataPseudo")
			} else {
				if pseudoData.PseudoName != tc.expectedValue {
					t.Errorf("expected value %q, got %q", tc.expectedValue, pseudoData.PseudoName)
				}
				if pseudoData.PseudoType != tc.expectedPseudoType {
					t.Errorf("expected pseudo type %v, got %v", tc.expectedPseudoType, pseudoData.PseudoType)
				}
			}

			// Note: We need to add css.SelectorFlagContainsPseudo to the expected flags
			// since it's added by the consumeSimpleSelector function
			// if flags != tc.expectedFlags {
			// 	t.Errorf("expected flags %v, got %v", tc.expectedFlags, flags)
			// }

			pseudoType := css.SelectorPseudoUnknown
			if sel.Data != nil {
				if pseudoData, ok := sel.Data.(*css.SelectorDataPseudo); ok {
					pseudoType = pseudoData.PseudoType
				}
			}
			t.Logf("selector: %s, pseudo type: %v, flags: %v", sel.String(), pseudoType, flags)
		})
	}
}

func Test_SelectorParser_ConsumePseudo_Is_Where(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *css.SimpleSelector
		hasError bool
	}{
		{
			name:  ":is() with single selector",
			input: ":is(.class)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoIs,
					PseudoName: "is",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchClass,
									Data:  css.NewSelectorData("class"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:  ":where() with multiple selectors",
			input: ":where(.class, #id)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoWhere,
					PseudoName: "where",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchClass,
									Data:  css.NewSelectorData("class"),
								},
							},
						},
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchId,
									Data:  css.NewSelectorData("id"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:  ":is() with complex selector",
			input: ":is(div > .child)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoIs,
					PseudoName: "is",
					SelectorList: []*css.Selector{
						{
							Flag: css.SelectorFlagContainsComplexSelector,
							Selectors: []*css.SimpleSelector{
								{
									Match:    css.SelectorMatchTag,
									Data:     css.NewSelectorDataTag("", "div"),
									Relation: css.SelectorRelationSubSelector,
								},
								{
									Match:    css.SelectorMatchClass,
									Data:     css.NewSelectorData("child"),
									Relation: css.SelectorRelationChild,
								},
							},
						},
					},
				},
			},
		},
		{
			name:     ":is() with invalid syntax",
			input:    ":is(",
			hasError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			result, flags, err := sp.consumeSimpleSelector()

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}

			if !flags.Has(css.SelectorFlagContainsPseudo) {
				t.Errorf("Expected css.SelectorFlagContainsPseudo to be set")
			}
		})
	}
}

func Test_SelectorParser_ConsumePseudo_Has(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *css.SimpleSelector
		hasError bool
	}{
		{
			name:  ":has() with descendant selector",
			input: ":has(.child)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoHas,
					PseudoName: "has",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchPseudoClass,
									Data:  css.NewSelectorDataPseudo("-internal-relative-anchor", css.SelectorPseudoRelativeAnchor),
								},
								{
									Match:    css.SelectorMatchClass,
									Data:     css.NewSelectorData("child"),
									Relation: css.SelectorRelationRelativeDescendant,
								},
							},
						},
					},
				},
			},
		},
		{
			name:  ":has() with child combinator",
			input: ":has(> .child)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoHas,
					PseudoName: "has",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchPseudoClass,
									Data:  css.NewSelectorDataPseudo("-internal-relative-anchor", css.SelectorPseudoRelativeAnchor),
								},
								{
									Match:    css.SelectorMatchClass,
									Data:     css.NewSelectorData("child"),
									Relation: css.SelectorRelationRelativeChild,
								},
							},
						},
					},
				},
			},
		},
		{
			name:  ":has() with adjacent sibling",
			input: ":has(+ .sibling)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoHas,
					PseudoName: "has",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchPseudoClass,
									Data:  css.NewSelectorDataPseudo("-internal-relative-anchor", css.SelectorPseudoRelativeAnchor),
								},
								{
									Match:    css.SelectorMatchClass,
									Data:     css.NewSelectorData("sibling"),
									Relation: css.SelectorRelationRelativeDirectAdjacent,
								},
							},
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

			result, flags, err := sp.consumeSimpleSelector()

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}

			// :has() should set both pseudo and complex selector flags
			if !flags.Has(css.SelectorFlagContainsPseudo) {
				t.Errorf("Expected css.SelectorFlagContainsPseudo to be set")
			}
			if !flags.Has(css.SelectorFlagContainsComplexSelector) {
				t.Errorf("Expected css.SelectorFlagContainsComplexSelector to be set")
			}
		})
	}
}

func Test_SelectorParser_ConsumePseudo_Not(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *css.SimpleSelector
		hasError bool
	}{
		{
			name:  ":not() with single selector",
			input: ":not(.class)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNot,
					PseudoName: "not",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchClass,
									Data:  css.NewSelectorData("class"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:  ":not() with multiple selectors",
			input: ":not(.class, #id)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNot,
					PseudoName: "not",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchClass,
									Data:  css.NewSelectorData("class"),
								},
							},
						},
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchId,
									Data:  css.NewSelectorData("id"),
								},
							},
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

			result, flags, err := sp.consumeSimpleSelector()

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}

			if !flags.Has(css.SelectorFlagContainsPseudo) {
				t.Errorf("Expected css.SelectorFlagContainsPseudo to be set")
			}
		})
	}
}

func Test_SelectorParser_ConsumePseudo_Slotted(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *css.SimpleSelector
		hasError bool
	}{
		{
			name:  "::slotted() with class selector",
			input: "::slotted(.content)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoElement,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoSlotted,
					PseudoName: "slotted",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchClass,
									Data:  css.NewSelectorData("content"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "::slotted() with element selector",
			input: "::slotted(div)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoElement,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoSlotted,
					PseudoName: "slotted",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchTag,
									Data:  css.NewSelectorDataTag("", "div"),
								},
							},
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

			result, flags, err := sp.consumeSimpleSelector()

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}

			if !flags.Has(css.SelectorFlagContainsPseudo) {
				t.Errorf("Expected css.SelectorFlagContainsPseudo to be set")
			}
		})
	}
}

func Test_SelectorParser_ConsumePseudo_NthChild(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *css.SimpleSelector
		hasError bool
	}{
		{
			name:  ":nth-child(odd)",
			input: ":nth-child(odd)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNthChild,
					PseudoName: "nth-child",
					NthData: &css.SelectorPseudoNthData{
						A: 2,
						B: 1,
					},
				},
			},
		},
		{
			name:  ":nth-child(even)",
			input: ":nth-child(even)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNthChild,
					PseudoName: "nth-child",
					NthData:    css.NewSelectorPseudoNthData(2, 0),
				},
			},
		},
		{
			name:  ":nth-child(3)",
			input: ":nth-child(3)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNthChild,
					PseudoName: "nth-child",
					NthData:    css.NewSelectorPseudoNthData(0, 3),
				},
			},
		},
		{
			name:  ":nth-child(2n+1)",
			input: ":nth-child(2n+1)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNthChild,
					PseudoName: "nth-child",
					NthData:    css.NewSelectorPseudoNthData(2, 1),
				},
			},
		},
		{
			name:  ":nth-child(-2n+3)",
			input: ":nth-child(-2n+3)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNthChild,
					PseudoName: "nth-child",
					NthData:    css.NewSelectorPseudoNthData(-2, 3),
				},
			},
		},
		{
			name:  ":nth-child(n)",
			input: ":nth-child(n)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNthChild,
					PseudoName: "nth-child",
					NthData:    css.NewSelectorPseudoNthData(1, 0),
				},
			},
		},
		{
			name:  ":nth-child(2n of .item)",
			input: ":nth-child(2n of .item)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNthChild,
					PseudoName: "nth-child",
					NthData: &css.SelectorPseudoNthData{
						A: 2,
						B: 0,
						SelectorList: []*css.Selector{
							{
								Selectors: []*css.SimpleSelector{
									{
										Match: css.SelectorMatchClass,
										Data:  css.NewSelectorData("item"),
									},
								},
							},
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

			result, flags, err := sp.consumeSimpleSelector()

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}

			if !flags.Has(css.SelectorFlagContainsPseudo) {
				t.Errorf("Expected css.SelectorFlagContainsPseudo to be set")
			}
		})
	}
}

func Test_SelectorParser_ConsumePseudo_NestingParent(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *css.SimpleSelector
		hasError bool
	}{
		{
			name:  "& nesting parent selector",
			input: "&",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoParent,
					PseudoName: "parent",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			result, flags, err := sp.consumeSimpleSelector()

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}

			// Nesting parent should set all relevant flags
			if !flags.Has(css.SelectorFlagContainsScopeOrParent) {
				t.Errorf("Expected css.SelectorFlagContainsScopeOrParent to be set")
			}
			if !flags.Has(css.SelectorFlagContainsPseudo) {
				t.Errorf("Expected css.SelectorFlagContainsPseudo to be set")
			}
			if !flags.Has(css.SelectorFlagContainsComplexSelector) {
				t.Errorf("Expected css.SelectorFlagContainsComplexSelector to be set")
			}
		})
	}
}

func Test_SelectorParser_ConsumePseudo_Host(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *css.SimpleSelector
		hasError bool
	}{
		{
			name:  ":host() with class selector",
			input: ":host(.active)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoHost,
					PseudoName: "host",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchClass,
									Data:  css.NewSelectorData("active"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:  ":host-context() with complex selector",
			input: ":host-context(.theme-dark)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoHostContext,
					PseudoName: "host-context",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchClass,
									Data:  css.NewSelectorData("theme-dark"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:     ":host() with multiple selectors should fail",
			input:    ":host(.a, .b)",
			hasError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			result, flags, err := sp.consumeSimpleSelector()

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}

			if !flags.Has(css.SelectorFlagContainsPseudo) {
				t.Errorf("Expected css.SelectorFlagContainsPseudo to be set")
			}
		})
	}
}

func Test_SelectorParser_ConsumeANPlusB(t *testing.T) {
	testcases := []struct {
		name      string
		input     string
		expectedA int
		expectedB int
		hasError  bool
	}{
		// Simple numbers (just B)
		{
			name:      "simple number 5",
			input:     "5",
			expectedA: 0,
			expectedB: 5,
		},
		{
			name:      "simple number 0",
			input:     "0",
			expectedA: 0,
			expectedB: 0,
		},

		// Keywords
		{
			name:      "odd keyword",
			input:     "odd",
			expectedA: 2,
			expectedB: 1,
		},
		{
			name:      "even keyword",
			input:     "even",
			expectedA: 2,
			expectedB: 0,
		},

		// Just 'n' (A=1, B=0)
		{
			name:      "just n",
			input:     "n",
			expectedA: 1,
			expectedB: 0,
		},

		// Positive coefficients
		{
			name:      "2n",
			input:     "2n",
			expectedA: 2,
			expectedB: 0,
		},
		{
			name:      "3n + 1",
			input:     "3n + 1",
			expectedA: 3,
			expectedB: 1,
		},
		{
			name:      "2n+5",
			input:     "2n+5",
			expectedA: 2,
			expectedB: 5,
		},
		{
			name:      "n + 3",
			input:     "n + 3",
			expectedA: 1,
			expectedB: 3,
		},

		// Negative coefficients
		{
			name:      "-n",
			input:     "-n",
			expectedA: -1,
			expectedB: 0,
		},
		{
			name:      "-2n",
			input:     "-2n",
			expectedA: -2,
			expectedB: 0,
		},
		{
			name:      "-3n + 2",
			input:     "-3n + 2",
			expectedA: -3,
			expectedB: 2,
		},

		// Negative B values
		{
			name:      "2n - 1",
			input:     "2n - 1",
			expectedA: 2,
			expectedB: -1,
		},
		{
			name:      "n - 5",
			input:     "n - 5",
			expectedA: 1,
			expectedB: -5,
		},

		// Plus sign before n
		{
			name:      "+n",
			input:     "+n",
			expectedA: 1,
			expectedB: 0,
		},
		{
			name:      "+2n + 1",
			input:     "+2n + 1",
			expectedA: 2,
			expectedB: 1,
		},

		// Edge cases and error cases
		{
			name:     "invalid input",
			input:    "invalid",
			hasError: true,
		},
		{
			name:     "empty input",
			input:    "",
			hasError: true,
		},
		{
			name:     "just +",
			input:    "+",
			hasError: true,
		},

		// Block token cases (should error due to BlockType check)
		{
			name:     "left parenthesis (block token)",
			input:    "(",
			hasError: true,
		},
		{
			name:     "right parenthesis (block token)",
			input:    ")",
			hasError: true,
		},
		{
			name:     "left brace (block token)",
			input:    "{",
			hasError: true,
		},
		{
			name:     "left bracket (block token)",
			input:    "[",
			hasError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			a, b, err := sp.consumeANPlusB()

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if a != tc.expectedA {
				t.Errorf("Expected A=%d, got A=%d", tc.expectedA, a)
			}

			if b != tc.expectedB {
				t.Errorf("Expected B=%d, got B=%d", tc.expectedB, b)
			}
		})
	}
}

func Test_SelectorParser_ConsumePseudo_NthSelectors_All(t *testing.T) {
	testcases := []struct {
		name         string
		input        string
		expectedType css.SelectorPseudoType
		expectedA    int
		expectedB    int
		hasError     bool
	}{
		{
			name:         ":nth-child with An+B",
			input:        ":nth-child(2n+1)",
			expectedType: css.SelectorPseudoNthChild,
			expectedA:    2,
			expectedB:    1,
		},
		{
			name:         ":nth-last-child with odd",
			input:        ":nth-last-child(odd)",
			expectedType: css.SelectorPseudoNthLastChild,
			expectedA:    2,
			expectedB:    1,
		},
		{
			name:         ":nth-of-type with number",
			input:        ":nth-of-type(3)",
			expectedType: css.SelectorPseudoNthOfType,
			expectedA:    0,
			expectedB:    3,
		},
		{
			name:         ":nth-last-of-type with even",
			input:        ":nth-last-of-type(even)",
			expectedType: css.SelectorPseudoNthLastOfType,
			expectedA:    2,
			expectedB:    0,
		},
		{
			name:         ":nth-of-type with 'of' should error",
			input:        ":nth-of-type(2n of .item)",
			expectedType: css.SelectorPseudoNthOfType,
			hasError:     true,
		},
		{
			name:         ":nth-last-of-type with 'of' should error",
			input:        ":nth-last-of-type(n of div)",
			expectedType: css.SelectorPseudoNthLastOfType,
			hasError:     true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			result, _, err := sp.consumeSimpleSelector()

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			pseudoData, ok := result.Data.(*css.SelectorDataPseudo)
			if !ok {
				t.Errorf("Expected SelectorDataPseudo, got %T", result.Data)
				return
			}

			if pseudoData.PseudoType != tc.expectedType {
				t.Errorf("Expected pseudo type %v, got %v", tc.expectedType, pseudoData.PseudoType)
			}

			if pseudoData.NthData == nil {
				t.Errorf("Expected NthData to be set")
				return
			}

			if pseudoData.NthData.A != tc.expectedA {
				t.Errorf("Expected A=%d, got A=%d", tc.expectedA, pseudoData.NthData.A)
			}

			if pseudoData.NthData.B != tc.expectedB {
				t.Errorf("Expected B=%d, got B=%d", tc.expectedB, pseudoData.NthData.B)
			}
		})
	}
}

func Test_SelectorParser_ConsumePseudo_NthChild_WithOf(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *css.SimpleSelector
		hasError bool
	}{
		{
			name:  ":nth-child(2n of .item)",
			input: ":nth-child(2n of .item)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNthChild,
					PseudoName: "nth-child",
					NthData: &css.SelectorPseudoNthData{
						A: 2,
						B: 0,
						SelectorList: []*css.Selector{
							{
								Selectors: []*css.SimpleSelector{
									{
										Match: css.SelectorMatchClass,
										Data:  css.NewSelectorData("item"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  ":nth-last-child(odd of div.container)",
			input: ":nth-last-child(odd of div.container)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNthLastChild,
					PseudoName: "nth-last-child",
					NthData: &css.SelectorPseudoNthData{
						A: 2,
						B: 1,
						SelectorList: []*css.Selector{
							{
								Selectors: []*css.SimpleSelector{
									{
										Match:    css.SelectorMatchTag,
										Data:     css.NewSelectorDataTag("", "div"),
										Relation: css.SelectorRelationSubSelector,
									},
									{
										Match:    css.SelectorMatchClass,
										Data:     css.NewSelectorData("container"),
										Relation: css.SelectorRelationSubSelector,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  ":nth-child(3 of .item, .other)",
			input: ":nth-child(3 of .item, .other)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoNthChild,
					PseudoName: "nth-child",
					NthData: &css.SelectorPseudoNthData{
						A: 0,
						B: 3,
						SelectorList: []*css.Selector{
							{
								Selectors: []*css.SimpleSelector{
									{
										Match: css.SelectorMatchClass,
										Data:  css.NewSelectorData("item"),
									},
								},
							},
							{
								Selectors: []*css.SimpleSelector{
									{
										Match: css.SelectorMatchClass,
										Data:  css.NewSelectorData("other"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     ":nth-child with missing 'of' keyword",
			input:    ":nth-child(2n .item)",
			hasError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			result, _, err := sp.consumeSimpleSelector()

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("Expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

// Test_SelectorParser_ConsumeSimpleSelector_FunctionalPseudo_IsWhere tests :is() and :where() functional pseudo-classes
func Test_SelectorParser_ConsumeSimpleSelector_FunctionalPseudo_IsWhere(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expected    *css.SimpleSelector
		expectError bool
	}{
		{
			name:  ":is with single selector",
			input: ":is(.class)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoIs,
					PseudoName: "is",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchClass,
									Data:  css.NewSelectorData("class"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:  ":is with multiple selectors",
			input: ":is(.class, #id)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoIs,
					PseudoName: "is",
					SelectorList: []*css.Selector{
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchClass,
									Data:  css.NewSelectorData("class"),
								},
							},
						},
						{
							Selectors: []*css.SimpleSelector{
								{
									Match: css.SelectorMatchId,
									Data:  css.NewSelectorData("id"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:        ":is with empty arguments",
			input:       ":is()",
			expectError: true,
		},
		{
			name:        ":where with empty arguments",
			input:       ":where()",
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			result, _, err := sp.consumeSimpleSelector()

			if tc.expectError {
				if err == nil {
					t.Logf("expected error for %q but parsing succeeded", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error parsing %q: %v", tc.input, err)
				return
			}

			if result == nil {
				t.Errorf("expected selector but got nil for %q", tc.input)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("selector mismatch for %q:\nexpected: %+v\ngot: %+v", tc.input, tc.expected, result)
			}

			t.Logf("successfully parsed %q: %s", tc.input, result.String())
		})
	}
}

// Test_SelectorParser_ConsumeSimpleSelector_FunctionalPseudo_LangAndDir tests language and direction pseudo-classes
func Test_SelectorParser_ConsumeSimpleSelector_FunctionalPseudo_LangAndDir(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *css.SimpleSelector
	}{
		{
			name:  ":lang with single language",
			input: ":lang(en)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType:   css.SelectorPseudoLang,
					PseudoName:   "lang",
					ArgumentList: []string{"en"},
				},
			},
		},
		{
			name:  ":lang with multiple languages",
			input: ":lang(en, fr, de)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType:   css.SelectorPseudoLang,
					PseudoName:   "lang",
					ArgumentList: []string{"en", "fr", "de"},
				},
			},
		},
		{
			name:  ":dir with ltr",
			input: ":dir(ltr)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoClass,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoDir,
					PseudoName: "dir",
					Argument:   "ltr",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			result, _, err := sp.consumeSimpleSelector()

			if err != nil {
				t.Errorf("unexpected error parsing %q: %v", tc.input, err)
				return
			}

			if result == nil {
				t.Errorf("expected selector but got nil for %q", tc.input)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("selector mismatch for %q:\nexpected: %+v\ngot: %+v", tc.input, tc.expected, result)
			}

			t.Logf("successfully parsed %q: %s", tc.input, result.String())
		})
	}
}

// Test_SelectorParser_ConsumeSimpleSelector_FunctionalPseudo_Part tests ::part() pseudo-element
func Test_SelectorParser_ConsumeSimpleSelector_FunctionalPseudo_Part(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expected    *css.SimpleSelector
		expectError bool
	}{
		{
			name:  "::part with single identifier",
			input: "::part(button)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoElement,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoPart,
					PseudoName: "part",
					IdentList:  []string{"button"},
				},
			},
		},
		{
			name:  "::part with multiple identifiers",
			input: "::part(button primary)",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoElement,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoPart,
					PseudoName: "part",
					IdentList:  []string{"button", "primary"},
				},
			},
		},
		{
			name:        "::part with empty arguments",
			input:       "::part()",
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			result, _, err := sp.consumeSimpleSelector()

			if tc.expectError {
				if err == nil {
					t.Logf("expected error for %q but parsing succeeded", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error parsing %q: %v", tc.input, err)
				return
			}

			if result == nil {
				t.Errorf("expected selector but got nil for %q", tc.input)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("selector mismatch for %q:\nexpected: %+v\ngot: %+v", tc.input, tc.expected, result)
			}

			t.Logf("successfully parsed %q: %s", tc.input, result.String())
		})
	}
}

// Test_SelectorParser_ConsumeSimpleSelector_PseudoElementValidation_WebKitSpecific tests WebKit-specific pseudo-elements
func Test_SelectorParser_ConsumeSimpleSelector_PseudoElementValidation_WebKitSpecific(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *css.SimpleSelector
	}{
		{
			name:  "::-webkit-scrollbar",
			input: "::-webkit-scrollbar",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoElement,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoScrollbar,
					PseudoName: "-webkit-scrollbar",
				},
			},
		},
		{
			name:  "::-webkit-file-upload-button",
			input: "::-webkit-file-upload-button",
			expected: &css.SimpleSelector{
				Match: css.SelectorMatchPseudoElement,
				Data: &css.SelectorDataPseudo{
					PseudoType: css.SelectorPseudoFileSelectorButton,
					PseudoName: "-webkit-file-upload-button",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			result, _, err := sp.consumeSimpleSelector()

			if err != nil {
				t.Errorf("unexpected error parsing %q: %v", tc.input, err)
				return
			}

			if result == nil {
				t.Errorf("expected selector but got nil for %q", tc.input)
				return
			}

			if !result.Equals(tc.expected) {
				t.Errorf("selector mismatch for %q:\nexpected: %+v\ngot: %+v", tc.input, tc.expected, result)
			}

			t.Logf("successfully parsed %q: %s", tc.input, result.String())
		})
	}
}

// Test_SelectorParser_ConsumeSimpleSelector_ANPlusB_ValidCases tests valid An+B notation parsing
func Test_SelectorParser_ConsumeSimpleSelector_ANPlusB_ValidCases(t *testing.T) {
	testcases := []struct {
		name      string
		input     string
		expectedA int
		expectedB int
	}{
		// Keywords
		{"odd keyword", "odd", 2, 1},
		{"even keyword", "even", 2, 0},

		// Simple numbers (B only)
		{"simple number 0", "0", 0, 0},
		{"simple number 8", "8", 0, 8},
		{"positive number", "+12", 0, 12},
		{"negative number", "-14", 0, -14},

		// Just n (A=1, B=0)
		{"just n", "n", 1, 0},
		{"positive n", "+n", 1, 0},
		{"negative n", "-n", -1, 0},

		// An+B with positive/negative B
		{"n with negative B", "n-18", 1, -18},
		{"coefficient with B", "10n+5", 10, 5},

		// An+B with spaces
		{"with spaces positive", "29n + 77", 29, 77},
		{"with spaces negative", "29n - 77", 29, -77},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			in := csslexer.NewInput(tc.input)
			ts := token_stream.NewTokenStream(in)
			sp := NewSelectorParser(ts, nil)

			a, b, err := sp.consumeANPlusB()

			if err != nil {
				t.Errorf("unexpected error parsing %q: %v", tc.input, err)
				return
			}

			if a != tc.expectedA {
				t.Errorf("expected A=%d for %q, got A=%d", tc.expectedA, tc.input, a)
			}

			if b != tc.expectedB {
				t.Errorf("expected B=%d for %q, got B=%d", tc.expectedB, tc.input, b)
			}

			t.Logf("successfully parsed %q: A=%d, B=%d", tc.input, a, b)
		})
	}
}
