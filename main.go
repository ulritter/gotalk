package main

/*
simple ad-hoc multi user communication program. communication is secured by tls over tcp.
the program can start in server mode or in client mode.

Parts of this program are based oh the great golang jump-start tutorial of David Byttow:
https://www.youtube.com/watch?v=oyKrI8e1-9U
*/

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"runtime"
)

// TODO: externalize strings
func main() {

	nl := Newline{}

	if runtime.GOOS == "windows" {
		nl.SetNewLine("\r\n")
	} else {
		nl.SetNewLine("\n")
	}

	whoami := WhoAmI{}

	getParams := checkArgs(&whoami)

	ch := make(chan ClientInput)

	if getParams == nil {
		if whoami.server {

			go serverDialogHandling(ch, nl)
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
				log.Fatal("failed to parse root certificate")
			}
			config := &tls.Config{RootCAs: roots, InsecureSkipVerify: true}
			connect := whoami.addr + whoami.port
			clientDialogHandling(connect, config, whoami.nick, nl)
		}
	}
}
