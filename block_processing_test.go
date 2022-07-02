package block_processing_test

import (
	bp "block_processing"
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
    rrp := bp.RangeResponseProcessor{}
	rrp.Initialize(4, 3)

	for count := 0; count < 3; count++ {
		blocks := []bp.Block{
			bp.Block(fmt.Sprintf("foo_%d", count)),
			bp.Block(fmt.Sprintf("bar_%d", count)),
			bp.Block(fmt.Sprintf("baz_%d", count)),
		}
		for i := 1; i < 10; i++ {
			rrp.ProcessRange(uint64(i), blocks)
		}
	}
	
	lo, hi := rrp.GetActiveRange()
	fmt.Printf("lo=%d, hi=%d", lo, hi)
}
