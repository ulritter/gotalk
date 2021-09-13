package main

import (
	"log"
)

func main() {

	whoami := WhoAmI{
		server: false,
		addr:   "localhost",
		port:   "8080",
		nick:   "J_Doe",
	}

	getParams := checkArgs(&whoami)

	var mport string
	if whoami.port[0] != ':' {
		mport = ":" + whoami.port
	} else {
		mport = whoami.port
	}
	ch := make(chan ClientInput)

	if getParams == nil {
		if whoami.server {
			go serverDialogHandling(ch)
			err := startServer(ch, mport)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			connect := whoami.addr + mport
			clientDialogHandling(connect, whoami.nick)
		}
	}
}
