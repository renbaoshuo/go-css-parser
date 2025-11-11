package css

import (
	"go.baoshuo.dev/cssutil"
)

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
