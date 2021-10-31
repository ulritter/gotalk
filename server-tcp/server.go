package main

import (
	"crypto/tls"
	"fmt"
	"gotalk/constants"
	"gotalk/models"
	"net"
	"strings"
	"time"
)

// read from connection, recognize request types and pass appropriate event types to the session handler (serverDialog())
func handleConnection(a *models.Application) error {
	buf := make([]byte, constants.BUFSIZE)

	session := &models.Session{Conn: a.Config.Conn}
	n, err := a.Config.Conn.Read(buf)
	if err != nil {
		a.Logger.Print(a.Lang.Lookup(a.Config.Locale, "Error reading from buffer")+a.Config.Newline, err)
		return err
	}
	msg := models.Message{}
	msg.Body = nil
	msg.UnmarshalMSG(buf[:n])

	if (msg.Action != constants.ACTION_INIT) || (len(msg.Body) != 2) {
		return fmt.Errorf(a.Lang.Lookup(a.Config.Locale, "Wrong connection initialization message."))
	} else {
		// expecting format {models.ACTION_INIT, [{<nickname>}, {<revision level>}]}
		if msg.Body[1] != constants.REVISION {
			models.SendJSONMessage(a.Config.Conn, constants.ACTION_REVISION, []string{constants.REVISION})
			return fmt.Errorf(a.Lang.Lookup(a.Config.Locale,
				"Connection request from ")+a.Config.Conn.RemoteAddr().(*net.TCPAddr).IP.String()+a.Lang.Lookup(a.Config.Locale,
				" rejected. ")+a.Lang.Lookup(a.Config.Locale,
				"Wrong client revision level. Should be: ")+" %s"+a.Lang.Lookup(a.Config.Locale, ", actual: ")+"%s", constants.REVISION, msg.Body[1])
		}
	}
	user := &models.User{Name: msg.Body[0], Session: session, Timejoined: time.Now().Format("2006.01.02 15:04:05")}
	a.Config.Ch <- models.ClientInput{
		User:  user,
		Event: &models.UserJoinedEvent{},
	}

	for {
		n, err1 := a.Config.Conn.Read(buf)
		if err1 != nil {
			a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "End condition, closing connection for:")+" %s"+a.Config.Newline, user.Name)
			a.Config.Ch <- models.ClientInput{
				User: user,
				Event: &models.UserLeftEvent{
					User: user, Msg: a.Lang.Lookup(a.Config.Locale, "Goodbye")},
			}
			return err1
		}

		msg.Action = ""
		msg.Body = nil
		err2 := msg.UnmarshalMSG(buf[:n])

		if err2 != nil {
			a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "Warning: Corrupt JSON Message from: ")+" %s"+a.Config.Newline, user.Name)
			a.Logger.Println(err2)
		}

		if msg.Action == constants.ACTION_EXIT {
			a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "End condition, closing connection for:")+" %s"+a.Config.Newline, user.Name)
			//echo exit condition for organized client shutdown
			models.SendJSONMessage(a.Config.Conn, constants.ACTION_EXIT, nil)
			a.Config.Ch <- models.ClientInput{
				User:  user,
				Event: &models.UserLeftEvent{User: user, Msg: a.Lang.Lookup(a.Config.Locale, "Goodbye")},
			}
			return err1
		}

		switch msg.Action {
		case constants.ACTION_CHANGENICK:
			if len(msg.Body) == 1 {
				a.Config.Ch <- models.ClientInput{
					User:  user,
					Event: &models.UserChangedNickEvent{User: user, Nick: msg.Body[0]},
				}
			}
		case constants.ACTION_LISTUSERS:
			a.Config.Ch <- models.ClientInput{
				User:  user,
				Event: &models.ListUsersEvent{User: user},
			}
		case constants.ACTION_SENDMESSAGE:
			if len(msg.Body) == 1 {
				sendmsg := strings.TrimSpace(msg.Body[0])
				e := models.ClientInput{User: user, Event: &models.MessageEvent{Msg: sendmsg}}
				a.Config.Ch <- e
			}
		default:
		}

	}
}

// this function is called by main() in the case the app needs to operate as server
// wait for connections and start a handler for each connection
func startServer(a *models.Application) error {
	a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "Starting server on port ")+"%s"+a.Config.Newline, a.Config.Port)
	ln, err := tls.Listen("tcp", a.Config.Port, a.Config.TLSconfig)
	if err != nil {
		// handle error
		return err
	}
	for {
		conn, err := ln.Accept()
		a.Config.Conn = conn
		if err != nil {
			// handle error
			a.Logger.Print(a.Lang.Lookup(a.Config.Locale, "Error accepting connection")+a.Config.Newline, err)
			continue
		}
		go func() {
			if err := handleConnection(a); err != nil {
				a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "Error handling connection or unexpected client exit")+": %v"+a.Config.Newline, err)
			}
		}()

	}
}
