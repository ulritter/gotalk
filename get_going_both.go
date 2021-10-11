//go:build !serveronly

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"

	"github.com/alecthomas/kong"
)

var cli struct {
	Client struct {
		Address string `help:"IP address or domain name." short:"a" default:"localhost"`
		Port    string `help:"Port number." short:"p" default:"8089"`
		Nick    string `help:"Nickname to be used." short:"n" default:"J_Doe"`
		Locale  string `help:"Language setting to be used." short:"l" `
	} `cmd:"" help:"Start gotalk client."`

	Server struct {
		Port   string `help:"Port number." short:"p" default:"8089"`
		Locale string `help:"Language setting to be used." short:"l" default:"en"`
	} `cmd:"" help:"Start gotalk server."`
}

func get_going() {

	nl := Newline{}
	nl.Init()

	whoami := WhoAmI{}

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
		whoami.server = false
		whoami.addr = cli.Client.Address
		whoami.port = cli.Client.Port
		whoami.nick = cli.Client.Nick
		if len(cli.Client.Locale) > 0 {
			actualLocale = cli.Client.Locale
		}

	case "server":
		whoami.server = true
		whoami.port = cli.Server.Port
		if len(cli.Server.Locale) > 0 {
			actualLocale = cli.Server.Locale
		}
	}

	if portOK(whoami.port) {

		if whoami.port[0] != ':' {
			whoami.port = ":" + whoami.port
		}

		ch := make(chan ClientInput)

		if whoami.server {
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
			roots := x509.NewCertPool()
			ok := roots.AppendCertsFromPEM([]byte(rootCert))
			if !ok {
				log.Fatal(lang.Lookup(actualLocale, "Failed to parse root certificate"))
			}
			config := &tls.Config{RootCAs: roots, InsecureSkipVerify: true}
			connect := whoami.addr + whoami.port
			handleClientSession(connect, config, whoami.nick, nl)
		}
	} else {
		fmt.Println(lang.Lookup(actualLocale, "Error in port number"))
	}
}
