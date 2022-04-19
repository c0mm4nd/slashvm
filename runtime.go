package slashvm

import (
	"github.com/c0mm4nd/slashvm/configs"
	"github.com/c0mm4nd/slashvm/core"
	"github.com/c0mm4nd/slashvm/core/exits"
)

type Runtime struct {
	Machine          *core.Machine
	Status           exits.Reason
	ReturnDataBuffer []byte
	Context          *Context
	Config           *configs.Config
}

func NewRuntime(code, data []byte, context *Context, config *configs.Config) *Runtime {
	return &Runtime{
		Machine:          core.NewMachine(code, data, config.StackLimit, config.MemoryLimit),
		Status:           nil,
		ReturnDataBuffer: nil,
		Context:          context,
		Config:           config,
	}
}

func (r *Runtime) Step(h Handler) exits.Reason {
	opcode, stack := r.Machine.Inspect()
	//emitEvent(Step,
	//	context: r.Context,
	//	opcode,
	//	position: r.Machine.Pos,
	//	stack,
	//	memory: r.Machine.Mem
	//)

	err := h.PreValidate(r.Context, opcode, stack)
	if err != nil {
		r.Machine.Exit(err)
		r.Status = err
	}

	if r.Status != nil {
		return err
	}

	exitReason, trap := r.Machine.Step()
	//emitEvent(result, r.Machine.ReturnValue())
	if exitReason != nil {
		r.Status = exitReason
		return exitReason
	}

	if trap != nil {
		ctrl := Eval(r, opcode, h)
		switch ctrl := ctrl.(type) {
		case *Continue:
			return nil
		case *CallInterrupt:
			r.Machine.Exit(exits.UnhandledInterrupt)
			r.Status = exits.UnhandledInterrupt
			return exits.UnhandledInterrupt
		case *CreateInterrupt:
			r.Machine.Exit(exits.UnhandledInterrupt)
			r.Status = exits.UnhandledInterrupt
			return exits.UnhandledInterrupt
		case *Exit:
			r.Machine.Exit(ctrl.Reason)
			r.Status = ctrl.Reason
			return ctrl.Reason
		default:
			panic(ctrl)
		}
	}

	return nil
}

func (r *Runtime) Run(h Handler) exits.Reason {
	for {
		exitReason := r.Step(h)
		if exitReason != nil {
			return exitReason
		}
	}
}
