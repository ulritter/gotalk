package main

import (
	"fmt"
	"gotalk/models"
	"time"
)

// dialog handling, broadcast user input to all users and status messages to all users or to a specific user depending on the request type
func handleServerSession(a *models.Application) {
	room := &models.Room{}
	for input := range a.Config.Ch {

		switch event := input.Event.(type) {
		case *models.MessageEvent:
			currentTime := time.Now()
			a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "Received Message at")+" %s "+a.Lang.Lookup(a.Config.Locale, "from")+" [%s]: %s"+a.Config.Newline, currentTime.Format("2006.01.02 15:04:05"), input.User.Name, event.Msg)
			for _, user := range room.Users {
				user.Session.WriteMessage([]string{fmt.Sprintf("[%s]: %s", input.User.Name, event.Msg)})
			}
		case *models.UserJoinedEvent:
			a.Logger.Printf("User joined: %s"+a.Config.Newline, input.User.Name)
			room.Users = append(room.Users, input.User)
			input.User.Session.WriteStatus([]string{
				fmt.Sprintf(a.Lang.Lookup(a.Config.Locale, "Welcome")+" %s", input.User.Name),
			})
			for _, user := range room.Users {
				if user != input.User {
					user.Session.WriteStatus([]string{
						" ",
						fmt.Sprintf("%s ", input.User.Name) + a.Lang.Lookup(a.Config.Locale, "entered the room"),
					})
				}
			}
		case *models.UserLeftEvent:
			a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "User left:")+" %s %s"+a.Config.Newline, event.Msg, event.User.Name)
			var users []*models.User
			for _, user := range room.Users {
				if user != input.User {
					user.Session.WriteStatus([]string{
						" ",
						fmt.Sprintf("%s "+a.Lang.Lookup(a.Config.Locale, "left the room"), input.User.Name),
					})
					users = append(users, user)
				}
			}
			room.Users = users

		case *models.UserChangedNickEvent:
			a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "User")+" %s "+a.Lang.Lookup(a.Config.Locale, "has changed his nick to:")+" %s", event.User.Name, event.Nick)
			for _, user := range room.Users {
				if user != input.User {
					user.Session.WriteStatus([]string{
						" ",
						fmt.Sprintf("[%s] "+a.Lang.Lookup(a.Config.Locale, "has changed his nick to:")+" [%s]", event.User.Name, event.Nick),
					})

				} else {
					user.Session.WriteStatus([]string{
						" ",
						a.Lang.Lookup(a.Config.Locale, "You have changed your nick"),
						fmt.Sprintf(a.Lang.Lookup(a.Config.Locale, "from")+" [%s] "+a.Lang.Lookup(a.Config.Locale, "to")+" [%s]", event.User.Name, event.Nick),
					})
				}
			}
			input.User.Name = event.Nick
		case *models.ListUsersEvent:
			a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "User")+" %s "+a.Lang.Lookup(a.Config.Locale, "has requested user list"), input.User.Name)
			var list []string
			list = nil
			list = append(list, " ")
			list = append(list, a.Lang.Lookup(a.Config.Locale, "User list:"))

			for _, user := range room.Users {
				list = append(list, fmt.Sprintf("[%s] - "+a.Lang.Lookup(a.Config.Locale, "joined at")+" [%s]", user.Name, user.Timejoined))
			}
			input.User.Session.WriteStatus(list)
		}
	}
}
