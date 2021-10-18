package script

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/hdkeychain"
)

type Key interface {
	Bytes() []byte
	String() string
}

type PubKey struct {
	key []byte
}

func NewPubKey(s string) (*PubKey, error) {
	key, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return &PubKey{key}, nil
}

func (pk *PubKey) Bytes() []byte  { return pk.key }
func (pk *PubKey) String() string { return hex.EncodeToString(pk.key) }

type XPub struct {
	key *hdkeychain.ExtendedKey
}

func NewXPub(s string) (*XPub, error) {
	key, err := hdkeychain.NewKeyFromString(s)
	if err != nil {
		return nil, err
	}

	return &XPub{key: key}, nil
}

func parsePath(path string, fn func(uint32) error) error {
	if !strings.HasPrefix(path, "m/") {
		return errors.New("xpub: invalid path prefix")
	}
	path = strings.TrimPrefix(path, "m/")
	levels := strings.Split(path, "/")
	for _, level := range levels {
		var v uint32

		// Verify if the level is hardened
		if strings.HasSuffix(level, "'") {
			v = 0x80000000
			level = strings.TrimSuffix(level, "'")
		}

		atoi, err := strconv.Atoi(level)
		if err != nil {
			return err
		}

		v += uint32(atoi)

		if err := fn(v); err != nil {
			return err
		}
	}

	return nil
}

func (xpub *XPub) Derive(path string) (*XPub, error) {
	key := xpub.key
	err := parsePath(path, func(i uint32) error {
		var err error
		key, err = key.Derive(i)
		return err
	})
	if err != nil {
		return xpub, err
	}

	return &XPub{key: key}, nil
}

func (xpub *XPub) Child(i uint32) (Key, error) {
	child, err := xpub.key.Derive(i)
	if err != nil {
		return nil, err
	}

	pub, err := child.ECPubKey()
	if err != nil {
		return nil, err
	}

	return &PubKey{key: pub.SerializeCompressed()}, nil
}

func (xpub *XPub) Key(n uint32) Key {
	return nil
}
