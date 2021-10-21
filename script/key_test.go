package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXPubDerivations(t *testing.T) {
	// Private key generated out of the BIP39 mnemonic seed:
	// `abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about`
	xpriv := "xprv9s21ZrQH143K3GJpoapnV8SFfukcVBSfeCficPSGfubmSFDxo1kuHnLisriDvSnRRuL2Qrg5ggqHKNVpxR86QEC8w35uxmGoggxtQTPvfUu"

	testCases := []struct {
		name         string
		key          string
		expectedKeys []string
	}{
		{
			name: "BIP44",
			key:  "[00000000/44'/0'/0'/0]" + xpriv,
			expectedKeys: []string{
				"03aaeb52dd7494c361049de67cc680e83ebcbbbdbeb13637d92cd845f70308af5e",
				"02dfcaec532010d704860e20ad6aff8cf3477164ffb02f93d45c552dadc70ed24f",
				"0338994349b3a804c44bbec55c2824443ebb9e475dfdad14f4b1a01a97d42751b3",
			},
		},
		{
			name: "BIP49",
			key:  "[00000000/49'/0'/0'/0]" + xpriv,
			expectedKeys: []string{
				"039b3b694b8fc5b5e07fb069c783cac754f5d38c3e08bed1960e31fdb1dda35c24",
				"022a421fa4a65a87d1c3e4238155d85f7bd2c5bb87632f331b5722f110586aa198",
				"02fdbd244eebd701270478af75ebb8894b963d61f2f686e366a626cb200ba13e45",
			},
		},
		{
			name: "BIP84",
			key:  "[00000000/84'/0'/0'/0]" + xpriv,
			expectedKeys: []string{
				"0330d54fd0dd420a6e5f8d3624f5f3482cae350f79d5f0753bf5beef9c2d91af3c",
				"03e775fd51f0dfb8cd865d9ff1cca2a158cf651fe997fdc9fee9c1d3b5e995ea77",
				"038ffea936b2df76bf31220ebd56a34b30c6b86f40d3bd92664e2f5f98488dddfa",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			key, err := NewXPub(test.key)
			require.NoError(t, err)

			for i, expected := range test.expectedKeys {
				child, err := key.Child(uint32(i))
				require.NoError(t, err)
				assert.Equal(t, expected, child.String())
			}
		})
	}
}
