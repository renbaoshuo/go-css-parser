package token_stream

import (
	"fmt"
	"maps"

	"go.baoshuo.dev/csslexer"
)

func (ts *TokenStream) SetBoundary(tokenType csslexer.TokenType, enable bool) {
	if enable {
		ts.b[tokenType] = true
	} else {
		delete(ts.b, tokenType)
	}
}

func (ts *TokenStream) ConsumeBlockToEndWithRestoring(
	endTokenType csslexer.TokenType,
	blockConsumer func(ts *TokenStream) (commit bool, err error)) error {
	initialState := ts.State()
	ob := maps.Clone(ts.b) // Clone the current boundaries to restore later if needed
	ts.b = make(map[csslexer.TokenType]bool)

	ts.Consume()                       // Consume the block start token
	ts.SetBoundary(endTokenType, true) // Set the boundary for the end token

	commit, err := blockConsumer(ts)
	if err != nil {
		initialState.Restore() // Restore the initial state in case of error
		return fmt.Errorf("error while consuming block: %w", err)
	}

	if !commit {
		initialState.Restore() // Restore the initial state if not committing
		return nil
	}

	ts.SkipUntil(endTokenType)
	endToken := ts.Peek()
	if endToken.Type != endTokenType {
		initialState.Restore() // Restore the initial state if the end token does not match
		return fmt.Errorf("expected end token %s, got %s", endTokenType, endToken.Type)
	}

	ts.Consume() // Consume the end token
	ts.b = ob    // Restore the boundaries
	return nil
}

func (ts *TokenStream) ConsumeBlockRestoring(
	blockConsumer func(ts *TokenStream) (commit bool, err error)) error {
	startToken := ts.Peek()

	if !isBlockStartToken(startToken.Type) {
		return fmt.Errorf("expected a block start token, got %s", startToken.Type)
	}

	endTokenType := getMatchingBlockEndToken(startToken.Type)

	return ts.ConsumeBlockToEndWithRestoring(endTokenType, blockConsumer)
}

func (ts *TokenStream) ConsumeBlock(blockConsumer func(ts *TokenStream) error) error {
	return ts.ConsumeBlockRestoring(func(ts *TokenStream) (commit bool, err error) {
		err = blockConsumer(ts)
		if err != nil {
			return false, err // If an error occurs, do not commit the block
		}
		return true, nil // Commit the block if no error
	})
}

func (ts *TokenStream) ConsumeBlockToEnd(
	endTokenType csslexer.TokenType,
	blockConsumer func(ts *TokenStream) error) error {
	return ts.ConsumeBlockToEndWithRestoring(endTokenType, func(ts *TokenStream) (commit bool, err error) {
		err = blockConsumer(ts)
		if err != nil {
			return false, err // If an error occurs, do not commit the block
		}
		return true, nil // Commit the block if no error
	})
}
