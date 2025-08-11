package selector

import (
	"go.baoshuo.dev/cssparser"
	"go.baoshuo.dev/cssparser/token_stream"
)

type SelectorParser struct {
	tokenStream          *token_stream.TokenStream
	parentRuleForNesting *cssparser.Rule
}

func NewSelectorParser(
	tokenStream *token_stream.TokenStream,
	parentRuleForNesting *cssparser.Rule,
) *SelectorParser {
	return &SelectorParser{
		tokenStream:          tokenStream,
		parentRuleForNesting: parentRuleForNesting,
	}
}

func ConsumeSelector(
	tokenStream *token_stream.TokenStream,
	nestingType cssparser.NestingTypeType,
	parentRuleForNesting *cssparser.Rule,
) ([]*Selector, error) {
	tokenStream.ConsumeWhitespace()
	return NewSelectorParser(tokenStream, parentRuleForNesting).consumeComplexSelectorList()
}
