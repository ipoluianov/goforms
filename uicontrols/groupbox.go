package uicontrols

import (
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type GroupBox struct {
	Container
	titleTextBlock *TextBlock
	panel          *Panel
}

func NewGroupBox(parent uiinterfaces.Widget, title string) *GroupBox {
	var c GroupBox
	c.InitControl(parent, &c)

	c.titleTextBlock = NewTextBlock(&c, title)
	c.titleTextBlock.TextHAlign = canvas.HAlignLeft
	c.AddWidgetOnGrid(c.titleTextBlock, 0, 0)

	c.panel = NewPanel(&c)
	c.panel.SetBorders(1, c.ForeColor())
	c.AddWidgetOnGrid(c.panel, 0, 1)

	c.SetCellPadding(0)
	c.SetPanelPadding(0)
	return &c
}

func (c *GroupBox) ControlType() string {
	return "GroupBox"
}

func (c *GroupBox) Dispose() {
	c.titleTextBlock = nil
	c.panel = nil
	c.Container.Dispose()
}

func (c *GroupBox) Panel() *Panel {
	return c.panel
}

func (c *GroupBox) SetTitle(title string) {
	c.titleTextBlock.SetText(title)
}

func (c *GroupBox) Title() string {
	return c.titleTextBlock.Text()
}
