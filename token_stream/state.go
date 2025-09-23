package token_stream

import (
	"maps"

	"go.baoshuo.dev/csslexer"
)

// TokenStreamState captures the state of a TokenStream at a specific point in time.
// It can be used to restore the TokenStream to this state later.
type TokenStreamState struct {
	tokenStream *TokenStream // Reference to the TokenStream for context.

	inputState  csslexer.InputState         // The state of the input.
	peekedToken *csslexer.Token             // The token that was peeked.
	boundaries  map[csslexer.TokenType]bool // Boundary tokens.
}

// State captures the current state of the TokenStream and returns it as a TokenStreamState.
func (ts *TokenStream) State() TokenStreamState {
	state := TokenStreamState{
		// Reference to the TokenStream for context.
		tokenStream: ts,

		// Capture the current input state.
		inputState:  ts.z.State(),
		peekedToken: nil,
		boundaries:  maps.Clone(ts.b),
	}

	if ts.p != nil {
		state.peekedToken = &csslexer.Token{
			Type:  ts.p.Type,
			Value: ts.p.Value,
			Raw:   ts.p.Raw,
		}
	}

	return state
}

// Restore restores the TokenStream to the state captured in the TokenStreamState.
func (tss TokenStreamState) Restore() {
	// Restore the input state.
	tss.inputState.Restore()

	// Restore the peeked token if it exists.
	if tss.peekedToken != nil {
		if tss.tokenStream.p != nil {
			tokenPool.Put(tss.tokenStream.p) // Return the old token to the pool.
		}
		p := tokenPool.Get().(*csslexer.Token)
		p.Type, p.Value, p.Raw = tss.peekedToken.Type, tss.peekedToken.Value, tss.peekedToken.Raw
		tss.tokenStream.p = p
	} else {
		tss.tokenStream.p = nil
	}

	// Restore the boundaries.
	tss.tokenStream.b = maps.Clone(tss.boundaries)
}
