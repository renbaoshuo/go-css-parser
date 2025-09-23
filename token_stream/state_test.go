package token_stream

import (
	"testing"

	"go.baoshuo.dev/csslexer"
)

func TestTokenStreamState(t *testing.T) {
	input := csslexer.NewInput("div { color: red; }")
	ts := NewTokenStream(input)

	// Consume a few tokens
	firstToken := ts.Consume() // 'div'
	if firstToken.Type != csslexer.IdentToken {
		t.Errorf("Expected IdentToken, got %v", firstToken.Type)
	}

	secondToken := ts.Consume() // whitespace
	if secondToken.Type != csslexer.WhitespaceToken {
		t.Errorf("Expected WhitespaceToken, got %v", secondToken.Type)
	}

	// Capture state before consuming the brace
	state := ts.State()

	// Verify state has correct reference
	if state.tokenStream != ts {
		t.Error("State should reference the original TokenStream")
	}

	// Consume more tokens
	braceToken := ts.Consume() // '{'
	if braceToken.Type != csslexer.LeftBraceToken {
		t.Errorf("Expected LeftBraceToken, got %v", braceToken.Type)
	}

	// Peek at next token to test peeked token state
	peekedToken := ts.Peek()
	if peekedToken.Type != csslexer.WhitespaceToken {
		t.Errorf("Expected WhitespaceToken, got %v", peekedToken.Type)
	}

	// Restore to previous state
	state.Restore()

	// Verify we're back to the state before consuming the brace
	nextToken := ts.Peek()
	if nextToken.Type != csslexer.LeftBraceToken {
		t.Errorf("Expected LeftBraceToken after restore, got %v", nextToken.Type)
	}
}

func TestTokenStreamStateWithPeekedToken(t *testing.T) {
	input := csslexer.NewInput("color: red;")
	ts := NewTokenStream(input)

	// Peek at a token
	peekedToken := ts.Peek()
	if peekedToken.Type != csslexer.IdentToken {
		t.Errorf("Expected IdentToken, got %v", peekedToken.Type)
	}
	if peekedToken.Value != "color" {
		t.Errorf("Expected 'color', got %q", peekedToken.Value)
	}

	// Capture state with peeked token
	state := ts.State()

	// Consume the peeked token
	consumed := ts.Consume()
	if consumed.Type != csslexer.IdentToken {
		t.Errorf("Expected IdentToken, got %v", consumed.Type)
	}

	// Consume more tokens
	ts.Consume() // ':'
	ts.Consume() // whitespace

	// Restore state
	state.Restore()

	// Verify the peeked token is restored
	restoredPeek := ts.Peek()
	if restoredPeek.Type != csslexer.IdentToken {
		t.Errorf("Expected IdentToken after restore, got %v", restoredPeek.Type)
	}
	if restoredPeek.Value != "color" {
		t.Errorf("Expected 'color' after restore, got %q", restoredPeek.Value)
	}
}

func TestTokenStreamStateWithBoundaries(t *testing.T) {
	input := csslexer.NewInput("div { color: red; }")
	ts := NewTokenStream(input)

	// Set some boundaries
	ts.b[csslexer.SemicolonToken] = true
	ts.b[csslexer.RightBraceToken] = true

	// Capture state
	state := ts.State()

	// Modify boundaries
	ts.b[csslexer.ColonToken] = true
	delete(ts.b, csslexer.SemicolonToken)

	// Verify boundaries are different
	if !ts.b[csslexer.ColonToken] {
		t.Error("Expected ColonToken boundary to be set")
	}
	if ts.b[csslexer.SemicolonToken] {
		t.Error("Expected SemicolonToken boundary to be removed")
	}

	// Restore state
	state.Restore()

	// Verify boundaries are restored
	if ts.b[csslexer.ColonToken] {
		t.Error("Expected ColonToken boundary to be removed after restore")
	}
	if !ts.b[csslexer.SemicolonToken] {
		t.Error("Expected SemicolonToken boundary to be restored")
	}
	if !ts.b[csslexer.RightBraceToken] {
		t.Error("Expected RightBraceToken boundary to be preserved")
	}
}

func TestTokenStreamStateNoPeekedToken(t *testing.T) {
	input := csslexer.NewInput("div")
	ts := NewTokenStream(input)

	// Capture state without peeking
	state := ts.State()

	// Verify no peeked token in state
	if state.peekedToken != nil {
		t.Error("Expected no peeked token in state")
	}

	// Peek to create a peeked token
	ts.Peek()

	// Restore state (should clear the peeked token)
	state.Restore()

	// Verify peeked token is cleared
	if ts.p != nil {
		t.Error("Expected peeked token to be cleared after restore")
	}
}