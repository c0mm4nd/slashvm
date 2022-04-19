package opcodes

func ToOpcodes(bytes []byte) []Opcode {
	ops := make([]Opcode, 0, len(bytes))
	for i := range bytes {
		ops = append(ops, Opcode(bytes[i]))
	}
	return ops
}

func ToBytes(opcodes []Opcode) []byte {
	data := make([]byte, 0, len(opcodes))
	for i := range opcodes {
		data = append(data, byte(opcodes[i]))
	}
	return data
}
