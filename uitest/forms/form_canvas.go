package forms

import (
	"github.com/ipoluianov/goforms/canvas"
	"github.com/ipoluianov/goforms/uicontrols"
	"github.com/ipoluianov/goforms/uiforms"
	"github.com/ipoluianov/goforms/uiinterfaces"
	"golang.org/x/image/colornames"
)

type FormCanvas struct {
	uiforms.Form
}

type CanvasCtrl struct {
	uicontrols.Control
}

func NewCanvasCtrl(parent uiinterfaces.Widget) *CanvasCtrl {
	var c CanvasCtrl
	//c.InitControl(parent, &c, 0, 0, 0, 0, 0)
	return &c
}

func (c *CanvasCtrl) Draw(ctx *canvas.CanvasDirect) {
	ctx.DrawRect(10, 10, 100, 100, colornames.Red, 1)
}

func (c *FormCanvas) OnInit() {
	c.Resize(600, 800)
	c.SetTitle("FormFont")

	//ctrl := NewCanvasCtrl(c.Panel())
	//c.Panel().AddWidgetOnGrid(ctrl, 0, 0)
}
