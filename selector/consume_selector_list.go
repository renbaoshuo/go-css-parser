package selector

import (
	"errors"

	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/css"
	"go.baoshuo.dev/cssparser/nesting"
)

func (sp *SelectorParser) consumeComplexSelectorList(nestingType nesting.NestingTypeType) ([]*css.Selector, error) {
	var selectors []*css.Selector
	firstInComplexSelectorList := true

	for {
		sel, err := sp.consumeComplexSelector(nesting.NestingTypeNone, firstInComplexSelectorList)
		if err != nil || !sp.atEndOfSelector() {
			sp.tokenStream.SkipUntil(csslexer.LeftBraceToken, csslexer.CommaToken)

			return nil, err
		}

		firstInComplexSelectorList = false
		selectors = append(selectors, sel)

		if sp.tokenStream.AtEnd() {
			break
		}

		nextToken := sp.tokenStream.Peek()
		if nextToken.Type == csslexer.LeftBraceToken {
			break
		}

		// at here, we assure the next token is a comma,
		// so we can safely consume it and its trailing whitespace.
		sp.tokenStream.ConsumeIncludingWhitespace()
	}

	return selectors, nil
}

// consumeCompoundSelectorList parses a comma-separated list of compound selectors
// Used by :is(), :where(), :any(), :host(), :host-context(), :cue()
func (sp *SelectorParser) consumeCompoundSelectorList() ([]*css.Selector, error) {
	var selectors []*css.Selector

	sel, err := sp.consumeCompoundSelectorAsComplexSelector()
	if err != nil {
		return nil, err
	}
	selectors = append(selectors, sel)
	sp.tokenStream.ConsumeWhitespace()

	for sp.tokenStream.Peek().Type == csslexer.CommaToken {
		sp.tokenStream.ConsumeIncludingWhitespace()

		sel, err := sp.consumeCompoundSelectorAsComplexSelector()
		if err != nil {
			return nil, err
		}
		selectors = append(selectors, sel)
		sp.tokenStream.ConsumeWhitespace()
	}

	return selectors, nil
}

// consumeCompoundSelectorAsComplexSelector wraps a compound selector as a complex selector
func (sp *SelectorParser) consumeCompoundSelectorAsComplexSelector() (*css.Selector, error) {
	compoundSelectors, flags := sp.consumeCompoundSelector(nesting.NestingTypeNone)
	if len(compoundSelectors) == 0 {
		return nil, errors.New("expected compound selector")
	}

	sel := &css.Selector{}
	sel.Flag.Set(flags)
	sel.Append(compoundSelectors...)

	return sel, nil
}

// consumeNestedSelectorList parses a nested selector list for :not()
func (sp *SelectorParser) consumeNestedSelectorList() ([]*css.Selector, error) {
	return sp.consumeComplexSelectorList(nesting.NestingTypeNone)
}

// consumeForgivingNestedSelectorList parses a forgiving nested selector list for :is(), :where()
func (sp *SelectorParser) consumeForgivingNestedSelectorList() ([]*css.Selector, error) {
	var selectors []*css.Selector
	firstInList := true

	for !sp.tokenStream.AtEnd() {
		// Save state to restore if parsing fails
		state := sp.tokenStream.State()

		sel, err := sp.consumeComplexSelector(nesting.NestingTypeNone, firstInList)
		if err != nil || !sp.atEndOfSelector() {
			// Restore state and skip to next comma or end
			state.Restore()
			sp.skipInvalidSelector()
		} else {
			selectors = append(selectors, sel)
		}

		if sp.tokenStream.Peek().Type != csslexer.CommaToken {
			break
		}
		sp.tokenStream.ConsumeIncludingWhitespace()
		firstInList = false
	}

	return selectors, nil
}

// consumeRelativeSelectorList parses a relative selector list for :has()
func (sp *SelectorParser) consumeRelativeSelectorList() ([]*css.Selector, error) {
	var selectors []*css.Selector

	sel, err := sp.consumeRelativeSelector()
	if err != nil {
		return nil, err
	}
	selectors = append(selectors, sel)

	for sp.tokenStream.Peek().Type == csslexer.CommaToken {
		sp.tokenStream.ConsumeIncludingWhitespace()

		sel, err := sp.consumeRelativeSelector()
		if err != nil {
			return nil, err
		}
		selectors = append(selectors, sel)
	}

	return selectors, nil
}

// consumeRelativeSelector parses a single relative selector
func (sp *SelectorParser) consumeRelativeSelector() (*css.Selector, error) {
	sel := &css.Selector{}

	// Create implicit relative anchor
	anchorSelector := &css.SimpleSelector{
		Match: css.SelectorMatchPseudoClass,
		Data:  css.NewSelectorDataPseudo("-internal-relative-anchor", css.SelectorPseudoRelativeAnchor),
	}
	sel.Append(anchorSelector)

	// Parse combinator if present
	combinator := convertRelationToRelative(sp.consumeCombinator())

	// Parse the rest of the complex selector
	rest, flags, err := sp.consumePartialComplexSelector(nesting.NestingTypeNone, combinator)
	if err != nil {
		return nil, err
	}

	sel.Flag.Set(flags)
	sel.Append(rest...)

	return sel, nil
}

// skipInvalidSelector skips tokens until we reach a comma or end of block
func (sp *SelectorParser) skipInvalidSelector() {
	for !sp.tokenStream.AtEnd() {
		token := sp.tokenStream.Peek()
		if token.Type == csslexer.CommaToken {
			break
		}
		sp.tokenStream.Consume()
	}
}
