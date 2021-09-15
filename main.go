package main

import (
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
			err := startServer(ch, whoami.port, nl)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			connect := whoami.addr + whoami.port
			clientDialogHandling(connect, whoami.nick, nl)
		}
	}
}
