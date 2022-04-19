package core

import (
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/c0mm4nd/slashvm/core/opcodes"
)

type Control interface {
	CoreCtrlType() string
}

type Jump struct{ N int }

func (j *Jump) CoreCtrlType() string { return "jump" }

type Continue struct{ N int }

func (c *Continue) CoreCtrlType() string { return "continue" }

type Exit struct {
	Result exits.Reason
}

func (c *Exit) CoreCtrlType() string { return "exit" }

type Trap struct {
	Opcode opcodes.Opcode
}

func (t Trap) CoreCtrlType() string { return "trap" }
