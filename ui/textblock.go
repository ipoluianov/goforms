package ui

import (
	"strings"

	"github.com/ipoluianov/goforms/canvas"
)

type TextBlock struct {
	Control
	pressed    bool
	OnClick    func(ev *Event)
	text       string
	TextVAlign canvas.VAlign
	TextHAlign canvas.HAlign
	underline  bool
}

func (c *TextBlock) ControlType() string {
	return "TextBlock"
}

func NewTextBlock(parent Widget, text string) *TextBlock {
	var c TextBlock
	c.text = text
	c.InitControl(parent, &c)
	c.TextVAlign = canvas.VAlignCenter
	c.TextHAlign = canvas.HAlignLeft
	return &c
}

func (c *TextBlock) Draw(ctx DrawContext) {
	ctx.SetColor(c.foregroundColor.Color())
	ctx.SetFontSize(c.fontSize.Float64())
	ctx.SetTextAlign(c.TextHAlign, c.TextVAlign)
	ctx.SetUnderline(c.underline)
	ctx.DrawText(0, 0, c.ClientWidth(), c.ClientHeight(), c.text)
}

func (c *TextBlock) SetUnderline(underline bool) {
	c.underline = underline
}

func (c *TextBlock) SetText(text string) {
	if c.text != text {
		c.text = text
		if c.Window() != nil {
			c.Window().UpdateLayout()
		}
	}
	c.Update("TextBlock")
}

func (c *TextBlock) MinWidth() int {
	if c.minWidth > 0 {
		return c.minWidth
	}

	result := 0
	for _, line := range c.lines() {
		lineW, _, _ := canvas.MeasureText(c.FontFamily(), c.FontSize(), c.FontBold(), c.FontItalic(), line, true)
		if lineW > result {
			result = lineW
		}
	}

	result += c.LeftBorderWidth()
	result += c.RightBorderWidth()

	return result
}

func (c *TextBlock) MinHeight() int {
	if c.minHeight > 0 {
		return c.minHeight
	}

	result := 0
	if len(c.lines()) > 0 {
		for _, line := range c.lines() {
			_, lineH, _ := canvas.MeasureText(c.FontFamily(), c.FontSize(), c.FontBold(), c.FontItalic(), line, true)
			result += lineH
		}
	} else {
		_, lineH, _ := canvas.MeasureText(c.FontFamily(), c.FontSize(), c.FontBold(), c.FontItalic(), "Qg", true)
		result += lineH
	}
	result += c.TopBorderWidth()
	result += c.BottomBorderWidth()
	return result
}

func (c *TextBlock) MouseDown(event *MouseDownEvent) {
	c.pressed = true
	c.Update("TextBlock")
}

func (c *TextBlock) MouseUp(event *MouseUpEvent) {
	if c.pressed {
		c.pressed = false
		c.Update("TextBlock")
		if event.X >= 0 && event.Y >= 0 && event.X < c.InnerWidth() && event.Y < c.InnerHeight() {
			if c.OnClick != nil {
				c.OnClick(NewEvent(c))
			}
		}
	}
}

func (c *TextBlock) Text() string {
	return c.text
}

func (c *TextBlock) lines() []string {
	if c == nil {
		panic("No text")
	}
	result := strings.Split(c.text, "\r\n")
	if len(result) == 1 && result[0] == "" {
		return make([]string, 0)
	}
	return result
}
