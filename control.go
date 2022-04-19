package slashvm

import "github.com/c0mm4nd/slashvm/core/exits"

type Control interface {
	RuntimeCtrlType() string
}

type Continue struct{}

func (*Continue) RuntimeCtrlType() string { return "continue" }

type CallInterrupt struct{}

func (*CallInterrupt) RuntimeCtrlType() string { return "CallInterrupt" }

type CreateInterrupt struct{}

func (*CreateInterrupt) RuntimeCtrlType() string { return "CreateInterrupt" }

type Exit struct {
	Reason exits.Reason
}

func (*Exit) RuntimeCtrlType() string { return "Exit" }
