package main

import (
	"encoding/json"
	"net"

	language "github.com/moemoe89/go-localization"
)

var lang *language.Config

var actualLocale string

const RAWFILE = "https://raw.githubusercontent.com/ulritter/gotalk/main/language.json"
const LANGFILE = "./language.json"

const ACTION_CHANGENICK = "changenick"
const ACTION_SENDMESSAGE = "message"
const ACTION_LISTUSERS = "listusers"
const ACTION_REVISION = "revision"
const ACTION_SENDSTATUS = "status"
const ACTION_EXIT = "exit"
const ACTION_INIT = "init"

const CMD_PREFIX = '/'
const CMD_EXIT1 = "exit"
const CMD_EXIT2 = "quit"
const CMD_EXIT3 = "q"
const CMD_CHANGENICK = "nick"
const CMD_LISTUSERS = "list"
const CMD_ERROR = "error"
const CMD_HELP = "help"
const CMD_HELP1 = "?"

const BUFSIZE = 4096

const REVISION = "0.8.1"

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

type Message struct {
	Action string   `json:"action"`
	Body   []string `json:"body"`
}

func (m Message) MarshalMSG() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) UnmarshalMSG(data []byte) error {
	return json.Unmarshal(data, &m)
}

// send Message {} type json message over a connection
func sendJSON(conn net.Conn, mtype string, str []string) error {
	msg := Message{}
	msg.Action = mtype
	msg.Body = nil
	for i := 0; i < len(str); i++ {
		msg.Body = append(msg.Body, str[i])
	}
	b, _ := msg.MarshalMSG()
	_, err := conn.Write(b)
	return err
}

// Holding new line flavours for either linux or windows type systems
type Newline struct {
	nl string
}

//sends a message string from server to client
func (s *Session) WriteMessage(str []string) error {
	err := sendJSON(s.conn, ACTION_SENDMESSAGE, str)
	return err
}

//sends a status string from server to client
func (s *Session) WriteStatus(str []string) error {
	err := sendJSON(s.conn, ACTION_SENDSTATUS, str)
	return err
}

//returns newline character for either linux or windows type systems
func (n *Newline) NewLine() string {
	return n.nl
}
