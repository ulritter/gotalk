//go:build !serveronly

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/alecthomas/kong"
)

var cli struct {
	Client struct {
		Address     string `help:"IP address or domain name." short:"a" default:"localhost"`
		Port        string `help:"Port number." short:"p" default:"8089"`
		Nick        string `help:"Nickname to be used." short:"n" default:"J_Doe"`
		Locale      string `help:"Language setting to be used." short:"l" `
		Environment string `help:"Application environment (development|production)." short:"e" default:"development"`
	} `cmd:"" help:"Start gotalk client."`

	Server struct {
		Port        string `help:"Port number." short:"p" default:"8089"`
		Locale      string `help:"Language setting to be used." short:"l"`
		Environment string `help:"Application environment (development|production)." short:"e" default:"development"`
	} `cmd:"" help:"Start gotalk server."`
}

func (a *application) get_going() {

	ctx := kong.Parse(&cli,
		kong.Name("gotalk"),
		kong.Description("An instant chat server and client."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	switch ctx.Command() {
	case "client":
		a.config.server = false
		a.config.addr = cli.Client.Address
		a.config.port = cli.Client.Port
		if len(cli.Client.Locale) > 0 {
			a.config.locale = cli.Client.Locale
		}
		a.config.env = cli.Client.Environment
	case "server":
		a.config.server = true
		a.config.port = cli.Server.Port
		if len(cli.Server.Locale) > 0 {
			a.config.locale = cli.Server.Locale
		}
		a.config.env = cli.Server.Environment
	}

	if portOK(a.config.port) {

		if a.config.port[0] != ':' {
			a.config.port = ":" + a.config.port
		}

		ch := make(chan ClientInput)

		if a.config.server {
			go a.handleServerSession(ch)
			cer, err := tls.X509KeyPair([]byte(rootCert), []byte(serverKey))
			config := &tls.Config{Certificates: []tls.Certificate{cer}}
			if err != nil {
				a.logger.Fatal(err)
			}
			err = a.startServer(ch, config, a.config.port)
			if err != nil {
				a.logger.Fatal(err)
			}
		} else {
			roots := x509.NewCertPool()
			ok := roots.AppendCertsFromPEM([]byte(rootCert))
			if !ok {
				a.logger.Fatal(a.lang.Lookup(a.config.locale, "Failed to parse root certificate"))
			}
			config := &tls.Config{RootCAs: roots, InsecureSkipVerify: true}
			connect := a.config.addr + a.config.port
			a.handleClientSession(connect, config, cli.Client.Nick)
		}
	} else {
		fmt.Println(a.lang.Lookup(a.config.locale, "Error in port number"))
	}
}
