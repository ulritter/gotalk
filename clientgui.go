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
	sty fyne.TextStyle
}

// data strcture to hold the ui elements
type Ui struct {
	mHeader *canvas.Text
	sHeader *canvas.Text
	mBox    *fyne.Container
	sBox    *fyne.Container
	mMsgs   []MessageLine
	sMsgs   []MessageLine
	input   *widget.Entry
	mScroll *container.Scroll
	sScroll *container.Scroll
	button  *widget.Button
	win     fyne.Window
	ui_ref  *Ui
	app     fyne.App
}

// create new ui with fyne elements
func (u *Ui) newUi(conn net.Conn, nl Newline) fyne.CanvasObject {

	actTheme = u.app.Settings().ThemeVariant()
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
			processInput(conn, u.input.Text, nl, u)
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
func (u *Ui) ShowMessage(msg string, test bool) {
	linecolor := MESSAGECOLOR
	linestyle := MESSAGESTYLE

	mWords := strings.Fields(msg)
	mb := container.NewHBox()

	// fill horizontal box
	for i := range mWords {
		w := mWords[i]
		if w[0] == '$' {
			c, m, coloronly := checkColor(w)
			if coloronly {
				linecolor = *c
			} else {
				t := canvas.NewText(m, *c)
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
	if test == false {
		u.win.Canvas().Focus(u.input)
	}
}

//display a status message in the (right hand) status area of the ui
func (u *Ui) ShowStatus(msg string, test bool) {
	b := canvas.NewText(msg, STATUSCOLOR)
	b.TextStyle = STATUSSTYLE
	refreshVBoxContent(msg, &u.sMsgs, u.sBox, container.NewHBox(b))
	u.sBox.Refresh()
	u.sScroll.Refresh()
	u.sScroll.ScrollToBottom()
	if test == false {
		u.win.Canvas().Focus(u.input)
	}
}

func refreshVBoxContent(msg string, m *[]MessageLine, b *fyne.Container, h *fyne.Container) {
	nlines := len((*m))
	t := MessageLine{
		txt: msg,
		obj: h,
	}
	if nlines < MAXLINES {
		(*m) = append((*m), t)
		(*b).Add(t.obj)
	} else {
		for i := 0; i < nlines-1; i++ {
			(*m)[i] = (*m)[i+1]
			(*b).Objects[i] = (*b).Objects[i+1]
		}
		(*m)[nlines-1] = t
		(*b).Objects[nlines-1] = t.obj
	}
}

//parse the given string for appearance of color code commands
//and strip color code from string
//return values are detected color (if any), else default color
//the new string and an indicator whether the string only consisted
//of only the color code or whether color code was anly precededing
//a message element
func checkColor(w string) (*color.RGBA, string, bool) {

	var c *color.RGBA
	var co bool
	c = &MESSAGECOLOR
	//unfortunately golang randomly interates on maps
	//thus we have to circle throug an outer loop so that we
	//can check on the key strings in descending length
	for i := 7; i > 1; i-- {
		for key, element := range cmap {
			if element.len == i {
				if (len(w) == element.len) && (w == key) {
					co = true
					c = &element.color
					w = ""
					return c, w, co
				} else {
					if (len(w) > element.len) && (w[:element.len] == key) {
						co = false
						c = &element.color
						w = w[element.len:]
						return c, w, co
					}
				}
			}
		}
	}
	return c, w, co
}
