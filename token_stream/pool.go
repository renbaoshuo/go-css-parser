package token_stream

import (
	"sync"

	"go.baoshuo.dev/csslexer"
)

// token represents a parsed token with its type and data.
//
// It is used internally by the token stream to represent
// the tokens extracted from the lexer.
type token struct {
	Type csslexer.TokenType
	Data []rune
}

// tokenPool is a sync.Pool for reusing token instances.
var tokenPool = sync.Pool{
	New: func() interface{} {
		return &token{
			Type: csslexer.DefaultToken,
			Data: nil,
		}
	},
}
