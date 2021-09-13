package main

import (
	"log"
)

func main() {

	whoami := WhoAmI{}
	getParams := checkArgs(&whoami)
	ch := make(chan ClientInput)

	if getParams == nil {
		if whoami.server {
			go serverDialogHandling(ch)
			err := startServer(ch, whoami.port)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			connect := whoami.addr + whoami.port
			clientDialogHandling(connect, whoami.nick)
		}
	}
}
