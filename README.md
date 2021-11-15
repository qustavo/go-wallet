# go-wallet
[![Build Status](https://github.com/qustavo/go-wallet/actions/workflows/test.yml/badge.svg)](https://github.com/qustavo/go-wallet/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/qustavo/go-wallet.svg)](https://pkg.go.dev/github.com/qustavo/go-wallet)

A work-in-progress Bitcoin wallet

## Descriptors
`go-wallet` is designed around [Bitcoin Descriptors](https://github.com/bitcoin/bips/blob/master/bip-0380.mediawiki).
It implements a Output Script Descriptors language parser which returns valid outputs scripts.
The script package also implements a Output Script DSL to build scripts with a fluent API.

### Supported scripts

| Script Type |Operator                   | |
|-------------|---------------------------|-|
| P2SH        | `sh(SCRIPT)`              |✓|
| P2WSH       | `wsh(SCRIPT)`             |✓|
| P2PK        | `pk(KEY)`                 |✗|
| P2PKH       | `pkh(KEY)`                |✓|
| P2WPKH      | `wpkh(KEY)`               |✓|
|             | `combo(KEY)`              |✗|
| Multi       | `multi(k,<keys>`          |✓|
| Sortedmulti | `sortedmulti(k,<keys>`    |✓|
| P2TR        | `tr()` or `tr(KEY, TREE)` |✗|
|             | `addr(ADDR)`              |✗|
|             | `hex(HEX)`                |✗|

### Example
The following example shows how generate addresses for an output descriptor using the cli tool:

```bash
$ wallet-cli newaddrs --num=1 "wsh(sortedmulti(2,
	0375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c,
	03a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff,
	03c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f880
))"

m/0/0: bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej
```

Or using the wallet API:

```go
package main

import (
	"fmt"
	"log"

	"github.com/qustavo/go-wallet"
	"github.com/qustavo/go-wallet/script"
)

func main() {
	w, err := wallet.NewWallet(`
		wsh(sortedmulti(2,
			0375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c,
			03a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff,
			03c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f880
	))`, script.Mainnet)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("> addr: %s\n", w.Address())
}
```

The same can be expresses using the script API.
```go
package main

import (
	"fmt"
	"log"

	. "github.com/qustavo/go-wallet/script"
)

func main() {
	script := Wsh(Sortedmulti(2,
		"0375e00eb72e29da82b89367947f29ef34afb75e8654f6ea368e0acdfd92976b7c",
		"03a1b26313f430c4b15bb1fdce663207659d8cac749a0e53d70eff01874496feff",
		"03c96d495bfdd5ba4145e3e046fee45e84a8a48ad05bd8dbb395c011a32cf9f880",
	))
	eval, err := script.Eval()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("> addr: %s\n", eval.Address(Mainnet))
}
```
