package cssparser

import (
	"testing"
)

func TestParseStylesheet(t *testing.T) {
	content := `
html, body {
	color: red;
}
@media screen and (max-width: 600px) {
	body {
		background-color: blue !important;
	}
}
`

	stylesheet, err := ParseStylesheet(content)
	if err != nil {
		t.Fatalf("Failed to parse stylesheet: %v", err)
	}

	if len(stylesheet.Rules) != 2 {
		t.Fatalf("Expected 2 rule, got %d", len(stylesheet.Rules))
	}

	rule := stylesheet.Rules[0]
	if len(rule.Selectors) != 2 || rule.Selectors[0] != "html" || rule.Selectors[1] != "body" {
		t.Errorf("Rule mismatch: %s with %d declarations", rule.Selectors[0], len(rule.Declarations))
	}

	println(stylesheet.String())
}

func TestParseDeclarations(t *testing.T) {
	content := "color: red; background-color: blue !important;"

	declarations, err := ParseDeclarations(content, WithInline(true))
	if err != nil {
		t.Fatalf("Failed to parse declarations: %v", err)
	}

	if len(declarations) != 2 {
		t.Fatalf("Expected 2 declarations, got %d", len(declarations))
	}

	if declarations[0].Property != "color" || declarations[0].Value != "red" {
		t.Errorf("First declaration mismatch: %s: %s", declarations[0].Property, declarations[0].Value)
	}

	if declarations[1].Property != "background-color" || declarations[1].Value != "blue" || !declarations[1].Important {
		t.Errorf("Second declaration mismatch: %s: %s !important=%v", declarations[1].Property, declarations[1].Value, declarations[1].Important)
	}

	for _, decl := range declarations {
		println(decl.String())
	}
}

func TestLooseMode(t *testing.T) {
	// Test CSS with syntax errors
	invalidCSS := `
.invalid {
	color: red;
	background-
	font-size: 16px;
	font-weight: bold;
	border:
	height: 100px;
	width: 200px;
}
	`

	// Test with strict mode (should fail)
	_, err := ParseStylesheet(invalidCSS)
	// In strict mode, we expect it to either fail or succeed - let's just check loose mode behavior
	if err == nil {
		t.Error("Expected error in strict mode, but got none")
	}

	// Test with loose mode (should succeed and skip errors)
	stylesheet, err := ParseStylesheet(invalidCSS, WithLooseParsing(true))
	if err != nil {
		t.Errorf("Expected no error in loose mode, but got: %v", err)
	}

	if stylesheet == nil {
		t.Error("Expected stylesheet to be parsed in loose mode")
	}

	// Should have at least some valid rules
	if len(stylesheet.Rules) < 1 {
		t.Error("Expected at least one valid rule to be parsed in loose mode")
	}

	println(stylesheet.String())
}

func TestLooseModeDeclarations(t *testing.T) {
	// Test declarations with syntax errors
	invalidDeclarations := `
color: red;
invalid-syntax;
background: blue;
: invalid-property;
font-size: 16px`

	// Test with strict mode (should fail)
	_, err := ParseDeclarations(invalidDeclarations)
	if err == nil {
		t.Error("Expected error in strict mode, but got none")
	}

	// Test with loose mode (should succeed)
	declarations, err := ParseDeclarations(invalidDeclarations, WithLooseParsing(true))
	if err != nil {
		t.Errorf("Expected no error in loose mode, but got: %v", err)
	}

	if declarations == nil {
		t.Error("Expected declarations to be parsed in loose mode")
	}

	// Should have at least some valid declarations
	if len(declarations) < 1 {
		t.Error("Expected at least one valid declaration to be parsed in loose mode")
	}

	for _, decl := range declarations {
		println(decl.String())
	}
}
