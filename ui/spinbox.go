package ui

import (
	"fmt"
	"strconv"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ipoluianov/goforms/uiresources"
)

type SpinBox struct {
	Container
	txt       *TextBox
	btnUp     *Button
	btnDown   *Button
	value     float64
	minValue  float64
	maxValue  float64
	precision int
	increment float64

	loadingTextBox bool

	OnValueChanged func(spinBox *SpinBox, value float64)
}

func NewSpinBox(parent Widget) *SpinBox {
	var c SpinBox

	c.InitControl(parent, &c)

	c.cellPadding = 0
	c.panelPadding = 0

	c.txt = NewTextBox(&c)
	c.AddWidgetOnGrid(c.txt, 0, 0)

	p := NewPanel(&c)
	p.cellPadding = 0
	p.panelPadding = 0
	c.AddWidgetOnGrid(p, 1, 0)

	c.btnUp = p.AddButtonOnGrid(0, 0, "", c.upClick)
	c.btnUp.SetImage(uiresources.ResImage(uiresources.R_icons_material4_png_navigation_expand_less_materialicons_48dp_1x_baseline_expand_less_black_48dp_png))
	c.btnUp.SetName("btnUp")
	c.btnUp.SetImageSize(16, 8)
	c.btnUp.SetPadding(0)

	c.btnDown = p.AddButtonOnGrid(0, 1, "", c.downClick)
	c.btnDown.SetImage(uiresources.ResImage(uiresources.R_icons_material4_png_navigation_expand_more_materialicons_48dp_1x_baseline_expand_more_black_48dp_png))
	c.btnDown.SetName("btnDown")
	c.btnDown.SetImageSize(16, 8)
	c.btnDown.SetPadding(0)

	c.txt.onValidateNeeded = c.onValidateNeeded
	c.txt.OnTextChanged = c.onTextChanged
	c.txt.onPreKeyDown = c.onPreTxtKeyDown

	c.precision = 3
	c.increment = 1
	c.minValue = -1000000000
	c.maxValue = 1000000000
	c.SetValue(c.value)

	return &c
}

func (c *SpinBox) Dispose() {
	c.txt = nil
	c.btnUp = nil
	c.btnDown = nil
	c.Container.Dispose()
}

func (c *SpinBox) SetPrecision(precision int) {
	c.precision = precision
	c.updateValue()
	c.Update("SpinBox")
}

func (c *SpinBox) SetIncrement(increment float64) {
	c.increment = increment
	c.updateValue()
	c.Update("SpinBox")
}

func (c *SpinBox) SetMinValue(minValue float64) {
	c.minValue = minValue
	c.updateValue()
	c.Update("SpinBox")
}

func (c *SpinBox) SetMaxValue(maxValue float64) {
	c.maxValue = maxValue
	c.updateValue()
	c.Update("SpinBox")
}

func (c *SpinBox) upClick(event *Event) {
	if c.validateValue(c.value+c.increment) == nil {
		c.value += c.increment
		c.SetValue(c.value)
	}
}

func (c *SpinBox) downClick(event *Event) {
	if c.validateValue(c.value-c.increment) == nil {
		c.value -= c.increment
		c.SetValue(c.value)
	}
}

func (c *SpinBox) validateValue(value float64) error {
	if value < c.minValue {
		return fmt.Errorf("< min value")
	}
	if value > c.maxValue {
		return fmt.Errorf("> max value")
	}
	return nil
}

func (c *SpinBox) SetValue(value float64) {
	changed := false
	if c.value != value {
		changed = true
	}

	c.value = value

	if c.value < c.minValue {
		value = c.minValue
	}

	if c.value > c.maxValue {
		value = c.maxValue
	}

	c.updateValue()

	if changed && c.OnValueChanged != nil {
		c.OnValueChanged(c, c.value)
	}
}

func (c *SpinBox) updateValue() {
	c.loadingTextBox = true
	c.txt.SetText(strconv.FormatFloat(c.value, 'f', c.precision, 64))
	c.loadingTextBox = false
}

func (c *SpinBox) Value() float64 {
	return c.value
}

func (c *SpinBox) onValidateNeeded(oldValue string, newValue string) bool {
	val, err := strconv.ParseFloat(newValue, 64)
	if err == nil {
		if c.validateValue(val) == nil {
			return true
		}
	}
	return false
}

func (c *SpinBox) onTextChanged(txtBox *TextBox, oldValue string, newValue string) {
	if c.loadingTextBox {
		return
	}
	val, err := strconv.ParseFloat(newValue, 64)
	if err == nil {
		c.SetValue(val)
	}
}

func (c *SpinBox) onPreTxtKeyDown(event *KeyDownEvent) {
	if event.Key == glfw.KeyUp {
		c.upClick(nil)
		event.Ignore = true
	}
	if event.Key == glfw.KeyDown {
		c.downClick(nil)
		event.Ignore = true
	}
	if event.Key == glfw.KeyEnter {
		c.txt.SelectAllText()
		event.Ignore = true
	}
}

func (c *SpinBox) MouseWheel(event *MouseWheelEvent) {
	//delta := event.Delta

	if event.Delta > 0 {
		c.SetValue(c.Value() + c.increment)
	} else {
		c.SetValue(c.Value() - c.increment)
	}

	c.Update("SpinBox")
}
