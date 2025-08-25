package selector

import (
	"go.baoshuo.dev/cssparser/nesting"
	"go.baoshuo.dev/cssparser/rule"
	"go.baoshuo.dev/cssparser/token_stream"
)

type SelectorParser struct {
	tokenStream          *token_stream.TokenStream
	parentRuleForNesting *rule.StyleRule
}

func NewSelectorParser(
	tokenStream *token_stream.TokenStream,
	parentRuleForNesting *rule.StyleRule,
) *SelectorParser {
	return &SelectorParser{
		tokenStream:          tokenStream,
		parentRuleForNesting: parentRuleForNesting,
	}
}

func ConsumeSelector(
	tokenStream *token_stream.TokenStream,
	nestingType nesting.NestingTypeType,
	parentRuleForNesting *rule.StyleRule,
) ([]*Selector, error) {
	tokenStream.ConsumeWhitespace()
	return NewSelectorParser(tokenStream, parentRuleForNesting).consumeComplexSelectorList(nestingType)
}
