package gasometer

import (
	"fmt"
	"math"

	"github.com/c0mm4nd/slashvm/configs"
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/c0mm4nd/slashvm/gasometer/gascosts"
	"github.com/holiman/uint256"
)

type Inner struct {
	memGas      uint64
	usedGas     uint64
	refundedGas int64
	config      *configs.Config
}

func (i *Inner) MemGas(mem MemCost) (*uint64, *exits.Error) {
	from := mem.Offset
	length := mem.Length

	if mem.Length.IsZero() {
		return &i.memGas, nil
	}

	end, overflow := new(uint256.Int).AddOverflow(from, length)
	if overflow {
		return nil, exits.OutOfGas
	}

	if end.GtUint64(math.MaxUint64) {
		return nil, exits.OutOfGas
	}

	size := end.Uint64()
	if size%32 == 0 {
		size = size / 32
	} else {
		size = size/32 + 1
	}

	gas, err := memoryGas(size)
	if err != nil {
		return nil, err
	}
	if *gas > i.memGas {
		*gas = i.memGas
	}

	return gas, nil
}

func (i *Inner) ExtraCheck(cost gascosts.GasCost, afterGas uint64) *exits.Error {
	switch cost := cost.(type) {
	case *gascosts.GasCostCall:
		return call_extra_check(cost.Gas, afterGas, i.config)
	case *gascosts.GasCostCallCode:
		return call_extra_check(cost.Gas, afterGas, i.config)
	case *gascosts.GasCostDelegateCall:
		return call_extra_check(cost.Gas, afterGas, i.config)
	case *gascosts.GasCostStaticCall:
		return call_extra_check(cost.Gas, afterGas, i.config)
	default:
		return nil // Ok
	}
}

func (i *Inner) GasCost(cost gascosts.GasCost, gas uint64) (uint64, *exits.Error) {
	switch cost := cost.(type) {
	case *gascosts.GasCostCall:
		return call_cost(cost.Value, cost.TargetIsCold, true, true, !cost.TargetExists, i.config), nil
	case *gascosts.GasCostCallCode:
		return call_cost(cost.Value, cost.TargetIsCold, true, false, !cost.TargetExists, i.config), nil
	case *gascosts.GasCostDelegateCall:
		return call_cost(&uint256.Int{}, cost.TargetIsCold, false, false, !cost.TargetExists, i.config), nil
	case *gascosts.GasCostStaticCall:
		return call_cost(&uint256.Int{}, cost.TargetIsCold, false, true, !cost.TargetExists, i.config), nil
	case *gascosts.GasCostSuicide:
		return suicide_cost(cost.Value, cost.TargetIsCold, cost.TargetExists, i.config), nil
	case *gascosts.GasCostSStore:
		return sstore_cost(cost.Original, cost.Current, cost.New, gas, cost.TargetIsCold, i.config)
	case *gascosts.GasCostSHA3:
		return sha3_cost(cost.Length)
	case *gascosts.GasCostLog:
		return log_cost(cost.N, cost.Length)
	case *gascosts.GasCostVeryLowCopy:
		return verylowcopy_cost(cost.Length)
	case *gascosts.GasCostExp:
		return exp_cost(cost.Power, i.config)
	case *gascosts.GasCostCreate:
		return G_CREATE, nil
	case *gascosts.GasCostCreate2:
		return create2_cost(cost.Length)
	case *gascosts.GasCostSLoad:
		return sload_cost(cost.TargetIsCold, i.config), nil
	case *gascosts.GasCostZero:
		return G_ZERO, nil
	case *gascosts.GasCostBase:
		return G_BASE, nil
	case *gascosts.GasCostVeryLow:
		return G_VERYLOW, nil
	case *gascosts.GasCostLow:
		return G_LOW, nil
	case *gascosts.GasCostInvalid:
		return 0, exits.OutOfGas
	case *gascosts.GasCostExtCodeSize:
		return address_access_cost(cost.TargetIsCold, i.config.GasExtCode, i.config), nil
	case *gascosts.GasCostExtCodeCopy:
		return extcodecopy_cost(cost.Length, cost.TargetIsCold, i.config)
	case *gascosts.GasCostBalance:
		return address_access_cost(cost.TargetIsCold, i.config.GasBalance, i.config), nil
	case *gascosts.GasCostBlockHash:
		return G_BLOCKHASH, nil
	case *gascosts.GasCostExtCodeHash:
		return address_access_cost(cost.TargetIsCold, i.config.GasExtCodeHash, i.config), nil
	default:
		panic(fmt.Sprintf("unknown cost %#v", cost))
	}
}

func (i *Inner) GasRefund(cost gascosts.GasCost) int64 {
	if i.config.Estimate {
		return 0
	}

	switch cost := cost.(type) {
	case *gascosts.GasCostSStore:
		return sstore_refund(cost.Original, cost.Current, cost.New, i.config)
	case *gascosts.GasCostSuicide:
		if i.config.DecreaseClearsRefund {
			return 0
		}
		return suicide_refund(cost.AlreadyRemoved)
	default:
		panic(fmt.Sprintf("unknown refund %#v", cost))
	}
}
