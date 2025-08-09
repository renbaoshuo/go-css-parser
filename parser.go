package cssparser

import (
	"go.baoshuo.dev/csslexer"

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

func (p *Parser) ParseStylesheet() ([]*Rule, error) {
	return p.consumeRuleList()
}
