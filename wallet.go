package wallet

import (
	"fmt"

	"github.com/qustavo/go-wallet/script"
)

type Wallet struct {
	desc    string
	script  *script.Script
	network script.Network
	child   uint32
}

func NewWallet(desc string, net script.Network) (*Wallet, error) {
	return newWallet(desc, net, "")
}

func newWallet(desc string, net script.Network, path string) (*Wallet, error) {
	script, err := script.ParseWithPath(desc, path)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		desc:    desc,
		script:  script,
		network: net,
	}, nil

}

func (m *Wallet) Address() string {
	return m.script.Address(m.network)
}

func (m *Wallet) Child(i uint32) (*Wallet, error) {
	return newWallet(m.desc, m.network, fmt.Sprintf("m/%d", i))
}
