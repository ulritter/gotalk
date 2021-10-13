package main

import (
	"net"
	"os"
	"runtime"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var testApp fyne.App
var testWindow fyne.Window
var testUi *Ui
var testBuf []byte
var testContent fyne.CanvasObject
var testMsg *Message
var testSnd *Message
var testSession *Session
var timeoutDuration time.Duration
var client net.Conn
var server net.Conn

func testSetUp() {

	if runtime.GOOS == "windows" {
		newLine = "\r\n"
	} else {
		newLine = "\n"
	}

	testBuf = make([]byte, BUFSIZE)

	server, client = net.Pipe()

	testApp = app.NewWithID(APPTITLE)
	setColors(testApp)
	testWindow = testApp.NewWindow(WINTITLE)

	testUi = &Ui{win: testWindow, app: testApp, conn: client}
	testContent = testUi.newUi()

	testMsg = &Message{}
	testSnd = &Message{}

	testSession = &Session{conn: server}

	timeoutDuration = 1 * time.Second

	testApp = app.NewWithID(APPTITLE)
	setColors(testApp)
	testWindow = testApp.NewWindow(WINTITLE)
	testUi = &Ui{win: testWindow, app: testApp}
}

func TestMain(m *testing.M) {

	testSetUp()
	exitVal := m.Run()

	os.Exit(exitVal)
}
