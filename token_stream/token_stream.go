package token_stream

import (
	"go.baoshuo.dev/csslexer"
)

type TokenStream struct {
	z *csslexer.Input
	l *csslexer.Lexer
}

func NewTokenStream(input *csslexer.Input) *TokenStream {
	return &TokenStream{
		z: input,
		l: csslexer.NewLexer(input),
	}
}

func (s *TokenStream) Peek() (csslexer.TokenType, []rune) {
	return s.l.Peek()
}

func (s *TokenStream) Next() (csslexer.TokenType, []rune) {
	return s.l.Next()
}
