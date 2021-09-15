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

func sendServerCommand(conn net.Conn, cmd string) error {
	_, err := fmt.Fprint(conn, cmd)
	return err
}

func printHelp(nl Newline) {
	fmt.Print("HELP TEXT" + nl.NewLine())
}

func printError(nl Newline) {
	fmt.Print("Command error, type /help of /? for command descriptions" + nl.NewLine())
}

func parseCommand(conn net.Conn, msg string, nl Newline) int {
	if msg[0] != CMD_PREFIX {
		return CODE_NOCMD
	} else {
		cmdstring := msg[1:]
		// fmt.Println("command recognized:", cmdstring)
		//TODO: split in words
		cmd := strings.Fields(cmdstring)
		switch cmd[0] {
		case CMD_EXIT:
			sendServerCommand(conn, string(CMD_ESCAPE_CHAR)+CMD_EXIT+string(CMD_ESCAPE_CHAR))
			return CODE_EXIT
		case CMD_HELP, CMD_HELP1:
			printHelp(nl)
			return CODE_DONOTHING
		case CMD_LISTUSERS:
			sendServerCommand(conn, string(CMD_ESCAPE_CHAR)+CMD_LISTUSERS+string(CMD_ESCAPE_CHAR))
			return CODE_DONOTHING
		case CMD_CHANGENICK:
			// TODO: discard spaces
			cmd_arguments := cmd[1:]
			if len(cmd_arguments) != 1 {
				printError(nl)
				return CODE_DONOTHING
			} else {
				new_nick := cmd_arguments[0]
				sendServerCommand(conn, string(CMD_ESCAPE_CHAR)+CMD_CHANGENICK+string(CMD_ESCAPE_CHAR)+new_nick+string(CMD_ESCAPE_CHAR))
				return CODE_DONOTHING
			}
		default:
			return CODE_DONOTHING
		}
	}
}

// TODO: error handling for whole function

func clientDialogHandling(connect string, nick string, nl Newline) error {
	buf := make([]byte, BUFSIZE)
	conn, err := net.Dial("tcp", connect)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Connected to: %s, Nickname: %s %s", connect, nick, nl.NewLine())
	fmt.Fprintf(conn, string(CMD_ESCAPE_CHAR)+nick+string(CMD_ESCAPE_CHAR))

	for {
		go func() {
			for { // TODO: error handling
				n, err := conn.Read(buf)
				if err != nil {
					log.Printf("Error reading from buffer, most likely server was terminated" + nl.NewLine())
					os.Exit(1)
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
					if len(msg) > 0 {
						switch cC := parseCommand(conn, msg, nl); cC {
						case CODE_NOCMD:
							fmt.Fprintln(conn, msg)
						case CODE_EXIT:
							fmt.Print("TCP client exiting..." + nl.NewLine())
							conn.Close()
							return nil
						case CODE_DONOTHING:
							fallthrough
						default:
							break
						}
					}
				}
			case <-time.After(1 * time.Second):
				// Do something when there is nothing read from stdin
			}
		}
	}
}
