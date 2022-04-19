package slashvm

import "github.com/c0mm4nd/slashvm/core/exits"

func evalCreate(runtime *Runtime, handler Handler) Control {
	runtime.ReturnDataBuffer = runtime.ReturnDataBuffer[len(runtime.ReturnDataBuffer):] // clear

	value, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	codeOffset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	length, err := runtime.Machine.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}
	err = runtime.Machine.Mem.ResizeOffset(codeOffset, length)
	if err != nil {
		return &Exit{err}
	}

	var code []byte
	if length.IsZero() {
		code = []byte{}
	} else {
		code = runtime.Machine.Mem.Get(codeOffset.Uint64(), length.Uint64())
	}

	scheme := &CreateSchemeLegacy{
		Caller: runtime.Context.Address,
	}

	exitReason, addr, data, createInterrupt := handler.Create(runtime.Context.Address, scheme, value, code)
	if exitReason != nil {
		runtime.ReturnDataBuffer = data
		createAddress := addr

		switch reason := exitReason.(type) {
		case *exits.Succeed:
			raw := createAddress.Bytes20()
			length.SetBytes20(raw[:])
			return &Continue{}
		case *exits.Revert:
			length.Clear()
			return &Continue{}
		case *exits.Error:
			length.Clear()
			return &Continue{}
		case *exits.Fatal:
			length.Clear()
			return &Exit{reason}
		default:
			panic("unknown reason")
		}
	}
	// Trap
	// createInterrupt != nil
	length.Clear()
	return createInterrupt
}
func evalCreate2(runtime *Runtime, handler Handler) Control {
	runtime.ReturnDataBuffer = runtime.ReturnDataBuffer[len(runtime.ReturnDataBuffer):] // clear

	value, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	codeOffset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	length, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	err = runtime.Machine.Mem.ResizeOffset(codeOffset, length)
	if err != nil {
		return &Exit{err}
	}

	var code []byte
	if length.IsZero() {
		code = []byte{}
	} else {
		code = runtime.Machine.Mem.Get(codeOffset.Uint64(), length.Uint64())
	}

	salt, err := runtime.Machine.Stack.Peek()
	codeHash := keccak256Hash(code)
	if err != nil {
		return &Exit{err}
	}
	scheme := &CreateSchemeCreate2{
		Caller:   runtime.Context.Address,
		Salt:     salt.Bytes(),
		CodeHash: codeHash,
	}

	exitReason, addr, data, createInterrupt := handler.Create2(runtime.Context.Address, scheme, value, code)
	if exitReason != nil {
		runtime.ReturnDataBuffer = data
		createAddress := addr

		switch reason := exitReason.(type) {
		case *exits.Succeed:
			raw := createAddress.Bytes20()
			length.SetBytes20(raw[:])
			return &Continue{}
		case *exits.Revert:
			length.Clear()
			return &Continue{}
		case *exits.Error:
			length.Clear()
			return &Continue{}
		case *exits.Fatal:
			length.Clear()
			return &Exit{reason}
		default:
			panic("unknown reason")
		}
	}
	// Trap
	// createInterrupt != nil
	length.Clear()
	return createInterrupt

}
