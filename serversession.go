package main

import (
	"fmt"
	"log"
	"time"
)

// dialog handling, broadcast user input to all users and status messages to all users or to a specific user depending on the request type
func handleServerSession(clientInputChannel <-chan ClientInput) {
	room := &Room{}
	for input := range clientInputChannel {
		switch event := input.event.(type) {
		case *MessageEvent:
			currentTime := time.Now()
			log.Printf(lang.Lookup(actualLocale, "Received Message at")+" %s "+lang.Lookup(actualLocale, "from")+" [%s]: %s"+newLine, currentTime.Format("2006.01.02 15:04:05"), input.user.name, event.msg)
			for _, user := range room.users {
				user.session.WriteMessage([]string{fmt.Sprintf("[%s]: %s", input.user.name, event.msg)})
			}
		case *UserJoinedEvent:
			log.Printf("User joined: %s"+newLine, input.user.name)
			room.users = append(room.users, input.user)
			input.user.session.WriteStatus([]string{
				fmt.Sprintf(lang.Lookup(actualLocale, "Welcome")+" %s", input.user.name),
			})
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus([]string{
						" ",
						fmt.Sprintf("%s ", input.user.name) + lang.Lookup(actualLocale, "entered the room"),
					})
				}
			}
		case *UserLeftEvent:
			log.Printf(lang.Lookup(actualLocale, "User left:")+" %s %s"+newLine, event.msg, event.user.name)
			var users []*User
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus([]string{
						" ",
						fmt.Sprintf("%s "+lang.Lookup(actualLocale, "left the room"), input.user.name),
					})
					users = append(users, user)
				}
			}
			room.users = users

		case *UserChangedNickEvent:
			log.Printf(lang.Lookup(actualLocale, "User")+" %s "+lang.Lookup(actualLocale, "has changed his nick to:")+" %s", event.user.name, event.nick)
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus([]string{
						" ",
						fmt.Sprintf("[%s] "+lang.Lookup(actualLocale, "has changed his nick to:")+" [%s]", event.user.name, event.nick),
					})

				} else {
					user.session.WriteStatus([]string{
						" ",
						lang.Lookup(actualLocale, "You have changed your nick"),
						fmt.Sprintf(lang.Lookup(actualLocale, "from")+" [%s] "+lang.Lookup(actualLocale, "to")+" [%s]", event.user.name, event.nick),
					})
				}
			}
			input.user.name = event.nick
		case *ListUsersEvent:
			log.Printf(lang.Lookup(actualLocale, "User")+" %s "+lang.Lookup(actualLocale, "has requested user list"), input.user.name)
			var list []string
			list = nil
			list = append(list, " ")
			list = append(list, lang.Lookup(actualLocale, "User list:"))

			for _, user := range room.users {
				list = append(list, fmt.Sprintf("[%s] - "+lang.Lookup(actualLocale, "joined at")+" [%s]", user.name, user.timejoined))
			}
			input.user.session.WriteStatus(list)
		}
	}
}
