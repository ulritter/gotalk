package main

import (
	"fmt"
	"time"
)

// dialog handling, broadcast user input to all users and status messages to all users or to a specific user depending on the request type
func (a *application) handleServerSession(clientInputChannel <-chan ClientInput) {
	room := &Room{}
	for input := range clientInputChannel {

		switch event := input.event.(type) {
		case *MessageEvent:
			currentTime := time.Now()
			a.logger.Printf(a.lang.Lookup(a.config.locale, "Received Message at")+" %s "+a.lang.Lookup(a.config.locale, "from")+" [%s]: %s"+a.config.newline, currentTime.Format("2006.01.02 15:04:05"), input.user.name, event.msg)
			for _, user := range room.users {
				user.session.WriteMessage([]string{fmt.Sprintf("[%s]: %s", input.user.name, event.msg)})
			}
		case *UserJoinedEvent:
			a.logger.Printf("User joined: %s"+a.config.newline, input.user.name)
			room.users = append(room.users, input.user)
			input.user.session.WriteStatus([]string{
				fmt.Sprintf(a.lang.Lookup(a.config.locale, "Welcome")+" %s", input.user.name),
			})
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus([]string{
						" ",
						fmt.Sprintf("%s ", input.user.name) + a.lang.Lookup(a.config.locale, "entered the room"),
					})
				}
			}
		case *UserLeftEvent:
			a.logger.Printf(a.lang.Lookup(a.config.locale, "User left:")+" %s %s"+a.config.newline, event.msg, event.user.name)
			var users []*User
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus([]string{
						" ",
						fmt.Sprintf("%s "+a.lang.Lookup(a.config.locale, "left the room"), input.user.name),
					})
					users = append(users, user)
				}
			}
			room.users = users

		case *UserChangedNickEvent:
			a.logger.Printf(a.lang.Lookup(a.config.locale, "User")+" %s "+a.lang.Lookup(a.config.locale, "has changed his nick to:")+" %s", event.user.name, event.nick)
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus([]string{
						" ",
						fmt.Sprintf("[%s] "+a.lang.Lookup(a.config.locale, "has changed his nick to:")+" [%s]", event.user.name, event.nick),
					})

				} else {
					user.session.WriteStatus([]string{
						" ",
						a.lang.Lookup(a.config.locale, "You have changed your nick"),
						fmt.Sprintf(a.lang.Lookup(a.config.locale, "from")+" [%s] "+a.lang.Lookup(a.config.locale, "to")+" [%s]", event.user.name, event.nick),
					})
				}
			}
			input.user.name = event.nick
		case *ListUsersEvent:
			a.logger.Printf(a.lang.Lookup(a.config.locale, "User")+" %s "+a.lang.Lookup(a.config.locale, "has requested user list"), input.user.name)
			var list []string
			list = nil
			list = append(list, " ")
			list = append(list, a.lang.Lookup(a.config.locale, "User list:"))

			for _, user := range room.users {
				list = append(list, fmt.Sprintf("[%s] - "+a.lang.Lookup(a.config.locale, "joined at")+" [%s]", user.name, user.timejoined))
			}
			input.user.session.WriteStatus(list)
		}
	}
}
