package main

import (
	"fmt"
	"log"
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
				user.session.WriteString(fmt.Sprintf("[%s]: %s", input.user.name, event.msg))
			}
		case *UserJoinedEvent:
			log.Print("User joined: "+nl.NewLine(), input.user.name)
			room.users = append(room.users, input.user)
			input.user.session.WriteStatus(fmt.Sprintf("Welcome %s", input.user.name))
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus(fmt.Sprintf(nl.NewLine()+"%s entered the room", input.user.name))
				}
			}
		case *UserLeftEvent:
			log.Printf("User left: %s: %s"+nl.NewLine(), event.msg, event.user.name)
			var users []*User
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus(fmt.Sprintf(nl.NewLine()+"%s left the room", input.user.name))
					users = append(users, user)
				}
			}
			room.users = users

		case *UserChangedNickEvent:
			log.Printf("User %s has changed his nick to: %s"+nl.NewLine(), event.user.name, event.nick)
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus(fmt.Sprintf(nl.NewLine()+"[%s] has changed his nick to: [%s]", event.user.name, event.nick))

				} else {
					user.session.WriteStatus(fmt.Sprintf(nl.NewLine()+"You have changed your nick from [%s] to [%s]", event.user.name, event.nick))
				}
			}
			input.user.name = event.nick
		case *ListUsersEvent:
			log.Printf("User %s has requested user list"+nl.NewLine(), input.user.name)
			input.user.session.WriteStatus(nl.NewLine() + "User list:" + nl.NewLine() + "==========")

			for _, user := range room.users {
				input.user.session.WriteStatus(fmt.Sprintf("[%s] - joined at [%s]", user.name, user.timejoined))
			}
		}
	}
}
