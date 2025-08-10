package token_stream

import (
	"go.baoshuo.dev/csslexer"
)

type TokenStream struct {
	z *csslexer.Input
	l *csslexer.Lexer
	p *token
}

func NewTokenStream(input *csslexer.Input) *TokenStream {
	return &TokenStream{
		z: input,
		l: csslexer.NewLexer(input),
		p: nil,
	}
}

func (s *TokenStream) Peek() (csslexer.TokenType, []rune) {
	if s.p == nil {
		tt, raw := s.l.Next()
		p := tokenPool.Get().(*token)
		p.Type = tt
		p.Data = raw
		s.p = p
	}
	return s.p.Type, s.p.Data
}

func (s *TokenStream) Consume() (csslexer.TokenType, []rune) {
	if s.p != nil {
		tt, raw := s.p.Type, s.p.Data
		tokenPool.Put(s.p)
		s.p = nil
		return tt, raw
	}
	return s.l.Next()
}
