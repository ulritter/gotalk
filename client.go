package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func clientDialogHandling(connect string, nick string) {
	c, err := net.Dial("tcp", connect)
	if err != nil {
		fmt.Println(err)
		return
	}
	// TODO: send user name nick
	fmt.Printf("Connected to: %s, Nickname: %s\n", connect, nick)
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			c.Close()
			return
		}
	}
}
