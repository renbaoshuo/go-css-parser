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

	nestingLevel := 0
	for {
		t := ts.Peek()

		if t.Type == csslexer.EOFToken {
			return // End of file, stop here.
		}

		if nestingLevel == 0 && stopTypes[t.Type] {
			return // Stop at the specified token type.
		}

		ts.Consume()

		// TODO: Handle unmatched block end tokens.
		// e.g., .foo { color: red); }
		//
		// Below is a simplified version, and may not handle all cases correctly.
		if isBlockStartToken(t.Type) {
			nestingLevel++ // Entering a block, increase nesting level.
		}
		if nestingLevel > 0 && isBlockEndToken(t.Type) {
			nestingLevel-- // Exiting a block, decrease nesting level.
		}
	}
}
