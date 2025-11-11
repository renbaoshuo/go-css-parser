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

// NewTokenStream creates a new TokenStream from the given input.
func NewTokenStream(input *csslexer.Input) *TokenStream {
	return &TokenStream{
		z: input,
		l: csslexer.NewLexer(input),
		p: nil,
		b: make(map[csslexer.TokenType]bool),
	}
}

// Peek returns the current token without consuming it.
func (s *TokenStream) Peek() csslexer.Token {
	if s.p == nil {
		// Skip comment tokens automatically
		for {
			token := s.l.Next()
			if token.Type != csslexer.CommentToken {
				p := tokenPool.Get().(*csslexer.Token)
				p.Type, p.Value, p.Raw = token.Type, token.Value, token.Raw
				s.p = p
				break
			}
		}
	}

	return csslexer.Token{
		Type:  s.p.Type,
		Value: s.p.Value,
		Raw:   s.p.Raw,
	}
}

// Consume consumes the current token and returns it.
func (s *TokenStream) Consume() csslexer.Token {
	if s.p != nil {
		tt, value, raw := s.p.Type, s.p.Value, s.p.Raw
		tokenPool.Put(s.p)
		s.p = nil

		return csslexer.Token{
			Type:  tt,
			Value: value,
			Raw:   raw,
		}
	} else {
		// Skip comment tokens automatically
		for {
			token := s.l.Next()
			if token.Type != csslexer.CommentToken {
				return token
			}
		}
	}
}

// AtEnd returns true if the current token is an EOF token or a boundary token.
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
