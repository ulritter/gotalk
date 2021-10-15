//go:build !serveronly

package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// display help text in the status are of the window (no server roudtrip required)
func showHelp(u *Ui) {
	u.ShowStatus([]string{" ",
		u.lang.Lookup(u.locale, "Available commands:"),
		u.lang.Lookup(u.locale, "  /exit, /quit, /q - exit program"),
		u.lang.Lookup(u.locale, "  /list - displays active users in room"),
		u.lang.Lookup(u.locale, "  /nick <nickname> - change nickname"),
		u.lang.Lookup(u.locale, "  /help, /?  - this list"),
		" ",
		u.lang.Lookup(u.locale, "Available color controls:"),
		u.lang.Lookup(u.locale, "General:"),
		u.lang.Lookup(u.locale, "A color control followed by space will change"),
		u.lang.Lookup(u.locale, "the color for the remainder of the line."),
		u.lang.Lookup(u.locale, "A color control attached to a word will change"),
		u.lang.Lookup(u.locale, "the color for the word."),
		u.lang.Lookup(u.locale, " "),
		u.lang.Lookup(u.locale, "Usage Example:"),
		u.lang.Lookup(u.locale, "$red this is my $ytext"),
		u.lang.Lookup(u.locale, " "),
		u.lang.Lookup(u.locale, "Color Controls: (long form and short form):"),
		u.lang.Lookup(u.locale, "$red $r $cyan $c $yellow $y $green $g"),
		u.lang.Lookup(u.locale, "$purple $p $white $w $black $b "),
		" "}, false)
}

// display error message in the status are of the window (no server roudtrip required)
func showError(u *Ui) {
	u.ShowStatus([]string{" ",
		u.lang.Lookup(u.locale, "Command error,"),
		u.lang.Lookup(u.locale, "type /help of /? for command descriptions"),
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
					sendMessage(conn, ACTION_EXIT, nil)
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
					return (sendMessage(conn, ACTION_LISTUSERS, nil))
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

func ciao(w fyne.Window, c net.Conn, e int) {
	c.Close()
	w.Close()
	os.Exit(e)
}

// this function is called by main() in the case the app needs to operate as client
// it starts the conenction to the server, listens to the server,
// creates the ui and starts the fyne ui loop
func (a *application) handleClientSession(connect string, config *tls.Config, nick string) error {
	buf := make([]byte, BUFSIZE)
	conn, err := tls.Dial("tcp", connect, config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	guiApp := app.NewWithID(APPTITLE)
	setColors(guiApp)
	myWindow := guiApp.NewWindow(WINTITLE)

	u := &Ui{win: myWindow, app: guiApp, conn: conn, locale: a.config.locale, lang: a.lang}
	content := u.newUi()
	rmsg := Message{}
	// sending init message, format {ACTION_INIT, [{<nickname>}, {<revision level>}]}
	err1 := sendMessage(conn, ACTION_INIT, []string{nick, REVISION})

	//try to catch ^C signals etc
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		//intercept signal and start closing roundtrip
		<-c
		sendMessage(conn, ACTION_EXIT, nil)
	}()

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
						case ACTION_EXIT:
							fmt.Println(a.config.newline + a.lang.Lookup(a.config.locale, "Goodbye") + a.config.newline)
							ciao(myWindow, conn, 0)
						case ACTION_REVISION:
							if rmsg.Body[0] != REVISION {
								fmt.Printf(a.lang.Lookup(a.config.locale, "Wrong client revision level. Should be: ")+" %s"+a.lang.Lookup(a.config.locale, ", actual: ")+"%s"+a.config.newline, rmsg.Body[0], REVISION)
								ciao(myWindow, conn, 1)
							}
						}
					}
				} else {
					a.logger.Printf(a.lang.Lookup(a.config.locale, "Error reading from network, most likely server was terminated") + a.config.newline)
					ciao(myWindow, conn, 1)
				}
			}
		}()

		myWindow.SetContent(content)
		u.ShowStatus([]string{fmt.Sprintf(a.lang.Lookup(a.config.locale, "Connected to:")+" %s, "+a.lang.Lookup(a.config.locale, "Nickname:")+" %s", connect, nick),
			" "}, false)

		myWindow.Canvas().Focus(u.input)

		//intercept quit and start closing roundtrip
		myWindow.SetCloseIntercept(func() {
			sendMessage(conn, ACTION_EXIT, nil)
		})

		myWindow.ShowAndRun()
	} else {
		a.logger.Printf(a.lang.Lookup(a.config.locale, "Send Message failed, error is ")+"%v", err1)
		return err1
	}

	return nil
}
