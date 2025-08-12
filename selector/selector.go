package selector

import (
	"strings"
)

// ===== SelectorListFlagType =====

type SelectorListFlagType int

const (
	SelectorFlagContainsPseudo SelectorListFlagType = 1 << iota
	SelectorFlagContainsComplexSelector
	SelectorFlagContainsScopeOrParent
)

func (s SelectorListFlagType) Has(flag SelectorListFlagType) bool {
	return s&flag != 0
}

func (s *SelectorListFlagType) Set(flag SelectorListFlagType) {
	*s |= flag
}

// ===== Selector =====

// Selector represents a list of simple selectors that can be combined
// to form a complex selector.
//
// It can contain multiple simple selectors and may include combinators
// to define relationships between them.
//
// It also has flags to indicate certain properties of the selector.
//
// The selectors are stored in the order they appear in the CSS.
type Selector struct {
	Flag      SelectorListFlagType // Flags for the selector
	Selectors []*SimpleSelector    // The list of selectors in this selector list
}

func (s *Selector) Append(sel ...*SimpleSelector) {
	if len(sel) == 0 {
		return
	}
	s.Selectors = append(s.Selectors, sel...)
}

func (s *Selector) InsertBefore(index int, sel *SimpleSelector) {
	if sel == nil || index < 0 || index > len(s.Selectors) {
		return
	}
	s.Selectors = append(s.Selectors[:index], append([]*SimpleSelector{sel}, s.Selectors[index:]...)...)
}

func (s *Selector) Prepend(sel *SimpleSelector) {
	if sel == nil {
		return
	}
	s.Selectors = append([]*SimpleSelector{sel}, s.Selectors...)
}

func (s *Selector) Equal(other *Selector) bool {
	if s.Flag != other.Flag || len(s.Selectors) != len(other.Selectors) {
		return false
	}

	for i, sel := range s.Selectors {
		if !sel.Equal(other.Selectors[i]) {
			return false
		}
	}

	return true
}

func (s *Selector) String() string {
	var result strings.Builder

	for _, sel := range s.Selectors {
		result.WriteString(sel.String())
	}

	return result.String()
}
