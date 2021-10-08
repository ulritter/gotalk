package main

import (
	"net"
	"runtime"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var testApp fyne.App
var testWindow fyne.Window
var testUi *Ui
var testConn net.Conn
var testNl Newline

func testNlInit(t *testing.T) {
	if runtime.GOOS == "windows" {
		testNl.nl = "\r\n"
	} else {
		testNl.nl = "\n"
	}
}

func testUiSetUp(t *testing.T) {

	testNlInit(t)
	testApp = app.NewWithID(APPTITLE)
	setColors(testApp)
	testWindow = testApp.NewWindow(WINTITLE)
	testUi = &Ui{win: testWindow, app: testApp}
}

/*
	inputline := container.NewBorder(nil, nil, nil, u.button, u.input)
	left := container.NewBorder(u.mHeader, inputline, nil, nil, u.mScroll)
	right := container.NewBorder(nil, nil, vSeparator, container.NewBorder(u.sHeader, nil, nil, nil, u.sScroll))
	content := container.NewBorder(nil, nil, nil, right, left)
	return container.New(layout.NewMaxLayout(), content)
*/
func TestUiLayout(t *testing.T) {

	testUiSetUp(t)
	test_content := testUi.newUi(testConn, testNl)
	o := test_content.(*fyne.Container).Objects
	if len(o) != 1 {
		t.Log("Wrong Ui Structure")
		t.Fail()
	}
	oo := test_content.(*fyne.Container).Objects[0].(*fyne.Container).Objects
	if len(oo) != 2 {
		t.Log("Wrong Ui Structure, expecting left and right containers")
		t.Fail()
	}
	ooo_r := test_content.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects
	if len(ooo_r) != 3 {
		t.Log("Wrong Ui Structure, expecting two elements on left half")
		t.Fail()
	}

	ooo_l := test_content.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects
	if len(ooo_l) != 2 {
		t.Log("Wrong Ui Structure, expecting three elements on right half")
		t.Fail()
	}

}

func TestMessage(t *testing.T) {
	for i := 0; i < MAXLINES+1; i++ {
		testUi.ShowMessage("test message", true)
	}
	if len(testUi.mMsgs) != MAXLINES ||
		len(testUi.mBox.Objects) != MAXLINES {
		t.Log("Message append and buffer handling failed")
		t.Fail()
	}
}

func TestStatus(t *testing.T) {
	for i := 0; i < MAXLINES+1; i++ {
		testUi.ShowStatus("test status", true)
	}
	if len(testUi.sMsgs) != MAXLINES ||
		len(testUi.sBox.Objects) != MAXLINES {
		t.Log("Status append and buffer handling failed")
		t.Fail()
	}
}
