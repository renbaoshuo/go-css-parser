package token_stream

import (
	"go.baoshuo.dev/csslexer"
)

func isBlockStartToken(tt csslexer.TokenType) bool {
	switch tt {
	case csslexer.LeftBraceToken, csslexer.LeftParenthesisToken, csslexer.LeftBracketToken:
		return true
	default:
		return false
	}
}

func isBlockEndToken(tt csslexer.TokenType) bool {
	switch tt {
	case csslexer.RightBraceToken, csslexer.RightParenthesisToken, csslexer.RightBracketToken:
		return true
	default:
		return false
	}
}
