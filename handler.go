package slashvm

import (
	"github.com/c0mm4nd/slashvm/core"
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/c0mm4nd/slashvm/core/opcodes"
	"github.com/holiman/uint256"
)

type Transfer struct {
	Source string // Address
	Target string // Address
	Value  *uint256.Int
}

// EVM context handler.
type Handler interface {
	Balance(addr *uint256.Int) *uint256.Int
	CodeSize(addr *uint256.Int) *uint256.Int
	CodeHash(addr *uint256.Int) *uint256.Int
	ChainID() *uint256.Int
	Origin() *uint256.Int
	GasPrice() *uint256.Int
	BlockBaseFeePerGas() *uint256.Int
	BlockHash(number uint64) *uint256.Int
	BlockCoinbase() *uint256.Int
	BlockTimestamp() *uint256.Int
	BlockNumber() *uint256.Int
	BlockDifficulty() *uint256.Int
	GasLimit() *uint256.Int
	SetStorage(address *uint256.Int, index *uint256.Int, value *uint256.Int) *exits.Error
	Storage(address *uint256.Int, index *uint256.Int) *uint256.Int
	GasLeft() *uint256.Int
	Create(address *uint256.Int, scheme *CreateSchemeLegacy, value *uint256.Int, code []byte) (exits.Reason, *uint256.Int, []byte, *CreateInterrupt)
	Create2(address *uint256.Int, scheme *CreateSchemeCreate2, value *uint256.Int, code []byte) (exits.Reason, *uint256.Int, []byte, *CreateInterrupt)
	Log(address *uint256.Int, topics []*uint256.Int, data []byte) *exits.Error
	MarkDelete(address *uint256.Int, target *uint256.Int) *exits.Error
	Code(address *uint256.Int) []byte
	Call(codeAddr *uint256.Int, transfer *Transfer, input []byte, gas uint64, context *Context) (exits.Reason, []byte, *CallInterrupt)
	CallCode(codeAddr *uint256.Int, transfer *Transfer, input []byte, gas uint64, context *Context) (exits.Reason, []byte, *CallInterrupt)
	DelegateCall(addr *uint256.Int, input []byte, gas uint64, context *Context) (exits.Reason, []byte, *CallInterrupt)
	StaticCode(addr *uint256.Int, input []byte, gas uint64, context *Context) (exits.Reason, []byte, *CallInterrupt)
	CallFeedback() *exits.Error

	PreValidate(*Context, opcodes.Opcode, *core.Stack) *exits.Error
	Other(opcodes.Opcode, *core.Machine) *exits.Error
}
