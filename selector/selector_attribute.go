package selector

import (
	"go.baoshuo.dev/cssutil"
)

// ===== SelectorDataAttribute =====

type SelectorDataAttr struct {
	AttrNamespace string                // The namespace of the attribute, if any.
	AttrName      string                // The name of the attribute.
	AttrValue     string                // The value of the attribute.
	AttrMatch     SelectorAttrMatchType // The match type for attribute selectors.
}

func NewSelectorDataAttr(name, value string, match SelectorAttrMatchType) *SelectorDataAttr {
	return &SelectorDataAttr{
		AttrName:  name,
		AttrValue: value,
		AttrMatch: match,
	}
}

func (d *SelectorDataAttr) String(match SelectorMatchType) string {
	attrName := cssutil.SerializeIdentifier(d.AttrName)
	attrValue := cssutil.SerializeString(d.AttrValue)

	switch match {
	case SelectorMatchAttributeExact:
		return "[" + attrName + "=" + attrValue + "]"
	case SelectorMatchAttributeSet:
		return "[" + attrName + "]"
	case SelectorMatchAttributeHyphen:
		return "[" + attrName + "|=" + attrValue + "]"
	case SelectorMatchAttributeList:
		return "[" + attrName + "~=" + attrValue + "]"
	case SelectorMatchAttributeContain:
		return "[" + attrName + "*=" + attrValue + "]"
	case SelectorMatchAttributeBegin:
		return "[" + attrName + "^=" + attrValue + "]"
	case SelectorMatchAttributeEnd:
		return "[" + attrName + "$=" + attrValue + "]"
	default:
		return "[UnknownAttributeMatchType]"
	}
}

func (d *SelectorDataAttr) Equals(other SelectorDataType) bool {
	otherData, ok := other.(*SelectorDataAttr)
	if !ok {
		return false
	}

	return d.AttrName == otherData.AttrName &&
		d.AttrValue == otherData.AttrValue &&
		d.AttrMatch == otherData.AttrMatch
}

// ===== SelectorAttrMatchType =====

type SelectorAttrMatchType int

const (
	SelectorAttrMatchCaseSensitive SelectorAttrMatchType = iota
	SelectorAttrMatchCaseInsensitive
	SelectorAttrMatchCaseSensitiveAlways
)
