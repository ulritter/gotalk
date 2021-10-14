package main

import (
	"log"
	"net"
	"os"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Xuanwo/go-locale"
	language "github.com/moemoe89/go-localization"
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
var newline string

func testSetUp() {
	cfg := language.New()
	cfg.BindPath(LANGFILE)
	cfg.BindMainLocale("en")
	lang, lerr := cfg.Init()
	if lerr != nil {
		panic(lerr)
	}

	tag, err := locale.Detect()
	appConfig := config{
		newline: newLine(),
	}

	if err != nil {
		log.Fatal(err)
		appConfig.locale = "en"
	} else {
		if len(tag.String()) > 2 {
			appConfig.locale = tag.String()[:2]
		} else {
			if len(tag.String()) == 2 {
				appConfig.locale = tag.String()
			}
		}
	}

	testBuf = make([]byte, BUFSIZE)

	server, client = net.Pipe()

	testApp = app.NewWithID(APPTITLE)
	setColors(testApp)
	testWindow = testApp.NewWindow(WINTITLE)

	testUi = &Ui{win: testWindow, app: testApp, conn: client, locale: appConfig.locale, lang: lang}
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
