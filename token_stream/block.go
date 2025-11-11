package token_stream

import (
	"fmt"
	"maps"

	"go.baoshuo.dev/csslexer"
)

// SetBoundary sets or unsets a token type as a boundary token.
// When the token stream encounters a boundary token, it may stop
// processing further tokens depending on the context.
func (ts *TokenStream) SetBoundary(tokenType csslexer.TokenType, enable bool) {
	if enable {
		ts.b[tokenType] = true
	} else {
		delete(ts.b, tokenType)
	}
}

// ConsumeBlockToEndWithRestoring consumes a block of tokens until the end token is reached,
// allowing the caller to process the tokens within the block using the provided blockConsumer
// function.
//
// If the blockConsumer function returns commit as false, the token stream state is restored
// to its initial state before the block was consumed. Note this behavior is not affected by
// whether an error is returned or not.
//
// The errors returned by blockConsumer are wrapped with additional context.
func (ts *TokenStream) ConsumeBlockToEndWithRestoring(
	endTokenType csslexer.TokenType,
	blockConsumer func(ts *TokenStream) (commit bool, err error),
) error {
	initialState := ts.State()
	ob := maps.Clone(ts.b) // Clone the current boundaries to restore later if needed
	ts.b = make(map[csslexer.TokenType]bool)
	defer func() { ts.b = ob }() // Always restore boundaries

	ts.Consume()                       // Consume the block start token
	ts.SetBoundary(endTokenType, true) // Set the boundary for the end token

	commit, err := blockConsumer(ts)
	if err != nil {
		if !commit {
			initialState.Restore() // Restore the initial state if not committing
		}
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
	return nil
}

// ConsumeBlockRestoring consumes a block of tokens, allowing the caller to process
// the tokens within the block using the provided blockConsumer function.
//
// If the blockConsumer function returns commit as false, the token stream state is restored
// to its initial state before the block was consumed. Note this behavior is not affected by
// whether an error is returned or not.
//
// The errors returned by blockConsumer are wrapped with additional context.
func (ts *TokenStream) ConsumeBlockRestoring(blockConsumer func(ts *TokenStream) (commit bool, err error)) error {
	startToken := ts.Peek()

	if !isBlockStartToken(startToken.Type) {
		return fmt.Errorf("expected a block start token, got %s", startToken.Type)
	}

	endTokenType := getMatchingBlockEndToken(startToken.Type)

	return ts.ConsumeBlockToEndWithRestoring(endTokenType, blockConsumer)
}

// ConsumeBlock consumes a block of tokens, allowing the caller to process
// the tokens within the block using the provided blockConsumer function.
//
// If an error occurs during processing, the token stream will skip to the end of the block
// and the error will be returned with additional context.
func (ts *TokenStream) ConsumeBlock(blockConsumer func(ts *TokenStream) error) error {
	return ts.ConsumeBlockRestoring(func(ts *TokenStream) (commit bool, err error) {
		return true, blockConsumer(ts)
	})
}

// ConsumeBlockToEnd consumes a block of tokens until the specified end token
// type is reached, allowing the caller to process the tokens within the block
// using the provided blockConsumer function.
//
// If an error occurs during processing, the token stream state is restored to
// its initial state before the block was consumed, and the error is returned
// with additional context.
func (ts *TokenStream) ConsumeBlockToEnd(
	endTokenType csslexer.TokenType,
	blockConsumer func(ts *TokenStream) error,
) error {
	return ts.ConsumeBlockToEndWithRestoring(endTokenType, func(ts *TokenStream) (commit bool, err error) {
		return true, blockConsumer(ts)
	})
}
