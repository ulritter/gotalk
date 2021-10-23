package utils

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

func NewLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	} else {
		return "\n"
	}
}

// check if a file exists
func FileExists(filename string) bool {
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

func PortOK(p string) bool {

	if p[0] == ':' && len(p) > 1 {
		p = p[1:]
	}
	if len(p) >= 1 {
		if _, err := strconv.Atoi(p); err == nil {
			return true
		}
	} else {
		return false
	}
	return false
}
