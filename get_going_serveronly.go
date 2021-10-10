//go:build serveronly

package main

import (
	"crypto/tls"
	"fmt"
	"github.com/alecthomas/kong"
	"log"
)

var cli struct {
	Port   string `help:"Port number." short:"p" default:"8089"`
	Locale string `help:"Language setting to be used." short:"l" `
}

func get_going() {

	nl := Newline{}
	nl.Init()

	whoami := WhoAmI{}

	kong.Parse(&cli,
		kong.Name("gotalk-server"),
		kong.Description("An instant chat server."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	whoami.server = true
	whoami.port = cli.Port

	if len(cli.Locale) > 0 {
		actualLocale = cli.Locale
	}

	if portOK(whoami.port) {
		if whoami.port[0] != ':' {
			whoami.port = ":" + whoami.port
		}

		ch := make(chan ClientInput)

		go handleServerSession(ch, nl)
		cer, err := tls.X509KeyPair([]byte(rootCert), []byte(serverKey))
		config := &tls.Config{Certificates: []tls.Certificate{cer}}
		if err != nil {
			log.Fatal(err)
		}
		err = startServer(ch, config, whoami.port, nl)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(lang.Lookup(actualLocale, "Error in port number"))
	}
}
