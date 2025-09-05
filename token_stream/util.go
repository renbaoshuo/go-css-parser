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

func getMatchingBlockEndToken(tt csslexer.TokenType) csslexer.TokenType {
	switch tt {
	case csslexer.LeftBraceToken:
		return csslexer.RightBraceToken
	case csslexer.LeftParenthesisToken:
		return csslexer.RightParenthesisToken
	case csslexer.LeftBracketToken:
		return csslexer.RightBracketToken
	default:
		return csslexer.DefaultToken // Should not happen
	}
}

// isBlockToken returns true if the token is a block-related token
func IsBlockToken(tt csslexer.TokenType) bool {
	return isBlockStartToken(tt) || isBlockEndToken(tt)
}
