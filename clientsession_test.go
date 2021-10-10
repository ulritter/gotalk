package main

import (
	"net"
	"testing"
)

type parse_test struct {
	code int
	str  string
}

//TODO extend tests
func TestClientSession(t *testing.T) {

	testBuf := make([]byte, BUFSIZE)
	server, client := net.Pipe()

	go func() {
		exits := 0
		for {
			n, err := server.Read(testBuf)
			if err == nil {
				pattern := string(testBuf[:n])
				if pattern[0] == CMD_ESCAPE_CHAR && (pattern[1:] == CMD_EXIT1 || pattern[1:] == CMD_EXIT2 || pattern[1:] == CMD_EXIT3) {
					exits++
					if exits >= 3 {
						server.Close()
						break
					}
				}

			}
		}
	}()

	processInput(client, "Test Input", testNl, testUi)

	codes := [5]parse_test{{CODE_NOCMD, "Test"}, {CODE_EXIT, "/" + CMD_EXIT1}, {CODE_EXIT, "/" + CMD_EXIT2}, {CODE_EXIT, "/" + CMD_EXIT3}, {CODE_DONOTHING, "/Test"}}
	for i := range codes {
		p := parseCommand(client, codes[i].str, testUi)
		if p != codes[i].code {
			t.Logf("parseCommand expected %d, got %d\n", codes[i].code, p)
			t.Fail()
		}
	}
	client.Close()
}
