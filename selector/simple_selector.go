package selector

import (
	"strings"

	"go.baoshuo.dev/cssutil"
)

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
	SelectorAttributeMatchUnknown SelectorAttributeMatchType = iota
	SelectorAttributeMatchCaseSensitive
	SelectorAttributeMatchCaseInsensitive
	SelectorAttributeMatchCaseSensitiveAlways
)

// ===== SimpleSelector =====

// SimpleSelector represents a single simple selector within a compound selector.
// It can represent various types of selectors such as tag, class, id, attribute, etc.
type SimpleSelector struct {
	Match    SelectorMatchType    // The type of selector match.
	Value    string               // The value of the selector, e.g., tag name, class name, id, etc.
	Relation SelectorRelationType // The relation to the previous selector in the list.

	// Below are optional fields that may be used for specific selector types.

	AttrValue string                     // The value of the attribute, if applicable.
	AttrMatch SelectorAttributeMatchType // The match type for attribute selectors, if applicable.
}

func (s *SimpleSelector) AttrValueString() string {
	return cssutil.SerializeString(s.AttrValue)
}

func (s *SimpleSelector) String() string {
	var result strings.Builder

	switch s.Relation {
	case SelectorRelationSubSelector:
		// No combinator, just append the selector.
	case SelectorRelationDescendant:
		result.WriteString(" ")
	case SelectorRelationChild:
		result.WriteString(" > ")
	case SelectorRelationDirectAdjacent:
		result.WriteString(" + ")
	case SelectorRelationIndirectAdjacent:
		result.WriteString(" ~ ")
	}

	switch s.Match {
	case SelectorMatchTag:
		result.WriteString(cssutil.SerializeIdentifier(s.Value))
	case SelectorMatchUniversalTag:
		if s.Value != "" {
			result.WriteString(cssutil.SerializeIdentifier(s.Value) + "|*")
		} else {
			result.WriteString("*")
		}
	case SelectorMatchId:
		result.WriteString("#" + cssutil.SerializeIdentifier(s.Value))
	case SelectorMatchClass:
		result.WriteString("." + cssutil.SerializeIdentifier(s.Value))
	case SelectorMatchPseudoClass:
		result.WriteString(":" + cssutil.SerializeIdentifier(s.Value))
	case SelectorMatchPseudoElement:
		result.WriteString("::" + cssutil.SerializeIdentifier(s.Value))
	case SelectorMatchPagePseudoClass:
		result.WriteString("@page :" + cssutil.SerializeIdentifier(s.Value))
	case SelectorMatchAttributeExact:
		result.WriteString("[" + cssutil.SerializeIdentifier(s.Value) + "=" + s.AttrValueString() + "]")
	case SelectorMatchAttributeSet:
		result.WriteString("[" + cssutil.SerializeIdentifier(s.Value) + "]")
	case SelectorMatchAttributeHyphen:
		result.WriteString("[" + cssutil.SerializeIdentifier(s.Value) + "|=" + s.AttrValueString() + "]")
	case SelectorMatchAttributeList:
		result.WriteString("[" + cssutil.SerializeIdentifier(s.Value) + "~=" + s.AttrValueString() + "]")
	case SelectorMatchAttributeContain:
		result.WriteString("[" + cssutil.SerializeIdentifier(s.Value) + "*=" + s.AttrValueString() + "]")
	case SelectorMatchAttributeBegin:
		result.WriteString("[" + cssutil.SerializeIdentifier(s.Value) + "^=" + s.AttrValueString() + "]")
	case SelectorMatchAttributeEnd:
		result.WriteString("[" + cssutil.SerializeIdentifier(s.Value) + "$=" + s.AttrValueString() + "]")
	default:
		result.WriteString("UnknownSelectorMatchType(" + cssutil.SerializeIdentifier(s.Value) + ")")
	}

	return result.String()
}

func (s *SimpleSelector) Equal(other *SimpleSelector) bool {
	return s.Match == other.Match &&
		s.Value == other.Value &&
		s.Relation == other.Relation &&
		s.AttrValue == other.AttrValue &&
		s.AttrMatch == other.AttrMatch
}
