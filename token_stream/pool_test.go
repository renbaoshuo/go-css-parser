package token_stream

import (
	"testing"

	"go.baoshuo.dev/csslexer"
)

func TestTokenPool(t *testing.T) {
	// Test getting token from pool
	token1 := tokenPool.Get().(*csslexer.Token)
	if token1 == nil {
		t.Error("Expected token from pool, got nil")
	}

	// Check initial values
	if token1.Type != csslexer.DefaultToken {
		t.Errorf("Expected DefaultToken, got %v", token1.Type)
	}
	if token1.Value != "" {
		t.Errorf("Expected empty string, got %q", token1.Value)
	}
	if len(token1.Raw) != 0 {
		t.Errorf("Expected empty Raw slice, got length %d", len(token1.Raw))
	}

	// Modify token and return to pool
	token1.Type = csslexer.IdentToken
	token1.Value = "test"
	token1.Raw = []rune("test")
	tokenPool.Put(token1)

	// Get another token - if it's the same instance from pool, it retains previous values
	token2 := tokenPool.Get().(*csslexer.Token)
	if token2 == nil {
		t.Error("Expected token from pool, got nil")
	}

	// If this is a fresh token from pool New(), it will have default values
	// If it's a reused token, it will have the previous values
	// Both cases are valid behavior for sync.Pool
	if token2 == token1 {
		// Same instance was reused, should have modified values
		if token2.Type != csslexer.IdentToken {
			t.Errorf("Reused token: expected IdentToken, got %v", token2.Type)
		}
		if token2.Value != "test" {
			t.Errorf("Reused token: expected 'test', got %q", token2.Value)
		}
	} else {
		// New instance, should have default values
		if token2.Type != csslexer.DefaultToken {
			t.Errorf("New token: expected DefaultToken, got %v", token2.Type)
		}
		if token2.Value != "" {
			t.Errorf("New token: expected empty string, got %q", token2.Value)
		}
		if len(token2.Raw) != 0 {
			t.Errorf("New token: expected empty Raw slice, got length %d", len(token2.Raw))
		}
	}

	tokenPool.Put(token2)
}

func TestTokenPoolConcurrency(t *testing.T) {
	const numTokens = 100
	tokens := make([]*csslexer.Token, numTokens)

	// Get multiple tokens
	for i := 0; i < numTokens; i++ {
		tokens[i] = tokenPool.Get().(*csslexer.Token)
		if tokens[i] == nil {
			t.Errorf("Expected token at index %d, got nil", i)
		}
	}

	// Return all tokens
	for i := 0; i < numTokens; i++ {
		tokenPool.Put(tokens[i])
	}
}
