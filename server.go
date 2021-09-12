package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func startServerDialogHandling(clientInputChannel <-chan ClientInput) {
	w := &World{}

	for input := range clientInputChannel {
		switch event := input.event.(type) {
		case *MessageEvent:
			currentTime := time.Now()
			fmt.Printf("Received Message at %s from [%s]: %s\n", currentTime.Format("2006.01.02 15:04:05"), input.user.name, event.msg)
			// TODO: error handling
			// input.user.session.WriteLine(fmt.Sprintf("You said \"%s\"\r\n", event.msg))

			for _, user := range w.users {
				if user != input.user {
					user.session.WriteLine(fmt.Sprintf("[%s]: %s", input.user.name, event.msg))
				}
			}
		case *UserJoinedEvent:
			fmt.Println("User joined: ", input.user.name)
			w.users = append(w.users, input.user)
			input.user.session.WriteLine(fmt.Sprintf("Welcome %s\n", input.user.name))
			for _, user := range w.users {
				if user != input.user {
					user.session.WriteLine(fmt.Sprintf("%s entered the room", input.user.name))
				}
			}
		case *UserLeftEvent:
			fmt.Printf("User: %s, %s", event.user.name, event.msg)
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
		if n == 0 {
			log.Println("Zero bytes, closing connection")
			inputChannel <- ClientInput{
				user,
				&UserLeftEvent{user, "Goodbye"},
			}
		}
		// TODO: check real empty imput like ^d
		msg := strings.TrimSpace(string(string(buf[:n])))
		// log.Println("Received message:", msg)

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
