package selector

import (
	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser"
)

func (sp *SelectorParser) consumeComplexSelectorList(nestingType cssparser.NestingTypeType) ([]*Selector, error) {
	var selectors []*Selector
	firstInComplexSelectorList := true

	for {
		sel, err := sp.consumeComplexSelector(cssparser.NestingTypeNone, firstInComplexSelectorList)
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
