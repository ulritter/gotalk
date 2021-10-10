package main

import (
	"log"

	"github.com/Xuanwo/go-locale"
)

/*
simple ad-hoc multi user communication program. communication is secured by tls over tcp.
the program can start in server mode or in client mode. Client is GUI using fyne.io as a graphics toolkit

*/

// TODO: improve command line parameter handling, use map for parser
// TODO: make it multi-room

func main() {

	//set actual locale to system default which can be overridden by cli flags
	tag, err := locale.Detect()

	if err != nil {
		log.Fatal(err)
		actualLocale = "en"
	} else {
		if len(tag.String()) > 2 {
			actualLocale = tag.String()[:2]
		} else {
			if len(tag.String()) == 2 {
				actualLocale = tag.String()
			}
		}
	}

	get_going()
}
