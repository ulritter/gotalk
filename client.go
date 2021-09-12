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
	c, err := net.Dial("tcp", connect)
	if err != nil {
		fmt.Println(err)
		return
	}
	// TODO: send user name nick
	fmt.Printf("Connected to: %s, Nickname: %s\n", connect, nick)
	fmt.Fprintf(c, "$$$"+nick+"$")

	for {

		go func() {
			for { // TODO: error handling
				n, err := c.Read(buf)
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
			/* Uncomment this block to actually read from stdin */
			reader := bufio.NewReader(os.Stdin)
			for {
				s, err := reader.ReadString('\n')
				if err != nil { // Maybe log non io.EOF errors, if you want
					close(ch)
					return
				}
				ch <- s
			}
			// Simulating stdin
			// ch <- "A line of text"
			close(ch)
		}(ch)

	stdinloop:
		for {
			select {
			case stdin, ok := <-ch:
				if !ok {
					break stdinloop
				} else {
					msg := strings.TrimSpace(string(stdin))
					fmt.Fprintln(c, msg)
					if strings.TrimSpace(string(stdin)) == "STOP" {
						fmt.Println("TCP client exiting...")
						c.Close()
						return
					}
				}
			case <-time.After(1 * time.Second):
				// Do something when there is nothing read from stdin
			}
		}

	}
}
