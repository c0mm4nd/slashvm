package exits

type Reason interface {
	IsExitReason()
}

type Succeed struct{ Reason string }

func (*Succeed) IsExitReason() {}

var (
	Stopped  = &Succeed{"stopped"}
	Returned = &Succeed{"returned"}
	Suicide  = &Succeed{"suicide"}
)

type Error struct{ Reason string }

func (*Error) IsExitReason() {}

// TODO: copy desc from https://github.com/rust-blockchain/evm/blob/master/core/src/error.rs
var (
	StackUnderflow      = &Error{"stack underflow"}
	StackOverflow       = &Error{"stack overflow"}
	InvalidJump         = &Error{"invalid jump destination"}
	InvalidRange        = &Error{"invalid range"}
	DesignatedInvalid   = &Error{"designated invalid"}
	CallTooDeep         = &Error{"call stack limit reached"}
	CreateCollision     = &Error{"create collision"}
	CreateContractLimit = &Error{"create contract limit"}
	InvalidCode         = &Error{"invalid code"}

	OutOfOffset = &Error{"out of offset"}
	OutOfGas    = &Error{"out of gas"}
	OutOfFund   = &Error{"out of fund"}

	PCUnderflow = &Error{"pc underflow"}
	CreateEmpty = &Error{"create empty"}

	OtherError = &Error{"other error"}
)

type Fatal struct{ Reason string }

func (*Fatal) IsExitReason() {}

var (
	NotSupported       = &Fatal{"not supported"}
	UnhandledInterrupt = &Fatal{"unhandled interrupt"}

	OtherFatal = &Fatal{"other fatal"}
)

type Revert struct{ Reason string }

func (*Revert) IsExitReason() {}

var (
	Reverted = &Revert{"reverted"}
)
