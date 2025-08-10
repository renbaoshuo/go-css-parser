package token_stream

import (
	"sync"

	"go.baoshuo.dev/csslexer"
)

// tokenPool is a sync.Pool for reusing token instances.
var tokenPool = sync.Pool{
	New: func() interface{} {
		return &csslexer.Token{
			Type: csslexer.DefaultToken,
			Data: nil,
		}
	},
}
