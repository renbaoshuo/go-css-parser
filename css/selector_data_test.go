package css

import (
	"testing"
)

func TestSelectorDataString(t *testing.T) {
	data := NewSelectorData("test-value")

	tests := []struct {
		name     string
		data     SelectorDataType
		match    SelectorMatchType
		expected string
	}{
		{"id selector", data, SelectorMatchId, "#test-value"},
		{"class selector", data, SelectorMatchClass, ".test-value"},
		{"pseudo class", data, SelectorMatchPseudoClass, ":test-value"},
		{"pseudo element", data, SelectorMatchPseudoElement, "::test-value"},
		{"page pseudo class", data, SelectorMatchPagePseudoClass, "@page :test-value"},
		{"tag selector", data, SelectorMatchTag, "test-value"},
		{"universal tag empty", NewSelectorData(""), SelectorMatchUniversalTag, "*"},
		{"unknown match", data, SelectorMatchUnknown, "test-value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.data.String(tt.match)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSelectorDataStringWithNamespace(t *testing.T) {
	data := NewSelectorData("svg")

	result := data.String(SelectorMatchUniversalTag)
	expected := "svg|*"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestSelectorDataEquals(t *testing.T) {
	data1 := NewSelectorData("test")
	data2 := NewSelectorData("test")
	data3 := NewSelectorData("different")

	tests := []struct {
		name     string
		data1    SelectorDataType
		data2    SelectorDataType
		expected bool
	}{
		{"identical values", data1, data2, true},
		{"same object", data1, data1, true},
		{"different values", data1, data3, false},
		{"different types", data1, NewSelectorDataTag("", "div"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.data1.Equals(tt.data2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSelectorDataTagString(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		tagName   string
		match     SelectorMatchType
		expected  string
	}{
		{
			name:     "simple tag",
			tagName:  "div",
			match:    SelectorMatchTag,
			expected: "div",
		},
		{
			name:     "universal tag",
			tagName:  "any",
			match:    SelectorMatchUniversalTag,
			expected: "*",
		},
		{
			name:      "namespaced tag",
			namespace: "svg",
			tagName:   "rect",
			match:     SelectorMatchTag,
			expected:  "svg|rect",
		},
		{
			name:      "namespaced universal",
			namespace: "xml",
			tagName:   "any",
			match:     SelectorMatchUniversalTag,
			expected:  "xml|*",
		},
		{
			name:     "complex tag name",
			tagName:  "custom-element",
			match:    SelectorMatchTag,
			expected: "custom-element",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := NewSelectorDataTag(tt.namespace, tt.tagName)
			result := data.String(tt.match)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSelectorDataTagEquals(t *testing.T) {
	tag1 := NewSelectorDataTag("", "div")
	tag2 := NewSelectorDataTag("", "div")
	tag3 := NewSelectorDataTag("", "span")
	tag4 := NewSelectorDataTag("svg", "rect")
	tag5 := NewSelectorDataTag("svg", "rect")
	tag6 := NewSelectorDataTag("xml", "rect")

	tests := []struct {
		name     string
		tag1     SelectorDataType
		tag2     SelectorDataType
		expected bool
	}{
		{"identical tags", tag1, tag2, true},
		{"same object", tag1, tag1, true},
		{"different tag names", tag1, tag3, false},
		{"namespaced tags identical", tag4, tag5, true},
		{"different namespaces", tag4, tag6, false},
		{"namespaced vs non-namespaced", tag1, tag4, false},
		{"different types", tag1, NewSelectorData("div"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tag1.Equals(tt.tag2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNewSelectorData(t *testing.T) {
	value := "test-value"
	data := NewSelectorData(value)

	if data == nil {
		t.Error("expected NewSelectorData to return non-nil")
	}
	if data.Value != value {
		t.Errorf("expected value %q, got %q", value, data.Value)
	}
}

func TestNewSelectorDataTag(t *testing.T) {
	namespace := "svg"
	tagName := "rect"
	data := NewSelectorDataTag(namespace, tagName)

	if data == nil {
		t.Fatal("expected NewSelectorDataTag to return non-nil")
	}
	if data.Namespace != namespace {
		t.Errorf("expected namespace %q, got %q", namespace, data.Namespace)
	}
	if data.TagName != tagName {
		t.Errorf("expected tagName %q, got %q", tagName, data.TagName)
	}
}
