package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func serverDialogHandling(clientInputChannel <-chan ClientInput, nl Newline) {
	room := &Room{}
	for input := range clientInputChannel {
		switch event := input.event.(type) {
		case *MessageEvent:
			currentTime := time.Now()
			log.Printf("Received Message at %s from [%s]: %s"+nl.NewLine(), currentTime.Format("2006.01.02 15:04:05"), input.user.name, event.msg)

			for _, user := range room.users {
				if user != input.user {
					user.session.WriteString(fmt.Sprintf("[%s]: %s"+nl.NewLine(), input.user.name, event.msg))
				}
			}
		case *UserJoinedEvent:
			log.Print("User joined: "+nl.NewLine(), input.user.name)
			room.users = append(room.users, input.user)
			input.user.session.WriteString(fmt.Sprintf("Welcome %s"+nl.NewLine(), input.user.name))
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteString(fmt.Sprintf("%s entered the room"+nl.NewLine(), input.user.name))
				}
			}
		case *UserLeftEvent:
			log.Printf("User left: %s: %s"+nl.NewLine(), event.msg, event.user.name)
			var users []*User
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteString(fmt.Sprintf("%s left the room"+nl.NewLine(), input.user.name))
					users = append(users, user)
				}
			}
			room.users = users

		case *UserChangedNickEvent:
			log.Printf("User %s has changed his nick to: %s"+nl.NewLine(), event.user.name, event.nick)
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteString(fmt.Sprintf("[%s] has changed his nick to: [%s]"+nl.NewLine(), event.user.name, event.nick))

				}
			}
			input.user.name = event.nick
		case *ListUsersEvent:
			log.Printf("User %s has requested user list"+nl.NewLine(), input.user.name)
			input.user.session.WriteString(fmt.Sprint("User list:" + nl.NewLine()))
			input.user.session.WriteString(fmt.Sprint("==========" + nl.NewLine()))
			for _, user := range room.users {
				input.user.session.WriteString(fmt.Sprintf("[%s] has joined at [%s]"+nl.NewLine(), user.name, user.timejoined))
			}
		}
	}
}

func handleConnection(conn net.Conn, inputChannel chan ClientInput, nl Newline) error {
	buf := make([]byte, BUFSIZE)

	session := &Session{conn}
	n, err := conn.Read(buf)
	if err != nil {
		log.Print("Error reading from buffer"+nl.NewLine(), err)
		return err
	}
	var nick string
	pattern := string(buf[:n])
	if (pattern[0] == CMD_ESCAPE_CHAR) && (pattern[n-1] == CMD_ESCAPE_CHAR) {
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
		if err != nil {
			log.Print("Error reading from buffer"+nl.NewLine(), err)
			return err
		}

		if buf[0] == CMD_ESCAPE_CHAR {
			pattern := strings.Fields(string(buf[:n]))
			if (len(pattern) == 1) && (pattern[0] == (CMD_EXIT)) {
				log.Printf("End condition, closing connection for %s"+nl.NewLine(), user.name)
				inputChannel <- ClientInput{
					user,
					&UserLeftEvent{user, "Goodbye"},
				}
				break
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
	return nil
}

func startServer(eventChannel chan ClientInput, config *tls.Config, port string, nl Newline) error {
	log.Printf("Starting server on %s"+nl.NewLine(), port)
	ln, err := tls.Listen("tcp", port, config)
	if err != nil {
		// handle error
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			log.Print("Error accepting connection"+nl.NewLine(), err)
			continue
		}
		go func() {
			if err := handleConnection(conn, eventChannel, nl); err != nil {
				log.Print("Error handling connection"+nl.NewLine(), err)
			}
		}()

	}
}
