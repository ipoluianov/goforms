package ui

import (
	"math"

	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/goforms/utils/uiproperties"
)

type ProgressBar struct {
	TextBlock
	minValue float64
	maxValue float64
	value    float64
	barColor *uiproperties.Property
}

func NewProgressBar(parent Widget) *ProgressBar {
	var c ProgressBar
	c.InitControl(parent, &c)
	c.TextVAlign = canvas.VAlignCenter
	c.TextHAlign = canvas.HAlignCenter
	c.minValue = 0
	c.maxValue = 100
	c.value = 42
	c.text = ""
	c.barColor = AddPropertyToWidget(&c, "barColor", uiproperties.PropertyTypeColor)
	InitDefaultStyle(&c)
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

func (c *ProgressBar) IsZero(value float64) bool {
	return math.Abs(value) <= math.SmallestNonzeroFloat64
}

func (c *ProgressBar) Draw(ctx DrawContext) {
	// Width of inner area of the control
	barMaxWidth := float64(c.InnerWidth())

	// Determine progress value in diapason [0:1]
	diff := (c.maxValue - c.minValue)
	effectiveValue := 0.0 // 0.0 - 1.0
	if !c.IsZero(diff) {
		effectiveValue = c.value / diff
	}

	// Calc width of the bar for display
	barWidth := barMaxWidth * effectiveValue
	if barWidth < 0 {
		barWidth = 0
	}
	if barWidth > barMaxWidth {
		barWidth = barMaxWidth
	}

	// Draw bar
	ctx.SetColor(c.barColor.Color())
	ctx.FillRect(0, 0, int(barWidth), c.ClientHeight())

	// Draw text
	c.TextBlock.Draw(ctx)
}

func (c *ProgressBar) SetMinValue(min float64) {
	c.minValue = min
	c.checkValueBounds()
	c.Window().UpdateLayout()
}

func (c *ProgressBar) SetMaxValue(max float64) {
	c.maxValue = max
	c.checkValueBounds()
	c.Window().UpdateLayout()
}

func (c *ProgressBar) SetMinMaxValue(min float64, max float64) {
	c.SetMinValue(min)
	c.SetMaxValue(max)
}

func (c *ProgressBar) SetValue(value float64) {
	c.value = value
	c.checkValueBounds()
	c.Window().UpdateLayout()
}

func (c *ProgressBar) SetValueAndText(value float64, text string) {
	c.value = value
	c.checkValueBounds()
	c.text = text
	c.Window().UpdateLayout()
}

func (c *ProgressBar) checkValueBounds() {
	if c.value < c.minValue {
		c.value = c.minValue
	}
	if c.value > c.maxValue {
		c.value = c.maxValue
	}
}
