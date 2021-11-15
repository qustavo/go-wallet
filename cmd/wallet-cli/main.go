package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/urfave/cli"

	"github.com/qustavo/go-wallet"
	"github.com/qustavo/go-wallet/script"
)

type Flags struct {
	Net  script.Network
	Path string
}

func (f *Flags) Args() []string {
	return flag.Args()
}

func newAddress(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return errors.New("`newaddress` requires exactly 1 argument")
	}
	desc := ctx.Args()[0]

	var net script.Network
	switch ctx.String("network") {
	case "mainnet":
		net = script.Mainnet
	case "testnet":
		net = script.Testnet
	case "regtest":
		net = script.Regtest
	default:
		return fmt.Errorf("net '%s' is invalid", ctx.String("network"))
	}

	basePath := "m/"
	// change addresses uses 1 in the path
	if !ctx.Bool("change") {
		basePath += "0"
	} else {
		basePath += "1"
	}
	offset := ctx.Uint("offset")

	for i := uint(0); i < ctx.Uint("num"); i++ {
		w, err := wallet.NewWallet(desc, net)
		if err != nil {
			return err
		}
		path := fmt.Sprintf("%s/%d", basePath, offset+i)
		w, err = w.Path(path)
		if err != nil {
			return err
		}

		fmt.Printf("%s: %s\n", path, w.Address())
	}

	return nil
}

func main() {
	(&cli.App{
		Name: "wallet-cli",
		Commands: []cli.Command{
			{
				Name:        "newaddrs",
				Usage:       "Generates new wallet addresses",
				ArgsUsage:   "<descriptor>",
				Description: "`newaddrs` generates address given a descriptor.",
				Flags: []cli.Flag{
					cli.StringFlag{Name: "network", Usage: "Sets the Bitcoin network [mainnet|testnet|regtest]", Value: "mainnet"},
					cli.UintFlag{Name: "num", Usage: "How many addresses to generate", Value: 10},
					cli.UintFlag{Name: "offset"},
					cli.BoolFlag{Name: "change", Usage: "Generate a change address"},
				},
				Action: newAddress,
			},
		},
	}).RunAndExitOnError()
}
