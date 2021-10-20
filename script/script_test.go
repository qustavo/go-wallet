package script_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/qustavo/go-wallet/script"
)

func TestScripts(t *testing.T) {
	newKey := func(s string) Key {
		require.NotEmpty(t, s)

		key, err := NewPubKey(s)
		require.NoError(t, err)

		return key
	}

	testCases := []struct {
		name         string
		script       Script
		expectedAddr string
		network      Network
	}{
		{
			name: "P2PKH",
			script: Pkh(
				newKey("03aaeb52dd7494c361049de67cc680e83ebcbbbdbeb13637d92cd845f70308af5e"),
			),
			expectedAddr: "1LqBGSKuX5yYUonjxT5qGfpUsXKYYWeabA",
		},
		{
			// This example has been taken from BDK docs
			name: "P2PKH (BDK)",
			script: Pkh(
				newKey("02e96fe52ef0e22d2f131dd425ce1893073a3c6ad20e8cac36726393dfb4856a4c"),
			),
			expectedAddr: "mrkwtj5xpYQjHeJe5wsweNjVeTKkvR5fCr",
			network:      Testnet,
		},
		{
			name: "P2SH-P2WPKH",
			script: Sh(Wpkh(
				newKey("039b3b694b8fc5b5e07fb069c783cac754f5d38c3e08bed1960e31fdb1dda35c24"),
			)),
			expectedAddr: "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf",
		},
		{
			name: "P2WPKH",
			script: Wpkh(
				newKey("0330d54fd0dd420a6e5f8d3624f5f3482cae350f79d5f0753bf5beef9c2d91af3c"),
			),
			expectedAddr: "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu",
		},
		{
			name: "P2WSH-Multi(not-sorted)",
			script: Wsh(Multi(2, []Key{
				newKey("03a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff"),
				newKey("0375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c"),
				newKey("03c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f880"),
			})),
			expectedAddr: "bc1qwhahvweerhg22ghn8ssqjl5e6r6rjj92jhca266ccmxts65840ks3pu0dp",
		},
		{
			name: "P2WSH-Multi(sorted)",
			script: Wsh(Multi(2, []Key{
				newKey("0375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c"),
				newKey("03a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff"),
				newKey("03c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f880"),
			})),
			expectedAddr: "bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej",
		},
		{
			name: "P2WSH-Sortedmulti",
			script: Wsh(Sortedmuti(2, []Key{
				newKey("0375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c"),
				newKey("03a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff"),
				newKey("03c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f880"),
			})),
			expectedAddr: "bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			net := test.network
			got := test.script.Address(net)
			require.Equal(t, test.expectedAddr, got)
		})
	}
}
