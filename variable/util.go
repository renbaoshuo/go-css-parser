package variable

import (
	"go.baoshuo.dev/csslexer"

	"go.baoshuo.dev/cssparser/token_stream"
)

func IsValidVariableName(token csslexer.Token) bool {
	if token.Type != csslexer.IdentToken {
		return false
	}

	return len(token.Value) >= 3 && token.Value[0] == '-' && token.Value[1] == '-'
}

func StartsCustomPropertyDeclaration(ts token_stream.TokenStream) bool {
	if !IsValidVariableName(ts.Peek()) {
		return false
	}
	state := ts.State()
	ts.ConsumeIncludingWhitespace() // <ident-token>
	result := ts.Peek().Type == csslexer.ColonToken
	state.Restore()
	return result
}
