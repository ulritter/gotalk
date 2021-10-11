package main

import (
	"log"
	"net"
	"os"
	"runtime"
	"testing"

	"fyne.io/fyne/v2/app"
)

type parse_test struct {
	code int
	str  string
}

//TODO extend tests
func TestClientSession(t *testing.T) {

	testBuf := make([]byte, BUFSIZE)

	server, client := net.Pipe()

	testApp := app.NewWithID(APPTITLE)
	setColors(testApp)
	testWindow := testApp.NewWindow(WINTITLE)
	testMsg := Message{}
	testUi := &Ui{win: testWindow, app: testApp}
	testContent := testUi.newUi(client, testNl)

	go func() {
		for {
			testMsg.Body = nil
			n, err := client.Read(testBuf)
			if err == nil {
				err := testMsg.UnmarshalMSG(testBuf[:n])
				if err == nil {
					switch testMsg.Action {
					case ACTION_SENDMESSAGE:
						for i := 0; i < len(testMsg.Body); i++ {
							testUi.ShowMessage(testMsg.Body[i])
						}
					case ACTION_SENDSTATUS:
						for i := 0; i < len(testMsg.Body); i++ {
							testUi.ShowStatus(testMsg.Body[i])
						}
					case ACTION_REVISION:
						if testMsg.Body[0] != REVISION {
							log.Printf(lang.Lookup(actualLocale, "Wrong client revision level. Should be: ")+" %s"+lang.Lookup(actualLocale, ", actual: ")+"%s", testMsg.Body[0], REVISION)
							client.Close()
							testUi.win.Close()
							os.Exit(1)
						}
					}
				}
			} else {
				log.Printf(lang.Lookup(actualLocale, "Error reading from buffer, most likely server was terminated") + testNl.NewLine())
				client.Close()
				server.Close()
				testUi.win.Close()
				os.Exit(1)

			}
		}
	}()

	if runtime.GOOS == "windows" {
		testNl.nl = "\r\n"
	} else {
		testNl.nl = "\n"
	}

	testWindow.SetContent(testContent)

	//TODO: Update tests to adapt to new json messages

	client.Close()
}
