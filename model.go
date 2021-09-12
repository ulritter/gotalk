package main

import (
	"fmt"
	"math/rand"
	"net"
)

type MessageEvent struct {
	msg string
}

type UserJoinedEvent struct {
}

type UserLeftEvent struct {
	user *User
	msg  string
}

type ClientInput struct {
	user  *User
	event interface{}
}

type User struct {
	name    string
	session *Session
}

type Session struct {
	conn net.Conn
}

type WhoAmI struct {
	server bool
	addr   string
	port   string
	nick   string
}

type World struct {
	users []*User
}

func (s *Session) WriteLine(str string) error {
	_, err := s.conn.Write([]byte(str + "\r\n"))
	return err
}

func generateName() string {
	return fmt.Sprintf("User %d", rand.Intn(100)+1)
}
