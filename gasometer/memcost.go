package gasometer

import (
	"github.com/c0mm4nd/slashvm/utils"
	"github.com/holiman/uint256"
)

type MemCost struct {
	Offset *uint256.Int
	Length *uint256.Int
}

func (memCost *MemCost) Join(other *MemCost) *MemCost {
	if other.Length.IsZero() {
		return memCost
	}

	if memCost.Length.IsZero() {
		return other
	}
	//saturating_add
	self_end := utils.SaturatingAdd(memCost.Offset, memCost.Length)
	other_end := utils.SaturatingAdd(other.Offset, other.Length)

	if self_end.Cmp(other_end) < 0 {
		return other
	}

	return memCost
}
