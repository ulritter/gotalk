package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// TODO: error handling for whole function
func clientDialogHandling(connect string, nick string) {
	buf := make([]byte, 4096)
	conn, err := net.Dial("tcp", connect)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Connected to: %s, Nickname: %s\n", connect, nick)
	fmt.Fprintf(conn, "$$$"+nick+"$")

	for {
		go func() {
			for { // TODO: error handling
				n, err := conn.Read(buf)
				if err != nil {
					log.Println("Error reading from buffer", err)
					return
				}
				msg := string(buf[:n])
				fmt.Print(msg)
			}
		}()

		ch := make(chan string)
		go func(ch chan string) {
			reader := bufio.NewReader(os.Stdin)
			for {
				s, err := reader.ReadString('\n')
				if err != nil { // Maybe log non io.EOF errors, if you want
					close(ch)
					return
				}
				ch <- s
			}
		}(ch)

	stdinloop:
		for {
			select {
			case stdin, ok := <-ch:
				if !ok {
					break stdinloop
				} else {
					msg := strings.TrimSpace(string(stdin))
					fmt.Fprintln(conn, msg)
					if msg == "STOP" {
						fmt.Println("TCP client exiting...")
						conn.Close()
						return
					}
				}
			case <-time.After(1 * time.Second):
				// Do something when there is nothing read from stdin
			}
		}

	}
}
