package uicontrols

import (
	"fmt"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uistyles"
	"image"
)

type MonthSelector struct {
	Container

	selectedMonth int

	pressed            bool
	mousePointerInRect bool

	enabled_ bool
	checked_ bool
	Image    image.Image

	buttons      []*Button
	buttonWidth  int
	buttonHeight int

	MonthChanged func()
}

func NewMonthSelector(parent uiinterfaces.Widget) *MonthSelector {
	var c MonthSelector
	c.InitControl(parent, &c)
	c.enabled_ = true
	c.selectedMonth = 1
	c.buttons = make([]*Button, 0)
	c.xExpandable = true

	c.fillMonths()
	c.SetMonth(c.Month())

	return &c
}

func (c *MonthSelector) Subclass() string {
	if c.pressed && c.mousePointerInRect {
		return "pressed"
	}
	return c.Control.Subclass()
}

func (c *MonthSelector) XExpandable() bool {
	return true
}

func (c *MonthSelector) YExpandable() bool {
	return false
}

func (c *MonthSelector) ControlType() string {
	return "MonthSelector"
}

func (c *MonthSelector) SetEnabled(enabled bool) {
	c.enabled_ = enabled
	c.Update("MonthSelector")
}

func (c *MonthSelector) fillMonths() {
	c.buttons = make([]*Button, 0)
	c.RemoveAllWidgets()

	for i := 1; i <= 12; i++ {
		monthText := fmt.Sprint(i)
		btn := NewButton(c, monthText, c.onClickButton)
		btn.SetMinWidth(30)
		btn.SetMinHeight(30)
		btn.SetUserData("month", i)
		c.AddWidgetOnGrid(btn, i, 0)
		c.buttons = append(c.buttons, btn)

	}
	c.updateButtonsColors()
}

func (c *MonthSelector) Month() int {
	return c.selectedMonth
}

func (c *MonthSelector) SetMonth(month int) {
	c.selectedMonth = month
	c.updateButtonsColors()
	if c.MonthChanged != nil {
		c.MonthChanged()
	}
	c.Update("MonthSelector")
}

func (c *MonthSelector) updateButtonsColors() {
	for _, btn := range c.buttons {
		if btn.UserData("month").(int) == c.selectedMonth {
			btn.SetForeColor(uistyles.DefaultBackColor)
			btn.SetBackColor(c.ForeColor())
		} else {
			btn.SetForeColor(c.ForeColor())
			btn.SetBackColor(uistyles.DefaultBackColor)
		}
	}
}

func (c *MonthSelector) onClickButton(ev *uievents.Event) {
	c.SetMonth(ev.Sender.(*Button).UserData("month").(int))
}

func (c *MonthSelector) UpdateStyle() {
	c.Container.UpdateStyle()
	c.updateButtonsColors()
}
