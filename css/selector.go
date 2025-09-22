package css

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

func (s *Selector) String() string {
	var result strings.Builder

	for _, sel := range s.Selectors {
		result.WriteString(sel.String())
	}

	return result.String()
}

func (s *Selector) Equals(other *Selector) bool {
	if s.Flag != other.Flag || len(s.Selectors) != len(other.Selectors) {
		return false
	}

	for i, sel := range s.Selectors {
		if !sel.Equals(other.Selectors[i]) {
			return false
		}
	}

	return true
}

// ===== SimpleSelector =====

// SimpleSelector represents a single simple selector within a compound selector.
// It can represent various types of selectors such as tag, class, id, attribute, etc.
type SimpleSelector struct {
	Match    SelectorMatchType    // The type of selector match.
	Relation SelectorRelationType // The relation to the previous selector in the list.
	Data     SelectorDataType
}

func (s *SimpleSelector) String() string {
	var result strings.Builder

	result.WriteString(s.Relation.String())

	if s.Data != nil {
		result.WriteString(s.Data.String(s.Match))
	} else {
		// This should not happen in the new implementation
		result.WriteString("[UnknownSelector]")
	}

	return result.String()
}

func (s *SimpleSelector) Equals(other *SimpleSelector) bool {
	if other == nil {
		return false
	}

	if s.Match != other.Match || s.Relation != other.Relation {
		return false
	}

	// Handle nil data cases
	if s.Data == nil && other.Data == nil {
		return true
	}
	if s.Data == nil || other.Data == nil {
		return false
	}

	return s.Data.Equals(other.Data)
}

// ===== SelectorMatchType =====

type SelectorMatchType int

const (
	SelectorMatchUnknown                     = iota
	SelectorMatchInvalidList                 // Used as a marker in CSSSelectorList.
	SelectorMatchTag                         // Example: div
	SelectorMatchUniversalTag                // Example: * (possibly with namespace)
	SelectorMatchId                          // Example: #id
	SelectorMatchClass                       // Example: .class
	SelectorMatchPseudoClass                 // Example: :nth-child(2)
	SelectorMatchPseudoElement               // Example: ::first-line
	SelectorMatchPagePseudoClass             // Example: @page :right
	SelectorMatchAttributeExact              // Example: E[foo="bar"]
	SelectorMatchAttributeSet                // Example: E[foo]
	SelectorMatchAttributeHyphen             // Example: E[foo|="bar"]
	SelectorMatchAttributeList               // Example: E[foo~="bar"]
	SelectorMatchAttributeContain            // css3: E[foo*="bar"]
	SelectorMatchAttributeBegin              // css3: E[foo^="bar"]
	SelectorMatchAttributeEnd                // css3: E[foo$="bar"]
	SelectorMatchFirstAttributeSelectorMatch = SelectorMatchAttributeExact
)

// ===== SelectorRelationType =====

type SelectorRelationType int

const (
	SelectorRelationSubSelector      SelectorRelationType = iota // No combinator. Used between simple selectors within the same compound.
	SelectorRelationDescendant                                   // "Space" combinator
	SelectorRelationChild                                        // > combinator
	SelectorRelationDirectAdjacent                               // + combinator
	SelectorRelationIndirectAdjacent                             // ~ combinator

	// The relation types below are implicit combinators inserted at parse time
	// before pseudo-elements which match another flat tree element than the
	// rest of the compound.

	// Implicit combinator inserted before pseudo-elements matching an element
	// inside a UA shadow tree. This combinator allows the selector matching to
	// cross a shadow root.
	//
	// Examples:
	// input::placeholder, video::cue(i), video::--webkit-media-controls-panel
	SelectorRelationUAShadow
	// Implicit combinator inserted before ::slotted() selectors.
	SelectorRelationShadowSlot
	// Implicit combinator inserted before ::part() selectors which allows
	// matching a ::part in shadow-including descendant tree for #host in
	// "#host::part(button)".
	SelectorRelationShadowPart

	// Relative selectors
	SelectorRelationRelativeDescendant       // leftmost "Space" combinator of relative selector
	SelectorRelationRelativeChild            // leftmost > combinator of relative selector
	SelectorRelationRelativeDirectAdjacent   // leftmost + combinator of relative selector
	SelectorRelationRelativeIndirectAdjacent // leftmost ~ combinator of relative selector
)

func (sr SelectorRelationType) String() string {
	switch sr {
	case SelectorRelationSubSelector:
		return ""
	case SelectorRelationDescendant:
		return " "
	case SelectorRelationChild:
		return " > "
	case SelectorRelationDirectAdjacent:
		return " + "
	case SelectorRelationIndirectAdjacent:
		return " ~ "
	case SelectorRelationRelativeDescendant:
		return " " // Same as regular descendant in string form
	case SelectorRelationRelativeChild:
		return " > " // Same as regular child in string form
	case SelectorRelationRelativeDirectAdjacent:
		return " + " // Same as regular direct adjacent in string form
	case SelectorRelationRelativeIndirectAdjacent:
		return " ~ " // Same as regular indirect adjacent in string form
	default:
		return ""
	}
}
