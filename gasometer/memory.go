package gasometer

import (
	"math"

	"github.com/c0mm4nd/slashvm/core/exits"
)

func memoryGas(a uint64) (*uint64, *exits.Error) {
	result := G_MEMORY
	if result > math.MaxUint64/a {
		return nil, exits.OutOfGas
	}
	result = a * result

	if a > math.MaxUint64/a {
		return nil, exits.OutOfGas
	}
	a = a * a

	if result > math.MaxUint64-a {
		return nil, exits.OutOfGas
	}
	result = result + a

	return &result, nil
}
