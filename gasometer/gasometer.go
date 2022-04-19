package gasometer

import (
	"github.com/c0mm4nd/slashvm/configs"
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/c0mm4nd/slashvm/gasometer/gascosts"
)

type Gasometer struct {
	GasLimit uint64
	Config   *configs.Config
	Inner    *Inner
	InnerErr *exits.Error // check err when Inner is nil
}

func NewGasometer(gasLimit uint64, config *configs.Config) *Gasometer {
	return &Gasometer{
		GasLimit: gasLimit,
		Config:   config,
		Inner: &Inner{
			memGas:      0,
			usedGas:     0,
			refundedGas: 0,
			config:      config,
		},
		InnerErr: nil,
	}
}

func (g *Gasometer) GasCost(cost gascosts.GasCost, gas uint64) (uint64, *exits.Error) {
	if g.InnerErr != nil {
		return 0, g.InnerErr
	}

	return g.Inner.GasCost(cost, gas)
}

// Remaining gas.
func (g *Gasometer) Gas() uint64 {
	if g.InnerErr != nil {
		return 0
	}

	return g.GasLimit - g.Inner.usedGas - g.Inner.memGas
}

func (g *Gasometer) TotalUsedGas() uint64 {
	if g.InnerErr != nil {
		return 0
	}

	return g.Inner.usedGas + g.Inner.memGas

}

func (g *Gasometer) RefundedGas() int64 {
	if g.InnerErr != nil {
		return 0
	}

	return g.Inner.refundedGas
}

// Explicitly fail the gasometer with out of gas. Return `OutOfGas` error.
func (g *Gasometer) Fail() *exits.Error {
	g.InnerErr = exits.OutOfGas
	return exits.OutOfGas
}

func (g *Gasometer) RecordCost(cost uint64) *exits.Error {
	// event!(RecordCost {
	// 	cost,
	// 	snapshot: self.snapshot(),
	// });

	all_gas_cost := g.TotalUsedGas() + cost
	if g.GasLimit < all_gas_cost {
		g.InnerErr = exits.OutOfGas
		return exits.OutOfGas
	}

	if g.Inner != nil {
		g.Inner.usedGas += cost
	}
	return nil
}

// Record an explicit refund.
func (g *Gasometer) RecordRefund(refund int64) *exits.Error {
	// event!(RecordRefund {
	// 	refund,
	// 	snapshot: self.snapshot(),
	// });
	if g.Inner != nil {
		g.Inner.refundedGas += refund
	}
	return nil
}

// Record `CREATE` code deposit.
func (g *Gasometer) RecordDeposit(length int64) *exits.Error {
	// event!(RecordRefund {
	// 	refund,
	// 	snapshot: self.snapshot(),
	// });
	cost := uint64(length) * G_CODEDEPOSIT
	return g.RecordCost(cost)
}

func (g *Gasometer) RecordDynamicCost(cost gascosts.GasCost, mem *MemCost) *exits.Error {
	return nil //TODO
}
