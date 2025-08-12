package selector

import (
	"errors"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser"
)

func (sp *SelectorParser) consumeComplexSelector(nestingType cssparser.NestingTypeType) (*Selector, error) {
	sel := &Selector{}

	if nestingType != cssparser.NestingTypeNone && sp.peekIsCombinator() {
		// Nested selectors that start with a combinator are to be
		// interpreted as relative selectors (with the anchor being
		// the parent selector, i.e., &).
		return sp.consumeNestedRelativeSelector(nestingType)
	}

	compoundSelectors := sp.consumeCompoundSelector(nestingType)
	if len(compoundSelectors) == 0 {
		return nil, errors.New("invalid selector: no compound selectors found")
	}
	sel.Append(compoundSelectors...)

	combinator := sp.consumeCombinator()
	if combinator != SelectorRelationSubSelector {
		sel.Flag.Set(SelectorFlagContainsComplexSelector)

		rest, err := sp.consumePartialComplexSelector(nestingType, combinator)
		if err != nil {
			return nil, err
		}

		sel.Append(rest...)
	}

	// TODO: handle if in nested top-level rules
	// if nestingType != cssparser.NestingTypeNone {
	// }

	return sel, nil
}

func (sp *SelectorParser) consumePartialComplexSelector(
	nestingType cssparser.NestingTypeType,
	combinator SelectorRelationType,
) ([]*SimpleSelector, error) {
	selectors := make([]*SimpleSelector, 0)
	for {
		compound := sp.consumeCompoundSelector(nestingType)
		if len(compound) == 0 {
			if combinator == SelectorRelationDescendant {
				break
			} else {
				return nil, errors.New("invalid selector: expected compound selector")
			}
		}

		compound[0].Relation = combinator // Set the relation for the first selector
		selectors = append(selectors, compound...)

		combinator = sp.consumeCombinator()
		if combinator == SelectorRelationSubSelector {
			break
		}
	}
	return selectors, nil
}

func (sp *SelectorParser) consumeCompoundSelector(nestingType cssparser.NestingTypeType) []*SimpleSelector {
	var selectors []*SimpleSelector

	// See if the compound selector starts with a tag name, universal selector
	// or the likes (these can only be at the beginning). Note that we don't
	// add this to output yet, because there are situations where it should
	// be ignored (like if we have a universal selector and don't need it;
	// e.g. *:hover is the same as :hover). Thus, we just keep its data around
	// and prepend it if needed.
	name, namespace, hasQName := sp.consumeName()

	// TODO: A tag name is not valid following a pseudo-element.

	for {
		selector, err := sp.consumeSimpleSelector()
		if err != nil {
			break
		}

		// TODO: handle pseudo-elements

		selector.Relation = SelectorRelationSubSelector
		selectors = append(selectors, selector)
	}

	selectors = prependTypeSelectorIfNeeded(selectors, name, namespace, hasQName)

	return selectors
}

// consumeName consumes a name token and returns the name and its namespace if applicable.
//
// Returns:
//   - The name as a string.
//   - The namespace as a string (empty if not applicable).
//   - Whether the name was successfully consumed.
func (sp *SelectorParser) consumeName() ([]rune, []rune, bool) {
	var name, namespace []rune

	first := sp.tokenStream.Peek()
	switch first.Type {
	case csslexer.IdentToken:
		name = first.Data
		sp.tokenStream.Consume()

	case csslexer.DelimiterToken:
		if len(first.Data) != 1 {
			return nil, nil, false // Invalid name
		}

		switch first.Data[0] {
		case '*':
			name = nil // This is a universal selector, no name.
			sp.tokenStream.Consume()

		case '|':
			// This is an empty namespace, no name.
			name = nil

		default:
			return nil, nil, false // Invalid name
		}

	default:
		return nil, nil, false // Invalid name
	}

	second := sp.tokenStream.Peek()
	if second.Type != csslexer.DelimiterToken || len(second.Data) != 1 || second.Data[0] != '|' {
		// No namespace, just a name.
		return name, nil, true
	}

	tss := sp.tokenStream.State()
	sp.tokenStream.Consume() // Consume the '|'

	if name == nil {
		namespace = []rune{'*'} // Universal selector with namespace
	} else {
		namespace = name // Use the name as the namespace
		name = nil       // Reset name to indicate that we are now looking for a namespace
	}

	third := sp.tokenStream.Peek()
	switch third.Type {
	case csslexer.IdentToken:
		name = third.Data
		sp.tokenStream.Consume()

	case csslexer.DelimiterToken:
		if len(third.Data) == 1 && third.Data[0] == '*' {
			name = nil // Universal selector, no name
			sp.tokenStream.Consume()
		} else {
			// Invalid name after namespace delimiter
			tss.Restore()
			return nil, nil, false
		}

	default:
		// Invalid token after namespace delimiter
		tss.Restore()
		return nil, nil, false
	}

	return name, namespace, true
}

func (sp *SelectorParser) consumeNestedRelativeSelector(nestingType cssparser.NestingTypeType) (*Selector, error) {
	return nil, errors.New("not implemented: SelectorParser.consumeNestedRelativeSelector")
}

func (sp *SelectorParser) consumeCombinator() SelectorRelationType {
	fallbackResult := SelectorRelationSubSelector
	for sp.tokenStream.Peek().Type == csslexer.WhitespaceToken {
		sp.tokenStream.Consume()
		fallbackResult = SelectorRelationDescendant
	}

	nextToken := sp.tokenStream.Peek()
	if nextToken.Type != csslexer.DelimiterToken {
		return fallbackResult // No combinator found, return fallback
	}
	if len(nextToken.Data) != 1 {
		return fallbackResult // Invalid combinator, return fallback
	}
	switch nextToken.Data[0] {
	case '>':
		sp.tokenStream.ConsumeIncludingWhitespace()
		return SelectorRelationChild
	case '+':
		sp.tokenStream.ConsumeIncludingWhitespace()
		return SelectorRelationDirectAdjacent
	case '~':
		sp.tokenStream.ConsumeIncludingWhitespace()
		return SelectorRelationIndirectAdjacent
	default:
		return fallbackResult // Invalid combinator, return fallback
	}
}
