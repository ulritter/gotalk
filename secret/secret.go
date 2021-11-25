package secret

import (
	"io/ioutil"
	"log"
)

func ServerKey(sk string) string {
	key, err := ioutil.ReadFile(sk)
	if err != nil {
		log.Fatal(err)
	}
	return string(key)
}

func RootCert(rc string) string {
	cert, err := ioutil.ReadFile(rc)
	if err != nil {
		log.Fatal(err)
	}
	return string(cert)
}
