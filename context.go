package slashvm

import "github.com/holiman/uint256"

type CreateScheme interface {
	IsCreateScheme()
}

type CreateSchemeLegacy struct {
	Caller *uint256.Int
}

func (CreateSchemeLegacy) IsCreateScheme() {}

type CreateSchemeCreate2 struct {
	Caller   *uint256.Int
	CodeHash []byte
	Salt     []byte
}

func (CreateSchemeCreate2) IsCreateScheme() {}

type Context struct {
	Address       *uint256.Int // Execution address.
	Caller        *uint256.Int // Caller of the EVM.
	ApparentValue *uint256.Int // Apparent value of the EVM.
}
