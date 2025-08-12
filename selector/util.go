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

func prependTypeSelectorIfNeeded(selectors []*SimpleSelector, name, namespace []rune, hasQName bool) []*SimpleSelector {
	if !hasQName {
		// If we don't have a qualified name, we don't need to prepend a type selector.
		return selectors
	}

	// TODO: Handle namespace uri
	nameData := name
	if namespace != nil {
		nameData = append(namespace, '|')
		nameData = append(nameData, name...)
	}

	// TODO: Check if has :host

	if !equalIgnoreCase(nameData, []rune("*")) {
		sel := &SimpleSelector{
			Match:    SelectorMatchTag,
			Data:     nameData,
			Relation: SelectorRelationSubSelector,
		}
		selectors = append([]*SimpleSelector{sel}, selectors...) // Prepend the type selector
	}

	return selectors
}
