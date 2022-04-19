package opcodes

type Opcode byte

// Core opcodes.
const (
	STOP Opcode = iota
	ADD
	MUL
	SUB
	DIV
	SDIV
	MOD
	SMOD
	ADDMOD
	MULMOD
	EXP
	SIGNEXTEND
)

const (
	LT Opcode = iota + 0x10
	GT
	SLT
	SGT
	EQ
	ISZERO
	AND
	OR
	XOR
	NOT
	BYTE
	SHL
	SHR
	SAR
)

const (
	CALLDATALOAD = iota + 0x35
	CALLDATASIZE
	CALLDATACOPY
	CODESIZE
	CODECOPY
)

const (
	POP      Opcode = 0x50
	MLOAD    Opcode = 0x51
	MSTORE   Opcode = 0x52
	MSTORE8  Opcode = 0x53
	JUMP     Opcode = 0x56
	JUMPI    Opcode = 0x57
	PC       Opcode = 0x58
	MSIZE    Opcode = 0x59
	JUMPDEST Opcode = 0x5b
)

const (
	PUSH1 Opcode = iota + 0x60
	PUSH2
	PUSH3
	PUSH4
	PUSH5
	PUSH6
	PUSH7
	PUSH8
	PUSH9
	PUSH10
	PUSH11
	PUSH12
	PUSH13
	PUSH14
	PUSH15
	PUSH16
	PUSH17
	PUSH18
	PUSH19
	PUSH20
	PUSH21
	PUSH22
	PUSH23
	PUSH24
	PUSH25
	PUSH26
	PUSH27
	PUSH28
	PUSH29
	PUSH30
	PUSH31
	PUSH32
	DUP1
	DUP2
	DUP3
	DUP4
	DUP5
	DUP6
	DUP7
	DUP8
	DUP9
	DUP10
	DUP11
	DUP12
	DUP13
	DUP14
	DUP15
	DUP16
	SWAP1
	SWAP2
	SWAP3
	SWAP4
	SWAP5
	SWAP6
	SWAP7
	SWAP8
	SWAP9
	SWAP10
	SWAP11
	SWAP12
	SWAP13
	SWAP14
	SWAP15
	SWAP16

	RETURN  Opcode = 0xf3
	REVERT  Opcode = 0xfd
	INVALID Opcode = 0xfe
)

// External opcodes
const (
	SHA3 Opcode = 0x20

	ADDRESS   Opcode = 0x30
	BALANCE   Opcode = 0x31
	ORIGIN    Opcode = 0x32
	CALLER    Opcode = 0x33
	CALLVALUE Opcode = 0x34

	GASPRICE       Opcode = 0x3a
	EXTCODESIZE    Opcode = 0x3b
	EXTCODECOPY    Opcode = 0x3c
	RETURNDATASIZE Opcode = 0x3d
	RETURNDATACOPY Opcode = 0x3e
	EXTCODEHASH    Opcode = 0x3f

	BLOCKHASH    Opcode = 0x40
	COINBASE     Opcode = 0x41
	TIMESTAMP    Opcode = 0x42
	NUMBER       Opcode = 0x43
	DIFFICULTY   Opcode = 0x44
	GASLIMIT     Opcode = 0x45
	CHAINID      Opcode = 0x46
	SELFBALANCE  Opcode = 0x47
	BASEFEE      Opcode = 0x48
	SLOAD        Opcode = 0x54
	SSTORE       Opcode = 0x55
	GAS          Opcode = 0x5a
	LOG0         Opcode = 0xa0
	LOG1         Opcode = 0xa1
	LOG2         Opcode = 0xa2
	LOG3         Opcode = 0xa3
	LOG4         Opcode = 0xa4
	CREATE       Opcode = 0xf0
	CALL         Opcode = 0xf1
	CALLCODE     Opcode = 0xf2
	DELEGATECALL Opcode = 0xf4
	STATICCALL   Opcode = 0xfa
	CREATE2      Opcode = 0xf5
	SUICIDE      Opcode = 0xff
)

// Whether the opcode is a PUSH opcode.
func (op Opcode) IsPush() int {
	if op >= 0x60 && op < 0x7f {
		return int(op - 0x60 + 1)
	} else {
		return 0
	}
}
