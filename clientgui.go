//go:build !serveronly

package main

import (
	"image/color"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const MAXLINES = 1024
const APPTITLE = "cooltide"
const WINTITLE = "gotalk"

const MESSAGEWIDTH = 618
const MESSAGEHEIGHT = 500
const STATUSWIDTH = 382
const STATUSHEIGHT = 540

var HEADERCOLOR color.RGBA
var HEADERSTYLE = fyne.TextStyle{Bold: true, Italic: false}

var MESSAGECOLOR color.RGBA
var MESSAGESTYLE = fyne.TextStyle{Bold: false, Italic: false}

var STATUSCOLOR color.RGBA
var STATUSSTYLE = fyne.TextStyle{Bold: false, Italic: true}

var actTheme fyne.ThemeVariant

type Colorfield struct {
	color color.RGBA
	len   int
}

//to add new colors, just add 'em here and they
//will be automatically both recognized and processed
var cmap_light = map[string]Colorfield{
	"$cyan":   {color.RGBA{20, 150, 220, 255}, 5},
	"$c":      {color.RGBA{20, 150, 220, 255}, 2},
	"$red":    {color.RGBA{210, 60, 18, 255}, 4},
	"$r":      {color.RGBA{210, 60, 18, 255}, 2},
	"$green":  {color.RGBA{50, 190, 40, 255}, 6},
	"$g":      {color.RGBA{50, 190, 40, 255}, 2},
	"$yellow": {color.RGBA{160, 150, 0, 255}, 7},
	"$y":      {color.RGBA{160, 150, 0, 255}, 2},
	"$purple": {color.RGBA{180, 90, 178, 255}, 7},
	"$p":      {color.RGBA{180, 90, 178, 255}, 2},
	"$white":  {color.RGBA{234, 234, 234, 255}, 6},
	"$w":      {color.RGBA{234, 234, 234, 255}, 2},
	"$black":  {color.RGBA{0, 0, 0, 255}, 6},
	"$b":      {color.RGBA{0, 0, 0, 255}, 2},
}

var cmap_dark = map[string]Colorfield{
	"$cyan":   {color.RGBA{30, 200, 234, 255}, 5},
	"$c":      {color.RGBA{30, 200, 234, 255}, 2},
	"$red":    {color.RGBA{255, 90, 25, 255}, 4},
	"$r":      {color.RGBA{255, 90, 25, 255}, 2},
	"$green":  {color.RGBA{90, 234, 81, 255}, 6},
	"$g":      {color.RGBA{90, 234, 81, 255}, 2},
	"$yellow": {color.RGBA{234, 195, 11, 255}, 7},
	"$y":      {color.RGBA{234, 195, 11, 255}, 2},
	"$purple": {color.RGBA{200, 100, 188, 255}, 7},
	"$p":      {color.RGBA{200, 100, 188, 255}, 2},
	"$white":  {color.RGBA{234, 234, 234, 255}, 6},
	"$w":      {color.RGBA{234, 234, 234, 255}, 2},
	"$black":  {color.RGBA{0, 0, 0, 255}, 6},
	"$b":      {color.RGBA{0, 0, 0, 255}, 2},
}

var cmap map[string]Colorfield

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

//set up custom theme to define backup colors for light / dark mode
func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return color.RGBA{200, 200, 200, 255}
		}
		return color.RGBA{70, 70, 70, 255}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// setup colors to be opaque in dark / light themes
func setColors(a fyne.App) {
	a.Settings().SetTheme(&myTheme{})

	switch a.Settings().ThemeVariant() {
	case theme.VariantDark:
		cmap = cmap_dark
		HEADERCOLOR = color.RGBA{255, 200, 100, 255}
		MESSAGECOLOR = color.RGBA{234, 234, 234, 255}
		STATUSCOLOR = color.RGBA{180, 180, 255, 255}
	case theme.VariantLight:
		cmap = cmap_light
		HEADERCOLOR = color.RGBA{200, 100, 0, 255}
		MESSAGECOLOR = color.RGBA{10, 10, 10, 255}
		STATUSCOLOR = color.RGBA{70, 70, 110, 255}
	default:
	}
}

// for future use
type MessageLine struct {
	txt string
	obj *fyne.Container
	//for future use
	sty fyne.TextStyle
}

// data strcture to hold the ui elements
type Ui struct {
	mHeader *canvas.Text
	sHeader *canvas.Text
	mBox    *fyne.Container
	sBox    *fyne.Container
	input   *widget.Entry
	mScroll *container.Scroll
	sScroll *container.Scroll
	button  *widget.Button
	win     fyne.Window
	ui_ref  *Ui
	app     fyne.App
	//for future use
	mMsgs []MessageLine
	sMsgs []MessageLine
}

//create new ui structure with fyne elements and
//set the callbacks
func (u *Ui) newUi(conn net.Conn, nl Newline) fyne.CanvasObject {

	actTheme = u.app.Settings().ThemeVariant()
	u.ui_ref = u

	u.mBox = container.NewVBox()
	u.mScroll = container.NewScroll(u.mBox)
	u.mScroll.SetMinSize(fyne.NewSize(MESSAGEWIDTH, MESSAGEHEIGHT))
	u.mHeader = canvas.NewText(lang.Lookup(actualLocale, " Messages"), HEADERCOLOR)
	u.mHeader.TextStyle = HEADERSTYLE

	u.sBox = container.NewVBox()
	u.sScroll = container.NewScroll(u.sBox)
	u.sScroll.SetMinSize(fyne.NewSize(STATUSWIDTH, STATUSHEIGHT))
	u.sHeader = canvas.NewText(lang.Lookup(actualLocale, " Status Info"), HEADERCOLOR)
	u.sHeader.TextStyle = HEADERSTYLE

	u.input = widget.NewEntry()

	vSeparator := canvas.NewRectangle(color.Gray{})
	vSeparator.Resize(fyne.NewSize(3, 400))

	u.button = widget.NewButton(">>", func() {
		if len(u.input.Text) > 0 {
			processInput(conn, u.input.Text, u)
			u.input.SetText("")
		}
		if actTheme != u.app.Settings().ThemeVariant() {
			actTheme = u.app.Settings().ThemeVariant()
			setColors(u.app)
		}
		u.mScroll.Refresh()
		u.mScroll.ScrollToBottom()
		u.win.Canvas().Focus(u.input)
	})

	u.input.OnSubmitted = func(text string) {
		if len(u.input.Text) > 0 {
			processInput(conn, u.input.Text, u)
			u.input.SetText("")
		}
		if actTheme != u.app.Settings().ThemeVariant() {
			actTheme = u.app.Settings().ThemeVariant()
			setColors(u.app)
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
//check for inline color commands and populate the horizontal box
//according to requested color values
func (u *Ui) ShowMessage(msg string) {
	linecolor := MESSAGECOLOR
	linestyle := MESSAGESTYLE

	mWords := strings.Fields(msg)
	mb := container.NewHBox()

	// fill horizontal box making up the message line
	for i := range mWords {
		w := mWords[i]
		// if [nickname:] needs color change
		if (i == 0) && (w[1] == '$') {
			returnColor, inputWord, _ := checkColor(w[1 : len(w)-2])
			t := canvas.NewText("["+inputWord+"]:", *returnColor)
			t.TextStyle = linestyle

			mb.Add(t)
		} else if w[0] == '$' {
			returnColor, inputWord, coloronly := checkColor(w)
			if coloronly {
				linecolor = *returnColor
			} else {
				t := canvas.NewText(inputWord, *returnColor)
				t.TextStyle = linestyle
				mb.Add(t)
			}
		} else {
			t := canvas.NewText(w, linecolor)
			t.TextStyle = linestyle
			mb.Add(t)
		}
	}

	refreshVBoxContent(msg, &u.mMsgs, u.mBox, mb)
	u.mBox.Refresh()
	u.mScroll.Refresh()
	u.mScroll.ScrollToBottom()
	u.win.Canvas().Focus(u.input)
}

//display a status message in the (right hand) status area of the ui
func (u *Ui) ShowStatus(msg string) {
	b := canvas.NewText(msg, STATUSCOLOR)
	b.TextStyle = STATUSSTYLE
	refreshVBoxContent(msg, &u.sMsgs, u.sBox, container.NewHBox(b))
	u.sBox.Refresh()
	u.sScroll.Refresh()
	u.sScroll.ScrollToBottom()
	u.win.Canvas().Focus(u.input)
}

//append new message / container to existing message slice / container
//if both message slice and container exceed MAXLINES, remove the oldest message / container
//and append the new one at the bottom
func refreshVBoxContent(msg string, messageLine *[]MessageLine, targetContainer *fyne.Container, newContainer *fyne.Container) {
	nlines := len((*messageLine))
	bufferLine := MessageLine{
		txt: msg,
		obj: newContainer,
	}
	if nlines < MAXLINES {
		(*messageLine) = append((*messageLine), bufferLine)
		(*targetContainer).Add(bufferLine.obj)
	} else {
		for i := 0; i < nlines-1; i++ {
			(*messageLine)[i] = (*messageLine)[i+1]
			(*targetContainer).Objects[i] = (*targetContainer).Objects[i+1]
		}
		(*messageLine)[nlines-1] = bufferLine
		(*targetContainer).Objects[nlines-1] = bufferLine.obj
	}
}

//parse the given string for appearance of color code commands
//and strip color code from string
//return values are detected color (if any), else default color
//the new string and an indicator whether the string only consisted
//of only the color code or whether color code was anly precededing
//a message element
func checkColor(returnString string) (*color.RGBA, string, bool) {

	var returnColor *color.RGBA
	var coloronly bool
	returnColor = &MESSAGECOLOR
	//since golang randomly iterates on maps,
	//we have to circle through an outer loop so that we
	//can check on the key strings in descending length.
	//While this might look clumsy, we can now add any new color
	//to the map definition(s) and it will be automatically
	//recognized
	for i := 7; i > 1; i-- {
		for colorKey, colorcode := range cmap {
			if colorcode.len == i {
				if (len(returnString) == colorcode.len) && (returnString == colorKey) {
					coloronly = true
					returnColor = &colorcode.color
					returnString = ""
					return returnColor, returnString, coloronly
				} else {
					if (len(returnString) > colorcode.len) && (returnString[:colorcode.len] == colorKey) {
						coloronly = false
						returnColor = &colorcode.color
						returnString = returnString[colorcode.len:]
						return returnColor, returnString, coloronly
					}
				}
			}
		}
	}
	return returnColor, returnString, coloronly
}
