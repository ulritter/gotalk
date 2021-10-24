package main

import (
	"gotalk/models"
	"gotalk/utils"
	"log"
	"os"
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

// app Config parameters and resources

func main() {

	//set actual locale to system default which can be overridden by cli flags

	appConfig := models.Config{
		Newline: utils.NewLine(),
		Locale:  utils.GetLocale(),
	}

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime)
	a := &models.Application{
		Logger: logger,
		Config: appConfig,
	}
	a.InitLocalization()
	get_going(a)
}
