package main

import (
	"image/color"
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

type Ui struct {
	mMsgs   []string
	sMsgs   []string
	win     fyne.Window
	input   *widget.Entry
	mLabel  *widget.Label
	mOutput *widget.Label
	mScroll *container.Scroll
	sLabel  *widget.Label
	sOutput *widget.Label
	sScroll *container.Scroll
}

func (u *Ui) makeUi(conn net.Conn, nl Newline) fyne.CanvasObject {
	u.input = widget.NewEntry()
	u.mLabel = widget.NewLabel("Messages")
	u.mOutput = widget.NewLabel("")
	u.mScroll = container.NewScroll(u.mOutput)
	u.sLabel = widget.NewLabel("Status Info")
	u.sOutput = widget.NewLabel("")
	u.sScroll = container.NewScroll(u.sOutput)

	hline1 := canvas.NewRectangle(color.Gray{})
	hline1.Resize(fyne.NewSize(800, 3))
	hline2 := canvas.NewRectangle(color.Gray{})
	hline2.Resize(fyne.NewSize(400, 3))
	vline := canvas.NewRectangle(color.Gray{})
	vline.Resize(fyne.NewSize(3, 400))

	u.mScroll.SetMinSize(fyne.NewSize(500, 400))
	u.sScroll.SetMinSize(fyne.NewSize(400, 440))

	left := container.NewVBox(u.mLabel, u.mScroll, u.input)
	right := container.NewHBox(vline, container.NewVBox(u.sLabel, u.sScroll))
	content := container.New(
		layout.NewBorderLayout(nil, nil, left, right),
		left,
		right)

	/* displayArea := container.NewHBox(
		container.NewVBox(u.mLabel,  u.mScroll),
		//layout.NewSpacer(),
		vline,
		container.NewVBox(u.sLabel,  u.sScroll))

	content := container.New(
		layout.NewBorderLayout(displayArea, u.input, nil, nil),
		displayArea,
		u.input) */

	u.input.OnSubmitted = func(text string) {
		processInput(conn, u.input.Text, nl)
		u.input.SetText("")
		u.mScroll.ScrollToBottom()
		u.win.Canvas().Focus(u.input)
	}
	return container.New(layout.NewMaxLayout(), content)
}

func (u *Ui) ShowMessage(msg string) {
	nlines := len(u.mMsgs)
	if nlines < MAXLINES {
		u.mMsgs = append(u.mMsgs, msg)
	} else {
		for i := 0; i < nlines-1; i++ {
			u.mMsgs[i] = u.mMsgs[i+1]
		}
		u.mMsgs[nlines-1] = msg
	}
	outMsg := ""
	nlines = len(u.mMsgs)
	for i := 0; i < nlines; i++ {
		outMsg = outMsg + "\n" + u.mMsgs[i]
	}
	u.mOutput.SetText(outMsg)
	u.mScroll.ScrollToBottom()
}

func (u *Ui) ShowStatus(msg string) {
	nlines := len(u.sMsgs)
	if nlines < MAXLINES {
		u.sMsgs = append(u.sMsgs, msg)
	} else {
		for i := 0; i < nlines-1; i++ {
			u.sMsgs[i] = u.sMsgs[i+1]
		}
		u.sMsgs[nlines-1] = msg
	}
	outMsg := ""
	nlines = len(u.sMsgs)
	for i := 0; i < nlines; i++ {
		outMsg = outMsg + "\n" + u.sMsgs[i]
	}
	u.sOutput.SetText(outMsg)
	u.sScroll.ScrollToBottom()
}
