package ui

import (
	"golang.org/x/image/colornames"
)

type HSpacer struct {
	Control
}

func NewHSpacer(parent Widget) *HSpacer {
	var c HSpacer
	c.InitControl(parent, &c)
	c.SetXExpandable(true)
	return &c
}

func (c *HSpacer) Subclass() string {
	return c.Control.Subclass()
}

func (c *HSpacer) ControlType() string {
	return "HSpacer"
}

func (c *HSpacer) Draw(ctx DrawContext) {
	ctx.SetColor(colornames.Lightblue)
	ctx.SetStrokeWidth(1)
	//ctx.DrawLine(0, c.ClientHeight()/2, c.ClientWidth(), c.ClientHeight()/2)
}

func (c *HSpacer) TabStop() bool {
	return false
}
