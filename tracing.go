package slashvm

import (
	"github.com/c0mm4nd/slashvm/core"
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/c0mm4nd/slashvm/core/opcodes"
)

type EventStep struct {
	Context    *Context
	Opcode     opcodes.Opcode
	Pos        uint
	ExitResult exits.Reason
	Stack      core.Stack
	Mem        core.Mem
}
