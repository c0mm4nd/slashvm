package gascosts

import "github.com/holiman/uint256"

type GasCost interface {
	IsGasCost()
}
type GasCostZero struct{}

func (*GasCostZero) IsGasCost() {}

type GasCostBase struct{}

func (*GasCostBase) IsGasCost() {}

type GasCostVeryLow struct{}

func (*GasCostVeryLow) IsGasCost() {}

type GasCostLow struct{}

func (*GasCostLow) IsGasCost() {}

type GasCostInvalid struct{}

func (*GasCostInvalid) IsGasCost() {}

type GasCostExtCodeSize struct {
	TargetIsCold bool
}

func (*GasCostExtCodeSize) IsGasCost() {}

type GasCostBalance struct {
	TargetIsCold bool
}

func (*GasCostBalance) IsGasCost() {}

type GasCostBlockHash struct{}

func (*GasCostBlockHash) IsGasCost() {}

type GasCostExtCodeHash struct {
	TargetIsCold bool
}

func (*GasCostExtCodeHash) IsGasCost() {}

type GasCostCall struct {
	Value        *uint256.Int
	Gas          *uint256.Int
	TargetIsCold bool
	TargetExists bool
}

func (*GasCostCall) IsGasCost() {}

type GasCostCallCode struct {
	Value        *uint256.Int
	Gas          *uint256.Int
	TargetIsCold bool
	TargetExists bool
}

func (*GasCostCallCode) IsGasCost() {}

type GasCostDelegateCall struct {
	Gas          *uint256.Int
	TargetIsCold bool
	TargetExists bool
}

func (*GasCostDelegateCall) IsGasCost() {}

type GasCostStaticCall struct {
	Gas          *uint256.Int
	TargetIsCold bool
	TargetExists bool
}

func (*GasCostStaticCall) IsGasCost() {}

type GasCostSuicide struct {
	Value          *uint256.Int
	TargetIsCold   bool
	TargetExists   bool
	AlreadyRemoved bool
}

func (*GasCostSuicide) IsGasCost() {}

type GasCostSStore struct {
	Original     *uint256.Int
	Current      *uint256.Int
	New          *uint256.Int
	TargetIsCold bool
}

func (*GasCostSStore) IsGasCost() {}

type GasCostSHA3 struct {
	Length *uint256.Int
}

func (*GasCostSHA3) IsGasCost() {}

type GasCostLog struct {
	N      uint8
	Length *uint256.Int
}

func (*GasCostLog) IsGasCost() {}

type GasCostExtCodeCopy struct {
	TargetIsCold bool
	Length       *uint256.Int
}

func (*GasCostExtCodeCopy) IsGasCost() {}

type GasCostVeryLowCopy struct {
	Length *uint256.Int
}

func (*GasCostVeryLowCopy) IsGasCost() {}

type GasCostExp struct {
	Power *uint256.Int
}

func (*GasCostExp) IsGasCost() {}

type GasCostCreate struct{}

func (*GasCostCreate) IsGasCost() {}

type GasCostCreate2 struct {
	Length *uint256.Int
}

func (*GasCostCreate2) IsGasCost() {}

type GasCostSLoad struct {
	TargetIsCold bool
}

func (*GasCostSLoad) IsGasCost() {}
