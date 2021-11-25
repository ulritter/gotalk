package secret

import (
	"io/ioutil"
	"log"
)

func GetKey(keyFile string) string {
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatal(err)
	}
	return string(key)
}
