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

// hard-wired key / certificate for test purposes
const serverKey = `-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEII1rgpj11/thNhne5jla7jJhdya34h3amxH+iYXduCC4oAoGCCqGSM49
AwEHoUQDQgAEeYUzwiEc8qzUWYOJPKEbKWlVvexG40pqsQA0eSaHRXXPV2gnrsWT
mUXfkjuYEpEREZgZH1HpiCjpy5hLAuQp7w==
-----END EC PRIVATE KEY-----
`
const rootCert = `-----BEGIN CERTIFICATE-----
MIIBCTCBsAIJAO2vPzVY2coAMAoGCCqGSM49BAMCMA0xCzAJBgNVBAYTAkRFMB4X
DTIxMDkxNjA4MzYwOFoXDTMxMDkxNDA4MzYwOFowDTELMAkGA1UEBhMCREUwWTAT
BgcqhkjOPQIBBggqhkjOPQMBBwNCAAR5hTPCIRzyrNRZg4k8oRspaVW97EbjSmqx
ADR5JodFdc9XaCeuxZOZRd+SO5gSkRERmBkfUemIKOnLmEsC5CnvMAoGCCqGSM49
BAMCA0gAMEUCIDKpwLhrMJWUJFcI5NC4YqQzcDaAHZTbOgRRIHsDZyCIAiEA+JZV
CgCeRmOCnFDNFFf9fl6ABui6hpRmwHj2dAK4e+U=
-----END CERTIFICATE-----
`

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
