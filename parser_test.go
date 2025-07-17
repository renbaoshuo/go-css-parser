package css_parser

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

	println("Parsed stylesheet:")
	println(stylesheet.String())
}

func TestParseDeclarations(t *testing.T) {
	content := "color: red; background-color: blue !important;"
	parser := NewParser(content, true)

	declarations, err := parser.ParseDeclarations()
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
		println("Parsed declaration:", decl.String())
	}
}
