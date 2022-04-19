package slashvm

import (
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/c0mm4nd/slashvm/utils"
	"github.com/holiman/uint256"
)

func evalCall(runtime *Runtime, handler Handler) Control {
	runtime.ReturnDataBuffer = runtime.ReturnDataBuffer[len(runtime.ReturnDataBuffer):] // clear
	gas, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	to, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	value, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	in_offset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	in_len, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	out_offset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	out_len, err := runtime.Machine.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	uint64Gas, overflow := gas.Uint64WithOverflow()
	if overflow {
		uint64Gas = 0
	}

	runtime.Machine.Mem.ResizeOffset(in_offset, in_len)
	runtime.Machine.Mem.ResizeOffset(out_offset, out_len)

	var input []byte
	if !in_len.IsZero() {
		offset, overflow := in_offset.Uint64WithOverflow()
		if overflow {
			return &Exit{exits.NotSupported}
		}
		l, overflow := in_len.Uint64WithOverflow()
		if overflow {
			return &Exit{exits.NotSupported}
		}
		input = runtime.Machine.Mem.Get(offset, l)
	}

	context := &Context{
		Address:       to,
		Caller:        runtime.Context.Address,
		ApparentValue: value,
	}

	transfer := &Transfer{
		Source: utils.U256ToHashHex(runtime.Context.Address),
		Target: utils.U256ToHashHex(to),
		Value:  value,
	}

	exitReason, data, callInterrupt := handler.Call(to, transfer, input, uint64Gas, context)
	if exitReason != nil {
		runtime.ReturnDataBuffer = data
		targetLen := out_len
		if uint64(len(data)) < targetLen.Uint64() {
			targetLen.SetUint64(uint64(len(data)))
		}
		switch exitReason := exitReason.(type) {
		case *exits.Succeed:
			fatal := runtime.Machine.Mem.CopyLarge(out_offset, new(uint256.Int), targetLen, data)
			if fatal != nil {
				out_len.SetUint64(0)
				return &Continue{}
			} else {
				out_len.SetUint64(1)
				return &Continue{}
			}
		case *exits.Revert:
			_ = runtime.Machine.Mem.CopyLarge(out_offset, new(uint256.Int), targetLen, data)

			return &Continue{}
		case *exits.Error:
			out_len.Clear()
			return &Continue{}
		case *exits.Fatal:
			out_len.Clear()
			return &Exit{exitReason}
		default:
			panic("unknown type")
		}
	}

	out_len.Clear()
	return callInterrupt

}
func evalCallCode(runtime *Runtime, handler Handler) Control {
	runtime.ReturnDataBuffer = runtime.ReturnDataBuffer[len(runtime.ReturnDataBuffer):] // clear
	gas, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	to, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	value, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	in_offset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	in_len, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	out_offset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	out_len, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	runtime.Machine.Mem.ResizeOffset(in_offset, in_len)
	runtime.Machine.Mem.ResizeOffset(out_offset, out_len)

	uint64Gas, overflow := gas.Uint64WithOverflow()
	if overflow {
		uint64Gas = 0
	}

	var input []byte
	if !in_len.IsZero() {
		offset, overflow := in_offset.Uint64WithOverflow()
		if overflow {
			return &Exit{exits.NotSupported}
		}
		l, overflow := in_len.Uint64WithOverflow()
		if overflow {
			return &Exit{exits.NotSupported}
		}
		input = runtime.Machine.Mem.Get(offset, l)
	}

	context := &Context{
		Address:       runtime.Context.Address,
		Caller:        runtime.Context.Address,
		ApparentValue: value,
	}

	transfer := &Transfer{
		Source: utils.U256ToHashHex(runtime.Context.Address),
		Target: utils.U256ToHashHex(runtime.Context.Address),
		Value:  value,
	}

	exitReason, data, callInterrupt := handler.CallCode(to, transfer, input, uint64Gas, context)
	if exitReason != nil {
		runtime.ReturnDataBuffer = data
		targetLen := out_len
		if uint64(len(data)) < targetLen.Uint64() {
			targetLen.SetUint64(uint64(len(data)))
		}
		switch exitReason := exitReason.(type) {
		case *exits.Succeed:
			fatal := runtime.Machine.Mem.CopyLarge(out_offset, new(uint256.Int), targetLen, data)
			if fatal != nil {
				out_len.SetUint64(0)
				return &Continue{}
			} else {
				out_len.SetUint64(1)
				return &Continue{}
			}
		case *exits.Revert:
			_ = runtime.Machine.Mem.CopyLarge(out_offset, new(uint256.Int), targetLen, data)

			return &Continue{}
		case *exits.Error:
			out_len.Clear()
			return &Continue{}
		case *exits.Fatal:
			out_len.Clear()
			return &Exit{exitReason}
		default:
			panic("unknown type")
		}
	}

	out_len.Clear()
	return callInterrupt
}

func evalDelegateCall(runtime *Runtime, handler Handler) Control {
	runtime.ReturnDataBuffer = runtime.ReturnDataBuffer[len(runtime.ReturnDataBuffer):] // clear
	gas, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	to, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	in_offset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	in_len, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	out_offset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	out_len, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	uint64Gas, overflow := gas.Uint64WithOverflow()
	if overflow {
		uint64Gas = 0
	}

	runtime.Machine.Mem.ResizeOffset(in_offset, in_len)
	runtime.Machine.Mem.ResizeOffset(out_offset, out_len)

	var input []byte
	if !in_len.IsZero() {
		offset, overflow := in_offset.Uint64WithOverflow()
		if overflow {
			return &Exit{exits.NotSupported}
		}
		l, overflow := in_len.Uint64WithOverflow()
		if overflow {
			return &Exit{exits.NotSupported}
		}
		input = runtime.Machine.Mem.Get(offset, l)
	}

	context := &Context{
		Address:       runtime.Context.Address,
		Caller:        runtime.Context.Caller,
		ApparentValue: runtime.Context.ApparentValue,
	}

	exitReason, data, callInterrupt := handler.DelegateCall(to, input, uint64Gas, context)
	if exitReason != nil {
		runtime.ReturnDataBuffer = data
		targetLen := out_len
		if uint64(len(data)) < targetLen.Uint64() {
			targetLen.SetUint64(uint64(len(data)))
		}
		switch exitReason := exitReason.(type) {
		case *exits.Succeed:
			fatal := runtime.Machine.Mem.CopyLarge(out_offset, new(uint256.Int), targetLen, data)
			if fatal != nil {
				out_len.SetUint64(0)
				return &Continue{}
			}
			out_len.SetUint64(1)
			return &Continue{}
		case *exits.Revert:
			_ = runtime.Machine.Mem.CopyLarge(out_offset, new(uint256.Int), targetLen, data)

			return &Continue{}
		case *exits.Error:
			out_len.Clear()
			return &Continue{}
		case *exits.Fatal:
			out_len.Clear()
			return &Exit{exitReason}
		default:
			panic("unknown type")
		}
	}
	out_len.Clear()
	return callInterrupt

}
func evalStaticCall(runtime *Runtime, handler Handler) Control {
	runtime.ReturnDataBuffer = runtime.ReturnDataBuffer[len(runtime.ReturnDataBuffer):] // clear
	gas, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	to, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	in_offset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	in_len, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	out_offset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	out_len, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	uint64Gas, overflow := gas.Uint64WithOverflow()
	if overflow {
		uint64Gas = 0
	}

	runtime.Machine.Mem.ResizeOffset(in_offset, in_len)
	runtime.Machine.Mem.ResizeOffset(out_offset, out_len)

	var input []byte
	if !in_len.IsZero() {
		offset, overflow := in_offset.Uint64WithOverflow()
		if overflow {
			return &Exit{exits.NotSupported}
		}
		l, overflow := in_len.Uint64WithOverflow()
		if overflow {
			return &Exit{exits.NotSupported}
		}
		input = runtime.Machine.Mem.Get(offset, l)
	}

	context := &Context{
		Address:       to,
		Caller:        runtime.Context.Address,
		ApparentValue: new(uint256.Int),
	}

	exitReason, data, callInterrupt := handler.StaticCode(to, input, uint64Gas, context)
	if exitReason != nil {
		runtime.ReturnDataBuffer = data
		targetLen := out_len
		if uint64(len(data)) < targetLen.Uint64() {
			targetLen.SetUint64(uint64(len(data)))
		}
		switch exitReason := exitReason.(type) {
		case *exits.Succeed:
			fatal := runtime.Machine.Mem.CopyLarge(out_offset, new(uint256.Int), targetLen, data)
			if fatal != nil {
				out_len.SetUint64(0)
				return &Continue{}
			}
			out_len.SetUint64(1)
			return &Continue{}

		case *exits.Revert:
			_ = runtime.Machine.Mem.CopyLarge(out_offset, new(uint256.Int), targetLen, data)

			return &Continue{}
		case *exits.Error:
			out_len.Clear()
			return &Continue{}
		case *exits.Fatal:
			out_len.Clear()
			return &Exit{exitReason}
		default:
			panic("unknown type")
		}
	}
	out_len.Clear()
	return callInterrupt
}
