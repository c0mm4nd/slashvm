package core

import (
	"testing"

	"github.com/holiman/uint256"
)

func Test_next_multiple_of_32(t *testing.T) {
	for i := 0; i < 32; i++ {
		x := uint256.NewInt(uint64(i) * 32)
		if result, overflow := next_multiple_of_32(x); overflow || x.Cmp(result) != 0 {
			if overflow {
				panic("overflow")
			}
			t.Logf("x(%s) != result(%s)", x, result)
			t.Fail()
		}
	}
	for x := 0; x < 1024; x++ {
		if x%32 == 0 {
			continue
		}

		next_multiple := uint256.NewInt(uint64(x + 32 - (x % 32)))
		if result, overflow := next_multiple_of_32(uint256.NewInt(uint64(x))); overflow || !next_multiple.Eq(result) {
			if overflow {
				panic("overflow")
			}
			t.Logf("x(%s) != result(%s)", next_multiple, result)
			t.Fail()
		}
	}
	maxUint64 := ^uint64(0)
	maxUint256 := &uint256.Int{maxUint64, maxUint64, maxUint64, maxUint64}

	last_multiple_of_32 := new(uint256.Int).Not(uint256.NewInt(31))
	last_multiple_of_32.And(last_multiple_of_32, maxUint256)
	for i := uint64(0); i < 63; i++ {
		var x = new(uint256.Int).Sub(maxUint256, uint256.NewInt(i))
		if x.Cmp(last_multiple_of_32) > 0 {
			if _, overflow := next_multiple_of_32(x); !overflow {
				t.Log("not overflow")
				t.Fail()
			}
		} else {
			if result, _ := next_multiple_of_32(x); !result.Eq(last_multiple_of_32) {
				t.Logf("x(%s) != result(%s)", last_multiple_of_32, result)
				t.Fail()
			}
		}
	}
}
