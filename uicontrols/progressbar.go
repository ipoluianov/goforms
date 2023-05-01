package uicontrols

import (
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uiproperties"
	"github.com/gazercloud/gazerui/uistyles"
)

type ProgressBar struct {
	TextBlock
	minValue int
	maxValue int
	value    int
	barColor *uiproperties.Property
}

func NewProgressBar(parent uiinterfaces.Widget) *ProgressBar {
	var c ProgressBar
	c.InitControl(parent, &c)
	c.TextVAlign = canvas.VAlignCenter
	c.TextHAlign = canvas.HAlignCenter
	c.minValue = 0
	c.maxValue = 100
	c.value = 42
	c.text = " 123 "
	c.barColor = AddPropertyToWidget(&c, "barColor", uiproperties.PropertyTypeColor)
	uistyles.InitDefaultStyle(&c)
	return &c
}

func (c *ProgressBar) OnInit() {
}

func (c *ProgressBar) ControlType() string {
	return "ProgressBar"
}

func (c *ProgressBar) Subclass() string {
	return "default"
}

func (c *ProgressBar) Draw(ctx ui.DrawContext) {

	maxWidth := float64(c.InnerWidth() - 6)
	perc := float64(c.value) / float64(c.maxValue-c.minValue)
	width := maxWidth * perc
	if width > maxWidth {
		width = maxWidth
	}

	ctx.SetColor(c.barColor.Color())
	ctx.FillRect(0, 0, int(width), c.ClientHeight())
	c.TextBlock.Draw(ctx)
}

func (c *ProgressBar) SetMin(min int) {
	c.minValue = min
	c.Window().UpdateLayout()
}

func (c *ProgressBar) SetMax(max int) {
	c.maxValue = max
	c.Window().UpdateLayout()
}

func (c *ProgressBar) SetMinMax(min int, max int) {
	c.SetMin(min)
	c.SetMax(max)
}

func (c *ProgressBar) SetValue(value int) {
	c.value = value
	c.Window().UpdateLayout()
}
