package gasometer

import (
	"math"

	"github.com/c0mm4nd/slashvm/configs"
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/c0mm4nd/slashvm/utils"
	"github.com/holiman/uint256"
)

func call_extra_check(gas *uint256.Int, after_gas uint64, config *configs.Config) *exits.Error {
	if config.ErrOnCallWithMoreGas && uint256.NewInt(after_gas).Cmp(gas) < 0 {
		return exits.OutOfGas
	}

	return nil
}

func suicide_refund(already_removed bool) int64 {
	if already_removed {
		return 0
	} else {
		return int64(R_SUICIDE)
	}
}

func sstore_refund(original, current, new *uint256.Int, config *configs.Config) int64 {
	if config.SStoreGasMetering {
		if current.Eq(new) {
			return 0
		} else {
			if original.Eq(current) && new.IsZero() {
				return int64(config.RefundSStoreClears)
			} else {
				refund := int64(0)

				if !original.IsZero() {
					if current.IsZero() {
						refund -= int64(config.RefundSStoreClears)
					} else if new.IsZero() {
						refund += int64(config.RefundSStoreClears)
					}
				}

				if original.Eq(new) {
					if original.IsZero() {
						refund += int64(config.GasSStoreSet - config.GasSLoad)
					} else {
						refund += int64(config.GasSStoreReset - config.GasSLoad)
					}
				}

				return refund
			}
		}
	} else {
		if !current.IsZero() && new.IsZero() {
			return int64(config.RefundSStoreClears)
		} else {
			return 0
		}
	}
}

func create2_cost(len *uint256.Int) (uint64, *exits.Error) {
	base := uint256.NewInt(G_CREATE)
	// ceil(len / 32.0)
	sha_addup_base := new(uint256.Int).Div(len, uint256.NewInt(32))
	if !new(uint256.Int).Mod(len, uint256.NewInt(32)).IsZero() {
		sha_addup_base.AddUint64(sha_addup_base, 1)
	}
	var overflow bool
	sha_addup := uint256.NewInt(G_SHA3WORD)
	sha_addup, overflow = sha_addup.MulOverflow(sha_addup_base, sha_addup)
	if overflow {
		return 0, exits.OutOfGas
	}

	gas, overflow := base.AddOverflow(base, sha_addup)
	if overflow {
		return 0, exits.OutOfGas
	}

	if gas.Cmp(uint256.NewInt(math.MaxUint64)) > 0 {
		return 0, exits.OutOfGas
	}

	return gas.Uint64(), nil
}

func exp_cost(power *uint256.Int, config *configs.Config) (uint64, *exits.Error) {
	if power.IsZero() {
		return G_EXP, nil
	} else {
		var overflow bool
		gas := uint256.NewInt(config.GasExpByte)
		gas, overflow = gas.MulOverflow(gas, uint256.NewInt(utils.Log2Floor(power)/8+1))
		if overflow {
			return 0, exits.OutOfGas
		}

		gas, overflow = gas.AddOverflow(gas, uint256.NewInt(G_EXP))
		if overflow {
			return 0, exits.OutOfGas
		}

		if gas.Cmp(uint256.NewInt(math.MaxUint64)) > 0 {
			return 0, exits.OutOfGas
		}
		return gas.Uint64(), nil
	}
}
func verylowcopy_cost(len *uint256.Int) (uint64, *exits.Error) {
	wordd := new(uint256.Int).Div(len, uint256.NewInt(32))
	wordr := new(uint256.Int).Mod(len, uint256.NewInt(32))
	if !wordr.IsZero() {
		wordd.AddUint64(wordd, 1)
	}

	var overflow bool
	gas := uint256.NewInt(G_COPY)
	gas, overflow = gas.MulOverflow(gas, wordd)
	if overflow {
		return 0, exits.OutOfGas
	}
	gas, overflow = gas.AddOverflow(gas, uint256.NewInt(G_COPY))
	if overflow {
		return 0, exits.OutOfGas
	}

	if gas.Cmp(uint256.NewInt(math.MaxUint64)) > 0 {
		return 0, exits.OutOfGas
	}
	return gas.Uint64(), nil
}
func extcodecopy_cost(len *uint256.Int, is_cold bool, config *configs.Config) (uint64, *exits.Error) {
	wordd := new(uint256.Int).Div(len, uint256.NewInt(32))
	wordr := new(uint256.Int).Mod(len, uint256.NewInt(32))
	if !wordr.IsZero() {
		wordd.AddUint64(wordd, 1)
	}
	var overflow bool
	gas := uint256.NewInt(G_COPY)
	gas, overflow = gas.MulOverflow(gas, wordd)
	if overflow {
		return 0, exits.OutOfGas
	}
	gas, overflow = gas.AddOverflow(gas, uint256.NewInt(address_access_cost(is_cold, config.GasExtCode, config)))
	if overflow {
		return 0, exits.OutOfGas
	}

	if gas.Cmp(uint256.NewInt(math.MaxUint64)) > 0 {
		return 0, exits.OutOfGas
	}
	return gas.Uint64(), nil
}

func log_cost(n uint8, len *uint256.Int) (uint64, *exits.Error) {
	var overflow bool
	gas := uint256.NewInt(G_LOGDATA)
	gas, overflow = gas.MulOverflow(gas, len)
	if overflow {
		return 0, exits.OutOfGas
	}
	gas, overflow = gas.AddOverflow(gas, uint256.NewInt(G_LOG))
	if overflow {
		return 0, exits.OutOfGas
	}
	gas, overflow = gas.AddOverflow(gas, uint256.NewInt(G_LOGTOPIC*uint64(n)))
	if overflow {
		return 0, exits.OutOfGas
	}

	if gas.Cmp(uint256.NewInt(math.MaxUint64)) > 0 {
		return 0, exits.OutOfGas
	}
	return gas.Uint64(), nil
}

func sha3_cost(len *uint256.Int) (uint64, *exits.Error) {
	wordd := new(uint256.Int).Div(len, uint256.NewInt(32))
	wordr := new(uint256.Int).Mod(len, uint256.NewInt(32))

	var overflow bool
	gas := uint256.NewInt(G_SHA3WORD)
	if !wordr.IsZero() {
		wordd.AddUint64(wordd, 1)
	}
	gas, overflow = gas.MulOverflow(gas, wordd)
	if overflow {
		return 0, exits.OutOfGas
	}
	gas, overflow = gas.AddOverflow(gas, uint256.NewInt(G_SHA3))
	if overflow {
		return 0, exits.OutOfGas
	}

	if gas.Cmp(uint256.NewInt(math.MaxUint64)) > 0 {
		return 0, exits.OutOfGas
	}

	return gas.Uint64(), nil
}

func sload_cost(is_cold bool, config *configs.Config) uint64 {
	if config.IncreaseStateAccessGas {
		if is_cold {
			return config.GasSLoadCold
		} else {
			return config.GasStorageReadWarm
		}
	} else {
		return config.GasSLoad
	}
}

func sstore_cost(
	original, current, new *uint256.Int,
	gas uint64,
	is_cold bool,
	config *configs.Config,
) (uint64, *exits.Error) {
	gas_cost := config.GasSStoreSet
	if !config.Estimate {
		if config.SStoreGasMetering {
			if config.SStoreRevertUnderStipend && gas <= config.CallStipend {
				return 0, exits.OutOfGas
			}

			if new == current {
				gas_cost = config.GasSLoad
			} else {
				if original == current {
					if original.IsZero() {
						gas_cost = config.GasSStoreSet
					} else {
						gas_cost = config.GasSStoreReset
					}
				} else {
					gas_cost = config.GasSLoad
				}
			}
		} else {
			if current.IsZero() && !new.IsZero() {
				gas_cost = config.GasSStoreSet
			} else {
				gas_cost = config.GasSStoreReset
			}
		}
	}
	// In EIP-2929 we charge extra if the slot has not been used yet in this transaction
	if is_cold {
		return gas_cost + config.GasSLoadCold, nil
	}
	return gas_cost, nil
}

func suicide_cost(value *uint256.Int, is_cold, target_exists bool, config *configs.Config) uint64 {
	eip161 := !config.EmptyConsideredExists
	should_charge_topup := !target_exists
	if eip161 {
		should_charge_topup = !value.IsZero() && !target_exists
	}

	suicide_gas_topup := uint64(0)
	if should_charge_topup {
		suicide_gas_topup = config.GasSuicideNewAccount
	}

	gas := config.GasSuicide + suicide_gas_topup
	if config.IncreaseStateAccessGas && is_cold {
		gas += config.GasAccountAccessCold
	}
	return gas
}

func call_cost(value *uint256.Int, is_cold bool,
	is_call_or_callcode bool,
	is_call_or_staticcall bool,
	new_account bool, config *configs.Config) uint64 {
	transfers_value := !value.IsZero()
	return address_access_cost(is_cold, config.GasCall, config) +
		xfer_cost(is_call_or_callcode, transfers_value) +
		new_cost(is_call_or_staticcall, new_account, transfers_value, config)
}
func address_access_cost(is_cold bool, regular_value uint64, config *configs.Config) uint64 {
	if config.IncreaseStateAccessGas {
		if is_cold {
			return config.GasAccountAccessCold
		} else {
			return config.GasStorageReadWarm
		}
	} else {
		return regular_value
	}
}
func xfer_cost(is_call_or_callcode bool, transfers_value bool) uint64 {
	if is_call_or_callcode && transfers_value {
		return G_CALLVALUE
	} else {
		return 0
	}
}
func new_cost(is_call_or_staticcall, new_account, transfers_value bool, config *configs.Config) uint64 {
	eip161 := !config.EmptyConsideredExists
	if is_call_or_staticcall {
		if eip161 {
			if transfers_value && new_account {
				return G_NEWACCOUNT
			} else {
				return 0
			}
		} else if new_account {
			return G_NEWACCOUNT
		} else {
			return 0
		}
	} else {
		return 0
	}
}
