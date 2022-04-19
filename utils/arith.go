package utils

import (
	"math/bits"

	"github.com/holiman/uint256"
	"golang.org/x/exp/constraints"
)

func SaturatingAdd(a, b *uint256.Int) (c *uint256.Int) {
	c, overflow := new(uint256.Int).AddOverflow(a, b)
	if overflow {
		c.SetAllOne()
	}

	return c
}

func SaturatingAddToA(a, b *uint256.Int) {
	a, overflow := new(uint256.Int).AddOverflow(a, b)
	if overflow {
		a.SetAllOne()
	}
}

func Log2Floor(a *uint256.Int) uint64 {
	if a.IsZero() {
		panic("Log2Floor(0)")
	}

	l := uint64(256)
	for i := 0; i < 4; i++ {
		n := uint64(3 - i)
		if a[n] == 0 {
			l -= 64
		} else {
			l -= uint64(bits.LeadingZeros64(a[n]))
			if l == 0 {
				return l
			} else {
				return l - 1
			}
		}
	}

	return l
}

func Max[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s {
		if m < v {
			m = v
		}
	}
	return m
}

func Min[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s {
		if m > v {
			m = v
		}
	}
	return m
}
