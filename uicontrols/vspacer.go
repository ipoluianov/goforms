package uicontrols

import (
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"golang.org/x/image/colornames"
)

type VSpacer struct {
	Control
}

func NewVSpacer(parent uiinterfaces.Widget) *VSpacer {
	var c VSpacer
	c.InitControl(parent, &c)
	c.SetYExpandable(true)
	c.SetXExpandable(false)
	return &c
}

func (c *VSpacer) Subclass() string {
	return c.Control.Subclass()
}

func (c *VSpacer) ControlType() string {
	return "VSpacer"
}

func (c *VSpacer) Draw(ctx ui.DrawContext) {
	ctx.SetColor(colornames.Lightblue)
	ctx.SetStrokeWidth(1)
	//ctx.DrawLine(c.InnerWidth()/2, 0, c.InnerWidth()/2, c.InnerHeight())
}

func (c *VSpacer) TabStop() bool {
	return false
}
