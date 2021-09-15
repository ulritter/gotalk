package main

import (
	"net"
)

const CMD_PREFIX = '/'
const CMD_EXIT = "exit"
const CMD_CHANGENICK = "nick"
const CMD_LISTUSERS = "list"
const CMD_ERROR = "error"
const CMD_HELP = "help"
const CMD_HELP1 = "?"
const CMD_ESCAPE_CHAR = '\f'

const CODE_NOCMD = 0
const CODE_EXIT = 1
const CODE_DONOTHING = 2

const BUFSIZE = 16384

type MessageEvent struct {
	msg string
}

type UserJoinedEvent struct {
}

type UserLeftEvent struct {
	user *User
	msg  string
}

type UserChangedNickEvent struct {
	user *User
	nick string
}

type ListUsersEvent struct {
	user *User
}

type ClientInput struct {
	user  *User
	event interface{}
}

type User struct {
	name       string
	session    *Session
	timejoined string
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

type Room struct {
	users []*User
}

type Newline struct {
	nl string
}

func (s *Session) WriteString(str string) error {
	_, err := s.conn.Write([]byte(str))
	return err
}
func (n *Newline) NewLine() string {
	return n.nl
}

func (n *Newline) SetNewLine(nline string) {
	n.nl = nline
}
