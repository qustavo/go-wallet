package script

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
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
		return nil, fmt.Errorf("invalid key format: %w", err)
	}
	return &PubKey{key}, nil
}

func (pk *PubKey) Bytes() []byte  { return pk.key }
func (pk *PubKey) String() string { return hex.EncodeToString(pk.key) }

type XPub struct {
	key *hdkeychain.ExtendedKey
}

var (
	keyOriginRegexp = regexp.MustCompile("\\[([0-9a-fA-F]{8})(.*)?\\](.+)")
)

// parseKeyOrigin returns the fingerprint and the derivation path (set)
func parseKeyOrigin(s string) (string, string, string, error) {
	submatch := keyOriginRegexp.FindStringSubmatch(s)
	switch len(submatch) {
	case 0:
		return "", "", "", errors.New("invalid key origin")
	case 3:
	case 4:
		return submatch[1], submatch[2], submatch[3], nil
	}
	panic("xxx")
}

func NewXPub(s string) (*XPub, error) {
	// check if key has fingerprint: "[" + <8-byte> + "]".
	if s[0] == '[' {
		_, path, key, err := parseKeyOrigin(s)
		if err != nil {
			return nil, err
		}

		xpub, err := newXPub(key)
		if err != nil {
			return nil, err
		}

		if path != "" {
			return xpub.Derive("m" + path)
		}

		return xpub, nil
	}

	return newXPub(s)
}

func newXPub(s string) (*XPub, error) {
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

	// Hardened levels can be defined as `'`, `h` or `H` so unify them into `'`
	path = strings.Map(func(r rune) rune {
		if r == 'h' || r == 'H' {
			return '\''
		}
		return r
	}, path)

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

func (xpub *XPub) String() string { return xpub.key.String() }
