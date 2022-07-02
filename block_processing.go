package block_processing

import (
	"sync"
)

type Block string

// Assumption: The genesis block (block 0) is considered fulfilled, and so is
// given n responses on initialization.

type RangeResponseProcessor struct {
	s uint64
	n uint64
	h uint64
	mu sync.Mutex
	responses map[uint64]uint64
}

func (rrp *RangeResponseProcessor) Initialize(s uint64, n uint64) {
	rrp.s = s
	rrp.n = n
	rrp.h = 1
	rrp.responses = make(map[uint64]uint64)
	rrp.responses[0] = n
}

func (rrp *RangeResponseProcessor) ProcessRange(startHeight uint64, blocks []Block) {
	lo, hi := rrp.GetActiveRange()
	blocklen := uint64(len(blocks))
	// Question: Is this the right place for the lock? If it's inside the if block that wraps
	// the responses read-update, it will potentially be called a lot more. Locks may be
	// expensive.
	rrp.mu.Lock()
	for h := startHeight; h < startHeight + blocklen; h++ {
		if lo <= h && h <= hi {
			rrp.responses[h]++
		}
	}
	rrp.mu.Unlock()
}

func (rrp *RangeResponseProcessor) GetActiveRange() (minHeight uint64, maxHeight uint64) {
	// h is the minimum height such that < n range responses containing h
	// have been received. Start at the last known h. rrp.h only increases, so there's
	// no need to revisit old heights. This also suggests that we could garbage collect
	// the 'responses' map and eliminate entries.

	h := uint64(rrp.h) // Start at the last known h
	for rrp.responses[h] >= rrp.n { // Skip heights that are fulfilled
		h++
	}

	if h > rrp.h { // Update h if it has increased
		rrp.mu.Lock()
		rrp.h = h
		rrp.mu.Unlock()
	}
	return h, h + rrp.s - uint64(1)
}
