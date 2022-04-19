package core

import "github.com/c0mm4nd/slashvm/core/opcodes"

type Valids struct{ bools []bool }

func NewValids(code []byte) *Valids {
	valids := make([]bool, len(code))

	i := 0
	for i < len(code) {
		opcode := opcodes.Opcode(code[i])
		if opcode == opcodes.JUMPDEST {
			valids[i] = true
			i += 1
		} else if v := opcode.IsPush(); v > 0 {
			i += v + 1
		} else {
			i += 1
		}
	}

	return &Valids{bools: valids}
}

func (vs Valids) IsEmpty() bool {
	return len(vs.bools) == 0
}

func (vs Valids) IsValid(pos int) bool {
	if pos > len(vs.bools) {
		return false
	}

	return vs.bools[pos]
}
