package models

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net"

	language "github.com/moemoe89/go-localization"
)

const RAWFILE = "https://raw.githubusercontent.com/ulritter/gotalk-app/main/language.json"
const LANGFILE = "../language.json"

// actions for client <-> server communication
const ACTION_CHANGENICK = "changenick"
const ACTION_SENDMESSAGE = "message"
const ACTION_LISTUSERS = "listusers"
const ACTION_REVISION = "revision"
const ACTION_SENDSTATUS = "status"
const ACTION_EXIT = "exit"
const ACTION_INIT = "init"

// end user commands on ui
const CMD_PREFIX = '/'
const CMD_EXIT1 = "exit"
const CMD_EXIT2 = "quit"
const CMD_EXIT3 = "q"
const CMD_CHANGENICK = "nick"
const CMD_LISTUSERS = "list"
const CMD_HELP = "help"
const CMD_HELP1 = "?"

const BUFSIZE = 4096

const REVISION = "0.8.3"

// Event type for messages
type MessageEvent struct {
	Msg string
}

// Event type for users joining
type UserJoinedEvent struct {
}

// Event type for users leaving
type UserLeftEvent struct {
	User *User
	Msg  string
}

// Event type for users changing their nick
type UserChangedNickEvent struct {
	User *User
	Nick string
}

// Event type for users requesting a room user list
type ListUsersEvent struct {
	User *User
}

// Commmunication structure between connection handler and user session
type ClientInput struct {
	User  *User
	Event interface{}
}

// User Info
type User struct {
	Name       string
	Session    *Session
	Timejoined string
}

// Structure holding the connection for the session
type Session struct {
	Conn net.Conn
}

// Structure holding the users the room
type Room struct {
	Users []*User
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
func SendJSONMessage(conn net.Conn, mtype string, str []string) error {
	msg := Message{}
	msg.Action = mtype
	msg.Body = nil
	for i := 0; i < len(str); i++ {
		msg.Body = append(msg.Body, str[i])
	}
	b, err := msg.MarshalMSG()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = conn.Write(b)
	return err
}

//sends a message string from server to client
func (s *Session) WriteMessage(str []string) error {
	err := SendJSONMessage(s.Conn, ACTION_SENDMESSAGE, str)
	return err
}

//sends a status string from server to client
func (s *Session) WriteStatus(str []string) error {
	err := SendJSONMessage(s.Conn, ACTION_SENDSTATUS, str)
	return err
}

type Config struct {
	Server    bool
	Addr      string
	Port      string
	Env       string
	Newline   string
	Locale    string
	Ch        chan ClientInput
	Conn      net.Conn
	TLSconfig *tls.Config
}

type Application struct {
	Config  Config
	Logger  *log.Logger
	Lang    *language.Config
	Version string
}
