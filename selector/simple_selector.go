package selector

import (
	"strings"

	"go.baoshuo.dev/cssutil"
)

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
	default:
		return ""
	}
}

// ===== SelectorDataType =====

type SelectorDataType interface {
	String(match SelectorMatchType) string
	Equals(other SelectorDataType) bool
}

// ===== SelectorData =====

type SelectorData struct {
	Value string // The value of the selector, e.g., tag name, class name, id, etc.
}

func NewSelectorData(value string) *SelectorData {
	return &SelectorData{Value: value}
}

func (d *SelectorData) String(match SelectorMatchType) string {
	switch match {
	case SelectorMatchId:
		return "#" + cssutil.SerializeIdentifier(d.Value)
	case SelectorMatchClass:
		return "." + cssutil.SerializeIdentifier(d.Value)
	case SelectorMatchPseudoClass:
		return ":" + cssutil.SerializeIdentifier(d.Value)
	case SelectorMatchPseudoElement:
		return "::" + cssutil.SerializeIdentifier(d.Value)
	case SelectorMatchPagePseudoClass:
		return "@page :" + cssutil.SerializeIdentifier(d.Value)
	case SelectorMatchTag:
		return cssutil.SerializeIdentifier(d.Value)
	case SelectorMatchUniversalTag:
		if d.Value != "" {
			return cssutil.SerializeIdentifier(d.Value) + "|*"
		} else {
			return "*"
		}
	default:
		return cssutil.SerializeIdentifier(d.Value)
	}
}

func (d *SelectorData) Equals(other SelectorDataType) bool {
	otherData, ok := other.(*SelectorData)
	if !ok {
		return false
	}
	return d.Value == otherData.Value
}

// ===== SelectorDataTag =====

type SelectorDataTag struct {
	Namespace string // The namespace of the tag, if any.
	TagName   string // The tag name.
}

func NewSelectorDataTag(namespace, tagName string) *SelectorDataTag {
	return &SelectorDataTag{
		Namespace: namespace,
		TagName:   tagName,
	}
}

func (d *SelectorDataTag) String(match SelectorMatchType) string {
	var tagName string
	if match == SelectorMatchUniversalTag {
		tagName = "*"
	} else {
		tagName = cssutil.SerializeIdentifier(d.TagName)
	}

	if d.Namespace != "" {
		return cssutil.SerializeIdentifier(d.Namespace) + "|" + tagName
	} else {
		return tagName
	}
}

func (d *SelectorDataTag) Equals(other SelectorDataType) bool {
	otherData, ok := other.(*SelectorDataTag)
	if !ok {
		return false
	}
	return d.Namespace == otherData.Namespace && d.TagName == otherData.TagName
}
