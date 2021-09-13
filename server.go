package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func serverDialogHandling(clientInputChannel <-chan ClientInput) {
	room := &Room{}
	for input := range clientInputChannel {
		switch event := input.event.(type) {
		case *MessageEvent:
			currentTime := time.Now()
			log.Printf("Received Message at %s from [%s]: %s\n", currentTime.Format("2006.01.02 15:04:05"), input.user.name, event.msg)

			for _, user := range room.users {
				if user != input.user {
					user.session.WriteLine(fmt.Sprintf("[%s]: %s", input.user.name, event.msg))
				}
			}
		case *UserJoinedEvent:
			log.Println("User joined: ", input.user.name)
			room.users = append(room.users, input.user)
			input.user.session.WriteLine(fmt.Sprintf("Welcome %s\n", input.user.name))
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteLine(fmt.Sprintf("%s entered the room", input.user.name))
				}
			}
		case *UserLeftEvent:
			log.Printf("User left: %s, %s\n", event.msg, event.user.name)
			var users []*User
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteLine(fmt.Sprintf("%s left the room", input.user.name))
					users = append(users, user)
				}
			}
			room.users = users
		}
	}
}

func handleConnection(conn net.Conn, inputChannel chan ClientInput) error {
	buf := make([]byte, 4096)

	session := &Session{conn}
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading from buffer", err)
		return err
	}
	var nick string
	pattern := string(buf[:n])
	if (pattern[0] == '$') && (pattern[1] == '$') && (pattern[2] == '$') && (pattern[n-1] == '$') {
		nick = string(buf[3 : n-1])
	} else {
		nick = "J_Doe"
	}

	user := &User{name: nick, session: session}
	inputChannel <- ClientInput{
		user,
		&UserJoinedEvent{},
	}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Error reading from buffer", err)
			return err
		}
		pattern := string(buf[:n])
		if (n == 0) || ((pattern[0] == 'S') && (pattern[1] == 'T') && (pattern[2] == 'O') && (pattern[3] == 'P')) {
			log.Println("End condition, closing connection")
			inputChannel <- ClientInput{
				user,
				&UserLeftEvent{user, "Goodbye"},
			}
			break
		}

		msg := strings.TrimSpace(string(string(buf[:n])))

		e := ClientInput{user, &MessageEvent{msg}}
		inputChannel <- e
	}
	return nil
}

func startServer(eventChannel chan ClientInput, port string) error {
	log.Println("Starting server")
	ln, err := net.Listen("tcp", port)
	if err != nil {
		// handle error
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			log.Println("Error accepting connection", err)
			continue
		}
		go func() {
			if err := handleConnection(conn, eventChannel); err != nil {
				log.Println("Error handling connection", err)
			}
		}()

	}
}
