package main

import (
	"net"

	language "github.com/moemoe89/go-localization"
)

var lang *language.Config

var actualLocale string

const RAWFILE = "https://raw.githubusercontent.com/ulritter/gotalk/main/language.json"
const LANGFILE = "./language.json"

const CMD_PREFIX = '/'
const CMD_EXIT1 = "exit"
const CMD_EXIT2 = "quit"
const CMD_EXIT3 = "q"
const CMD_CHANGENICK = "nick"
const CMD_LISTUSERS = "list"
const CMD_ERROR = "error"
const CMD_HELP = "help"
const CMD_HELP1 = "?"
const CMD_ESCAPE_CHAR = '\f'

const CODE_NOCMD = 0
const CODE_EXIT = 1
const CODE_DONOTHING = 2

const BUFSIZE = 4096

// Event type for messages
type MessageEvent struct {
	msg string
}

// Event type for users joining
type UserJoinedEvent struct {
}

// Event type for users leaving
type UserLeftEvent struct {
	user *User
	msg  string
}

// Event type for users changing their nick
type UserChangedNickEvent struct {
	user *User
	nick string
}

// Event type for users requesting a room user list
type ListUsersEvent struct {
	user *User
}

// Commmunication structure between connection handler and user session
type ClientInput struct {
	user  *User
	event interface{}
}

// User Info
type User struct {
	name       string
	session    *Session
	timejoined string
}

// Structure holding the connection for the session
type Session struct {
	conn net.Conn
}

// Structure holding the users the room
type Room struct {
	users []*User
}

// Structure holding the command line parameters (are filled with defaults on startup)
type WhoAmI struct {
	server bool
	addr   string
	port   string
	nick   string
}

// Holding new line flavours for either linux or windows type systems
type Newline struct {
	nl string
}

//sends a message string from server to client
func (s *Session) WriteMessage(str string) error {
	_, err := s.conn.Write([]byte(str))
	return err
}

//sends a status string from server to client
func (s *Session) WriteStatus(str string) error {
	str = string(CMD_ESCAPE_CHAR) + str
	_, err := s.conn.Write([]byte(str))
	return err
}

//returns newline character for either linux or windows type systems
func (n *Newline) NewLine() string {
	return n.nl
}
