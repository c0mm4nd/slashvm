package slashvm

import (
	"hash"

	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/sha3"
)

var keccak256 hash.Hash

func keccak256Hash(b []byte) []byte {
	if keccak256 == nil {
		keccak256 = sha3.NewLegacyKeccak256()
	} else {
		keccak256.Reset()
	}

	keccak256.Write(b)
	return keccak256.Sum(nil)
}

func evalKeccak256(runtime *Runtime) Control {
	from, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	length, err := runtime.Machine.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	if length.IsZero() {
		length.Clear()
		return nil
	}

	length.SetBytes(keccak256Hash(from.Bytes()))

	return &Continue{}
}

func evalChainID(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(handler.ChainID())

	return &Continue{}
}

func evalAddress(runtime *Runtime) Control {
	runtime.Machine.Stack.Push(new(uint256.Int).SetBytes(runtime.Context.Address.Bytes()))

	return &Continue{}
}
func evalBalance(runtime *Runtime, handler Handler) Control {
	addr, err := runtime.Machine.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	addr.Set(handler.Balance(addr))
	return &Continue{}
}
func evalSelfBalance(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(handler.Balance(runtime.Context.Address))

	return &Continue{}
}
func evalOrigin(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(handler.Origin())

	return &Continue{}
}
func evalCaller(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(new(uint256.Int).SetBytes(runtime.Context.Caller.Bytes()))

	return &Continue{}
}
func evalCallValue(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(new(uint256.Int).Set(runtime.Context.ApparentValue))

	return &Continue{}
}
func evalGasPrice(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(handler.GasPrice())

	return &Continue{}
}
func evalBaseFee(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(handler.BlockBaseFeePerGas())

	return &Continue{}
}
func evalExtCodeSize(runtime *Runtime, handler Handler) Control {
	addr, err := runtime.Machine.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	addr.Set(handler.CodeSize(addr))

	return &Continue{}
}
func evalExtCodeHash(runtime *Runtime, handler Handler) Control {
	addr, err := runtime.Machine.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}

	addr.Set(handler.CodeHash(addr))

	return &Continue{}
}
func evalExtCodeCopy(runtime *Runtime, handler Handler) Control {
	addr, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	memoryOffset, err := runtime.Machine.Stack.Pop()
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

	fatal := runtime.Machine.Mem.CopyLarge(memoryOffset, codeOffset, length, handler.Code(addr))
	if fatal != nil {
		return &Exit{fatal}
	}
	return &Continue{}
}
func evalReturnDataSize(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(uint256.NewInt(uint64(len(runtime.ReturnDataBuffer))))

	return &Continue{}
}
func evalReturnDataCopy(runtime *Runtime, handler Handler) Control {
	memoryOffset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	dataOffset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	length, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	if new(uint256.Int).Add(dataOffset, length).Cmp(uint256.NewInt(uint64(len(runtime.ReturnDataBuffer)))) > 0 {
		return &Exit{exits.OutOfOffset}
	}

	fatal := runtime.Machine.Mem.CopyLarge(memoryOffset, dataOffset, length, runtime.ReturnDataBuffer)
	if fatal != nil {
		return &Exit{fatal}
	}

	return &Continue{}
}

func evalBlockHash(runtime *Runtime, handler Handler) Control {
	number, err := runtime.Machine.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}
	number.Set(handler.BlockHash(number.Uint64()))
	return &Continue{}
}
func evalCoinbase(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(handler.BlockCoinbase())

	return &Continue{}
}
func evalTimestamp(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(handler.BlockTimestamp())

	return &Continue{}
}
func evalNumber(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(handler.BlockNumber())

	return &Continue{}
}
func evalDifficulty(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(handler.BlockDifficulty())

	return &Continue{}
}
func evalGasLimit(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(handler.GasLimit())

	return &Continue{}
}
func evalSLoad(runtime *Runtime, handler Handler) Control {
	index, err := runtime.Machine.Stack.Peek()
	if err != nil {
		return &Exit{err}
	}
	value := handler.Storage(runtime.Context.Address, index)
	index.Set(value)

	//emitEvent(EventSLoad{}, runtime.Context.Address, index, value)

	return &Continue{}
}

func evalSStore(runtime *Runtime, handler Handler) Control {
	index, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	value, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	//emitEvent(EventSStore{}, runtime.Context.Address, index, value)

	fatal := handler.SetStorage(runtime.Context.Address, index, value)
	if fatal != nil {
		return &Exit{fatal}
	}
	return &Continue{}
}
func evalGas(runtime *Runtime, handler Handler) Control {
	runtime.Machine.Stack.Push(handler.GasLeft())

	return &Continue{}
}

func evalLog(runtime *Runtime, n uint8, handler Handler) Control {
	offset, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	length, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}
	err = runtime.Machine.Mem.ResizeOffset(offset, length)
	if err != nil {
		return &Exit{err}
	}
	var data []byte
	if length.IsZero() {
		data = make([]byte, 0)
	} else {
		data = runtime.Machine.Mem.Get(offset.Uint64(), length.Uint64())
	}

	topics := make([]*uint256.Int, 0, n)
	for i := uint8(0); i < n; i++ {
		val, err := runtime.Machine.Stack.Pop()
		if err != nil {
			return &Exit{err}
		}
		topics = append(topics, val)
	}

	err = handler.Log(runtime.Context.Address, topics, data)
	if err != nil {
		return &Exit{err}
	}

	return &Continue{}
}

func evalSuicide(runtime *Runtime, handler Handler) Control {
	target, err := runtime.Machine.Stack.Pop()
	if err != nil {
		return &Exit{err}
	}

	err = handler.MarkDelete(runtime.Context.Address, target)
	if err != nil {
		return &Exit{err}
	}

	return &Continue{}
}
