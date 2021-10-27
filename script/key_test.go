package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXPubDerivations(t *testing.T) {
	// Private key generated out of the BIP39 mnemonic seed:
	// `abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about`
	// "xprv9s21ZrQH143K3GJpoapnV8SFfukcVBSfeCficPSGfubmSFDxo1kuHnLisriDvSnRRuL2Qrg5ggqHKNVpxR86QEC8w35uxmGoggxtQTPvfUu"

	testCases := []struct {
		name         string
		key          string
		expectedKeys []string
	}{
		{
			name: "BIP44",
			key:  "xpub6ELHKXNimKbxMCytPh7EdC2QXx46T9qLDJWGnTraz1H9kMMFdcduoU69wh9cxP12wDxqAAfbaESWGYt5rREsX1J8iR2TEunvzvddduAPYcY",
			expectedKeys: []string{
				"03aaeb52dd7494c361049de67cc680e83ebcbbbdbeb13637d92cd845f70308af5e",
				"02dfcaec532010d704860e20ad6aff8cf3477164ffb02f93d45c552dadc70ed24f",
				"0338994349b3a804c44bbec55c2824443ebb9e475dfdad14f4b1a01a97d42751b3",
			},
		},
		{
			name: "BIP49",
			key:  "ypub6Ynvx7RLNYgWzFGM8aeU43hFNjTh7u5Grrup7Ryu2nKZ1Y8FWKaJZXiUrkJSnMmGVNBoVH1DNDtQ32tR4YFDRSpSUXjjvsiMnCvoPHVWXJP/42/*",
			expectedKeys: []string{
				"021c4be1736ca2f364962244bba47d54dd569daefcb522d5116df82c56903dc599",
				"02c1440e470a04822ea0731e92b570c57cc1f241f6d5fb11c7ad496a3217a7bf70",
				"03eb2430b67655df8b1ffb497fa4a92d28cea4698f7690dd457f09a37fa52e7be6",
			},
		},
		{
			name: "BIP84",
			key:  "zpub6u4KbU8TSgNuZSxzv7HaGq5Tk361gMHdZxnM4UYuwzg5CMLcNytzhobitV4Zq6vWtWHpG9QijsigkxAzXvQWyLRfLq1L7VxPP1tky1hPfD4/1/2/3/*",
			expectedKeys: []string{
				"0316c81f32166ce834e22095b843728ee6d4483f9f3035ff03c353921ff42fa095",
				"03cf40b9f31ab1f2b5fc30bb81f28d50587ad28841e35bbd92f7ab48820678730b",
				"03984d879ee70ffa31c6b2bf6d44f7fbc380ad1e2cc410a14d72cd2a994daf2bf0",
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
