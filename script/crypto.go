package script

import (
	"crypto/sha256"
	"hash"

	"github.com/btcsuite/btcutil/bech32"
	"golang.org/x/crypto/ripemd160" // nolint:staticcheck // SA1019 ripem160 is deprecated but it is used by Bitcoin
)

func calcHash(buf []byte, hasher hash.Hash) []byte {
	hasher.Write(buf)
	return hasher.Sum(nil)
}

func Sha256(buf []byte) []byte {
	return calcHash(buf, sha256.New())
}

func Hash160(buf []byte) []byte {
	return calcHash(
		Sha256(buf), ripemd160.New(),
	)
}

// Taken from btcutil
func encodeSegWitAddress(hrp string, witnessVersion byte, witnessProgram []byte) (string, error) {
	// Group the address bytes into 5 bit groups, as this is what is used to
	// encode each character in the address string.
	converted, err := bech32.ConvertBits(witnessProgram, 8, 5, true)
	if err != nil {
		return "", err
	}

	// Concatenate the witness version and program, and encode the resulting
	// bytes using bech32 encoding.
	combined := make([]byte, len(converted)+1)
	combined[0] = witnessVersion
	copy(combined[1:], converted)
	bech, err := bech32.Encode(hrp, combined)
	if err != nil {
		return "", err
	}

	return bech, nil
}
