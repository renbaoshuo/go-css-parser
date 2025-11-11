package token_stream

import (
	"go.baoshuo.dev/csslexer"
)

// isBlockStartToken returns true if the token is a block start token.
func isBlockStartToken(tt csslexer.TokenType) bool {
	switch tt {
	case csslexer.LeftBraceToken, csslexer.LeftParenthesisToken, csslexer.LeftBracketToken, csslexer.FunctionToken:
		return true
	default:
		return false
	}
}

// isBlockEndToken returns true if the token is a block end token.
func isBlockEndToken(tt csslexer.TokenType) bool {
	switch tt {
	case csslexer.RightBraceToken, csslexer.RightParenthesisToken, csslexer.RightBracketToken:
		return true
	default:
		return false
	}
}

// getMatchingBlockEndToken returns the matching end token type for a given start token type.
func getMatchingBlockEndToken(tt csslexer.TokenType) csslexer.TokenType {
	switch tt {
	case csslexer.LeftBraceToken:
		return csslexer.RightBraceToken
	case csslexer.LeftParenthesisToken:
		return csslexer.RightParenthesisToken
	case csslexer.LeftBracketToken:
		return csslexer.RightBracketToken
	case csslexer.FunctionToken:
		return csslexer.RightParenthesisToken
	default:
		return csslexer.DefaultToken // Should not happen
	}
}

// isBlockToken returns true if the token is a block-related token.
func IsBlockToken(tt csslexer.TokenType) bool {
	return isBlockStartToken(tt) || isBlockEndToken(tt)
}
