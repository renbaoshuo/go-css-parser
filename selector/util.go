package selector

func equalIgnoreCase(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] && a[i] != b[i]+32 && a[i] != b[i]-32 { // ASCII case-insensitive comparison
			return false
		}
	}
	return true
}

func prependTypeSelectorIfNeeded(selectors []*SimpleSelector, name, namespace string, hasQName bool) []*SimpleSelector {
	if !hasQName {
		// If we don't have a qualified name, we don't need to prepend a type selector.
		return selectors
	}

	// TODO: Handle namespace uri
	nameStr := name
	if namespace != "" {
		nameStr = namespace + "|" + name
	}

	// TODO: Check if has :host

	if nameStr != "" {
		sel := &SimpleSelector{
			Match:    SelectorMatchTag,
			Value:    nameStr,
			Relation: SelectorRelationSubSelector,
		}
		selectors = append([]*SimpleSelector{sel}, selectors...) // Prepend the type selector
	} else if len(selectors) == 0 {
		// If we only have a universal selector, we still need to return it.
		sel := &SimpleSelector{
			Match:    SelectorMatchUniversalTag,
			Value:    namespace,
			Relation: SelectorRelationSubSelector,
		}
		selectors = append([]*SimpleSelector{sel}, selectors...) // Prepend the universal selector
	}

	return selectors
}
