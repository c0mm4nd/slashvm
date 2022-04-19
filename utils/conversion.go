package utils

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/holiman/uint256"
)

func U256ToHashHex(u256 *uint256.Int) string {
	b := make([]byte, 32)
	for i, u64 := range u256 {
		binary.BigEndian.PutUint64(b[i*8:], u64)
	}

	return hex.EncodeToString(b)
}

func HashHexToU256(hexString string) *uint256.Int {
	b, err := hex.DecodeString(hexString)
	if err != nil {
		panic(err.Error() + ": " + hexString)
	}

	if len(b) != 32 {
		panic("invalid hash hex length: " + hexString)
	}

	u256 := new(uint256.Int)
	for i := 0; i < len(b); i += 8 {
		u256[i/8] = binary.BigEndian.Uint64(b[i : i+8])
	}

	return u256
}
