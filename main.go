package main

/*
simple ad-hoc multi user communication program. communication is secured by tls over tcp.
the program can start in server mode or in client mode. Client is GUI using fyne.io as a graphics toolkit

*/

import (
	"crypto/tls"
	"crypto/x509"
	"log"

	language "github.com/moemoe89/go-localization"
)

// TODO: externalize strings
// TODO: make it multi-room

func init() {
	cfg := language.New()
	cfg.BindPath("./language.json")
	cfg.BindMainLocale("en")
	var err error
	lang, err = cfg.Init()
	if err != nil {
		panic(err)
	}
}

func main() {

	locale = "de"
	nl := Newline{}
	nl.Init()

	whoami := WhoAmI{}

	getParams := checkArgs(&whoami)

	ch := make(chan ClientInput)

	if getParams == nil {
		if whoami.server {
			go handleServerDialog(ch, nl)
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
				log.Fatal(lang.Lookup(locale, "Failed to parse root certificate"))
			}
			config := &tls.Config{RootCAs: roots, InsecureSkipVerify: true}
			connect := whoami.addr + whoami.port
			handleClientDialog(connect, config, whoami.nick, nl)
		}
	}
}
