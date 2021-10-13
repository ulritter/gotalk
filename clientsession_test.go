package main

import (
	"net"
	"runtime"
	"testing"
	"time"

	"fyne.io/fyne/v2/app"
)

//TODO extend tests
func TestClientSession(t *testing.T) {

	testBuf := make([]byte, BUFSIZE)

	server, client := net.Pipe()

	testApp := app.NewWithID(APPTITLE)
	setColors(testApp)
	testWindow := testApp.NewWindow(WINTITLE)

	testUi := &Ui{win: testWindow, app: testApp, conn: client}
	testContent := testUi.newUi()

	testMsg := Message{}
	testSnd := Message{}

	testSession := &Session{conn: server}

	quit := make(chan bool)

	timeoutDuration := 1 * time.Second

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				testMsg.Body = nil
				client.SetReadDeadline(time.Now().Add(timeoutDuration))
				n, err := client.Read(testBuf)
				if err == nil {
					err := testMsg.UnmarshalMSG(testBuf[:n])
					if err == nil {
						switch testMsg.Action {
						case ACTION_SENDMESSAGE:
							if len(testMsg.Body) == 0 {
								t.Errorf("bad test user message")
								t.Fail()
							} else {
								t.Log("ACTION_SENDMESSAGE passed")
							}

						case ACTION_SENDSTATUS:
							if len(testMsg.Body) == 0 {
								t.Errorf("bad test status mesage")
								t.Fail()
							} else {
								t.Log("ACTION_SENDSTATUS passed")
							}
						case ACTION_REVISION:
							if len(testMsg.Body) != 1 {
								t.Errorf("bad test revision message")
								t.Fail()
							} else {
								t.Log("ACTION_REVISION passed")
							}
						}
					}
				}
			}

		}
	}()

	if runtime.GOOS == "windows" {
		newLine = "\r\n"
	} else {
		newLine = "\n"
	}

	testWindow.SetContent(testContent)

	testSnd.Body = nil
	testSnd.Body = append(testSnd.Body, "Testmessage")
	testSession.WriteMessage(testSnd.Body)

	testSnd.Body = nil
	testSnd.Body = append(testSnd.Body, "Test status")
	testSession.WriteStatus(testSnd.Body)

	sendMessage(testSession.conn, ACTION_REVISION, []string{REVISION})

	quit <- true

	client.Close()
	server.Close()
}
