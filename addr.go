package wallet

import (
	"fmt"

	"github.com/qustavo/go-wallet/script"
)

type AddrManager struct {
	desc    string
	script  *script.Script
	network script.Network
	child   uint32
}

func NewAddrManager(desc string, net script.Network) (*AddrManager, error) {
	return newAddrManager(desc, net, "")
}

func newAddrManager(desc string, net script.Network, path string) (*AddrManager, error) {
	script, err := script.ParseWithPath(desc, path)
	if err != nil {
		return nil, err
	}

	return &AddrManager{
		desc:    desc,
		script:  script,
		network: net,
	}, nil

}

func (m *AddrManager) Address() string {
	return m.script.Address(m.network)
}

func (m *AddrManager) Child(i uint32) (*AddrManager, error) {
	return newAddrManager(m.desc, m.network, fmt.Sprintf("m/%d", i))
}
