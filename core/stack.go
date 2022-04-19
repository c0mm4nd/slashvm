package core

import (
	"github.com/c0mm4nd/slashvm/core/exits"
	"github.com/holiman/uint256"
)

type Stack struct {
	Data  []*uint256.Int
	Limit uint
}

// Create a new stack with given limit.
func NewStack(limit uint) *Stack {
	return &Stack{Data: make([]*uint256.Int, 0, limit), Limit: limit}
}

func (s *Stack) Pop() (*uint256.Int, *exits.Error) {
	// s.Data.pop
	l := len(s.Data)
	if l == 0 {
		return nil, exits.StackUnderflow
	}

	ret := s.Data[l-1]
	s.Data = s.Data[:l-1]

	return ret, nil
}

func (s *Stack) Push(v *uint256.Int) *exits.Error {
	if v == nil {
		panic("pushing nil")
	}
	if len(s.Data)+1 > int(s.Limit) {
		return exits.StackOverflow
	}

	s.Data = append(s.Data, v)
	return nil
}

func (s *Stack) PeekN(numFromTop int) (*uint256.Int, *exits.Error) {
	if len(s.Data) > numFromTop {
		return s.Data[len(s.Data)-numFromTop-1], nil
	} else {
		return nil, exits.StackUnderflow
	}
}

func (s *Stack) Peek() (*uint256.Int, *exits.Error) {
	if len(s.Data) > 0 {
		return s.Data[len(s.Data)-1], nil
	} else {
		return nil, exits.StackUnderflow
	}
}

func (s *Stack) Set(numFromTop int, v *uint256.Int) *exits.Error {
	if len(s.Data) > numFromTop {
		s.Data[len(s.Data)-numFromTop-1] = v
		return nil
	} else {
		return exits.StackUnderflow
	}
}
