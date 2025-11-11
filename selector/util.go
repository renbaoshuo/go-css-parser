package selector

import (
	"strings"

	"go.baoshuo.dev/cssparser/css"
)

func prependTypeSelectorIfNeeded(selectors []*css.SimpleSelector, name, namespace string, hasQName bool) []*css.SimpleSelector {
	if !hasQName {
		// If we don't have a qualified name, we don't need to prepend a type selector.
		return selectors
	}

	// TODO: Check if has :host

	if name != "" {
		sel := &css.SimpleSelector{
			Match:    css.SelectorMatchTag,
			Data:     css.NewSelectorDataTag(namespace, name),
			Relation: css.SelectorRelationSubSelector,
		}
		selectors = append([]*css.SimpleSelector{sel}, selectors...) // Prepend the type selector
	} else if namespace != "" {
		sel := &css.SimpleSelector{
			Match:    css.SelectorMatchUniversalTag,
			Data:     css.NewSelectorDataTag(namespace, ""),
			Relation: css.SelectorRelationSubSelector,
		}
		selectors = append([]*css.SimpleSelector{sel}, selectors...) // Prepend the universal selector with namespace
	} else if len(selectors) == 0 {
		// If we only have a universal selector, we still need to return it.
		sel := &css.SimpleSelector{
			Match:    css.SelectorMatchUniversalTag,
			Data:     css.NewSelectorDataTag("", ""),
			Relation: css.SelectorRelationSubSelector,
		}
		selectors = append([]*css.SimpleSelector{sel}, selectors...) // Prepend the universal selector
	}

	return selectors
}

func parsePseudoType(name string, hasArguments bool) css.SelectorPseudoType {
	if hasArguments {
		pseudoType, ok := PseudoTypeWithArgumentsMap[name]
		if ok {
			return pseudoType
		}
	} else {
		pseudoType, ok := PseudoTypeWithoutArgumentsMap[name]
		if ok {
			return pseudoType
		}
	}

	if strings.HasPrefix(name, "-webkit-") {
		return css.SelectorPseudoWebKitCustomElement
	}
	if strings.HasPrefix(name, "-internal-") {
		return css.SelectorPseudoBlinkInternalElement
	}

	return css.SelectorPseudoUnknown
}

// convertRelationToRelative converts regular relations to relative relations
func convertRelationToRelative(relation css.SelectorRelationType) css.SelectorRelationType {
	switch relation {
	case css.SelectorRelationChild:
		return css.SelectorRelationRelativeChild
	case css.SelectorRelationDescendant:
		return css.SelectorRelationRelativeDescendant
	case css.SelectorRelationDirectAdjacent:
		return css.SelectorRelationRelativeDirectAdjacent
	case css.SelectorRelationIndirectAdjacent:
		return css.SelectorRelationRelativeIndirectAdjacent
	default:
		return css.SelectorRelationRelativeDescendant // Default for :has()
	}
}
