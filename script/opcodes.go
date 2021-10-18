package script

const (
	// Push value onto stack
	OP_0     = 0x00
	OP_FALSE = OP_0
	OP_1     = 0x51
	OP_TRUE  = 0x01

	// Stack Operation
	OP_DUP = 0x76

	// Binary arithmetic and conditionals
	OP_EQUAL       = 0x87
	OP_EQUALVERIFY = 0x88

	// Cryptographic and hashing operations
	OP_HASH160       = 0xA9
	OP_CHECKSIG      = 0xAC
	OP_CHECKMULTISIG = 0xAE
)

func OP_PUSH_BYTES(b int) byte {
	if b < 0x01 || b > 0x4B {
		panic("OP_PUSH_BYTES value MUST be between 0x01 and 0x4B")
	}
	return byte(b)
}

func OP_N(n int) byte {
	if n < 0x02 || n > 0x10 {
		panic("OP_N value MUST be between 0x02 and 0x10")
	}
	return byte(0x50 + n)
}
