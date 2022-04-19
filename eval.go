package slashvm

import (
	"github.com/c0mm4nd/slashvm/core/opcodes"
)

func handleOther(state *Runtime, opcode opcodes.Opcode, handler Handler) Control {
	err := handler.Other(opcode, state.Machine)
	if err != nil {
		return &Exit{err}
	}

	return &Continue{}
}

func Eval(r *Runtime, opcode opcodes.Opcode, h Handler) Control {
	switch opcode {
	case opcodes.SHA3:
		return evalKeccak256(r)
	case opcodes.ADDRESS:
		return evalAddress(r)
	case opcodes.BALANCE:
		return evalBalance(r, h)
	case opcodes.SELFBALANCE:
		return evalSelfBalance(r, h)
	case opcodes.ORIGIN:
		return evalOrigin(r, h)
	case opcodes.CALLER:
		return evalCaller(r, h)
	case opcodes.CALLVALUE:
		return evalCallValue(r, h)
	case opcodes.GASPRICE:
		return evalGasPrice(r, h)
	case opcodes.EXTCODESIZE:
		return evalExtCodeSize(r, h)
	case opcodes.EXTCODEHASH:
		return evalExtCodeHash(r, h)
	case opcodes.EXTCODECOPY:
		return evalExtCodeCopy(r, h)
	case opcodes.RETURNDATASIZE:
		return evalReturnDataSize(r, h)
	case opcodes.RETURNDATACOPY:
		return evalReturnDataCopy(r, h)
	case opcodes.BLOCKHASH:
		return evalBlockHash(r, h)
	case opcodes.COINBASE:
		return evalCoinbase(r, h)
	case opcodes.TIMESTAMP:
		return evalTimestamp(r, h)
	case opcodes.NUMBER:
		return evalNumber(r, h)
	case opcodes.DIFFICULTY:
		return evalDifficulty(r, h)
	case opcodes.GASLIMIT:
		return evalGasLimit(r, h)
	case opcodes.SLOAD:
		return evalSLoad(r, h)
	case opcodes.SSTORE:
		return evalSStore(r, h)
	case opcodes.GAS:
		return evalGas(r, h)
	case opcodes.LOG0:
		return evalLog(r, 0, h)
	case opcodes.LOG1:
		return evalLog(r, 1, h)
	case opcodes.LOG2:
		return evalLog(r, 2, h)
	case opcodes.LOG3:
		return evalLog(r, 3, h)
	case opcodes.LOG4:
		return evalLog(r, 4, h)
	case opcodes.SUICIDE:
		return evalSuicide(r, h)
	case opcodes.CREATE:
		return evalCreate(r, h)
	case opcodes.CREATE2:
		return evalCreate2(r, h)
	case opcodes.CALL:
		return evalCall(r, h)
	case opcodes.CALLCODE:
		return evalCallCode(r, h)
	case opcodes.DELEGATECALL:
		return evalDelegateCall(r, h)
	case opcodes.STATICCALL:
		return evalStaticCall(r, h)
	case opcodes.CHAINID:
		return evalChainID(r, h)
	case opcodes.BASEFEE:
		return evalBaseFee(r, h)
	default:
		return handleOther(r, opcode, h)
	}
}
