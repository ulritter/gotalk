package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gotalk/models"
	"gotalk/secret"
	"gotalk/utils"

	"github.com/alecthomas/kong"
)

var cli struct {
	Address     string `help:"IP address or domain name." short:"a" default:"localhost"`
	Port        string `help:"Port number." short:"p" default:"8089"`
	Nick        string `help:"Nickname to be used." short:"n" default:"J_Doe"`
	Locale      string `help:"Language setting to be used." short:"l" `
	Environment string `help:"Application environment (development|production)." short:"e" default:"development"`
}

func get_going(a *models.Application) {

	kong.Parse(&cli,
		kong.Name("gotalk"),
		kong.Description("An instant chat client."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)

	a.Config.Server = false
	a.Config.Addr = cli.Address
	a.Config.Port = cli.Port
	if a.Config.Port[0] != ':' {
		a.Config.Port = ":" + a.Config.Port
	}
	if len(cli.Locale) > 0 {
		a.Config.Locale = cli.Locale
	}
	a.Config.Env = cli.Environment
	if utils.PortOK(a.Config.Port) {
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM([]byte(secret.RootCert))
		if !ok {
			a.Logger.Fatal(a.Lang.Lookup(a.Config.Locale, "Failed to parse root certificate"))
		}
		a.Config.TLSconfig = &tls.Config{RootCAs: roots, InsecureSkipVerify: true}

		handleClientSession(a, cli.Nick)

	} else {
		fmt.Println(a.Lang.Lookup(a.Config.Locale, "Error in port number"))
	}
}
