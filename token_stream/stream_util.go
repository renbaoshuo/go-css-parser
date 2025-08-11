package token_stream

import (
	"go.baoshuo.dev/csslexer"
)

// ConsumeWhitespace consumes all upcoming whitespace tokens
// until a non-whitespace token is encountered.
func (ts *TokenStream) ConsumeWhitespace() {
	// Consume whitespace tokens until we find a non-whitespace token.
	for {
		t := ts.Peek()
		if t.Type != csslexer.WhitespaceToken {
			break
		}
		ts.Consume()
	}
}

// ConsumeIncludingWhitespace consumes the next token and also consumes
// any whitespace that follows it, returning the token that was consumed.
// This is useful when you want to process a token and ensure that any
// whitespace after it is also consumed, so that the next call to Peek()
// will return the next non-whitespace token.
func (ts *TokenStream) ConsumeIncludingWhitespace() csslexer.Token {
	token := ts.Consume()
	ts.ConsumeWhitespace()
	return token
}

// Skip tokens until one of these is true:
//
//   - EOF is reached.
//   - The next token would signal a premature end of the current block
//     (an unbalanced } or similar).
//   - The next token is of any of the given types, except if it occurs
//     within a block.
//
// The tokens that we consume are discarded. So e.g., if we ask
// to stop at semicolons, and the rest of the input looks like
// “.foo { color; } bar ; baz”, we would skip “.foo { color; } bar ”
// and stop there.
func (ts *TokenStream) SkipUntil(types ...csslexer.TokenType) {
	stopTypes := make(map[csslexer.TokenType]bool, len(types))
	for _, t := range types {
		stopTypes[t] = true
	}

	nestingStack := make([]csslexer.TokenType, 0)
	for {
		t := ts.Peek()

		if t.Type == csslexer.EOFToken {
			return // End of file, stop here.
		}

		if len(nestingStack) == 0 && stopTypes[t.Type] {
			return // Stop at the specified token type.
		}

		ts.Consume()

		if isBlockStartToken(t.Type) {
			nestingStack = append(nestingStack, getMatchingBlockEndToken(t.Type)) // Push the expected end token onto the stack.
		} else if isBlockEndToken(t.Type) {
			if len(nestingStack) == 0 {
				// Unmatched block end token, we should stop here.
				return
			}

			if t.Type == nestingStack[len(nestingStack)-1] {
				nestingStack = nestingStack[:len(nestingStack)-1] // Pop the stack.
			} else {
				// Unmatched end token, ignore it and continue.
				continue
			}
		}
	}
}
