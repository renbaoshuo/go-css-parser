package token_stream

import (
	"go.baoshuo.dev/csslexer"
)

type TokenStream struct {
	z *csslexer.Input
	l *csslexer.Lexer
	p *csslexer.Token
}

func NewTokenStream(input *csslexer.Input) *TokenStream {
	return &TokenStream{
		z: input,
		l: csslexer.NewLexer(input),
		p: nil,
	}
}

func (s *TokenStream) Peek() csslexer.Token {
	if s.p == nil {
		token := s.l.Next()
		p := tokenPool.Get().(*csslexer.Token)
		p.Type = token.Type
		p.Data = token.Data
		s.p = p
	}

	return csslexer.Token{
		Type: s.p.Type,
		Data: s.p.Data,
	}
}

func (s *TokenStream) Consume() csslexer.Token {
	if s.p != nil {
		tt, raw := s.p.Type, s.p.Data
		tokenPool.Put(s.p)
		s.p = nil

		return csslexer.Token{
			Type: tt,
			Data: raw,
		}
	} else {
		return s.l.Next()
	}
}
