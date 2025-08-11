package selector

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
)

// ===== SelectorAttributeMatchType =====

type SelectorAttributeMatchType int

const (
	SelectorAttributeMatchUnknown       SelectorAttributeMatchType = iota
	SelectorAttributeMatchCaseSensitive SelectorAttributeMatchType = iota
	SelectorAttributeMatchCaseInsensitive
	SelectorAttributeMatchCaseSensitiveAlways
)

// ===== SimpleSelector =====

// SimpleSelector represents a single simple selector within a compound selector.
// It can represent various types of selectors such as tag, class, id, attribute, etc.
type SimpleSelector struct {
	Match    SelectorMatchType    // The type of selector match.
	Data     []rune               // The raw selector data.
	Relation SelectorRelationType // The relation to the previous selector in the list.

	// Below are optional fields that may be used for specific selector types.

	AttrValue []rune                     // The value of the attribute, if applicable.
	AttrMatch SelectorAttributeMatchType // The match type for attribute selectors, if applicable.
}

func (s *SimpleSelector) Equal(other *SimpleSelector) bool {
	return s.Match == other.Match &&
		string(s.Data) == string(other.Data) &&
		s.Relation == other.Relation &&
		string(s.AttrValue) == string(other.AttrValue) &&
		s.AttrMatch == other.AttrMatch
}
