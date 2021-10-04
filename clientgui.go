package main

import (
	"image/color"
	"log"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const MAXLINES = 1024
const APPTITLE = "cooltide"
const WINTITLE = "gotalk"

const MESSAGEWIDTH = 618
const MESSAGEHEIGHT = 400
const STATUSWIDTH = 382
const STATUSHEIGHT = 440

var HEADERCOLOR = color.RGBA{255, 200, 100, 255}
var HEADERSTYLE = fyne.TextStyle{Bold: true, Italic: false}

var MESSAGECOLOR = color.RGBA{200, 255, 100, 255}
var MESSAGESTYLE = fyne.TextStyle{Bold: false, Italic: false}

var STATUSCOLOR = color.RGBA{180, 180, 255, 255}
var STATUSSTYLE = fyne.TextStyle{Bold: false, Italic: true}

type MessageLine struct {
	txt string
	obj *canvas.Text
	col color.Color
}

// data strcture to hold the ui elements
type Ui struct {
	mHeader *canvas.Text
	sHeader *canvas.Text
	mBox    *fyne.Container
	sBox    *fyne.Container
	mMsgs   []MessageLine
	sMsgs   []MessageLine
	win     fyne.Window
	input   *widget.Entry
	mScroll *container.Scroll
	sScroll *container.Scroll
	button  *widget.Button
	ui_ref  *Ui
}

// create new ui with fyne elements
func (u *Ui) newUi(conn net.Conn, nl Newline) fyne.CanvasObject {
	u.ui_ref = u

	u.mBox = container.NewVBox()
	u.mScroll = container.NewScroll(u.mBox)
	u.mScroll.SetMinSize(fyne.NewSize(MESSAGEWIDTH, MESSAGEHEIGHT))
	u.mHeader = canvas.NewText(lang.Lookup(locale, " Messages"), HEADERCOLOR)
	u.mHeader.TextStyle = HEADERSTYLE

	u.sBox = container.NewVBox()
	u.sScroll = container.NewScroll(u.sBox)
	u.sScroll.SetMinSize(fyne.NewSize(STATUSWIDTH, STATUSHEIGHT))
	u.sHeader = canvas.NewText(lang.Lookup(locale, " Status Info"), HEADERCOLOR)
	u.sHeader.TextStyle = HEADERSTYLE

	u.input = widget.NewEntry()

	vSeparator := canvas.NewRectangle(color.Gray{})
	vSeparator.Resize(fyne.NewSize(3, 400))

	u.button = widget.NewButton(">>", func() {
		if len(u.input.Text) > 0 {
			processInput(conn, u.input.Text, nl, u)
			u.input.SetText("")
		}
		u.mScroll.Refresh()
		u.mScroll.ScrollToBottom()
		u.win.Canvas().Focus(u.input)
	})

	u.input.OnSubmitted = func(text string) {
		if len(u.input.Text) > 0 {
			processInput(conn, u.input.Text, nl, u)
			u.input.SetText("")
		}
		u.mScroll.Refresh()
		u.mScroll.ScrollToBottom()
		u.win.Canvas().Focus(u.input)
	}

	inputline := container.NewBorder(nil, nil, nil, u.button, u.input)
	left := container.NewBorder(u.mHeader, inputline, nil, nil, u.mScroll)
	right := container.NewBorder(nil, nil, vSeparator, container.NewBorder(u.sHeader, nil, nil, nil, u.sScroll))
	content := container.NewBorder(nil, nil, nil, right, left)

	return container.New(layout.NewMaxLayout(), content)
}

//display a user message in the (left hand) message area of the ui
func (u *Ui) ShowMessage(msg string) {
	refreshBoxContent(&u.mMsgs, u.mBox, msg, MESSAGESTYLE, MESSAGECOLOR)
	u.mMsgs[len(u.mMsgs)-1].obj.SetMinSize(fyne.NewSize(MESSAGEWIDTH, MESSAGEHEIGHT))
	u.mScroll.SetMinSize(fyne.NewSize(MESSAGEWIDTH, MESSAGEHEIGHT))
	u.mBox.Refresh()
	u.mScroll.Refresh()
	u.mScroll.ScrollToBottom()
	u.win.Canvas().Focus(u.input)
}

//display a status message in the (right hand) status area of the ui
func (u *Ui) ShowStatus(msg string) {
	refreshBoxContent(&u.sMsgs, u.sBox, msg, STATUSSTYLE, STATUSCOLOR)
	u.sBox.Refresh()
	u.sScroll.Refresh()
	u.sScroll.ScrollToBottom()
	u.win.Canvas().Focus(u.input)
}

func refreshBoxContent(m *[]MessageLine, b *fyne.Container, msg string, s fyne.TextStyle, c color.Color) {
	nlines := len((*m))
	t := MessageLine{
		txt: msg,
		col: c,
	}
	t.obj = canvas.NewText(msg, t.col)
	t.obj.TextStyle = s

	cs := (*b).Size()
	os := fyne.MeasureText(t.obj.Text, t.obj.TextSize, t.obj.TextStyle)
	as := cs.Width - os.Width

	if as < 0 {
		log.Printf("String size exceeds container size: os: %f, cs: %f", os.Width, cs.Width)
	}

	if nlines < MAXLINES {
		(*m) = append((*m), t)
		(*b).Add(t.obj)
		(*b).Objects[nlines].Resize(cs)
	} else {
		for i := 0; i < nlines-1; i++ {
			(*m)[i] = (*m)[i+1]
			(*b).Objects[i] = (*b).Objects[i+1]
		}
		(*m)[nlines-1] = t
		(*b).Objects[nlines-1] = t.obj
		(*b).Objects[nlines-1].Resize(cs)
	}
	(*b).Resize(cs)
}
