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
