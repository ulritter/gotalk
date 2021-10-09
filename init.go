package main

import (
	"io"
	"net/http"
	"os"
	"runtime"

	language "github.com/moemoe89/go-localization"
)

// init localization environment
func init() {
	err := initLocalization()
	if err != nil {
		panic(err)
	}
}

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

// set newline representation for wither linux or windows systems
func (n *Newline) Init() {
	if runtime.GOOS == "windows" {
		n.nl = "\r\n"
	} else {
		n.nl = "\n"
	}
}
