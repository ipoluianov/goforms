package example13

import (
	"image/color"

	"github.com/ipoluianov/goforms/ui"
)

type CustomControl struct {
	ui.Control
}

func NewCustomControl(parent ui.Widget) *CustomControl {
	var c CustomControl
	c.InitControl(parent, &c)
	return &c
}

func (c *CustomControl) Draw(ctx ui.DrawContext) {
	ctx.SetColor(color.White)
	ctx.SetFontSize(16)
	ctx.DrawText(0, 0, c.Width(), c.Height(), "This is a Custom Control")
}
