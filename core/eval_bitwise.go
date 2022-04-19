package core

import "github.com/c0mm4nd/slashvm/core/opcodes"

func evalAnd(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.And(op1, op2)

	return &Continue{1}
}
func evalOr(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.Or(op1, op2)

	return &Continue{1}
}
func evalXor(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.Xor(op1, op2)

	return &Continue{1}
}
func evalNot(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op1.Not(op1)

	return &Continue{1}
}
func evalByte(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op1.Byte(op1)

	return &Continue{1}
}
func evalShl(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}
	if op1.LtUint64(256) {
		op2.Lsh(op2, uint(op1.Uint64()))
	} else {
		op2.Clear()
	}
	return &Continue{1}
}
func evalShr(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}
	if op1.LtUint64(256) {
		op2.Rsh(op2, uint(op1.Uint64()))
	} else {
		op2.Clear()
	}
	return &Continue{1}
}
func evalSar(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	if op1.GtUint64(256) {
		if op2.Sign() >= 0 {
			op2.Clear()
		} else {
			// Max negative shift: all bits set
			op2.SetAllOne()
		}
		return nil
	}

	op2.SRsh(op2, uint(op1.Uint64()))
	return &Continue{1}
}
