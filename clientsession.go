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

// this function is called by ui events and starts to process the user input
func processInput(conn net.Conn, msg string, u *Ui) error {
	if len(msg) > 0 {
		if msg[0] != CMD_PREFIX {
			return (sendJSON(conn, ACTION_SENDMESSAGE, []string{msg}))
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
					sendJSON(conn, ACTION_EXIT, []string{""})
					fmt.Println(lang.Lookup(actualLocale, "Goodbye"))
					u.win.Close()
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
					return (sendJSON(conn, ACTION_LISTUSERS, []string{""}))
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
						return (sendJSON(conn, ACTION_CHANGENICK, []string{cmd_arguments[0]}))
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
	rmsg := Message{}

	err1 := sendJSON(conn, ACTION_INIT, []string{nick, REVISION})

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
							for i := 0; i < len(rmsg.Body); i++ {
								u.ShowMessage(rmsg.Body[i])
							}
						case ACTION_SENDSTATUS:
							for i := 0; i < len(rmsg.Body); i++ {
								u.ShowStatus(rmsg.Body[i])
							}
						case ACTION_REVISION:
							if rmsg.Body[0] != REVISION {
								log.Printf(lang.Lookup(actualLocale, "Wrong client revision level. Should be: ")+" %s"+lang.Lookup(actualLocale, ", actual: ")+"%s", rmsg.Body[0], REVISION)
								conn.Close()
								u.win.Close()
								os.Exit(1)
							}
						}
					}
				} else {
					log.Printf(lang.Lookup(actualLocale, "Error reading from buffer, most likely server was terminated") + nl.NewLine())
					conn.Close()
					u.win.Close()
					os.Exit(1)

				}
			}
		}()

		myWindow.SetContent(content)
		u.ShowStatus(fmt.Sprintf(lang.Lookup(actualLocale, "Connected to:")+" %s, "+lang.Lookup(actualLocale, "Nickname:")+" %s", connect, nick))
		u.ShowStatus(" ")
		myWindow.Canvas().Focus(u.input)
		myWindow.ShowAndRun()
	} else {
		log.Printf("Send Message failed, error is %v", err1)
		return err1
	}

	return nil
}
