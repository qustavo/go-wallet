package script

import (
	"sort"

	"github.com/btcsuite/btcutil/base58"
)

type Network int

const (
	Mainnet Network = iota
	Testnet
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
}

type Script interface {
	Bytes() []byte
	Address(Network) string
}

type p2Sh struct {
	hash160 []byte
}

func Sh(s Script) Script {
	// Execute the nested script
	script := s.Bytes()

	return &p2Sh{
		hash160: Hash160(script),
	}
}

func (s *p2Sh) Bytes() []byte {
	return NewBytes(
		[]byte{OP_HASH160, OP_PUSH_BYTES(20)},
		s.hash160,
		[]byte{OP_EQUAL},
	)
}

func (s *p2Sh) Address(net Network) string {
	return base58.CheckEncode(s.hash160, networks[net].p2sh)
}

type p2Wsh struct {
	sha256 []byte
}

func Wsh(s Script) Script {
	return &p2Wsh{
		sha256: Sha256(s.Bytes()),
	}
}

func (s *p2Wsh) Bytes() []byte {
	return NewBytes(
		[]byte{OP_0, OP_PUSH_BYTES(32)},
		s.sha256,
	)
}

func (s *p2Wsh) Address(net Network) string {
	addr, err := encodeSegWitAddress(networks[net].bech32, 0x00, s.sha256)
	if err != nil {
		panic(err)
	}
	return addr
}

type p2Pkh struct {
	hash160 []byte
}

func Pkh(key Key) Script {
	return &p2Pkh{
		hash160: Hash160(key.Bytes()),
	}
}

func (s *p2Pkh) Bytes() []byte {
	return NewBytes(
		[]byte{OP_DUP, OP_HASH160, OP_PUSH_BYTES(20)},
		Hash160(s.hash160),
		[]byte{OP_EQUALVERIFY, OP_CHECKSIG},
	)
}

func (s *p2Pkh) Address(net Network) string {
	return base58.CheckEncode(s.hash160, networks[net].p2pkh)
}

type p2Wpkh struct {
	hash160 []byte
}

func Wpkh(key Key) Script {
	return &p2Wpkh{
		hash160: Hash160(key.Bytes()),
	}
}

func (s *p2Wpkh) Bytes() []byte {
	return NewBytes(
		[]byte{OP_0, OP_PUSH_BYTES(20)},
		s.hash160,
	)
}

func (s *p2Wpkh) Address(net Network) string {
	// TODO: refactor encodeSegWitAddress so that we can compute the computed
	// bits on Wpkh constructor and this method never returns error.
	addr, err := encodeSegWitAddress(networks[net].bech32, 0x00, s.hash160)
	if err != nil {
		panic(err)
	}
	return addr
}

type multi struct {
	m    int
	keys []Key
}

func Multi(m int, keys []Key) Script {
	return &multi{m, keys}
}

func Sortedmuti(m int, keys []Key) Script {
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})
	return &multi{m, keys}
}

func (s *multi) Bytes() []byte {
	pushedKeys := []byte{}
	for _, key := range s.keys {
		pushedKey := NewBytes(
			[]byte{OP_PUSH_BYTES(33)},
			key.Bytes(),
		)
		pushedKeys = append(pushedKeys, pushedKey...)
	}

	return NewBytes(
		[]byte{OP_N(s.m)}, // required keys
		pushedKeys,        // [ 33 <key_1> ... 33 <key_N> ]
		[]byte{
			OP_N(len(s.keys)), // total keys
			OP_CHECKMULTISIG,
		},
	)
}

func (m *multi) Address(Network) string { return "" }

type Tree interface {
}
