package wallet

import (
	"testing"

	"github.com/qustavo/go-wallet/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalletAddrManager(t *testing.T) {
	w, err := NewWallet(
		"wpkh([00000000/84'/0'/0'/0]zpub6u4KbU8TSgNuZSxzv7HaGq5Tk361gMHdZxnM4UYuwzg5CMLcNytzhobitV4Zq6vWtWHpG9QijsigkxAzXvQWyLRfLq1L7VxPP1tky1hPfD4/*)",
		script.Mainnet,
	)
	require.NoError(t, err)

	addrs := []string{
		"bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu",
		"bc1qnjg0jd8228aq7egyzacy8cys3knf9xvrerkf9g",
		"bc1qp59yckz4ae5c4efgw2s5wfyvrz0ala7rgvuz8z",
	}

	for i, addr := range addrs {
		t.Run("", func(t *testing.T) {
			child, err := w.Child(uint32(i))
			require.NoError(t, err)

			assert.Equal(t, addr, child.Address())
		})
	}
}
