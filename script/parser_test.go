package script_test

import (
	"testing"

	. "github.com/qustavo/go-wallet/script"
	"github.com/stretchr/testify/require"
)

func TestScriptParser(t *testing.T) {
	testCases := []struct {
		name         string
		script       string
		expectedAddr string
	}{
		{
			name:         "P2PKH",
			script:       "pkh(03aaeb52dd7494c361049de67cc680e83ebcbbbdbeb13637d92cd845f70308af5e)",
			expectedAddr: "1LqBGSKuX5yYUonjxT5qGfpUsXKYYWeabA",
		},
		{
			name:         "P2SH-P2PWPKH",
			script:       "sh(wpkh(039b3b694b8fc5b5e07fb069c783cac754f5d38c3e08bed1960e31fdb1dda35c24))",
			expectedAddr: "37VucYSaXLCAsxYyAPfbSi9eh4iEcbShgf",
		},
		{
			name:         "P2WPKH",
			script:       "wpkh(0330d54fd0dd420a6e5f8d3624f5f3482cae350f79d5f0753bf5beef9c2d91af3c)",
			expectedAddr: "bc1qcr8te4kr609gcawutmrza0j4xv80jy8z306fyu",
		},
		{
			name: "P2WSH-Multi-not-sorted",
			script: `
				wsh(multi(2,
					03a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff,
					0375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c,
					03c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f880
				))
			`,
			expectedAddr: "bc1qwhahvweerhg22ghn8ssqjl5e6r6rjj92jhca266ccmxts65840ks3pu0dp",
		},
		{
			name: "P2WSH-Multi-sorted",
			script: `
				wsh(multi(2,
					0375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c,
					03a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff,
					03c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f880
				))
			`,
			expectedAddr: "bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej",
		},
		{
			name: "P2WSH-Sortedmulti",
			script: `
				wsh(sortedmulti(2,
					0375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c,
					03a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff,
					03c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f880
				))
				`,
			expectedAddr: "bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			script, err := New(test.script)
			require.NoError(t, err)

			addr := script.Address(Mainnet)
			require.Equal(t,
				test.expectedAddr,
				addr,
			)
		})
	}
}
