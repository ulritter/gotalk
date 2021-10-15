package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

// read from connection, recognize request types and pass appropriate event types to the session handler (serverDialog())
func (a *application) handleConnection() error {
	buf := make([]byte, BUFSIZE)

	session := &Session{a.config.conn}
	n, err := a.config.conn.Read(buf)
	if err != nil {
		a.logger.Print(a.lang.Lookup(a.config.locale, "Error reading from buffer")+a.config.newline, err)
		return err
	}
	msg := Message{}
	msg.Body = nil
	msg.UnmarshalMSG(buf[:n])

	if (msg.Action != ACTION_INIT) || (len(msg.Body) != 2) {
		return fmt.Errorf(a.lang.Lookup(a.config.locale, "Wrong connection initialization message."))
	} else {
		// expecting format {ACTION_INIT, [{<nickname>}, {<revision level>}]}
		if msg.Body[1] != REVISION {
			sendMessage(a.config.conn, ACTION_REVISION, []string{REVISION})
			return fmt.Errorf(a.lang.Lookup(a.config.locale,
				"Connection request from ")+a.config.conn.RemoteAddr().(*net.TCPAddr).IP.String()+a.lang.Lookup(a.config.locale,
				" rejected. ")+a.lang.Lookup(a.config.locale,
				"Wrong client revision level. Should be: ")+" %s"+a.lang.Lookup(a.config.locale, ", actual: ")+"%s", REVISION, msg.Body[1])
		}
	}
	user := &User{name: msg.Body[0], session: session, timejoined: time.Now().Format("2006.01.02 15:04:05")}
	a.config.ch <- ClientInput{
		user,
		&UserJoinedEvent{},
	}

	for {
		n, err1 := a.config.conn.Read(buf)
		if err1 != nil {
			a.logger.Printf(a.lang.Lookup(a.config.locale, "End condition, closing connection for:")+" %s"+a.config.newline, user.name)
			a.config.ch <- ClientInput{
				user,
				&UserLeftEvent{user, a.lang.Lookup(a.config.locale, "Goodbye")},
			}
			return err1
		}

		msg.Action = ""
		msg.Body = nil
		err2 := msg.UnmarshalMSG(buf[:n])

		if err2 != nil {
			a.logger.Printf(a.lang.Lookup(a.config.locale, "Warning: Corrupt JSON Message from: ")+" %s"+a.config.newline, user.name)
			a.logger.Println(err2)
		}

		if msg.Action == ACTION_EXIT {
			a.logger.Printf(a.lang.Lookup(a.config.locale, "End condition, closing connection for:")+" %s"+a.config.newline, user.name)
			//echo exit condition for organized client shutdown
			sendMessage(a.config.conn, ACTION_EXIT, []string{""})
			a.config.ch <- ClientInput{
				user,
				&UserLeftEvent{user, a.lang.Lookup(a.config.locale, "Goodbye")},
			}
			return err1
		}

		switch msg.Action {
		case ACTION_CHANGENICK:
			if len(msg.Body) == 1 {
				a.config.ch <- ClientInput{
					user,
					&UserChangedNickEvent{user, msg.Body[0]},
				}
			}
		case ACTION_LISTUSERS:
			if msg.Action == ACTION_LISTUSERS {
				a.config.ch <- ClientInput{
					user,
					&ListUsersEvent{user},
				}
			}
		case ACTION_SENDMESSAGE:
			if len(msg.Body) == 1 {
				sendmsg := strings.TrimSpace(msg.Body[0])
				e := ClientInput{user, &MessageEvent{sendmsg}}
				a.config.ch <- e
			}
		default:
		}

	}
}

// this function is called by main() in the case the app needs to operate as server
// wait for connections and start a handler for each connection
func (a *application) startServer() error {
	a.logger.Printf(a.lang.Lookup(a.config.locale, "Starting server on port ")+"%s"+a.config.newline, a.config.port)
	ln, err := tls.Listen("tcp", a.config.port, a.config.tlsConfig)
	if err != nil {
		// handle error
		return err
	}
	for {
		conn, err := ln.Accept()
		a.config.conn = conn
		if err != nil {
			// handle error
			a.logger.Print(a.lang.Lookup(a.config.locale, "Error accepting connection")+a.config.newline, err)
			continue
		}
		go func() {
			if err := a.handleConnection(); err != nil {
				a.logger.Printf(a.lang.Lookup(a.config.locale, "Error handling connection or unexpected client exit")+": %v"+a.config.newline, err)
			}
		}()

	}
}
