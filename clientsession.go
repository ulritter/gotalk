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

// display help text in the status are of the window (no server roudtrip required)
func showHelp(u *Ui) {
	u.ShowStatus([]string{" ",
		lang.Lookup(actualLocale, "Available commands:"),
		lang.Lookup(actualLocale, "  /exit, /quit, /q - exit program"),
		lang.Lookup(actualLocale, "  /list - displays active users in room"),
		lang.Lookup(actualLocale, "  /nick <nickname> - change nickname"),
		lang.Lookup(actualLocale, "  /help, /?  - this list"),
		" ",
		lang.Lookup(actualLocale, "Available color controls:"),
		lang.Lookup(actualLocale, "General:"),
		lang.Lookup(actualLocale, "A color control followed by space will change"),
		lang.Lookup(actualLocale, "the color for the remainder of the line."),
		lang.Lookup(actualLocale, "A color control attached to a word will change"),
		lang.Lookup(actualLocale, "the color for the word."),
		lang.Lookup(actualLocale, " "),
		lang.Lookup(actualLocale, "Usage Example:"),
		lang.Lookup(actualLocale, "$red this is my $ytext"),
		lang.Lookup(actualLocale, " "),
		lang.Lookup(actualLocale, "Color Controls: (long form and short form):"),
		lang.Lookup(actualLocale, "$red $r $cyan $c $yellow $y $green $g"),
		lang.Lookup(actualLocale, "$purple $p $white $w $black $b "),
		" "}, false)
}

// display error message in the status are of the window (no server roudtrip required)
func showError(u *Ui) {
	u.ShowStatus([]string{" ",
		lang.Lookup(actualLocale, "Command error,"),
		lang.Lookup(actualLocale, "type /help of /? for command descriptions"),
	}, false)
}

// this function is called by ui events and starts to process the user input
func parseInput(conn net.Conn, msg string, u *Ui) error {
	if len(msg) > 0 {
		if msg[0] != CMD_PREFIX {
			return (sendMessage(conn, ACTION_SENDMESSAGE, []string{msg}))
		} else {
			cmd := strings.Fields(msg)
			lc := len(cmd)
			cmd[0] = cmd[0][1:] // strip leading command symbol

			switch cmd[0] {
			case CMD_EXIT1:
				fallthrough
			case CMD_EXIT2:
				fallthrough
			case CMD_EXIT3:
				if lc == 1 {
					sendMessage(conn, ACTION_EXIT, []string{""})
					fmt.Println(lang.Lookup(actualLocale, "Goodbye"))
					u.win.Close()
					conn.Close()
					os.Exit(1)
				} else {
					showError(u)
					return nil
				}
			case CMD_HELP:
				fallthrough
			case CMD_HELP1:
				if lc == 1 {
					showHelp(u)
					return nil
				} else {
					showError(u)
					return nil
				}
			case CMD_LISTUSERS:
				if lc == 1 {
					return (sendMessage(conn, ACTION_LISTUSERS, []string{""}))
				} else {
					showError(u)
					return nil
				}
			case CMD_CHANGENICK:
				cmdErr := false
				if lc == 2 {
					cmd_arguments := cmd[1:]
					if len(cmd_arguments) != 1 || len(cmd_arguments[0]) == 0 {
						cmdErr = true
					} else {
						return (sendMessage(conn, ACTION_CHANGENICK, []string{cmd_arguments[0]}))
					}
				} else {
					cmdErr = true
				}
				if cmdErr {
					showError(u)
					return nil
				}

			default:
				showError(u)
				return nil
			}
		}
	}
	return nil
}

// this function is called by main() in the case the app needs to operate as client
// it starts the conenction to the server, listens to the server,
// creates the ui and starts the fyne ui loop
func handleClientSession(connect string, config *tls.Config, nick string) error {
	buf := make([]byte, BUFSIZE)
	conn, err := tls.Dial("tcp", connect, config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	myApp := app.NewWithID(APPTITLE)
	setColors(myApp)
	myWindow := myApp.NewWindow(WINTITLE)

	u := &Ui{win: myWindow, app: myApp, conn: conn}
	content := u.newUi()
	rmsg := Message{}
	// sending format {ACTION_INIT, [{<nickname>}, {<revision level>}]}
	err1 := sendMessage(conn, ACTION_INIT, []string{nick, REVISION})

	if err1 == nil {
		go func() {
			for {
				rmsg.Body = nil
				n, err := conn.Read(buf)
				if err == nil {
					err := rmsg.UnmarshalMSG(buf[:n])
					if err == nil {
						switch rmsg.Action {
						case ACTION_SENDMESSAGE:
							u.ShowMessage(rmsg.Body, false)
						case ACTION_SENDSTATUS:
							u.ShowStatus(rmsg.Body, false)
						case ACTION_REVISION:
							if rmsg.Body[0] != REVISION {
								log.Printf(lang.Lookup(actualLocale, "Wrong client revision level. Should be: ")+" %s"+lang.Lookup(actualLocale, ", actual: ")+"%s", rmsg.Body[0], REVISION)
								u.win.Close()
								conn.Close()
								os.Exit(1)
							}
						}
					}
				} else {
					log.Printf(lang.Lookup(actualLocale, "Error reading from buffer, most likely server was terminated") + newLine)
					u.win.Close()
					conn.Close()
					os.Exit(1)
				}
			}
		}()

		myWindow.SetContent(content)
		u.ShowStatus([]string{fmt.Sprintf(lang.Lookup(actualLocale, "Connected to:")+" %s, "+lang.Lookup(actualLocale, "Nickname:")+" %s", connect, nick),
			" "}, false)
		myWindow.Canvas().Focus(u.input)
		myWindow.ShowAndRun()
	} else {
		log.Printf("Send Message failed, error is %v", err1)
		return err1
	}

	return nil
}
