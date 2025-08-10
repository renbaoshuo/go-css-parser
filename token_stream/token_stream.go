package token_stream

import (
	"go.baoshuo.dev/csslexer"
)

type TokenStream struct {
	z *csslexer.Input             // The input to the lexer.
	l *csslexer.Lexer             // The lexer that reads from the input.
	p *csslexer.Token             // The current token being processed.
	b map[csslexer.TokenType]bool // Boundary tokens, used to determine if the current token is a boundary token.
}

func NewTokenStream(input *csslexer.Input) *TokenStream {
	return &TokenStream{
		z: input,
		l: csslexer.NewLexer(input),
		p: nil,
		b: make(map[csslexer.TokenType]bool),
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

func (ts *TokenStream) AtEnd() bool {
	token := ts.Peek()
	if token.Type == csslexer.EOFToken {
		return true
	}
	if ts.b[token.Type] {
		return true
	}
	return false
}
