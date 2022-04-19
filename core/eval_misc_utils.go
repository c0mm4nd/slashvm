package core

import (
	"github.com/c0mm4nd/slashvm/core/opcodes"
	"github.com/holiman/uint256"
)

func push(state *Machine, n, pos int) Control {
	end := min(pos+1+n, len(state.Code))
	slice := state.Code[pos+1 : end]
	val := make([]byte, 32)
	copy(val[32-len(slice):32], opcodes.ToBytes(slice))

	state.Stack.Push(new(uint256.Int).SetBytes32(val))
	return &Continue{1 + n}
}
func dup(state *Machine, n int) Control {
	value, err := state.Stack.PeekN(n - 1)
	if err != nil {
		return &Exit{err}
	}
	err = state.Stack.Push(new(uint256.Int).Set(value))
	if err != nil {
		return &Exit{err}
	}
	return &Continue{1}
}
func swap(state *Machine, n int) Control {
	val1, err := state.Stack.PeekN(0)
	if err != nil {
		return &Exit{err}
	}
	val2, err := state.Stack.PeekN(n)
	if err != nil {
		return &Exit{err}
	}
	err = state.Stack.Set(0, val2)
	if err != nil {
		return &Exit{err}
	}
	state.Stack.Set(n, val1)
	if err != nil {
		return &Exit{err}
	}
	return &Continue{1}
}
