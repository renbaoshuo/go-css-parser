package cssparser

import (
	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/nesting"
	"go.baoshuo.dev/cssparser/rule"
	"go.baoshuo.dev/cssparser/token_stream"
)

type Parser struct {
	s *token_stream.TokenStream
}

func NewParser(input *csslexer.Input) *Parser {
	return &Parser{
		s: token_stream.NewTokenStream(input),
	}
}

func (p *Parser) ParseStylesheet() ([]*rule.StyleRule, error) {
	return p.consumeRuleList(
		topLevelAllowedRules,
		true,
		nesting.NestingTypeNone,
		nil,
	)
}
