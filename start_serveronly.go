//go:build serveronly

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
)

// print usage message in case of wrong parameters given
func printUsage(appname string) {
	fmt.Printf(lang.Lookup(locale, "Usage:")+" %s  [<port>]+\n", appname)
}

// parse command line arguments
func checkArgs(whoami *WhoAmI) error {

	whoami.server = true
	whoami.addr = "localhost"
	whoami.port = ":8089"
	whoami.nick = "J_Doe"

	arguments := os.Args
	if len(arguments) == 1 {
		return nil
	} else if len(arguments) == 2 {
		whoami.port = arguments[1]
	} else if len(arguments) > 2 {
		printUsage(arguments[0])
		return fmt.Errorf("parameter error")
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
	}
}
