//go:build serveronly

package main

import (
	"crypto/tls"
	"fmt"
	"github.com/alecthomas/kong"
)

var cli struct {
	Port        string `help:"Port number." short:"p" default:"8089"`
	Locale      string `help:"Language setting to be used." short:"l"`
	Environment string `help:"Application environment (development|production)." short:"e" default:"development"`
}

func (a *application) get_going() {

	kong.Parse(&cli,
		kong.Name("gotalk-server"),
		kong.Description("An instant chat server."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	a.config.server = true
	a.config.port = cli.Port
	a.config.env = cli.Environment

	if len(cli.Locale) > 0 {
		a.config.locale = cli.Locale
	}

	if portOK(a.config.port) {
		if a.config.port[0] != ':' {
			a.config.port = ":" + a.config.port
		}

		ch := make(chan ClientInput)

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
		fmt.Println(a.lang.Lookup(a.config.locale, "Error in port number"))
	}
}
