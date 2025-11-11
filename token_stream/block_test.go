package token_stream

import (
	"testing"

	"go.baoshuo.dev/csslexer"
)

func Test_TokenStream_ConsumeBlock(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType csslexer.TokenType
		consumer     func(ts *TokenStream) error
	}{
		{
			"simple block skipping",
			`{ foo: bar; }`,
			csslexer.EOFToken,
			func(ts *TokenStream) error {
				return nil
			},
		},
		{
			"simple block consuming",
			`{ foo: bar; }`,
			csslexer.EOFToken,
			func(ts *TokenStream) error {
				// Consume the block content
				for !ts.AtEnd() {
					ts.Consume()
				}
				return nil
			},
		},
		{
			"nested block skipping",
			`{ foo: bar; { nested: value; } baz: qux; }`,
			csslexer.EOFToken,
			func(ts *TokenStream) error {
				return nil
			},
		},
		{
			"nested bad block consuming",
			`{ foo: bar; { nested: value); } baz: qux; }`,
			csslexer.EOFToken,
			func(ts *TokenStream) error {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := csslexer.NewInput(tt.input)
			ts := NewTokenStream(input)

			err := ts.ConsumeBlock(tt.consumer)
			if err != nil {
				t.Errorf("ConsumeBlock() error = %v", err)
				return
			}

			endToken := ts.Peek()
			if endToken.Type != tt.expectedType {
				t.Errorf("Expected end token type %s, got %s", tt.expectedType, endToken.Type)
			}
		})
	}
}
