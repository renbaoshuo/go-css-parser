package token_stream

import (
	"maps"

	"go.baoshuo.dev/csslexer"
)

type TokenStreamState struct {
	tokenStream *TokenStream // Reference to the TokenStream for context.

	inputState  csslexer.InputState         // The state of the input.
	peekedToken *csslexer.Token             // The token that was peeked.
	boundaries  map[csslexer.TokenType]bool // Boundary tokens.
}

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
			Type: ts.p.Type,
			Data: ts.p.Data,
		}
	}

	return state
}

func (tss TokenStreamState) Restore() {
	// Restore the input state.
	tss.tokenStream.z.RestoreState(tss.inputState)

	// Restore the peeked token if it exists.
	if tss.peekedToken != nil {
		tss.tokenStream.p = &csslexer.Token{
			Type: tss.peekedToken.Type,
			Data: tss.peekedToken.Data,
		}
	} else {
		tss.tokenStream.p = nil
	}

	// Restore the boundaries.
	tss.tokenStream.b = maps.Clone(tss.boundaries)
}
