package css_parser

import (
	"fmt"
)

type CssDeclaration struct {
	Property  string
	Value     string
	Important bool
}

func NewCssDeclaration() *CssDeclaration {
	return &CssDeclaration{}
}

func (d *CssDeclaration) String() string {
	return d.StringWithImportant(true)
}

func (d *CssDeclaration) StringWithImportant(option bool) string {
	result := fmt.Sprintf("%s: %s", d.Property, d.Value)

	if option && d.Important {
		result += " !important"
	}

	result += ";"

	return result
}

func (d *CssDeclaration) Equal(other *CssDeclaration) bool {
	return (d.Property == other.Property) && (d.Value == other.Value) && (d.Important == other.Important)
}

type DeclarationsByProperty []*CssDeclaration

// Implements sort.Interface
func (ds DeclarationsByProperty) Len() int {
	return len(ds)
}

func (ds DeclarationsByProperty) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

func (ds DeclarationsByProperty) Less(i, j int) bool {
	return ds[i].Property < ds[j].Property
}
