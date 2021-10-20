package main

import (
	"log"
	"os"

	"github.com/Xuanwo/go-locale"
)

/*
simple ad-hoc multi user communication program. communication is secured by tls over tcp.
the program can start in server mode or in client mode. Client is GUI using fyne.io as a graphics toolkit

*/

// TODO: multi-room
// TODO: server: intercept signal and exit confirmation dialogue
// TODO: client web interface (React)
// TODO: server admin gui
// TODO: server admin web interface
// TODO: mobile versions (iOS / Android)
// TODO: login / user management
// TODO: IAM

func main() {

	//set actual locale to system default which can be overridden by cli flags
	tag, err := locale.Detect()
	appConfig := config{
		newline: newLine(),
	}

	if err != nil {
		log.Fatal(err)
		appConfig.locale = "en"
	} else {
		if len(tag.String()) > 2 {
			appConfig.locale = tag.String()[:2]
		} else {
			if len(tag.String()) == 2 {
				appConfig.locale = tag.String()
			}
		}
	}
	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime)
	a := &application{
		logger: logger,
		config: appConfig,
	}
	a.initLocalization()
	a.get_going()
}
