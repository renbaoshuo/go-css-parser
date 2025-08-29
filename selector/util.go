package selector

import (
	"strings"
)

func prependTypeSelectorIfNeeded(selectors []*SimpleSelector, name, namespace string, hasQName bool) []*SimpleSelector {
	if !hasQName {
		// If we don't have a qualified name, we don't need to prepend a type selector.
		return selectors
	}

	// TODO: Check if has :host

	if name != "" {
		sel := &SimpleSelector{
			Match:    SelectorMatchTag,
			Data:     NewSelectorDataTag(namespace, name),
			Relation: SelectorRelationSubSelector,
		}
		selectors = append([]*SimpleSelector{sel}, selectors...) // Prepend the type selector
	} else if namespace != "" {
		sel := &SimpleSelector{
			Match:    SelectorMatchUniversalTag,
			Data:     NewSelectorDataTag(namespace, ""),
			Relation: SelectorRelationSubSelector,
		}
		selectors = append([]*SimpleSelector{sel}, selectors...) // Prepend the universal selector with namespace
	} else if len(selectors) == 0 {
		// If we only have a universal selector, we still need to return it.
		sel := &SimpleSelector{
			Match:    SelectorMatchUniversalTag,
			Data:     NewSelectorDataTag("", ""),
			Relation: SelectorRelationSubSelector,
		}
		selectors = append([]*SimpleSelector{sel}, selectors...) // Prepend the universal selector
	}

	return selectors
}

func parsePseudoType(name string, hasArguments bool) SelectorPseudoType {
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
		return SelectorPseudoWebKitCustomElement
	}
	if strings.HasPrefix(name, "-internal-") {
		return SelectorPseudoBlinkInternalElement
	}

	return SelectorPseudoUnknown
}
