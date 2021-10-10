package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// read from connection, recognize request types and pass appropriate event types to the session handler (serverDialog())
func handleConnection(conn net.Conn, inputChannel chan ClientInput, nl Newline) error {
	buf := make([]byte, BUFSIZE)

	session := &Session{conn}
	n, err := conn.Read(buf)
	if err != nil {
		log.Print(lang.Lookup(actualLocale, "Error reading from buffer")+nl.NewLine(), err)
		return err
	}
	var nick string

	rawData := string(buf[:n])
	rawDataFields := strings.Fields(rawData)

	if len(rawDataFields) != 2 {
		return fmt.Errorf(lang.Lookup(actualLocale, "Wrong connection initialization message."))
	} else if rawDataFields[1] != REVISION {
		str := string(CMD_ESCAPE_CHAR) + string(CMD_ESCAPE_CHAR) + REVISION
		conn.Write([]byte(str))
		return fmt.Errorf(lang.Lookup(actualLocale, "Connection request from ")+conn.RemoteAddr().(*net.TCPAddr).IP.String()+lang.Lookup(actualLocale, " rejected. ")+lang.Lookup(actualLocale, "Wrong client revision level. Should be: ")+" %s"+lang.Lookup(actualLocale, ", actual: ")+"%s", REVISION, rawDataFields[1])
	}

	assumedNick := rawDataFields[0]

	if (assumedNick[0] == CMD_ESCAPE_CHAR) && (assumedNick[n-1] == CMD_ESCAPE_CHAR) {
		nick = string(buf[1 : n-1])
	} else {
		nick = "J_Doe"
	}

	user := &User{name: nick, session: session, timejoined: time.Now().Format("2006.01.02 15:04:05")}
	inputChannel <- ClientInput{
		user,
		&UserJoinedEvent{},
	}

	for {
		n, err := conn.Read(buf)

		if (buf[0] == CMD_ESCAPE_CHAR) || (err != nil) {
			pattern := strings.Fields(string(buf[:n]))
			if (len(pattern) == 1) && ((pattern[0] == (CMD_EXIT1)) || (pattern[0] == (CMD_EXIT2)) || (pattern[0] == (CMD_EXIT3))) || (err != nil) {
				log.Printf(lang.Lookup(actualLocale, "End condition, closing connection for:")+" %s"+nl.NewLine(), user.name)
				inputChannel <- ClientInput{
					user,
					&UserLeftEvent{user, lang.Lookup(actualLocale, "Goodbye")},
				}
				return err
			} else if (len(pattern) == 2) && (pattern[0] == (CMD_CHANGENICK)) {
				inputChannel <- ClientInput{
					user,
					&UserChangedNickEvent{user, pattern[1]},
				}
			} else if (len(pattern) == 1) && (pattern[0] == (CMD_LISTUSERS)) {
				inputChannel <- ClientInput{
					user,
					&ListUsersEvent{user},
				}
			}
		} else {
			msg := strings.TrimSpace(string(string(buf[:n])))
			e := ClientInput{user, &MessageEvent{msg}}
			inputChannel <- e
		}

	}
}

// this function is called by main() in the case the app needs to operate as server
// wait for connections and start a handler for each connection
func startServer(eventChannel chan ClientInput, config *tls.Config, port string, nl Newline) error {
	log.Printf(lang.Lookup(actualLocale, "Starting server on port ")+"%s"+nl.NewLine(), port)
	ln, err := tls.Listen("tcp", port, config)
	if err != nil {
		// handle error
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			log.Print(lang.Lookup(actualLocale, "Error accepting connection")+nl.NewLine(), err)
			continue
		}
		go func() {
			if err := handleConnection(conn, eventChannel, nl); err != nil {
				log.Print(lang.Lookup(actualLocale, "Error handling connection or unexpected client exit")+nl.NewLine(), err)
			}
		}()

	}
}
