package selector

import (
	"errors"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/css"
	"go.baoshuo.dev/cssparser/nesting"
)

func (sp *SelectorParser) consumeComplexSelector(
	nestingType nesting.NestingTypeType,
	firstInComplexSelector bool,
) (*css.Selector, error) {
	if nestingType != nesting.NestingTypeNone && sp.peekIsCombinator() {
		// Nested selectors that start with a combinator are to be
		// interpreted as relative selectors (with the anchor being
		// the parent selector, i.e., &).
		return sp.consumeNestedRelativeSelector(nestingType)
	}

	sel := &css.Selector{}

	compoundSelectors, firstFlags := sp.consumeCompoundSelector(nestingType)
	if len(compoundSelectors) == 0 {
		return nil, errors.New("invalid selector: no compound selectors found")
	}
	sel.Flag.Set(firstFlags)
	sel.Append(compoundSelectors...)

	if combinator := sp.consumeCombinator(); combinator != css.SelectorRelationSubSelector {
		rest, restFlags, err := sp.consumePartialComplexSelector(nestingType, combinator)
		if err != nil {
			return nil, err
		}

		sel.Flag.Set(css.SelectorFlagContainsComplexSelector)
		sel.Flag.Set(restFlags)
		sel.Append(rest...)
	}

	// TODO: handle if in nested top-level rules
	// if nestingType != nesting.NestingTypeNone {
	// }

	return sel, nil
}

func (sp *SelectorParser) consumePartialComplexSelector(
	nestingType nesting.NestingTypeType,
	combinator css.SelectorRelationType,
) ([]*css.SimpleSelector, css.SelectorListFlagType, error) {
	var flags css.SelectorListFlagType
	selectors := make([]*css.SimpleSelector, 0)

	for {
		compound, compoundFlags := sp.consumeCompoundSelector(nestingType)
		if len(compound) == 0 {
			if combinator == css.SelectorRelationDescendant {
				flags.Set(compoundFlags)
				break
			} else {
				return nil, 0, errors.New("invalid selector: expected compound selector")
			}
		}

		compound[0].Relation = combinator // Set the relation for the first selector
		flags.Set(compoundFlags)
		selectors = append(selectors, compound...)

		combinator = sp.consumeCombinator()
		if combinator == css.SelectorRelationSubSelector {
			break
		}
	}

	return selectors, flags, nil
}

func (sp *SelectorParser) consumeCompoundSelector(nestingType nesting.NestingTypeType) ([]*css.SimpleSelector, css.SelectorListFlagType) {
	var selectors []*css.SimpleSelector
	var flags css.SelectorListFlagType

	// See if the compound selector starts with a tag name, universal selector
	// or the likes (these can only be at the beginning). Note that we don't
	// add this to output yet, because there are situations where it should
	// be ignored (like if we have a universal selector and don't need it;
	// e.g. *:hover is the same as :hover). Thus, we just keep its data around
	// and prepend it if needed.
	name, namespace, hasQName := sp.consumeName()

	// TODO: A tag name is not valid following a pseudo-element.

	for {
		selector, selectorFlags, err := sp.consumeSimpleSelector()
		if err != nil {
			break
		}

		// TODO: handle pseudo-elements

		selector.Relation = css.SelectorRelationSubSelector
		flags.Set(selectorFlags)
		selectors = append(selectors, selector)
	}

	selectors = prependTypeSelectorIfNeeded(selectors, name, namespace, hasQName)

	return selectors, flags
}

// consumeName consumes a name token and returns the name and its namespace if applicable.
//
// Returns:
//   - The name as a string.
//   - The namespace as a string (empty if not applicable).
//   - Whether the name was successfully consumed.
func (sp *SelectorParser) consumeName() (string, string, bool) {
	var name, namespace string

	first := sp.tokenStream.Peek()
	switch first.Type {
	case csslexer.IdentToken:
		name = first.Value
		sp.tokenStream.Consume()

	case csslexer.DelimiterToken:
		switch first.Value {
		case "*":
			name = "" // Universal selector, no name
			sp.tokenStream.Consume()

		case "|":
			// This is an empty namespace, no name.
			name = ""

		default:
			return "", "", false // Invalid name
		}

	default:
		return "", "", false // Invalid name
	}

	second := sp.tokenStream.Peek()
	if second.Type != csslexer.DelimiterToken || second.Value != "|" {
		// No namespace, just a name.
		return name, "", true
	}

	tss := sp.tokenStream.State()
	sp.tokenStream.Consume() // Consume the '|'

	if name == "" {
		namespace = "*" // Universal selector with namespace
	} else {
		namespace = name // Use the name as the namespace
		name = ""        // Reset name to indicate that we are now looking for a namespace
	}

	third := sp.tokenStream.Peek()
	switch third.Type {
	case csslexer.IdentToken:
		name = third.Value
		sp.tokenStream.Consume()

	case csslexer.DelimiterToken:
		if third.Value == "*" {
			name = "" // Universal selector, no name
			sp.tokenStream.Consume()
		} else {
			// Invalid name after namespace delimiter
			tss.Restore()
			return "", "", false
		}

	default:
		// Invalid token after namespace delimiter
		tss.Restore()
		return "", "", false
	}

	return name, namespace, true
}

func (sp *SelectorParser) consumeNestedRelativeSelector(nestingType nesting.NestingTypeType) (*css.Selector, error) {
	return nil, errors.New("not implemented: SelectorParser.consumeNestedRelativeSelector")
}

func (sp *SelectorParser) consumeCombinator() css.SelectorRelationType {
	fallbackResult := css.SelectorRelationSubSelector
	for sp.tokenStream.Peek().Type == csslexer.WhitespaceToken {
		sp.tokenStream.Consume()
		fallbackResult = css.SelectorRelationDescendant
	}

	nextToken := sp.tokenStream.Peek()
	if nextToken.Type != csslexer.DelimiterToken {
		return fallbackResult // No combinator found, return fallback
	}
	switch nextToken.Value {
	case ">":
		sp.tokenStream.ConsumeIncludingWhitespace()
		return css.SelectorRelationChild
	case "+":
		sp.tokenStream.ConsumeIncludingWhitespace()
		return css.SelectorRelationDirectAdjacent
	case "~":
		sp.tokenStream.ConsumeIncludingWhitespace()
		return css.SelectorRelationIndirectAdjacent
	default:
		return fallbackResult // Invalid combinator, return fallback
	}
}
