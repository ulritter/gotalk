package main

import (
	"fmt"
	"log"
	"net"
)

func startServerDialogHandling(clientInputChannel <-chan ClientInput) {
	w := &World{}

	for input := range clientInputChannel {
		switch event := input.event.(type) {
		case *MessageEvent:
			fmt.Println("Received Message: ", event.msg)
			// TODO: error handling
			input.user.session.WriteLine(fmt.Sprintf("You said \"%s\"\r\n", event.msg))

			for _, user := range w.users {
				if user != input.user {
					user.session.WriteLine(fmt.Sprintf("%s said, \"%s\"", input.user.name, event.msg))
				}
			}
		case *UserJoinedEvent:
			fmt.Println("User joined: ", input.user.name)
			w.users = append(w.users, input.user)
			input.user.session.WriteLine(fmt.Sprintf("Welcome %s", input.user.name))
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
	user := &User{name: generateName(), session: session}
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
		msg := string(buf[:n-2])
		log.Println("Received message:", msg)

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
