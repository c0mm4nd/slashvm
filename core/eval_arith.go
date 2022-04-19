package core

import "github.com/c0mm4nd/slashvm/core/opcodes"

// all have done inside "github.com/holiman/uint256"
func evalAdd(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.Add(op1, op2)
	return &Continue{1}
}

func evalMul(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.Mul(op1, op2)
	return &Continue{1}
}

func evalSub(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.Sub(op1, op2)

	return &Continue{1}
}
func evalDiv(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.Div(op1, op2)

	return &Continue{1}
}
func evalSDiv(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.SDiv(op1, op2)

	return &Continue{1}
}
func evalMod(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.Mod(op1, op2)

	return &Continue{1}
}
func evalSMod(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.SMod(op1, op2)

	return &Continue{1}
}
func evalAddMod(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op3, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op3.AddMod(op1, op2, op3)

	return &Continue{1}
}
func evalMulMod(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op3, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op3.MulMod(op1, op2, op3)

	return &Continue{1}
}
func evalExp(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.Exp(op1, op2)

	return &Continue{1}
}
func evalSignExtend(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	op2.ExtendSign(op1, op2)

	return &Continue{1}
}
func evalLt(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	if op1.Lt(op2) {
		op2.SetOne()
	} else {
		op2.Clear()
	}
	return &Continue{1}
}
func evalGt(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	if op1.Gt(op2) {
		op2.SetOne()
	} else {
		op2.Clear()
	}
	return &Continue{1}
}
func evalSLt(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	if op1.Slt(op2) {
		op2.SetOne()
	} else {
		op2.Clear()
	}
	return &Continue{1}
}
func evalSGt(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	if op1.Sgt(op2) {
		op2.SetOne()
	} else {
		op2.Clear()
	}
	return &Continue{1}
}
func evalEq(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	op2, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	if op1.Eq(op2) {
		op2.SetOne()
	} else {
		op2.Clear()
	}
	return &Continue{1}
}
func evalIsZero(state *Machine, _ opcodes.Opcode, _ int) Control {
	op1, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	if op1.IsZero() {
		op1.SetOne()
	} else {
		op1.Clear()
	}
	return &Continue{1}
}
