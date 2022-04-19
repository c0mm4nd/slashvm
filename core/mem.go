package core

import (
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/holiman/uint256"
)

type Mem struct {
	Data         []byte
	EffectiveLen *uint256.Int // u256
	Limit        uint
}

func NewMem(limit uint) *Mem {
	return &Mem{
		Data:         nil,
		EffectiveLen: uint256.NewInt(0),
		Limit:        limit,
	}
}

func (m *Mem) Len() uint {
	return uint(len(m.Data))
}

func (m *Mem) IsEmpty() bool {
	return m.Len() == 0
}

func (m *Mem) ResizeOffset(offset, length *uint256.Int) *exits.Error {
	if length.IsZero() {
		return nil
	}

	if end, overflow := offset.AddOverflow(offset, length); !overflow {
		return m.ResizeEnd(end)
	} else {
		return exits.InvalidRange
	}
}

func next_multiple_of_32(x *uint256.Int) (*uint256.Int, bool) {
	r := uint32(x[0] & 0xFFFFFFFF)
	r = (^(r & 31) + 1) & 31
	return new(uint256.Int).AddOverflow(x, uint256.NewInt(uint64(r)))
}

func (m *Mem) ResizeEnd(end *uint256.Int) *exits.Error {
	if end.Cmp(m.EffectiveLen) > 0 {
		newEnd, overflow := next_multiple_of_32(end)
		if overflow {
			return exits.InvalidRange
		}

		m.EffectiveLen = newEnd
	}

	return nil
}

func (m *Mem) Get(offset, size uint64) []byte {
	ret := make([]byte, size)

	for index := uint64(0); index < size; index++ {
		pos := offset + index
		if pos >= uint64(m.Len()) {
			break
		}

		ret[index] = m.Data[pos]
	}

	return ret
}

func (m *Mem) Set(offset int, value []byte, targetSize int) *exits.Fatal {
	if targetSize == 0 {
		return nil
	}

	if uint(offset+targetSize) > m.Limit {
		return exits.NotSupported
	}

	if m.Len() < uint(offset+targetSize) {
		m.Data = append(m.Data, make([]byte, offset+targetSize-int(m.Len()))...)
	}

	if targetSize > len(value) {
		copy(m.Data[offset:offset+len(value)], value)
		for index := len(value); index < targetSize; index++ {
			m.Data[offset+index] = 0
		}
	} else {
		copy(m.Data[offset:offset+targetSize], value[:targetSize])
	}

	return nil
}

func (m *Mem) CopyLarge(memOffset, dataOffset, length *uint256.Int, data []byte) *exits.Fatal {
	if length.IsZero() {
		return nil
	}

	if memOffset.Cmp(uint256.NewInt(uint64(^uint(0)))) > 0 {
		return exits.NotSupported
	}
	uMemOffset := int(memOffset[0])

	if length.Cmp(uint256.NewInt(uint64(^uint(0)))) > 0 {
		return exits.NotSupported
	}
	ulen := int(length.Uint64())

	if end, overflow := dataOffset.AddOverflow(dataOffset, length); !overflow {
		if end.Cmp(uint256.NewInt(uint64(^uint(0)))) > 0 {
			data = data[:0]
		} else {
			uDataOffset := int(dataOffset.Uint64())
			uEnd := int(end.Uint64())
			if uDataOffset > len(data) {
				data = data[:0]
			} else {
				data = data[uDataOffset:min(uEnd, len(data))]
			}
		}
	} else {
		data = data[:0]
	}

	m.Set(uMemOffset, data, ulen)

	return nil
}
