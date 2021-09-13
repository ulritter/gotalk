package main

import (
	"fmt"
	"os"
)

func printUsage(appname string) {
	fmt.Printf("Usage: %s server [<port>] or\n", appname)
	fmt.Printf("Usage: %s client [<nickname> [<address>] [<port>]] \n", appname)
}

func checkArgs(whoami *WhoAmI) error {
	// TODO: beautify parameter handling

	whoami.server = false
	whoami.addr = "localhost"
	whoami.port = ":8080"
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
