package main

import (
	"fmt"
	"log"
	"math/rand"
)

func generateName() string {
	return fmt.Sprintf("User %d", rand.Intn(100)+1)
}

func startDialog(clientInputChannel <-chan ClientInput) {
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

func main() {

	ch := make(chan ClientInput)

	go startDialog(ch)

	err := startServer(ch)
	if err != nil {
		log.Fatal(err)
	}

}
