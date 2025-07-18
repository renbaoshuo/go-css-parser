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

func (d *CssDeclaration) ValueString() string {
	return d.Value
}

func (d *CssDeclaration) ValueStringWithImportant(option bool) string {
	if option && d.Important {
		return fmt.Sprintf("%s !important", d.Value)
	}
	return d.Value
}

func (d *CssDeclaration) String() string {
	return d.StringWithImportant(true)
}

func (d *CssDeclaration) StringWithImportant(option bool) string {
	return fmt.Sprintf("%s: %s;", d.Property, d.ValueStringWithImportant(option))
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

func (ds DeclarationsByProperty) ToObject() map[string]string {
	obj := make(map[string]string, len(ds))
	for _, d := range ds {
		if _, exists := obj[d.Property]; !exists {
			obj[d.Property] = d.ValueString()
		}
	}
	return obj
}
