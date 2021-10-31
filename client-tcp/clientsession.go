package main

import (
	"crypto/tls"
	"fmt"
	"gotalk/constants"
	"gotalk/models"
	"net"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

//wrap up and exit
func ciao(w fyne.Window, c net.Conn, e int) {
	c.Close()
	w.Close()
	os.Exit(e)
}

// this function is called by main() in the case the app needs to operate as client
// it starts the conenction to the server, listens to the server,
// creates the ui and starts the fyne ui loop
func handleClientSession(a *models.Application, nick string) error {
	buf := make([]byte, constants.BUFSIZE)
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

	// send init message, format {models.ACTION_INIT, [{<nickname>}, {<revision level>}]}
	err = models.SendJSONMessage(conn, constants.ACTION_INIT, []string{nick, constants.REVISION})

	if err == nil {
		//try to intercept ^C signals etc
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
		go func() {
			<-c
			models.SendJSONMessage(conn, constants.ACTION_EXIT, nil)
		}()

		go func() {
			for {
				rmsg.Body = nil
				n, err := conn.Read(buf)
				if err == nil {
					err := rmsg.UnmarshalMSG(buf[:n])
					if err == nil {
						switch rmsg.Action {
						case constants.ACTION_SENDMESSAGE:
							u.ShowMessage(rmsg.Body, false)
						case constants.ACTION_SENDSTATUS:
							u.ShowStatus(rmsg.Body, false)
						case constants.ACTION_EXIT:
							fmt.Println(a.Config.Newline + a.Lang.Lookup(a.Config.Locale, "Goodbye") + a.Config.Newline)
							ciao(myWindow, conn, 0)
						case constants.ACTION_REVISION:
							if rmsg.Body[0] != constants.REVISION {
								fmt.Printf(a.Lang.Lookup(a.Config.Locale,
									"Wrong client revision level. Should be: ")+" %s"+a.Lang.Lookup(a.Config.Locale,
									", actual: ")+"%s"+a.Config.Newline, rmsg.Body[0], constants.REVISION)
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
			models.SendJSONMessage(conn, constants.ACTION_EXIT, nil)
		})

		myWindow.ShowAndRun()
	} else {
		a.Logger.Printf(a.Lang.Lookup(a.Config.Locale, "Send Message failed, error is ")+"%v", err)
		return err
	}

	return nil
}
