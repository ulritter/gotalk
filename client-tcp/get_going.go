package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gotalk/models"
	"gotalk/secret"
	"gotalk/utils"
	"os"
	"path"

	"github.com/alecthomas/kong"
)

var cli struct {
	Address     string `help:"IP address or domain name." short:"a" default:"localhost"`
	Port        string `help:"Port number." short:"p" default:"8089"`
	Locale      string `help:"Language setting to be used." short:"l" `
	RootCert    string `help:"Path to root certificate for TLS." short:"c" default:"./root_cert.pem"`
	Nick        string `help:"Nickname to be used." short:"n" default:"J_Doe"`
	Environment string `help:"Application environment (development|production)." short:"e" default:"development"`
	Version     bool   `help:"Show Version." short:"v"`
}

func get_going(a *models.Application) {

	kong.Parse(&cli,
		kong.Name(os.Args[0]),
		kong.Description("An instant chat client."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)
	if cli.Version {
		fmt.Printf(
			a.Config.Newline+
				"%s (client) "+a.Lang.Lookup(a.Config.Locale, "version")+
				": %s"+
				a.Config.Newline+a.Config.Newline,
			path.Base(os.Args[0]), a.Version)
		return
	}
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
		ok := roots.AppendCertsFromPEM([]byte(secret.RootCert(cli.RootCert)))
		if !ok {
			a.Logger.Fatal(a.Lang.Lookup(a.Config.Locale, "Failed to parse root certificate"))
		}
		a.Config.TLSconfig = &tls.Config{RootCAs: roots, InsecureSkipVerify: true}

		handleClientSession(a, cli.Nick)

	} else {
		fmt.Println(a.Lang.Lookup(a.Config.Locale, "Error in port number"))
	}
}
