package token_stream

import (
	"go.baoshuo.dev/csslexer"
)

func (ts *TokenStream) SetBoundary(tokenType csslexer.TokenType, enable bool) {
	if enable {
		ts.b[tokenType] = true
	} else {
		delete(ts.b, tokenType)
	}
}
