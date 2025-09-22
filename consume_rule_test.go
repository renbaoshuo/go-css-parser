package cssparser

import (
	"testing"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/css"
	"go.baoshuo.dev/cssparser/nesting"
	"go.baoshuo.dev/cssparser/token_stream"
)

func TestParser_ConsumeQualifiedRule(t *testing.T) {
	testcases := []struct {
		name         string
		input        string
		allowedRules allowedRuleType
		nestingType  nesting.NestingTypeType
		expectError  bool
		expected     *css.StyleRule
	}{
		{
			name:         "style rule allowed with class selector",
			input:        ".test { color: red; }",
			allowedRules: qualifiedRuleTypeStyle,
			nestingType:  nesting.NestingTypeNone,
			expectError:  false,
			expected: &css.StyleRule{
				Type: css.StyleRuleTypeQualifiedRule,
				Selectors: []*css.Selector{
					{
						Flag: 2,
						Selectors: []*css.SimpleSelector{
							{
								Match:    css.SelectorMatchTag,
								Data:     css.NewSelectorDataTag("", "test"),
								Relation: css.SelectorRelationSubSelector,
							},
						},
					},
				},
				Declarations: []*css.Declaration{
					{Property: "color", Value: "red", Important: false},
				},
				Rules: []*css.GenericRule{},
			},
		},
		{
			name:         "keyframes rule allowed",
			input:        "0% { opacity: 0; }",
			allowedRules: qualifiedRuleTypeKeyframes,
			nestingType:  nesting.NestingTypeNone,
			expectError:  true, // not implemented yet
		},
		{
			name:         "no allowed rules",
			input:        ".test { color: red; }",
			allowedRules: 0,
			nestingType:  nesting.NestingTypeNone,
			expectError:  true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			input := csslexer.NewInput(tc.input)
			parser := NewParser(input)

			styleRule, err := parser.consumeQualifiedRule(tc.allowedRules, tc.nestingType, nil)

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

			if styleRule == nil {
				t.Errorf("expected styleRule but got nil")
				return
			}

			if !styleRule.Equals(tc.expected) {
				t.Errorf("styleRule mismatch:\nexpected: %+v\ngot: %+v", tc.expected, styleRule)
			}
		})
	}
}

func TestParser_ConsumeStyleRule(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		nestingType nesting.NestingTypeType
		expectError bool
		expected    *css.StyleRule
	}{
		{
			name:        "class selector",
			input:       ".container { margin: 10px; padding: 5px; }",
			nestingType: nesting.NestingTypeNone,
			expectError: false,
			expected: &css.StyleRule{
				Type: css.StyleRuleTypeQualifiedRule,
				Selectors: []*css.Selector{
					{
						Flag: 2,
						Selectors: []*css.SimpleSelector{
							{
								Match:    css.SelectorMatchTag,
								Data:     css.NewSelectorDataTag("", "container"),
								Relation: css.SelectorRelationSubSelector,
							},
						},
					},
				},
				Declarations: []*css.Declaration{
					{Property: "margin", Value: "10px", Important: false},
					{Property: "padding", Value: "5px", Important: false},
				},
				Rules: []*css.GenericRule{},
			},
		},
		{
			name:        "missing opening brace",
			input:       "div color: red; }",
			nestingType: nesting.NestingTypeNone,
			expectError: true,
		},
		{
			name:        "custom property ambiguity",
			input:       "--property: value;",
			nestingType: nesting.NestingTypeNone,
			expectError: true,
		},
		{
			name:        "nested custom property ambiguity",
			input:       "--property: value;",
			nestingType: nesting.NestingTypeNesting,
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			input := csslexer.NewInput(tc.input)
			parser := NewParser(input)

			styleRule, err := parser.consumeStyleRule(tc.nestingType, nil, false)

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

			if styleRule == nil {
				t.Errorf("expected styleRule but got nil")
				return
			}

			if !styleRule.Equals(tc.expected) {
				t.Errorf("styleRule mismatch:\nexpected: %+v\ngot: %+v", tc.expected, styleRule)
			}
		})
	}
}

func TestParser_ConsumeDeclaration(t *testing.T) {
	testcases := []struct {
		name        string
		input       string
		expectError bool
		expected    *css.Declaration
	}{
		{
			name:        "simple property",
			input:       "color: red",
			expectError: false,
			expected: &css.Declaration{
				Property:  "color",
				Value:     "red",
				Important: false,
			},
		},
		{
			name:        "property with important",
			input:       "margin: 10px !important",
			expectError: false,
			expected: &css.Declaration{
				Property:  "margin",
				Value:     "10px",
				Important: true,
			},
		},
		{
			name:        "property with semicolon",
			input:       "padding: 5px;",
			expectError: false,
			expected: &css.Declaration{
				Property:  "padding",
				Value:     "5px",
				Important: false,
			},
		},
		{
			name:        "complex value",
			input:       "background: url('image.jpg') no-repeat center",
			expectError: false,
			expected: &css.Declaration{
				Property:  "background",
				Value:     "urlimage.jpg) no-repeat center",
				Important: false,
			},
		},
		{
			name:        "custom property",
			input:       "--main-color: #ff0000",
			expectError: false,
			expected: &css.Declaration{
				Property:  "--main-color",
				Value:     "ff0000",
				Important: false,
			},
		},
		{
			name:        "missing colon",
			input:       "color red",
			expectError: true,
		},
		{
			name:        "missing value",
			input:       "color:",
			expectError: true,
		},
		{
			name:        "not an identifier",
			input:       "123: red",
			expectError: true,
		},
		{
			name:        "exclamation without important",
			input:       "color: red ! notimportant",
			expectError: false,
			expected: &css.Declaration{
				Property:  "color",
				Value:     "red !notimportant",
				Important: false,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			input := csslexer.NewInput(tc.input)
			parser := NewParser(input)

			decl, err := parser.consumeDeclaration()

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

			if decl == nil {
				t.Errorf("expected declaration but got nil")
				return
			}

			if !decl.Equals(tc.expected) {
				t.Errorf("declaration mismatch:\nexpected: %+v\ngot: %+v", tc.expected, decl)
			}
		})
	}
}

func TestParser_ConsumeStyleRuleContents(t *testing.T) {
	testcases := []struct {
		name                 string
		input                string
		nestingType          nesting.NestingTypeType
		expectedDeclarations int
		expectedChildRules   int
	}{
		{
			name:                 "simple declarations",
			input:                "color: red; margin: 10px;",
			nestingType:          nesting.NestingTypeNone,
			expectedDeclarations: 2,
			expectedChildRules:   0,
		},
		{
			name:                 "declarations with semicolons",
			input:                "color: blue;; padding: 5px;;;",
			nestingType:          nesting.NestingTypeNone,
			expectedDeclarations: 2,
			expectedChildRules:   0,
		},
		{
			name:                 "empty content",
			input:                "",
			nestingType:          nesting.NestingTypeNone,
			expectedDeclarations: 0,
			expectedChildRules:   0,
		},
		{
			name:                 "only whitespace",
			input:                "   \n  \t  ",
			nestingType:          nesting.NestingTypeNone,
			expectedDeclarations: 0,
			expectedChildRules:   0,
		},
		{
			name:                 "declarations with comments",
			input:                "/* comment */ color: red; /* another comment */ margin: 10px;",
			nestingType:          nesting.NestingTypeNone,
			expectedDeclarations: 2,
			expectedChildRules:   0,
		},
		{
			name:                 "mixed valid and invalid declarations",
			input:                "color: red; invalid-property-without-colon; padding: 5px;",
			nestingType:          nesting.NestingTypeNone,
			expectedDeclarations: 2,
			expectedChildRules:   0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			input := csslexer.NewInput("{ " + tc.input + " }")
			parser := NewParser(input)

			styleRule := &css.StyleRule{
				Type:      css.StyleRuleTypeQualifiedRule,
				Selectors: []*css.Selector{},
			}

			err := parser.s.ConsumeBlock(func(ts *token_stream.TokenStream) error {
				return parser.consumeStyleRuleContents(styleRule, tc.nestingType)
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(styleRule.Declarations) != tc.expectedDeclarations {
				t.Errorf("expected %d declarations, got %d", tc.expectedDeclarations, len(styleRule.Declarations))
			}

			if len(styleRule.Rules) != tc.expectedChildRules {
				t.Errorf("expected %d child rules, got %d", tc.expectedChildRules, len(styleRule.Rules))
			}
		})
	}
}

func TestParser_StartsCustomPropertyDeclaration(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "custom property",
			input:    "--main-color: red;",
			expected: true,
		},
		{
			name:     "custom property with underscore",
			input:    "--my_var: blue;",
			expected: true,
		},
		{
			name:     "regular property",
			input:    "color: red;",
			expected: false,
		},
		{
			name:     "property starting with dash",
			input:    "-webkit-transform: rotate(45deg);",
			expected: false,
		},
		{
			name:     "empty input",
			input:    "",
			expected: false,
		},
		{
			name:     "custom property with spaces",
			input:    "  --spacing: 10px;",
			expected: false, // whitespace before the property
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			input := csslexer.NewInput(tc.input)
			parser := NewParser(input)

			result := parser.startsCustomPropertyDeclaration()

			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestParser_ConsumeRuleList(t *testing.T) {
	testcases := []struct {
		name         string
		input        string
		allowedRules allowedRuleType
		expectError  bool
		expected     []*css.StyleRule
	}{
		{
			name:         "single rule",
			input:        "div { color: red; }",
			allowedRules: qualifiedRuleTypeStyle,
			expectError:  false,
			expected: []*css.StyleRule{
				{
					Type: css.StyleRuleTypeQualifiedRule,
					Selectors: []*css.Selector{
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
					Declarations: []*css.Declaration{
						{Property: "color", Value: "red", Important: false},
					},
					Rules: []*css.GenericRule{},
				},
			},
		},
		{
			name:         "multiple rules",
			input:        "div { color: red; } .class { margin: 10px; }",
			allowedRules: qualifiedRuleTypeStyle,
			expectError:  false,
			expected: []*css.StyleRule{
				{
					Type: css.StyleRuleTypeQualifiedRule,
					Selectors: []*css.Selector{
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
					Declarations: []*css.Declaration{
						{Property: "color", Value: "red", Important: false},
					},
					Rules: []*css.GenericRule{},
				},
				{
					Type: css.StyleRuleTypeQualifiedRule,
					Selectors: []*css.Selector{
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
					Declarations: []*css.Declaration{
						{Property: "margin", Value: "10px", Important: false},
					},
					Rules: []*css.GenericRule{},
				},
			},
		},
		{
			name:         "empty input",
			input:        "",
			allowedRules: qualifiedRuleTypeStyle,
			expectError:  false,
			expected:     []*css.StyleRule{},
		},
		{
			name:         "whitespace and comments",
			input:        "/* comment */ div { color: red; } /* another comment */",
			allowedRules: qualifiedRuleTypeStyle,
			expectError:  false,
			expected: []*css.StyleRule{
				{
					Type: css.StyleRuleTypeQualifiedRule,
					Selectors: []*css.Selector{
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
					Declarations: []*css.Declaration{
						{Property: "color", Value: "red", Important: false},
					},
					Rules: []*css.GenericRule{},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			input := csslexer.NewInput(tc.input)
			parser := NewParser(input)

			rules, err := parser.consumeRuleList(tc.allowedRules, true, nesting.NestingTypeNone, nil)

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

			if len(rules) != len(tc.expected) {
				t.Errorf("expected %d rules, got %d", len(tc.expected), len(rules))
				return
			}

			for i, rule := range rules {
				if !rule.Equals(tc.expected[i]) {
					t.Errorf("rule %d mismatch:\nexpected: %+v\ngot: %+v", i, tc.expected[i], rule)
				}
			}
		})
	}
}
