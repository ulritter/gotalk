package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	language "github.com/moemoe89/go-localization"
)

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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

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

func printUsage(appname string) {
	fmt.Printf("Usage: %s server [<port>] or\n", appname)
	fmt.Printf("Usage: %s client [<nickname> [<address>] [<port>]] \n", appname)
}

func checkArgs(whoami *WhoAmI) error {
	// TODO: beautify parameter handling

	whoami.server = false
	whoami.addr = "localhost"
	whoami.port = ":8089"
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
