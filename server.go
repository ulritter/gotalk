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
func handleConnection(conn net.Conn, inputChannel chan ClientInput) error {
	buf := make([]byte, BUFSIZE)

	session := &Session{conn}
	n, err := conn.Read(buf)
	if err != nil {
		log.Print(lang.Lookup(actualLocale, "Error reading from buffer")+newLine, err)
		return err
	}
	msg := Message{}
	msg.Body = nil
	msg.UnmarshalMSG(buf[:n])

	if (msg.Action != ACTION_INIT) || (len(msg.Body) != 2) {
		return fmt.Errorf(lang.Lookup(actualLocale, "Wrong connection initialization message."))
	} else {
		// expecting format {ACTION_INIT, [{<nickname>}, {<revision level>}]}
		if msg.Body[1] != REVISION {
			sendMessage(conn, ACTION_REVISION, []string{REVISION})
			return fmt.Errorf(lang.Lookup(actualLocale,
				"Connection request from ")+conn.RemoteAddr().(*net.TCPAddr).IP.String()+lang.Lookup(actualLocale,
				" rejected. ")+lang.Lookup(actualLocale,
				"Wrong client revision level. Should be: ")+" %s"+lang.Lookup(actualLocale, ", actual: ")+"%s", REVISION, msg.Body[1])
		}
	}
	user := &User{name: msg.Body[0], session: session, timejoined: time.Now().Format("2006.01.02 15:04:05")}
	inputChannel <- ClientInput{
		user,
		&UserJoinedEvent{},
	}

	for {
		n, err1 := conn.Read(buf)
		if err1 != nil {
			log.Printf(lang.Lookup(actualLocale, "End condition, closing connection for:")+" %s"+newLine, user.name)
			inputChannel <- ClientInput{
				user,
				&UserLeftEvent{user, lang.Lookup(actualLocale, "Goodbye")},
			}
			return err1
		}

		msg.Action = ""
		msg.Body = nil
		err2 := msg.UnmarshalMSG(buf[:n])

		if err2 != nil {
			log.Printf(lang.Lookup(actualLocale, "Warning: Corrupt JSON Message from: ")+" %s"+newLine, user.name)
			log.Println(err2)
		}

		if msg.Action == ACTION_EXIT {
			log.Printf(lang.Lookup(actualLocale, "End condition, closing connection for:")+" %s"+newLine, user.name)
			inputChannel <- ClientInput{
				user,
				&UserLeftEvent{user, lang.Lookup(actualLocale, "Goodbye")},
			}
			return err1
		}

		switch msg.Action {
		case ACTION_CHANGENICK:
			if len(msg.Body) == 1 {
				inputChannel <- ClientInput{
					user,
					&UserChangedNickEvent{user, msg.Body[0]},
				}
			}
		case ACTION_LISTUSERS:
			if msg.Action == ACTION_LISTUSERS {
				inputChannel <- ClientInput{
					user,
					&ListUsersEvent{user},
				}
			}
		case ACTION_SENDMESSAGE:
			if len(msg.Body) == 1 {
				sendmsg := strings.TrimSpace(msg.Body[0])
				e := ClientInput{user, &MessageEvent{sendmsg}}
				inputChannel <- e
			}
		default:
		}

	}
}

// this function is called by main() in the case the app needs to operate as server
// wait for connections and start a handler for each connection
func startServer(eventChannel chan ClientInput, config *tls.Config, port string) error {
	log.Printf(lang.Lookup(actualLocale, "Starting server on port ")+"%s"+newLine, port)
	ln, err := tls.Listen("tcp", port, config)
	if err != nil {
		// handle error
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			log.Print(lang.Lookup(actualLocale, "Error accepting connection")+newLine, err)
			continue
		}
		go func() {
			if err := handleConnection(conn, eventChannel); err != nil {
				log.Print(lang.Lookup(actualLocale, "Error handling connection or unexpected client exit")+newLine, err)
			}
		}()

	}
}
