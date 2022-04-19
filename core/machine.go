package core

import (
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/c0mm4nd/slashvm/core/opcodes"
	"github.com/holiman/uint256"
)

type Machine struct {
	Data        []byte
	Code        []opcodes.Opcode
	Pos         int
	FinalResult exits.Reason
	ReturnRange *RangeU256
	Valids      *Valids
	Mem         *Mem
	Stack       *Stack
}

func NewMachine(code, data []byte, stackLimit, memLimit uint) *Machine {
	return &Machine{
		Data:        data,
		Code:        opcodes.ToOpcodes(code),
		Pos:         0,
		ReturnRange: &RangeU256{uint256.NewInt(0), uint256.NewInt(0)},
		Valids:      NewValids(code),
		Mem:         NewMem(memLimit),
		Stack:       NewStack(stackLimit),
	}
}

func (m *Machine) Exit(exitReason exits.Reason) {
	m.Pos = -1
	m.FinalResult = exitReason
}

// Inspect the machine's next opcode and current stack.
func (m *Machine) Inspect() (opcodes.Opcode, *Stack) {
	return opcodes.Opcode(m.Code[m.Pos]), m.Stack
}

func (m *Machine) ReturnValue() []byte {
	d := new(uint256.Int).Sub(m.ReturnRange.End, m.ReturnRange.Start)
	if m.ReturnRange.Start.Cmp(uint256.NewInt(uint64(^uint(0)))) > 0 {
		return make([]byte, d.Uint64())
	} else if m.ReturnRange.End.Cmp(uint256.NewInt(uint64(^uint(0)))) > 0 {
		ret := m.Mem.Get(
			m.ReturnRange.Start.Uint64(),
			new(uint256.Int).Sub(uint256.NewInt(uint64(^uint(0))), m.ReturnRange.Start).Uint64(),
		)
		if d.Cmp(uint256.NewInt(uint64(len(ret)))) > 0 {
			d.Sub(d, uint256.NewInt(uint64(len(ret))))
			ret = append(ret, make([]byte, d.Sub(d, uint256.NewInt(uint64(len(ret)))).Uint64())...)
		}
		return ret
	} else {
		return m.Mem.Get(m.ReturnRange.Start.Uint64(), d.Uint64())
	}
}
func (m *Machine) Run() (exits.Reason, *Trap) {
	for {
		err, trap := m.Step()
		if err != nil {
			return err, nil
		}
		if trap != nil {
			return nil, trap
		}
	}
}

func (m *Machine) Step() (exits.Reason, *Trap) {
	pos := m.Pos

	opcode := opcodes.Opcode(m.Code[pos])
	ctrl := Eval(m, opcode, pos)
	if ctrl != nil {
		switch ctrl := ctrl.(type) {
		case *Exit:
			m.FinalResult = ctrl.Result
			return ctrl.Result, nil
		case *Jump:
			m.Pos = ctrl.N
			return nil, nil
		case *Trap:
			m.Pos = pos + 1
			err := &Trap{Opcode: ctrl.Opcode}
			return nil, err
		case *Continue:
			m.Pos = pos + ctrl.N
			return nil, nil
		default:
			panic("unknown type")
		}
	} else {
		m.FinalResult = exits.Stopped
		return m.FinalResult, nil
	}
}
