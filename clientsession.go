//go:build !serveronly

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"fyne.io/fyne/v2/app"
)

// send a string to the server
func sendToServer(conn net.Conn, str string) error {
	_, err := fmt.Fprint(conn, str)
	return err
}

// display help text in the status are of the window (no server roudtrip required)
func showHelp(u *Ui) {
	u.ShowStatus(" ")
	u.ShowStatus(lang.Lookup(actualLocale, "Available commands:"))
	u.ShowStatus(lang.Lookup(actualLocale, "  /exit, /quit, /q - exit program"))
	u.ShowStatus(lang.Lookup(actualLocale, "  /list - displays active users in room"))
	u.ShowStatus(lang.Lookup(actualLocale, "  /nick <nickname> - change nickname"))
	u.ShowStatus(lang.Lookup(actualLocale, "  /help, /?  - this list"))
	u.ShowStatus(" ")
	u.ShowStatus(lang.Lookup(actualLocale, "Available color controls:"))
	u.ShowStatus(lang.Lookup(actualLocale, "General:"))
	u.ShowStatus(lang.Lookup(actualLocale, "A color control followed by space will change"))
	u.ShowStatus(lang.Lookup(actualLocale, "the color for the remainder of the line."))
	u.ShowStatus(lang.Lookup(actualLocale, "A color control attached to a word will change"))
	u.ShowStatus(lang.Lookup(actualLocale, "the color for the word."))
	u.ShowStatus(lang.Lookup(actualLocale, " "))
	u.ShowStatus(lang.Lookup(actualLocale, "Usage Example:"))
	u.ShowStatus(lang.Lookup(actualLocale, "$red this is my $ytext"))
	u.ShowStatus(lang.Lookup(actualLocale, " "))
	u.ShowStatus(lang.Lookup(actualLocale, "Color Controls: (long form and short form):"))
	u.ShowStatus(lang.Lookup(actualLocale, "$red $r $cyan $c $yellow $y $green $g"))
	u.ShowStatus(lang.Lookup(actualLocale, "$purple $p $white $w $black $b "))
	u.ShowStatus(" ")
}

// display error message in the status are of the window (no server roudtrip required)
func showError(u *Ui) {
	u.ShowStatus(" ")
	u.ShowStatus(lang.Lookup(actualLocale, "Command error,"))
	u.ShowStatus(lang.Lookup(actualLocale, "type /help of /? for command descriptions"))
}

//parse given string whether it is a command or not and take respective action
func parseCommand(conn net.Conn, msg string, u *Ui) int {
	if msg[0] != CMD_PREFIX {
		return CODE_NOCMD
	} else {
		cmdstring := msg[1:]
		cmd := strings.Fields(cmdstring)
		lc := len(cmd)
		switch cmd[0] {
		case CMD_EXIT1:
			if lc == 1 {
				sendToServer(conn, string(CMD_ESCAPE_CHAR)+CMD_EXIT1+string(CMD_ESCAPE_CHAR))
				return CODE_EXIT
			} else {
				showError(u)
				return CODE_DONOTHING
			}
		case CMD_EXIT2:
			if lc == 1 {
				sendToServer(conn, string(CMD_ESCAPE_CHAR)+CMD_EXIT1+string(CMD_ESCAPE_CHAR))
				return CODE_EXIT
			} else {
				showError(u)
				return CODE_DONOTHING
			}
		case CMD_EXIT3:
			if lc == 1 {
				sendToServer(conn, string(CMD_ESCAPE_CHAR)+CMD_EXIT1+string(CMD_ESCAPE_CHAR))
				return CODE_EXIT
			} else {
				showError(u)
				return CODE_DONOTHING
			}
		case CMD_HELP, CMD_HELP1:
			if lc == 1 {
				showHelp(u)
				return CODE_DONOTHING
			} else {
				showError(u)
				return CODE_DONOTHING
			}
		case CMD_LISTUSERS:
			if lc == 1 {
				sendToServer(conn, string(CMD_ESCAPE_CHAR)+CMD_LISTUSERS+string(CMD_ESCAPE_CHAR))
				return CODE_DONOTHING
			} else {
				showError(u)
				return CODE_DONOTHING
			}
		case CMD_CHANGENICK:
			cmd_arguments := cmd[1:]
			if len(cmd_arguments) != 1 {
				showError(u)
				return CODE_DONOTHING
			} else {
				new_nick := cmd_arguments[0]
				sendToServer(conn, string(CMD_ESCAPE_CHAR)+CMD_CHANGENICK+string(CMD_ESCAPE_CHAR)+new_nick+string(CMD_ESCAPE_CHAR))
				return CODE_DONOTHING
			}
		default:
			showError(u)
			return CODE_DONOTHING
		}
	}
}

// TODO: error handling for whole function

// this function is called by ui events and starts to process the user input
func processInput(conn net.Conn, msg string, nl Newline, u *Ui) error {

	if len(msg) > 0 {
		switch cC := parseCommand(conn, msg, u); cC {
		case CODE_NOCMD:
			sendToServer(conn, msg+nl.NewLine())
		case CODE_EXIT:
			fmt.Println(lang.Lookup(actualLocale, "Goodbye"))
			u.win.Close()
			os.Exit(0)
		case CODE_DONOTHING:
			fallthrough
		default:
			break
		}
	}
	return nil
}

// this function is called by main() in the case the app needs to operate as client
// it starts the conenction to the server, listens to the server,
// creates the ui and starts the fyne ui loop
func handleClientSession(connect string, config *tls.Config, nick string, nl Newline) error {
	buf := make([]byte, BUFSIZE)
	conn, err := tls.Dial("tcp", connect, config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	myApp := app.NewWithID(APPTITLE)
	setColors(myApp)
	myWindow := myApp.NewWindow(WINTITLE)

	u := &Ui{win: myWindow, app: myApp}
	content := u.newUi(conn, nl)

	//send user nick and revision level
	fmt.Fprintf(conn, string(CMD_ESCAPE_CHAR)+nick+string(CMD_ESCAPE_CHAR)+" "+REVISION)

	go func() {
		for { // TODO: error handling
			n, err := conn.Read(buf)
			if err != nil {
				log.Printf(lang.Lookup(actualLocale, "Error reading from buffer, most likely server was terminated") + nl.NewLine())
				conn.Close()
				u.win.Close()
				os.Exit(1)
			}
			if buf[0] != CMD_ESCAPE_CHAR {
				msg := string(buf[:n])
				u.ShowMessage(msg)
			} else {
				if n > 1 && buf[1] == CMD_ESCAPE_CHAR {
					//2 escape characters means wrong revision level
					expectedRevision := string(buf[2:])
					log.Printf(lang.Lookup(actualLocale, "Wrong client revision level. Should be: ")+" %s"+lang.Lookup(actualLocale, ", actual: ")+"%s", expectedRevision, REVISION)
					conn.Close()
					u.win.Close()
					os.Exit(1)
				} else {
					msg := string(buf[1:n])
					u.ShowStatus(msg)
				}
			}
		}
	}()

	myWindow.SetContent(content)
	u.ShowStatus(fmt.Sprintf(lang.Lookup(actualLocale, "Connected to:")+" %s, "+lang.Lookup(actualLocale, "Nickname:")+" %s", connect, nick))
	u.ShowStatus(" ")
	myWindow.Canvas().Focus(u.input)
	myWindow.ShowAndRun()

	return nil
}
