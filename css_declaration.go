package cssparser

import (
	"fmt"
)

type Declaration struct {
	Property  string
	Value     string
	Important bool
}

func NewCssDeclaration() *Declaration {
	return &Declaration{}
}

func (d *Declaration) ValueString() string {
	return d.Value
}

func (d *Declaration) ValueStringWithImportant(option bool) string {
	if option && d.Important {
		return fmt.Sprintf("%s !important", d.Value)
	}
	return d.Value
}

func (d *Declaration) String() string {
	return d.StringWithImportant(true)
}

func (d *Declaration) StringWithImportant(option bool) string {
	return fmt.Sprintf("%s: %s;", d.Property, d.ValueStringWithImportant(option))
}

func (d *Declaration) Equal(other *Declaration) bool {
	return (d.Property == other.Property) && (d.Value == other.Value) && (d.Important == other.Important)
}

type DeclarationsByProperty []*Declaration

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

// ToObject converts DeclarationsByProperty to a map[string]string
func (ds DeclarationsByProperty) ToObject() map[string]string {
	obj := make(map[string]string, len(ds))
	for _, d := range ds {
		if _, exists := obj[d.Property]; !exists {
			obj[d.Property] = d.ValueString()
		}
	}
	return obj
}
