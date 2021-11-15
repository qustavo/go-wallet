package script

import (
	"sort"

	"github.com/btcsuite/btcutil/base58"
)

type Network int

const (
	Mainnet Network = iota
	Testnet
	Regtest
)

var DefaultNetwork = Mainnet

type netParams struct {
	p2pkh  byte
	p2sh   byte
	bech32 string
}

var networks = map[Network]netParams{
	Mainnet: {
		p2pkh:  0x00,
		p2sh:   0x05,
		bech32: "bc",
	},
	Testnet: {
		p2pkh:  0x6F,
		p2sh:   0xC4,
		bech32: "tb",
	},
	Regtest: {
		p2pkh:  0x6F,
		p2sh:   0xC4,
		bech32: "bcrt",
	},
}

type Script struct {
	bytes  []byte
	addrFn func(Network) string
}

func (s *Script) Bytes() []byte {
	return s.bytes
}

func (s *Script) Address(net Network) string {
	if s.addrFn == nil {
		return "<not implemented>"
	}
	return s.addrFn(net)
}

type ScriptExpr interface {
	Eval() (*Script, error)
}

type p2Sh struct {
	expr ScriptExpr
}

func Sh(expr ScriptExpr) ScriptExpr {
	return &p2Sh{
		expr: expr,
	}
}

func (s *p2Sh) Eval() (*Script, error) {
	eval, err := s.expr.Eval()
	if err != nil {
		return nil, err
	}

	hash160 := Hash160(eval.Bytes())
	return &Script{
		bytes: NewBytes(
			[]byte{OP_HASH160, OP_PUSH_BYTES(20)},
			hash160,
			[]byte{OP_EQUAL},
		),
		addrFn: func(net Network) string {
			return base58.CheckEncode(hash160, networks[net].p2sh)
		},
	}, nil
}

type p2Wsh struct {
	expr ScriptExpr
}

func Wsh(expr ScriptExpr) ScriptExpr {
	return &p2Wsh{
		expr: expr,
	}
}

func (s *p2Wsh) Eval() (*Script, error) {
	eval, err := s.expr.Eval()
	if err != nil {
		return nil, err
	}

	hash256 := Sha256(eval.Bytes())

	return &Script{
		bytes: NewBytes(
			[]byte{OP_0, OP_PUSH_BYTES(32)},
			hash256,
		),
		addrFn: func(net Network) string {
			addr, err := encodeSegWitAddress(networks[net].bech32, 0x00, hash256)
			if err != nil {
				panic(err)
			}
			return addr
		},
	}, nil
}

type p2Pkh struct {
	key string
}

func Pkh(key string) ScriptExpr {
	return &p2Pkh{key: key}
}

func (s *p2Pkh) Eval() (*Script, error) {
	key, err := NewPubKey(s.key)
	if err != nil {
		return nil, err
	}

	hash160 := Hash160(key.Bytes())
	script := &Script{
		bytes: NewBytes(
			[]byte{OP_DUP, OP_HASH160, OP_PUSH_BYTES(20)},
			hash160,
			[]byte{OP_EQUALVERIFY, OP_CHECKSIG},
		),
		addrFn: func(net Network) string {
			return base58.CheckEncode(hash160, networks[net].p2pkh)
		},
	}

	return script, nil
}

type p2Wpkh struct {
	key string
}

func Wpkh(key string) ScriptExpr {
	return &p2Wpkh{key: key}
}

func (s *p2Wpkh) Eval() (*Script, error) {
	key, err := NewPubKey(s.key)
	if err != nil {
		return nil, err
	}

	hash160 := Hash160(key.Bytes())
	script := &Script{
		bytes: NewBytes(
			[]byte{OP_0, OP_PUSH_BYTES(20)},
			hash160,
		),
		addrFn: func(net Network) string {
			// TODO: refactor encodeSegWitAddress so that we can compute the computed
			// bits on Wpkh constructor and this method never returns error.
			addr, err := encodeSegWitAddress(networks[net].bech32, 0x00, hash160)
			if err != nil {
				panic(err)
			}
			return addr
		},
	}

	return script, nil
}

type multi struct {
	m    int
	keys []string
}

func Multi(m int, keys ...string) ScriptExpr {
	return &multi{m, keys}
}

func Sortedmulti(m int, keys ...string) ScriptExpr {
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return &multi{m, keys}
}

func (s *multi) Eval() (*Script, error) {
	// Convert input keys from string into PubKey.
	keys := make([]Key, len(s.keys))
	for i, str := range s.keys {
		key, err := NewPubKey(str)
		if err != nil {
			return nil, err
		}
		keys[i] = key
	}

	pushedKeys := []byte{}
	for _, key := range keys {
		pushedKey := NewBytes(
			[]byte{OP_PUSH_BYTES(33)},
			key.Bytes(),
		)
		pushedKeys = append(pushedKeys, pushedKey...)
	}

	return &Script{
		bytes: NewBytes(
			[]byte{OP_N(s.m)}, // required keys
			pushedKeys,        // [ 33 <key_1> ... 33 <key_N> ]
			[]byte{
				OP_N(len(s.keys)), // total keys
				OP_CHECKMULTISIG,
			},
		),
		addrFn: func(net Network) string { return "" },
	}, nil
}

type Tree interface {
}

type tr struct {
	key  string
	tree Tree
}

func Tr(key string, tree Tree) ScriptExpr {
	return &tr{
		key:  key,
		tree: tree,
	}
}

func (s *tr) Eval() (*Script, error) {
	return &Script{}, nil
}
