package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func sendServerCommand(conn net.Conn, cmd string) error {
	_, err := fmt.Fprint(conn, cmd)
	return err
}

func printHelp(nl Newline) {
	// TODO: create help text
	fmt.Print("Available commands:" + nl.NewLine())
	fmt.Print("- `/exit` - terminate connection and exit" + nl.NewLine())
	fmt.Print("- `/list` - displays active users in room" + nl.NewLine())
	fmt.Print("- `/nick <nickname>` - change nickname" + nl.NewLine())
}

func printError(nl Newline) {
	fmt.Print("Command error, type /help of /? for command descriptions" + nl.NewLine())
}

func parseCommand(conn net.Conn, msg string, nl Newline) int {
	if msg[0] != CMD_PREFIX {
		return CODE_NOCMD
	} else {
		cmdstring := msg[1:]
		cmd := strings.Fields(cmdstring)
		switch cmd[0] {
		case CMD_EXIT:
			sendServerCommand(conn, string(CMD_ESCAPE_CHAR)+CMD_EXIT+string(CMD_ESCAPE_CHAR))
			return CODE_EXIT
		case CMD_HELP, CMD_HELP1:
			printHelp(nl)
			return CODE_DONOTHING
		case CMD_LISTUSERS:
			sendServerCommand(conn, string(CMD_ESCAPE_CHAR)+CMD_LISTUSERS+string(CMD_ESCAPE_CHAR))
			return CODE_DONOTHING
		case CMD_CHANGENICK:
			cmd_arguments := cmd[1:]
			if len(cmd_arguments) != 1 {
				printError(nl)
				return CODE_DONOTHING
			} else {
				new_nick := cmd_arguments[0]
				sendServerCommand(conn, string(CMD_ESCAPE_CHAR)+CMD_CHANGENICK+string(CMD_ESCAPE_CHAR)+new_nick+string(CMD_ESCAPE_CHAR))
				return CODE_DONOTHING
			}
		default:
			return CODE_DONOTHING
		}
	}
}

// TODO: error handling for whole function

func processInput(conn net.Conn, msg string, nl Newline) error {

	if len(msg) > 0 {
		switch cC := parseCommand(conn, msg, nl); cC {
		case CODE_NOCMD:
			fmt.Fprintln(conn, msg)
		case CODE_EXIT:
			conn.Close()
			os.Exit(0)
		case CODE_DONOTHING:
			fallthrough
		default:
			break
		}
	}
	return nil
}

func clientDialogHandling(connect string, config *tls.Config, nick string, nl Newline) error {
	buf := make([]byte, BUFSIZE)
	conn, err := tls.Dial("tcp", connect, config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	myApp := app.NewWithID(APPTITLE)
	myWindow := myApp.NewWindow(WINTITLE)
	myWindow.Resize(fyne.NewSize(1200, 800))

	ui := &Ui{win: myWindow}
	ui_content := ui.makeUi(conn, nl)

	ui.ShowStatus(fmt.Sprintf("Connected to: %s, Nickname: %s %s", connect, nick, nl.NewLine()))

	fmt.Fprintf(conn, string(CMD_ESCAPE_CHAR)+nick+string(CMD_ESCAPE_CHAR))

	go func() {
		for { // TODO: error handling
			n, err := conn.Read(buf)
			if err != nil {
				log.Printf("Error reading from buffer, most likely server was terminated" + nl.NewLine())
				conn.Close()
				os.Exit(1)
			}
			if buf[0] != CMD_ESCAPE_CHAR {
				msg := string(buf[:n])
				ui.ShowMessage(msg)
			} else {
				msg := string(buf[1:n])
				ui.ShowStatus(msg)
			}
		}
	}()

	myWindow.SetContent(ui_content)
	myWindow.Canvas().Focus(ui.input)
	myWindow.ShowAndRun()

	return nil
}
