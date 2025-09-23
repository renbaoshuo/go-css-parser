package token_stream

import (
	"testing"

	"go.baoshuo.dev/csslexer"
)

func TestNewTokenStream(t *testing.T) {
	input := csslexer.NewInput("div { color: red; }")
	ts := NewTokenStream(input)

	if ts == nil {
		t.Fatal("Expected TokenStream, got nil")
	}
	if ts.z != input {
		t.Error("Expected input to be set")
	}
	if ts.l == nil {
		t.Error("Expected lexer to be created")
	}
	if ts.p != nil {
		t.Error("Expected peeked token to be nil initially")
	}
	if ts.b == nil {
		t.Error("Expected boundaries map to be initialized")
	}
	if len(ts.b) != 0 {
		t.Error("Expected boundaries map to be empty initially")
	}
}

func TestPeek(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedType  csslexer.TokenType
		expectedValue string
	}{
		{
			name:          "identifier",
			input:         "div",
			expectedType:  csslexer.IdentToken,
			expectedValue: "div",
		},
		{
			name:          "left brace",
			input:         "{",
			expectedType:  csslexer.LeftBraceToken,
			expectedValue: "{",
		},
		{
			name:          "string",
			input:         "\"hello\"",
			expectedType:  csslexer.StringToken,
			expectedValue: "hello",
		},
		{
			name:          "number",
			input:         "123",
			expectedType:  csslexer.NumberToken,
			expectedValue: "123",
		},
		{
			name:          "whitespace",
			input:         "   div",
			expectedType:  csslexer.WhitespaceToken,
			expectedValue: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := csslexer.NewInput(tt.input)
			ts := NewTokenStream(input)

			token := ts.Peek()
			if token.Type != tt.expectedType {
				t.Errorf("Expected type %v, got %v", tt.expectedType, token.Type)
			}
			if token.Value != tt.expectedValue {
				t.Errorf("Expected value %q, got %q", tt.expectedValue, token.Value)
			}

			// Peek again should return the same token
			token2 := ts.Peek()
			if token2.Type != token.Type || token2.Value != token.Value {
				t.Error("Peek should return the same token on subsequent calls")
			}
		})
	}
}

func TestPeekSkipsComments(t *testing.T) {
	input := csslexer.NewInput("/* comment */ div")
	ts := NewTokenStream(input)

	token := ts.Peek()
	if token.Type != csslexer.WhitespaceToken {
		t.Errorf("Expected WhitespaceToken (after comment), got %v", token.Type)
	}
}

func TestConsume(t *testing.T) {
	input := csslexer.NewInput("div { color: red; }")
	ts := NewTokenStream(input)

	// First consume
	token1 := ts.Consume()
	if token1.Type != csslexer.IdentToken {
		t.Errorf("Expected IdentToken, got %v", token1.Type)
	}
	if token1.Value != "div" {
		t.Errorf("Expected 'div', got %q", token1.Value)
	}

	// Second consume
	token2 := ts.Consume()
	if token2.Type != csslexer.WhitespaceToken {
		t.Errorf("Expected WhitespaceToken, got %v", token2.Type)
	}

	// Third consume
	token3 := ts.Consume()
	if token3.Type != csslexer.LeftBraceToken {
		t.Errorf("Expected LeftBraceToken, got %v", token3.Type)
	}
}

func TestConsumeSkipsComments(t *testing.T) {
	input := csslexer.NewInput("/* comment */ div")
	ts := NewTokenStream(input)

	token := ts.Consume()
	if token.Type != csslexer.WhitespaceToken {
		t.Errorf("Expected WhitespaceToken (after comment), got %v", token.Type)
	}
}

func TestConsumeAfterPeek(t *testing.T) {
	input := csslexer.NewInput("div { color: red; }")
	ts := NewTokenStream(input)

	// Peek first
	peeked := ts.Peek()
	if peeked.Type != csslexer.IdentToken {
		t.Errorf("Expected IdentToken, got %v", peeked.Type)
	}

	// Consume should return the same token
	consumed := ts.Consume()
	if consumed.Type != peeked.Type || consumed.Value != peeked.Value {
		t.Error("Consume after Peek should return the same token")
	}

	// Next consume should advance
	next := ts.Consume()
	if next.Type != csslexer.WhitespaceToken {
		t.Errorf("Expected WhitespaceToken, got %v", next.Type)
	}
}

func TestAtEnd(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		boundaries map[csslexer.TokenType]bool
		consumeN   int
		expected   bool
	}{
		{
			name:     "at EOF",
			input:    "div",
			consumeN: 1,
			expected: true,
		},
		{
			name:     "not at end",
			input:    "div { color: red; }",
			consumeN: 0,
			expected: false,
		},
		{
			name:       "at boundary token",
			input:      "div ; color: red;",
			boundaries: map[csslexer.TokenType]bool{csslexer.SemicolonToken: true},
			consumeN:   2, // consume 'div' and whitespace
			expected:   true,
		},
		{
			name:       "not at boundary token",
			input:      "div : color",
			boundaries: map[csslexer.TokenType]bool{csslexer.SemicolonToken: true},
			consumeN:   0,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := csslexer.NewInput(tt.input)
			ts := NewTokenStream(input)

			if tt.boundaries != nil {
				ts.b = tt.boundaries
			}

			for i := 0; i < tt.consumeN; i++ {
				ts.Consume()
			}

			result := ts.AtEnd()
			if result != tt.expected {
				t.Errorf("Expected AtEnd() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTokenStreamIntegration(t *testing.T) {
	input := csslexer.NewInput("div.class { color: red; margin: 10px; }")
	ts := NewTokenStream(input)

	// Test the flow of parsing
	tokens := []struct {
		expectedType  csslexer.TokenType
		expectedValue string
	}{
		{csslexer.IdentToken, "div"},
		{csslexer.DelimiterToken, "."},
		{csslexer.IdentToken, "class"},
		{csslexer.WhitespaceToken, " "},
		{csslexer.LeftBraceToken, "{"},
		{csslexer.WhitespaceToken, " "},
		{csslexer.IdentToken, "color"},
		{csslexer.ColonToken, ":"},
		{csslexer.WhitespaceToken, " "},
		{csslexer.IdentToken, "red"},
		{csslexer.SemicolonToken, ";"},
		{csslexer.WhitespaceToken, " "},
		{csslexer.IdentToken, "margin"},
		{csslexer.ColonToken, ":"},
		{csslexer.WhitespaceToken, " "},
		{csslexer.DimensionToken, "10px"},
		{csslexer.SemicolonToken, ";"},
		{csslexer.WhitespaceToken, " "},
		{csslexer.RightBraceToken, "}"},
		{csslexer.EOFToken, ""},
	}

	for i, expected := range tokens {
		token := ts.Consume()
		if token.Type != expected.expectedType {
			t.Errorf("Token %d: expected type %v, got %v", i, expected.expectedType, token.Type)
		}
		if token.Value != expected.expectedValue {
			t.Errorf("Token %d: expected value %q, got %q", i, expected.expectedValue, token.Value)
		}
	}
}
