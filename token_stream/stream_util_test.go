package token_stream

import (
	"testing"

	"go.baoshuo.dev/csslexer"
)

func TestConsumeWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected csslexer.TokenType
	}{
		{
			name:     "no whitespace",
			input:    "div",
			expected: csslexer.IdentToken,
		},
		{
			name:     "single whitespace",
			input:    " div",
			expected: csslexer.IdentToken,
		},
		{
			name:     "multiple whitespace",
			input:    "   \t\n  div",
			expected: csslexer.IdentToken,
		},
		{
			name:     "only whitespace",
			input:    "   \t\n  ",
			expected: csslexer.EOFToken,
		},
		{
			name:     "whitespace before EOF",
			input:    "div   ",
			expected: csslexer.IdentToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := csslexer.NewInput(tt.input)
			ts := NewTokenStream(input)

			ts.ConsumeWhitespace()

			token := ts.Peek()
			if token.Type != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, token.Type)
			}
		})
	}
}

func TestConsumeIncludingWhitespace(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedConsumed  csslexer.TokenType
		expectedConsumedValue string
		expectedNext      csslexer.TokenType
	}{
		{
			name:              "token with following whitespace",
			input:             "div   { color: red; }",
			expectedConsumed:  csslexer.IdentToken,
			expectedConsumedValue: "div",
			expectedNext:      csslexer.LeftBraceToken,
		},
		{
			name:              "token without following whitespace",
			input:             "div{ color: red; }",
			expectedConsumed:  csslexer.IdentToken,
			expectedConsumedValue: "div",
			expectedNext:      csslexer.LeftBraceToken,
		},
		{
			name:              "token at end",
			input:             "div   ",
			expectedConsumed:  csslexer.IdentToken,
			expectedConsumedValue: "div",
			expectedNext:      csslexer.EOFToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := csslexer.NewInput(tt.input)
			ts := NewTokenStream(input)

			consumed := ts.ConsumeIncludingWhitespace()

			if consumed.Type != tt.expectedConsumed {
				t.Errorf("Expected consumed token type %v, got %v", tt.expectedConsumed, consumed.Type)
			}
			if consumed.Value != tt.expectedConsumedValue {
				t.Errorf("Expected consumed token value %q, got %q", tt.expectedConsumedValue, consumed.Value)
			}

			next := ts.Peek()
			if next.Type != tt.expectedNext {
				t.Errorf("Expected next token type %v, got %v", tt.expectedNext, next.Type)
			}
		})
	}
}

func TestSkipUntil(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		stopTypes  []csslexer.TokenType
		expectedAt csslexer.TokenType
	}{
		{
			name:       "skip until semicolon",
			input:      "color: red; margin: 10px;",
			stopTypes:  []csslexer.TokenType{csslexer.SemicolonToken},
			expectedAt: csslexer.SemicolonToken,
		},
		{
			name:       "skip until EOF",
			input:      "color: red",
			stopTypes:  []csslexer.TokenType{csslexer.SemicolonToken},
			expectedAt: csslexer.EOFToken,
		},
		{
			name:       "skip with nested blocks",
			input:      ".foo { color: red; } bar ; baz",
			stopTypes:  []csslexer.TokenType{csslexer.SemicolonToken},
			expectedAt: csslexer.SemicolonToken,
		},
		{
			name:       "skip with multiple stop types",
			input:      "div, p { color: red; }",
			stopTypes:  []csslexer.TokenType{csslexer.CommaToken, csslexer.LeftBraceToken},
			expectedAt: csslexer.CommaToken,
		},
		{
			name:       "skip until unmatched closing brace",
			input:      "color: red }",
			stopTypes:  []csslexer.TokenType{csslexer.SemicolonToken},
			expectedAt: csslexer.EOFToken, // Unmatched brace is consumed, then EOF is reached
		},
		{
			name:       "skip with nested parentheses",
			input:      "calc(100% - (20px + 5px)) ; more",
			stopTypes:  []csslexer.TokenType{csslexer.SemicolonToken},
			expectedAt: csslexer.SemicolonToken,
		},
		{
			name:       "skip with nested brackets",
			input:      "attr[data-value='test[nested]'] ; more",
			stopTypes:  []csslexer.TokenType{csslexer.SemicolonToken},
			expectedAt: csslexer.SemicolonToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := csslexer.NewInput(tt.input)
			ts := NewTokenStream(input)

			ts.SkipUntil(tt.stopTypes...)

			token := ts.Peek()
			if token.Type != tt.expectedAt {
				t.Errorf("Expected to stop at %v, got %v", tt.expectedAt, token.Type)
			}
		})
	}
}

func TestSkipUntilWithMismatchedBlocks(t *testing.T) {
	input := csslexer.NewInput("{ color: red; ) more")
	ts := NewTokenStream(input)

	// Start inside a brace block
	ts.Consume() // consume '{'

	ts.SkipUntil(csslexer.SemicolonToken)

	// Should continue past the mismatched ')' and stop at ';'
	token := ts.Peek()
	if token.Type != csslexer.SemicolonToken {
		t.Errorf("Expected SemicolonToken, got %v", token.Type)
	}
}