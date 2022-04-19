package core

import (
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/c0mm4nd/slashvm/core/opcodes"
	"github.com/holiman/uint256"
)

func evalCodeSize(state *Machine, _ opcodes.Opcode, _ int) Control {
	state.Stack.Push(uint256.NewInt(uint64(len(state.Code))))
	return &Continue{1}
}
func evalCodeCopy(state *Machine, _ opcodes.Opcode, _ int) Control {
	memory_offset, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	code_offset, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	length, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	fatal := state.Mem.CopyLarge(
		memory_offset, code_offset,
		length, opcodes.ToBytes(state.Code))
	if fatal != nil {
		return &Exit{fatal}
	}

	return &Continue{1}
}

func evalCallDataLoad(state *Machine, _ opcodes.Opcode, _ int) Control {
	index, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	load := make([]byte, 32)
	var p = uint256.NewInt(0)
	for i := 0; i < 32; i++ {
		p.Add(index, uint256.NewInt(uint64(i)))
		if p.Lt(uint256.NewInt(uint64(^uint(0)))) {
			up := p.Uint64()
			if up < uint64(len(state.Data)) {
				load[i] = state.Data[up]
			}
		}
	}

	index.SetBytes32(load)
	return &Continue{1}
}
func evalCallDataSize(state *Machine, _ opcodes.Opcode, _ int) Control {
	state.Stack.Push(uint256.NewInt(uint64(len(state.Data))))
	return &Continue{1}
}
func evalCallDataCopy(state *Machine, _ opcodes.Opcode, _ int) Control {
	memory_offset, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	data_offset, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	length, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	state.Mem.ResizeOffset(memory_offset, length)
	if length.IsZero() {
		return &Continue{1}
	}

	fatal := state.Mem.CopyLarge(memory_offset, data_offset, length, state.Data)
	if fatal != nil {
		return &Exit{fatal}
	}

	return &Continue{1}
}
func evalPop(state *Machine, _ opcodes.Opcode, _ int) Control {
	_, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	return &Continue{1}
}
func evalMLoad(state *Machine, _ opcodes.Opcode, _ int) Control {
	index, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	err = state.Mem.ResizeOffset(index, uint256.NewInt(32))
	if err != nil {
		return &Exit{err}
	}

	i := index.Uint64()
	index.SetBytes32(state.Mem.Get(i, 32))

	return &Continue{1}
}
func evalMStore(state *Machine, _ opcodes.Opcode, _ int) Control {
	index, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	value, err := state.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	err = state.Mem.ResizeOffset(index, uint256.NewInt(32))
	if err != nil {
		return &Exit{err}
	}

	bytes := value.Bytes32()
	fatal := state.Mem.Set(int(index.Uint64()), bytes[:], 32)
	if fatal != nil {
		return &Exit{fatal}
	}

	return &Continue{1}
}
func evalMStore8(state *Machine, _ opcodes.Opcode, _ int) Control {
	index, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	value, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	err = state.Mem.ResizeOffset(index, uint256.NewInt(1))
	if err != nil {
		return &Exit{err}
	}

	b := byte(value.Uint64() & 0xff)
	fatal := state.Mem.Set(int(index.Uint64()), []byte{b}, 1)
	if fatal != nil {
		return &Exit{fatal}
	}

	return &Continue{1}
}
func evalJump(state *Machine, _ opcodes.Opcode, _ int) Control {
	dest, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	if state.Valids.IsValid(int(dest.Uint64())) {
		return &Jump{int(dest.Uint64())}
	}
	return &Exit{exits.InvalidJump}
}
func evalJumpI(state *Machine, _ opcodes.Opcode, _ int) Control {
	dest, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	value, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	if !value.IsZero() {
		d, of := dest.Uint64WithOverflow()
		if of {
			return &Exit{exits.InvalidJump}
		}
		if state.Valids.IsValid(int(d)) {
			return &Jump{int(d)}
		}

		return &Exit{exits.InvalidJump}
	}

	return &Continue{1}
}

func evalPC(state *Machine, _ opcodes.Opcode, pos int) Control {
	state.Stack.Push(uint256.NewInt(uint64(pos)))
	return &Continue{1}
}
func evalMSize(state *Machine, _ opcodes.Opcode, pos int) Control {
	state.Stack.Push(new(uint256.Int).Set(state.Mem.EffectiveLen))
	return &Continue{1}
}
func evalJumpDest(state *Machine, _ opcodes.Opcode, _ int) Control {
	return &Continue{1}
}

func evalPush1(state *Machine, _ opcodes.Opcode, pos int) Control  { return push(state, 1, pos) }
func evalPush2(state *Machine, _ opcodes.Opcode, pos int) Control  { return push(state, 2, pos) }
func evalPush3(state *Machine, _ opcodes.Opcode, pos int) Control  { return push(state, 3, pos) }
func evalPush4(state *Machine, _ opcodes.Opcode, pos int) Control  { return push(state, 4, pos) }
func evalPush5(state *Machine, _ opcodes.Opcode, pos int) Control  { return push(state, 5, pos) }
func evalPush6(state *Machine, _ opcodes.Opcode, pos int) Control  { return push(state, 6, pos) }
func evalPush7(state *Machine, _ opcodes.Opcode, pos int) Control  { return push(state, 7, pos) }
func evalPush8(state *Machine, _ opcodes.Opcode, pos int) Control  { return push(state, 8, pos) }
func evalPush9(state *Machine, _ opcodes.Opcode, pos int) Control  { return push(state, 9, pos) }
func evalPush10(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 10, pos) }
func evalPush11(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 11, pos) }
func evalPush12(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 12, pos) }
func evalPush13(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 13, pos) }
func evalPush14(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 14, pos) }
func evalPush15(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 15, pos) }
func evalPush16(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 16, pos) }
func evalPush17(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 17, pos) }
func evalPush18(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 18, pos) }
func evalPush19(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 19, pos) }
func evalPush20(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 20, pos) }
func evalPush21(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 21, pos) }
func evalPush22(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 22, pos) }
func evalPush23(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 23, pos) }
func evalPush24(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 24, pos) }
func evalPush25(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 25, pos) }
func evalPush26(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 26, pos) }
func evalPush27(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 27, pos) }
func evalPush28(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 28, pos) }
func evalPush29(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 29, pos) }
func evalPush30(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 30, pos) }
func evalPush31(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 31, pos) }
func evalPush32(state *Machine, _ opcodes.Opcode, pos int) Control { return push(state, 32, pos) }

func evalDup1(state *Machine, _ opcodes.Opcode, _ int) Control  { return dup(state, 1) }
func evalDup2(state *Machine, _ opcodes.Opcode, _ int) Control  { return dup(state, 2) }
func evalDup3(state *Machine, _ opcodes.Opcode, _ int) Control  { return dup(state, 3) }
func evalDup4(state *Machine, _ opcodes.Opcode, _ int) Control  { return dup(state, 4) }
func evalDup5(state *Machine, _ opcodes.Opcode, _ int) Control  { return dup(state, 5) }
func evalDup6(state *Machine, _ opcodes.Opcode, _ int) Control  { return dup(state, 6) }
func evalDup7(state *Machine, _ opcodes.Opcode, _ int) Control  { return dup(state, 7) }
func evalDup8(state *Machine, _ opcodes.Opcode, _ int) Control  { return dup(state, 8) }
func evalDup9(state *Machine, _ opcodes.Opcode, _ int) Control  { return dup(state, 9) }
func evalDup10(state *Machine, _ opcodes.Opcode, _ int) Control { return dup(state, 10) }
func evalDup11(state *Machine, _ opcodes.Opcode, _ int) Control { return dup(state, 11) }
func evalDup12(state *Machine, _ opcodes.Opcode, _ int) Control { return dup(state, 12) }
func evalDup13(state *Machine, _ opcodes.Opcode, _ int) Control { return dup(state, 13) }
func evalDup14(state *Machine, _ opcodes.Opcode, _ int) Control { return dup(state, 14) }
func evalDup15(state *Machine, _ opcodes.Opcode, _ int) Control { return dup(state, 15) }
func evalDup16(state *Machine, _ opcodes.Opcode, _ int) Control { return dup(state, 16) }

func evalSwap1(state *Machine, _ opcodes.Opcode, _ int) Control  { return swap(state, 1) }
func evalSwap2(state *Machine, _ opcodes.Opcode, _ int) Control  { return swap(state, 2) }
func evalSwap3(state *Machine, _ opcodes.Opcode, _ int) Control  { return swap(state, 3) }
func evalSwap4(state *Machine, _ opcodes.Opcode, _ int) Control  { return swap(state, 4) }
func evalSwap5(state *Machine, _ opcodes.Opcode, _ int) Control  { return swap(state, 5) }
func evalSwap6(state *Machine, _ opcodes.Opcode, _ int) Control  { return swap(state, 6) }
func evalSwap7(state *Machine, _ opcodes.Opcode, _ int) Control  { return swap(state, 7) }
func evalSwap8(state *Machine, _ opcodes.Opcode, _ int) Control  { return swap(state, 8) }
func evalSwap9(state *Machine, _ opcodes.Opcode, _ int) Control  { return swap(state, 9) }
func evalSwap10(state *Machine, _ opcodes.Opcode, _ int) Control { return swap(state, 10) }
func evalSwap11(state *Machine, _ opcodes.Opcode, _ int) Control { return swap(state, 11) }
func evalSwap12(state *Machine, _ opcodes.Opcode, _ int) Control { return swap(state, 12) }
func evalSwap13(state *Machine, _ opcodes.Opcode, _ int) Control { return swap(state, 13) }
func evalSwap14(state *Machine, _ opcodes.Opcode, _ int) Control { return swap(state, 14) }
func evalSwap15(state *Machine, _ opcodes.Opcode, _ int) Control { return swap(state, 15) }
func evalSwap16(state *Machine, _ opcodes.Opcode, _ int) Control { return swap(state, 16) }

func evalReturn(state *Machine, _ opcodes.Opcode, _ int) Control {
	start, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	length, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	state.Mem.ResizeOffset(start, length)
	state.ReturnRange = &RangeU256{Start: start, End: new(uint256.Int).Add(start, length)}
	return &Exit{exits.Returned}
}
func evalRevert(state *Machine, _ opcodes.Opcode, _ int) Control {
	start, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	length, err := state.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	state.ReturnRange = &RangeU256{Start: start, End: new(uint256.Int).Add(start, length)}
	return &Exit{exits.Reverted}
}
func evalInvalid(state *Machine, _ opcodes.Opcode, _ int) Control {
	return &Exit{exits.DesignatedInvalid}
}

func evalExternel(state *Machine, opcode opcodes.Opcode, _ int) Control { return &Trap{opcode} }
