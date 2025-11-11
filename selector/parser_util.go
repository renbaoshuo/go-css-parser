package selector

import (
	"go.baoshuo.dev/csslexer"
)

func (sp *SelectorParser) atEndOfSelector() bool {
	if sp.tokenStream.AtEnd() {
		return true
	}

	t := sp.tokenStream.Peek()

	return t.Type == csslexer.LeftBraceToken || t.Type == csslexer.CommaToken
}

func (sp *SelectorParser) peekIsCombinator() bool {
	sp.tokenStream.ConsumeWhitespace()

	t := sp.tokenStream.Peek()

	if t.Type != csslexer.DelimiterToken {
		return false
	}

	switch t.Value {
	case ">", "+", "~":
		return true
	default:
		return false
	}
}
