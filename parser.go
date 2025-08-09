package cssparser

import (
	"go.baoshuo.dev/csslexer"
)

type Parser struct {
	lexer *csslexer.Lexer
}

func NewParser(lexer *csslexer.Lexer) *Parser {
	return &Parser{
		lexer: lexer,
	}
}

func (p *Parser) ParseStylesheet() ([]*Rule, error) {
	return p.consumeRuleList()
}
