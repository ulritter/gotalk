package main

import (
	"crypto/tls"
	"fmt"
	"gotalk/models"
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
		if msg[0] != models.CMD_PREFIX {
			return (models.SendMessage(conn, models.ACTION_SENDMESSAGE, []string{msg}))
		} else {
			cmd := strings.Fields(msg)
			lc := len(cmd)
			cmd[0] = cmd[0][1:] // strip leading command symbol

			switch cmd[0] {
			case models.CMD_EXIT1:
				fallthrough
			case models.CMD_EXIT2:
				fallthrough
			case models.CMD_EXIT3:
				if lc == 1 {
					models.SendMessage(conn, models.ACTION_EXIT, nil)
				} else {
					showError(u)
					return nil
				}
			case models.CMD_HELP:
				fallthrough
			case models.CMD_HELP1:
				if lc == 1 {
					showHelp(u)
					return nil
				} else {
					showError(u)
					return nil
				}
			case models.CMD_LISTUSERS:
				if lc == 1 {
					return (models.SendMessage(conn, models.ACTION_LISTUSERS, nil))
				} else {
					showError(u)
					return nil
				}
			case models.CMD_CHANGENICK:
				cmdErr := false
				if lc == 2 {
					cmd_arguments := cmd[1:]
					if len(cmd_arguments) != 1 || len(cmd_arguments[0]) == 0 {
						cmdErr = true
					} else {
						return (models.SendMessage(conn, models.ACTION_CHANGENICK, []string{cmd_arguments[0]}))
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
func handleClientSession(a *models.Application, nick string) error {
	buf := make([]byte, models.BUFSIZE)
	adr := a.Config.Addr + a.Config.Port
	a.Logger.Println("address: ", adr)
	conn, err := tls.Dial("tcp", adr, a.Config.TLSconfig)
	if err != nil {
		fmt.Println(err)
		return err
	}

	guiApp := app.NewWithID(APPTITLE)
	setColors(guiApp)
	myWindow := guiApp.NewWindow(WINTITLE)

	u := &Ui{win: myWindow, app: guiApp, conn: conn, locale: a.Config.Locale, lang: a.Lang}
	content := u.newUi()
	rmsg := models.Message{}
	// sending init message, format {models.ACTION_INIT, [{<nickname>}, {<revision level>}]}
	err1 := models.SendMessage(conn, models.ACTION_INIT, []string{nick, models.REVISION})

	//try to catch ^C signals etc
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		//intercept signal and start closing roundtrip
		<-c
		models.SendMessage(conn, models.ACTION_EXIT, nil)
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
						case models.ACTION_SENDMESSAGE:
							u.ShowMessage(rmsg.Body, false)
						case models.ACTION_SENDSTATUS:
							u.ShowStatus(rmsg.Body, false)
						case models.ACTION_EXIT:
							fmt.Println(a.Config.Newline + a.Lang.Lookup(a.Config.Locale, "Goodbye") + a.Config.Newline)
							ciao(myWindow, conn, 0)
						case models.ACTION_REVISION:
							if rmsg.Body[0] != models.REVISION {
								fmt.Printf(a.Lang.Lookup(a.Config.Locale,
									"Wrong client revision level. Should be: ")+" %s"+a.Lang.Lookup(a.Config.Locale,
									", actual: ")+"%s"+a.Config.Newline, rmsg.Body[0], models.REVISION)
								ciao(myWindow, conn, 1)
							}
						}
					}
				} else {
					a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "Error reading from network, most likely server was terminated") + a.Config.Newline)
					ciao(myWindow, conn, 1)
				}
			}
		}()

		myWindow.SetContent(content)
		u.ShowStatus([]string{fmt.Sprintf(a.Lang.Lookup(a.Config.Locale, "Connected to:")+" %s, "+a.Lang.Lookup(a.Config.Locale, "Nickname:")+" %s", adr, nick),
			" "}, false)

		myWindow.Canvas().Focus(u.input)

		//intercept quit and start closing roundtrip
		myWindow.SetCloseIntercept(func() {
			models.SendMessage(conn, models.ACTION_EXIT, nil)
		})

		myWindow.ShowAndRun()
	} else {
		a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "Send Message failed, error is ")+"%v", err1)
		return err1
	}

	return nil
}
