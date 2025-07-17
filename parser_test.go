package css_parser

import (
	"testing"
)

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
}
