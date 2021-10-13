package main

import (
	"net"
	"os"
	"runtime"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var testApp fyne.App
var testWindow fyne.Window
var testUi *Ui
var testConn net.Conn

func testNlInit() {
	if runtime.GOOS == "windows" {
		newLine = "\r\n"
	} else {
		newLine = "\n"
	}
}

func testUiSetUp() {

	testNlInit()
	testApp = app.NewWithID(APPTITLE)
	setColors(testApp)
	testWindow = testApp.NewWindow(WINTITLE)
	testUi = &Ui{win: testWindow, app: testApp}
}

func TestMain(m *testing.M) {

	testUiSetUp()
	exitVal := m.Run()

	os.Exit(exitVal)
}
