package css

import (
	"testing"
)

func TestSelectorDataAttrString(t *testing.T) {
	tests := []struct {
		name      string
		attrName  string
		attrValue string
		attrMatch SelectorAttrMatchType
		match     SelectorMatchType
		expected  string
	}{
		{
			name:      "exact attribute match",
			attrName:  "class",
			attrValue: "button",
			attrMatch: SelectorAttrMatchCaseSensitive,
			match:     SelectorMatchAttributeExact,
			expected:  "[class=\"button\"]",
		},
		{
			name:      "attribute set",
			attrName:  "disabled",
			attrValue: "",
			attrMatch: SelectorAttrMatchCaseSensitive,
			match:     SelectorMatchAttributeSet,
			expected:  "[disabled]",
		},
		{
			name:      "attribute hyphen match",
			attrName:  "lang",
			attrValue: "en",
			attrMatch: SelectorAttrMatchCaseSensitive,
			match:     SelectorMatchAttributeHyphen,
			expected:  "[lang|=\"en\"]",
		},
		{
			name:      "attribute list match",
			attrName:  "class",
			attrValue: "primary",
			attrMatch: SelectorAttrMatchCaseSensitive,
			match:     SelectorMatchAttributeList,
			expected:  "[class~=\"primary\"]",
		},
		{
			name:      "attribute contains match",
			attrName:  "href",
			attrValue: "example",
			attrMatch: SelectorAttrMatchCaseSensitive,
			match:     SelectorMatchAttributeContain,
			expected:  "[href*=\"example\"]",
		},
		{
			name:      "attribute begins match",
			attrName:  "href",
			attrValue: "https",
			attrMatch: SelectorAttrMatchCaseSensitive,
			match:     SelectorMatchAttributeBegin,
			expected:  "[href^=\"https\"]",
		},
		{
			name:      "attribute ends match",
			attrName:  "href",
			attrValue: ".pdf",
			attrMatch: SelectorAttrMatchCaseSensitive,
			match:     SelectorMatchAttributeEnd,
			expected:  "[href$=\".pdf\"]",
		},
		{
			name:      "unknown attribute match",
			attrName:  "data-test",
			attrValue: "value",
			attrMatch: SelectorAttrMatchCaseSensitive,
			match:     SelectorMatchUnknown,
			expected:  "[UnknownAttributeMatchType]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := NewSelectorDataAttr(tt.attrName, tt.attrValue, tt.attrMatch)
			result := attr.String(tt.match)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSelectorDataAttrEquals(t *testing.T) {
	attr1 := NewSelectorDataAttr("class", "button", SelectorAttrMatchCaseSensitive)
	attr2 := NewSelectorDataAttr("class", "button", SelectorAttrMatchCaseSensitive)
	attr3 := NewSelectorDataAttr("id", "button", SelectorAttrMatchCaseSensitive)
	attr4 := NewSelectorDataAttr("class", "link", SelectorAttrMatchCaseSensitive)
	attr5 := NewSelectorDataAttr("class", "button", SelectorAttrMatchCaseInsensitive)

	tests := []struct {
		name     string
		attr1    SelectorDataType
		attr2    SelectorDataType
		expected bool
	}{
		{"identical attributes", attr1, attr2, true},
		{"same object", attr1, attr1, true},
		{"different attribute names", attr1, attr3, false},
		{"different attribute values", attr1, attr4, false},
		{"different match types", attr1, attr5, false},
		{"different data types", attr1, NewSelectorData("class"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.attr1.Equals(tt.attr2)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNewSelectorDataAttr(t *testing.T) {
	name := "data-test"
	value := "example"
	match := SelectorAttrMatchCaseInsensitive

	attr := NewSelectorDataAttr(name, value, match)

	if attr == nil {
		t.Error("expected NewSelectorDataAttr to return non-nil")
	}
	if attr.AttrName != name {
		t.Errorf("expected AttrName %q, got %q", name, attr.AttrName)
	}
	if attr.AttrValue != value {
		t.Errorf("expected AttrValue %q, got %q", value, attr.AttrValue)
	}
	if attr.AttrMatch != match {
		t.Errorf("expected AttrMatch %v, got %v", match, attr.AttrMatch)
	}
	if attr.AttrNamespace != "" {
		t.Errorf("expected empty AttrNamespace, got %q", attr.AttrNamespace)
	}
}

func TestSelectorDataAttrWithNamespace(t *testing.T) {
	attr := &SelectorDataAttr{
		AttrNamespace: "xml",
		AttrName:      "lang",
		AttrValue:     "en",
		AttrMatch:     SelectorAttrMatchCaseSensitive,
	}

	// Test that the namespace field is properly stored
	if attr.AttrNamespace != "xml" {
		t.Errorf("expected namespace %q, got %q", "xml", attr.AttrNamespace)
	}

	// Test equality with namespace consideration
	attr2 := &SelectorDataAttr{
		AttrNamespace: "xml",
		AttrName:      "lang",
		AttrValue:     "en",
		AttrMatch:     SelectorAttrMatchCaseSensitive,
	}

	attr3 := &SelectorDataAttr{
		AttrNamespace: "html",
		AttrName:      "lang",
		AttrValue:     "en",
		AttrMatch:     SelectorAttrMatchCaseSensitive,
	}

	if !attr.Equals(attr2) {
		t.Error("expected attributes with same namespace to be equal")
	}

	if attr.Equals(attr3) {
		t.Error("expected attributes with different namespaces to not be equal")
	}
}

func TestSelectorAttrMatchTypes(t *testing.T) {
	// Test that all SelectorAttrMatchType constants are defined
	cases := []SelectorAttrMatchType{
		SelectorAttrMatchCaseSensitive,
		SelectorAttrMatchCaseInsensitive,
		SelectorAttrMatchCaseSensitiveAlways,
	}

	for i, matchType := range cases {
		if int(matchType) != i {
			t.Errorf("expected SelectorAttrMatchType constant %d to have value %d, got %d", i, i, int(matchType))
		}
	}
}
