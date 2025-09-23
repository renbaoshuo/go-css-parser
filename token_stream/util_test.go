package token_stream

import (
	"testing"

	"go.baoshuo.dev/csslexer"
)

func TestIsBlockStartToken(t *testing.T) {
	tests := []struct {
		tokenType csslexer.TokenType
		expected  bool
	}{
		{csslexer.LeftBraceToken, true},
		{csslexer.LeftParenthesisToken, true},
		{csslexer.LeftBracketToken, true},
		{csslexer.FunctionToken, true},
		{csslexer.RightBraceToken, false},
		{csslexer.RightParenthesisToken, false},
		{csslexer.RightBracketToken, false},
		{csslexer.IdentToken, false},
		{csslexer.StringToken, false},
		{csslexer.NumberToken, false},
		{csslexer.EOFToken, false},
	}

	for _, tt := range tests {
		t.Run(tt.tokenType.String(), func(t *testing.T) {
			result := isBlockStartToken(tt.tokenType)
			if result != tt.expected {
				t.Errorf("isBlockStartToken(%v) = %v, expected %v", tt.tokenType, result, tt.expected)
			}
		})
	}
}

func TestIsBlockEndToken(t *testing.T) {
	tests := []struct {
		tokenType csslexer.TokenType
		expected  bool
	}{
		{csslexer.RightBraceToken, true},
		{csslexer.RightParenthesisToken, true},
		{csslexer.RightBracketToken, true},
		{csslexer.LeftBraceToken, false},
		{csslexer.LeftParenthesisToken, false},
		{csslexer.LeftBracketToken, false},
		{csslexer.IdentToken, false},
		{csslexer.StringToken, false},
		{csslexer.NumberToken, false},
		{csslexer.EOFToken, false},
	}

	for _, tt := range tests {
		t.Run(tt.tokenType.String(), func(t *testing.T) {
			result := isBlockEndToken(tt.tokenType)
			if result != tt.expected {
				t.Errorf("isBlockEndToken(%v) = %v, expected %v", tt.tokenType, result, tt.expected)
			}
		})
	}
}

func TestGetMatchingBlockEndToken(t *testing.T) {
	tests := []struct {
		startToken csslexer.TokenType
		expected   csslexer.TokenType
	}{
		{csslexer.LeftBraceToken, csslexer.RightBraceToken},
		{csslexer.LeftParenthesisToken, csslexer.RightParenthesisToken},
		{csslexer.LeftBracketToken, csslexer.RightBracketToken},
		{csslexer.FunctionToken, csslexer.RightParenthesisToken},
		{csslexer.IdentToken, csslexer.DefaultToken}, // Should not happen case
		{csslexer.NumberToken, csslexer.DefaultToken}, // Should not happen case
	}

	for _, tt := range tests {
		t.Run(tt.startToken.String(), func(t *testing.T) {
			result := getMatchingBlockEndToken(tt.startToken)
			if result != tt.expected {
				t.Errorf("getMatchingBlockEndToken(%v) = %v, expected %v", tt.startToken, result, tt.expected)
			}
		})
	}
}

func TestIsBlockToken(t *testing.T) {
	tests := []struct {
		tokenType csslexer.TokenType
		expected  bool
	}{
		// Block start tokens
		{csslexer.LeftBraceToken, true},
		{csslexer.LeftParenthesisToken, true},
		{csslexer.LeftBracketToken, true},
		{csslexer.FunctionToken, true},
		// Block end tokens
		{csslexer.RightBraceToken, true},
		{csslexer.RightParenthesisToken, true},
		{csslexer.RightBracketToken, true},
		// Non-block tokens
		{csslexer.IdentToken, false},
		{csslexer.StringToken, false},
		{csslexer.NumberToken, false},
		{csslexer.EOFToken, false},
		{csslexer.ColonToken, false},
		{csslexer.SemicolonToken, false},
		{csslexer.CommaToken, false},
	}

	for _, tt := range tests {
		t.Run(tt.tokenType.String(), func(t *testing.T) {
			result := IsBlockToken(tt.tokenType)
			if result != tt.expected {
				t.Errorf("IsBlockToken(%v) = %v, expected %v", tt.tokenType, result, tt.expected)
			}
		})
	}
}

func TestBlockTokenPairs(t *testing.T) {
	// Test that block start tokens correctly map to their end tokens
	pairs := map[csslexer.TokenType]csslexer.TokenType{
		csslexer.LeftBraceToken:      csslexer.RightBraceToken,
		csslexer.LeftParenthesisToken: csslexer.RightParenthesisToken,
		csslexer.LeftBracketToken:    csslexer.RightBracketToken,
		csslexer.FunctionToken:       csslexer.RightParenthesisToken,
	}

	for start, expectedEnd := range pairs {
		t.Run(start.String()+"_to_"+expectedEnd.String(), func(t *testing.T) {
			if !isBlockStartToken(start) {
				t.Errorf("Expected %v to be a block start token", start)
			}
			if !isBlockEndToken(expectedEnd) {
				t.Errorf("Expected %v to be a block end token", expectedEnd)
			}
			actualEnd := getMatchingBlockEndToken(start)
			if actualEnd != expectedEnd {
				t.Errorf("Expected %v to map to %v, got %v", start, expectedEnd, actualEnd)
			}
			if !IsBlockToken(start) {
				t.Errorf("Expected %v to be a block token", start)
			}
			if !IsBlockToken(expectedEnd) {
				t.Errorf("Expected %v to be a block token", expectedEnd)
			}
		})
	}
}