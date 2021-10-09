//go:build !serveronly

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
)

// print usage message in case of wrong parameters given
func printUsage(appname string) {
	fmt.Printf(lang.Lookup(locale, "Usage:"+"%s "+lang.Lookup(locale, "server [<port>] or")+"\n"), appname)
	fmt.Printf(lang.Lookup(locale, "Usage:"+"%s "+lang.Lookup(locale, "client [<nickname> [<address>] [<port>]]")+"\n"), appname)
}

// parse command line arguments
func checkArgs(whoami *WhoAmI) error {
	// TODO: beautify parameter handling

	whoami.server = false
	whoami.addr = "localhost"
	whoami.port = ":8089"
	whoami.nick = "J_Doe"

	arguments := os.Args
	if len(arguments) == 1 {
		printUsage(arguments[0])
		// TODO: error handling
		return fmt.Errorf("parameter error")
	} else if arguments[1] == "client" {
		whoami.server = false
		if len(arguments) >= 3 {
			whoami.nick = arguments[2]
			if len(arguments) >= 4 {
				whoami.addr = arguments[3]
			}
			if len(arguments) == 5 {
				whoami.port = arguments[4]
			} else {
				printUsage(arguments[0])
				// TODO: error handling
				return fmt.Errorf("parameter error")
			}
		}
	} else if arguments[1] == "server" {
		whoami.server = true
		if len(arguments) == 3 {
			whoami.port = arguments[2]
		} else if len(arguments) > 3 {
			printUsage(arguments[0])
			// TODO: error handling
			return fmt.Errorf("parameter error")
		}
	} else {
		printUsage(arguments[0])
		// TODO: error handling
		return fmt.Errorf("parameter error")
	}
	if whoami.port[0] != ':' {
		whoami.port = ":" + whoami.port
	}
	return nil
}

func get_going() {
	locale = "en"
	nl := Newline{}
	nl.Init()

	whoami := WhoAmI{}

	getParams := checkArgs(&whoami)

	ch := make(chan ClientInput)

	if getParams == nil {
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
				log.Fatal(lang.Lookup(locale, "Failed to parse root certificate"))
			}
			config := &tls.Config{RootCAs: roots, InsecureSkipVerify: true}
			connect := whoami.addr + whoami.port
			handleClientSession(connect, config, whoami.nick, nl)
		}
	}
}
