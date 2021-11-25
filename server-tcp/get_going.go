package main

import (
	"crypto/tls"
	"fmt"
	"gotalk/models"
	"gotalk/secret"
	"gotalk/utils"
	"os"
	"path"

	"github.com/alecthomas/kong"
)

var cli struct {
	Port        string `help:"Port number." short:"p" default:"8089"`
	Locale      string `help:"Language setting to be used." short:"l"`
	RootCert    string `help:"Path to root certificate for TLS." short:"c" default:"./root_cert.pem"`
	ServerKey   string `help:"Path to server key for TLS." short:"k" default:"./server.key"`
	Environment string `help:"Application environment (development|production)." short:"e" default:"development"`
	Version     bool   `help:"Show Version." short:"v"`
}

func get_going(a *models.Application) {

	kong.Parse(&cli,
		kong.Name(os.Args[0]),
		kong.Description("An instant chat server."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	if cli.Version {
		fmt.Printf(
			a.Config.Newline+
				"%s (server) "+a.Lang.Lookup(a.Config.Locale, "version")+
				": %s"+
				a.Config.Newline+a.Config.Newline,
			path.Base(os.Args[0]), a.Version)
		return
	}

	a.Config.Server = true
	a.Config.Port = cli.Port
	a.Config.Env = cli.Environment

	if len(cli.Locale) > 0 {
		a.Config.Locale = cli.Locale
	}

	if utils.PortOK(a.Config.Port) {
		if a.Config.Port[0] != ':' {
			a.Config.Port = ":" + a.Config.Port
		}

		a.Config.Ch = make(chan models.ClientInput)

		go handleServerSession(a)
		cer, err := tls.X509KeyPair([]byte(secret.RootCert(cli.RootCert)), []byte(secret.ServerKey(cli.ServerKey)))
		a.Config.TLSconfig = &tls.Config{Certificates: []tls.Certificate{cer}}
		if err != nil {
			a.Logger.Fatal(err)
		}
		err = startServer(a)
		if err != nil {
			a.Logger.Fatal(err)
		}
	} else {
		fmt.Println(a.Lang.Lookup(a.Config.Locale, "Error in port number"))
	}
}
