package css

import (
	"strings"
)

// Declaration represents a CSS property declaration (property: value)
type Declaration struct {
	Property  string // CSS property name
	Value     string // CSS property value (unparsed)
	Important bool   // Whether the declaration has !important
}

// String returns the string representation of the declaration
func (d *Declaration) String() string {
	result := d.Property + ": " + d.Value
	if d.Important {
		result += " !important"
	}
	return result
}

// IsCustomProperty returns true if this is a custom CSS property (starts with --)
func (d *Declaration) IsCustomProperty() bool {
	return strings.HasPrefix(d.Property, "--")
}

// Equals compares two Declaration instances
func (d *Declaration) Equals(other *Declaration) bool {
	if other == nil {
		return false
	}
	return d.Property == other.Property &&
		d.Value == other.Value &&
		d.Important == other.Important
}

// DeclarationList represents a list of CSS declarations
type DeclarationList struct {
	Declarations []*Declaration
}

// Append adds a declaration to the list
func (dl *DeclarationList) Append(decl *Declaration) {
	if decl != nil {
		dl.Declarations = append(dl.Declarations, decl)
	}
}

// Size returns the number of declarations
func (dl *DeclarationList) Size() int {
	return len(dl.Declarations)
}

// IsEmpty returns true if there are no declarations
func (dl *DeclarationList) IsEmpty() bool {
	return len(dl.Declarations) == 0
}

// String returns the string representation of the declaration list
func (dl *DeclarationList) String() string {
	var result []string
	for _, decl := range dl.Declarations {
		result = append(result, decl.String())
	}
	return strings.Join(result, "; ")
}
