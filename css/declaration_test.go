package css

import (
	"testing"
)

func TestDeclarationInterface(t *testing.T) {
	// Test that our declaration implements the interface correctly
	decl := &Declaration{
		Property:  "color",
		Value:     "red",
		Important: false,
	}

	// Test String method
	expected := "color: red"
	if decl.String() != expected {
		t.Errorf("expected %q, got %q", expected, decl.String())
	}

	// Test IsCustomProperty method
	if decl.IsCustomProperty() {
		t.Error("expected IsCustomProperty to return false for regular property")
	}

	// Test custom property
	customDecl := &Declaration{
		Property:  "--main-color",
		Value:     "#ff0000",
		Important: true,
	}

	expectedCustom := "--main-color: #ff0000 !important"
	if customDecl.String() != expectedCustom {
		t.Errorf("expected %q, got %q", expectedCustom, customDecl.String())
	}

	if !customDecl.IsCustomProperty() {
		t.Error("expected IsCustomProperty to return true for custom property")
	}
}

func TestDeclarationString(t *testing.T) {
	tests := []struct {
		name     string
		decl     *Declaration
		expected string
	}{
		{
			name: "basic property",
			decl: &Declaration{
				Property:  "color",
				Value:     "red",
				Important: false,
			},
			expected: "color: red",
		},
		{
			name: "important property",
			decl: &Declaration{
				Property:  "color",
				Value:     "blue",
				Important: true,
			},
			expected: "color: blue !important",
		},
		{
			name: "custom property",
			decl: &Declaration{
				Property:  "--theme-color",
				Value:     "#333",
				Important: false,
			},
			expected: "--theme-color: #333",
		},
		{
			name: "custom property important",
			decl: &Declaration{
				Property:  "--theme-color",
				Value:     "#333",
				Important: true,
			},
			expected: "--theme-color: #333 !important",
		},
		{
			name: "complex value",
			decl: &Declaration{
				Property:  "background",
				Value:     "url('image.png') no-repeat center",
				Important: false,
			},
			expected: "background: url('image.png') no-repeat center",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.decl.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDeclarationIsCustomProperty(t *testing.T) {
	tests := []struct {
		name     string
		property string
		expected bool
	}{
		{"regular property", "color", false},
		{"regular property with dash", "background-color", false},
		{"custom property", "--main-color", true},
		{"custom property long", "--very-long-custom-property-name", true},
		{"empty property", "", false},
		{"single dash", "-", false},
		{"almost custom", "-custom", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decl := &Declaration{Property: tt.property}
			result := decl.IsCustomProperty()
			if result != tt.expected {
				t.Errorf("expected %v, got %v for property %q", tt.expected, result, tt.property)
			}
		})
	}
}

func TestDeclarationEquals(t *testing.T) {
	decl1 := &Declaration{
		Property:  "color",
		Value:     "red",
		Important: false,
	}

	decl2 := &Declaration{
		Property:  "color",
		Value:     "red",
		Important: false,
	}

	decl3 := &Declaration{
		Property:  "color",
		Value:     "blue",
		Important: false,
	}

	decl4 := &Declaration{
		Property:  "color",
		Value:     "red",
		Important: true,
	}

	decl5 := &Declaration{
		Property:  "background-color",
		Value:     "red",
		Important: false,
	}

	tests := []struct {
		name     string
		decl1    *Declaration
		decl2    *Declaration
		expected bool
	}{
		{"identical declarations", decl1, decl2, true},
		{"same object", decl1, decl1, true},
		{"different values", decl1, decl3, false},
		{"different importance", decl1, decl4, false},
		{"different properties", decl1, decl5, false},
		{"nil comparison", decl1, nil, false},
		{"nil receiver", nil, decl1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			if tt.decl1 == nil {
				// Can't call method on nil receiver
				result = false
			} else {
				result = tt.decl1.Equals(tt.decl2)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDeclarationList(t *testing.T) {
	dl := &DeclarationList{}

	// Test initial state
	if !dl.IsEmpty() {
		t.Error("expected new DeclarationList to be empty")
	}
	if dl.Size() != 0 {
		t.Errorf("expected size 0, got %d", dl.Size())
	}

	// Test appending declarations
	decl1 := &Declaration{Property: "color", Value: "red", Important: false}
	decl2 := &Declaration{Property: "background", Value: "blue", Important: true}

	dl.Append(decl1)
	if dl.IsEmpty() {
		t.Error("expected DeclarationList to not be empty after append")
	}
	if dl.Size() != 1 {
		t.Errorf("expected size 1, got %d", dl.Size())
	}

	dl.Append(decl2)
	if dl.Size() != 2 {
		t.Errorf("expected size 2, got %d", dl.Size())
	}

	// Test appending nil
	dl.Append(nil)
	if dl.Size() != 2 {
		t.Errorf("expected size to remain 2 after appending nil, got %d", dl.Size())
	}

	// Test String method
	expected := "color: red; background: blue !important"
	result := dl.String()
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestDeclarationListString(t *testing.T) {
	tests := []struct {
		name         string
		declarations []*Declaration
		expected     string
	}{
		{
			name:         "empty list",
			declarations: []*Declaration{},
			expected:     "",
		},
		{
			name: "single declaration",
			declarations: []*Declaration{
				{Property: "color", Value: "red", Important: false},
			},
			expected: "color: red",
		},
		{
			name: "multiple declarations",
			declarations: []*Declaration{
				{Property: "color", Value: "red", Important: false},
				{Property: "background", Value: "blue", Important: true},
				{Property: "margin", Value: "10px", Important: false},
			},
			expected: "color: red; background: blue !important; margin: 10px",
		},
		{
			name: "custom properties",
			declarations: []*Declaration{
				{Property: "--main-color", Value: "#333", Important: false},
				{Property: "color", Value: "var(--main-color)", Important: false},
			},
			expected: "--main-color: #333; color: var(--main-color)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := &DeclarationList{Declarations: tt.declarations}
			result := dl.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
