package uicontrols

import (
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type Space struct {
	Control

	onPress func(event *uievents.Event)
}

func NewSpace(parent uiinterfaces.Widget) *Space {
	var c Space
	c.InitControl(parent, &c)
	return &c
}

func (c *Space) Subclass() string {
	return c.Control.Subclass()
}

func (c *Space) ControlType() string {
	return "Space"
}

func (c *Space) Draw(ctx ui.DrawContext) {
	//ctx.DrawRect(0, 0, c.InnerWidth(), c.InnerHeight(), c.rightBorderColor.Color(), 1)
}

func (c *Space) TabStop() bool {
	return false
}
