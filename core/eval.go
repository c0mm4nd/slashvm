package core

import (
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/c0mm4nd/slashvm/core/opcodes"
)

func evalStop(_ *Machine, _ opcodes.Opcode, _ int) Control {
	return &Exit{exits.Stopped}
}

var evalTable = map[opcodes.Opcode]func(*Machine, opcodes.Opcode, int) Control{
	opcodes.STOP:         evalStop,
	opcodes.ADD:          evalAdd,
	opcodes.MUL:          evalMul,
	opcodes.SUB:          evalSub,
	opcodes.DIV:          evalDiv,
	opcodes.SDIV:         evalSDiv,
	opcodes.MOD:          evalMod,
	opcodes.SMOD:         evalSMod,
	opcodes.ADDMOD:       evalAddMod,
	opcodes.MULMOD:       evalMulMod,
	opcodes.EXP:          evalExp,
	opcodes.SIGNEXTEND:   evalSignExtend,
	opcodes.LT:           evalLt,
	opcodes.GT:           evalGt,
	opcodes.SLT:          evalSLt,
	opcodes.SGT:          evalSGt,
	opcodes.EQ:           evalEq,
	opcodes.ISZERO:       evalIsZero,
	opcodes.AND:          evalAnd,
	opcodes.OR:           evalOr,
	opcodes.XOR:          evalXor,
	opcodes.NOT:          evalNot,
	opcodes.BYTE:         evalByte,
	opcodes.SHL:          evalShl,
	opcodes.SHR:          evalShr,
	opcodes.SAR:          evalSar,
	opcodes.CODESIZE:     evalCodeSize,
	opcodes.CODECOPY:     evalCodeCopy,
	opcodes.CALLDATALOAD: evalCallDataLoad,
	opcodes.CALLDATASIZE: evalCallDataSize,
	opcodes.CALLDATACOPY: evalCallDataCopy,
	opcodes.POP:          evalPop,
	opcodes.MLOAD:        evalMLoad,
	opcodes.MSTORE:       evalMStore,
	opcodes.MSTORE8:      evalMStore8,
	opcodes.JUMP:         evalJump,
	opcodes.JUMPI:        evalJumpI,
	opcodes.PC:           evalPC,
	opcodes.MSIZE:        evalMSize,
	opcodes.JUMPDEST:     evalJumpDest,

	opcodes.PUSH1:  evalPush1,
	opcodes.PUSH2:  evalPush2,
	opcodes.PUSH3:  evalPush3,
	opcodes.PUSH4:  evalPush4,
	opcodes.PUSH5:  evalPush5,
	opcodes.PUSH6:  evalPush6,
	opcodes.PUSH7:  evalPush7,
	opcodes.PUSH8:  evalPush8,
	opcodes.PUSH9:  evalPush9,
	opcodes.PUSH10: evalPush10,
	opcodes.PUSH11: evalPush11,
	opcodes.PUSH12: evalPush12,
	opcodes.PUSH13: evalPush13,
	opcodes.PUSH14: evalPush14,
	opcodes.PUSH15: evalPush15,
	opcodes.PUSH16: evalPush16,
	opcodes.PUSH17: evalPush17,
	opcodes.PUSH18: evalPush18,
	opcodes.PUSH19: evalPush19,
	opcodes.PUSH20: evalPush20,
	opcodes.PUSH21: evalPush21,
	opcodes.PUSH22: evalPush22,
	opcodes.PUSH23: evalPush23,
	opcodes.PUSH24: evalPush24,
	opcodes.PUSH25: evalPush25,
	opcodes.PUSH26: evalPush26,
	opcodes.PUSH27: evalPush27,
	opcodes.PUSH28: evalPush28,
	opcodes.PUSH29: evalPush29,
	opcodes.PUSH30: evalPush30,
	opcodes.PUSH31: evalPush31,
	opcodes.PUSH32: evalPush32,

	opcodes.DUP1:  evalDup1,
	opcodes.DUP2:  evalDup2,
	opcodes.DUP3:  evalDup3,
	opcodes.DUP4:  evalDup4,
	opcodes.DUP5:  evalDup5,
	opcodes.DUP6:  evalDup6,
	opcodes.DUP7:  evalDup7,
	opcodes.DUP8:  evalDup8,
	opcodes.DUP9:  evalDup9,
	opcodes.DUP10: evalDup10,
	opcodes.DUP11: evalDup11,
	opcodes.DUP12: evalDup12,
	opcodes.DUP13: evalDup13,
	opcodes.DUP14: evalDup14,
	opcodes.DUP15: evalDup15,
	opcodes.DUP16: evalDup16,

	opcodes.SWAP1:  evalSwap1,
	opcodes.SWAP2:  evalSwap2,
	opcodes.SWAP3:  evalSwap3,
	opcodes.SWAP4:  evalSwap4,
	opcodes.SWAP5:  evalSwap5,
	opcodes.SWAP6:  evalSwap6,
	opcodes.SWAP7:  evalSwap7,
	opcodes.SWAP8:  evalSwap8,
	opcodes.SWAP9:  evalSwap9,
	opcodes.SWAP10: evalSwap10,
	opcodes.SWAP11: evalSwap11,
	opcodes.SWAP12: evalSwap12,
	opcodes.SWAP13: evalSwap13,
	opcodes.SWAP14: evalSwap14,
	opcodes.SWAP15: evalSwap15,
	opcodes.SWAP16: evalSwap16,

	opcodes.RETURN:  evalReturn,
	opcodes.REVERT:  evalRevert,
	opcodes.INVALID: evalInvalid,
}

func Eval(state *Machine, opcode opcodes.Opcode, pos int) Control {
	evalMethod := evalTable[opcode]
	if evalMethod == nil {
		evalMethod = evalExternel
	}

	return evalMethod(state, opcode, pos)
}
