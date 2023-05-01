package ui

import (
	"image/color"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/golang-collections/collections/stack"
	"github.com/ipoluianov/goforms/canvas"
	"golang.org/x/image/colornames"
)

type DrawContextBase struct {
	// Current settings
	CurrentColor color.Color
	StrokeWidth  int
	FontFamily   string
	FontSize     float64
	UnderLine    bool
	TextHAlign   canvas.HAlign
	TextVAlign   canvas.VAlign

	CurrentClipSettings ClipSettings
	StackClipSettings   stack.Stack
	WindowWidth         int
	WindowHeight        int

	Window *glfw.Window
}

type StateStruct struct {
	CurrentColor color.Color
	StrokeWidth  int
	FontFamily   string
	FontSize     float64
	UnderLine    bool
	TextHAlign   canvas.HAlign
	TextVAlign   canvas.VAlign

	CurrentClipSettings ClipSettings
}

func (c *DrawContextBase) InitBase() {
	c.TextHAlign = canvas.HAlignLeft
	c.TextVAlign = canvas.VAlignCenter
	c.FontFamily = "Roboto"
	c.FontSize = 12
	c.StrokeWidth = 1
	c.CurrentColor = colornames.Black
	for c.StackClipSettings.Len() > 0 {
		c.StackClipSettings.Pop()
	}

	c.CurrentClipSettings.x = 0
	c.CurrentClipSettings.y = 0
	c.CurrentClipSettings.width = c.WindowWidth
	c.CurrentClipSettings.height = c.WindowHeight
}

func (c *DrawContextBase) Save() {
	var ss StateStruct
	ss.CurrentColor = c.CurrentColor
	ss.StrokeWidth = c.StrokeWidth
	ss.FontFamily = c.FontFamily
	ss.FontSize = c.FontSize
	ss.UnderLine = c.UnderLine
	ss.TextHAlign = c.TextHAlign
	ss.TextVAlign = c.TextVAlign
	ss.CurrentClipSettings = c.CurrentClipSettings

	c.StackClipSettings.Push(ss)
}

func (c *DrawContextBase) TranslateAndClip(x, y, width, height int) {
	c.CurrentClipSettings.x += x
	c.CurrentClipSettings.y += y
	c.CurrentClipSettings.width = width
	c.CurrentClipSettings.height = height
	//c.setViewport()
}

func (c *DrawContextBase) Load() {
	ss := c.StackClipSettings.Peek().(StateStruct)
	c.CurrentColor = ss.CurrentColor
	c.StrokeWidth = ss.StrokeWidth
	c.FontFamily = ss.FontFamily
	c.FontSize = ss.FontSize
	c.UnderLine = ss.UnderLine
	c.TextHAlign = ss.TextHAlign
	c.TextVAlign = ss.TextVAlign
	c.CurrentClipSettings = ss.CurrentClipSettings
	c.StackClipSettings.Pop()
	//c.setViewport()
}

func (c *DrawContextBase) Finish() {
}

func (c *DrawContextBase) SetColor(col color.Color) {
	c.CurrentColor = col
}

func (c *DrawContextBase) SetStrokeWidth(w int) {
	c.StrokeWidth = w
}

func (c *DrawContextBase) SetFontFamily(fontFamily string) {
	c.FontFamily = fontFamily
}

func (c *DrawContextBase) SetFontSize(s float64) {
	c.FontSize = s
}

func (c *DrawContextBase) SetTextAlign(h canvas.HAlign, v canvas.VAlign) {
	c.TextHAlign = h
	c.TextVAlign = v
}

func (c *DrawContextBase) SetUnderline(underline bool) {
	c.UnderLine = underline
}
