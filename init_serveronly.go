// +build serveronly

package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	language "github.com/moemoe89/go-localization"
)

// Initialize localization environment, if localization file is not present, create one by downloading it from github
func initLocalization() error {
	if !fileExists(LANGFILE) {
		fileUrl := RAWFILE
		err := GetFileFromGithub(LANGFILE, fileUrl)
		if err != nil {
			panic(err)
		}
	}
	var err error
	cfg := language.New()
	cfg.BindPath(LANGFILE)
	cfg.BindMainLocale("en")
	lang, err = cfg.Init()
	if err != nil {
		panic(err)
	}
	return err
}

// check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// download a file from github (raw format)
func GetFileFromGithub(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

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

// set newline representation for wither linux or windows systems
func (n *Newline) Init() {
	if runtime.GOOS == "windows" {
		n.nl = "\r\n"
	} else {
		n.nl = "\n"
	}
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
