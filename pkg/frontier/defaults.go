package frontier

import "context"

type StandardBufferPolicy struct{}

func (b StandardBufferPolicy) ShouldPrefetch(ctx context.Context, state PrefetchState) PrefetchDecision {
	if state.Size < state.Capacity%2 {
		return PrefetchDecision{
			ShouldFetch: true,
		}
	}

	return PrefetchDecision{
		ShouldFetch: false,
	}
}

// Flush on every insertion
func (b StandardBufferPolicy) ShouldFlush(ctx context.Context, state FlushState) FlushDecision {
	return FlushDecision{
		ShouldFlush: true,
	}
}
