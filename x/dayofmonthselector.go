package x

/*
import (
	"image"
	"time"
)

type DayOfMonthSelector struct {
	Container

	year  int
	month int

	selectedDay int

	pressed            bool
	mousePointerInRect bool

	enabled_ bool
	//checked_ bool
	Image image.Image

	buttons []*Button

	DayChanged func()
}

func NewDayOfMonthSelector(parent Widget) *DayOfMonthSelector {
	var c DayOfMonthSelector

	c.InitControl(parent, &c)
	c.enabled_ = true
	c.month = 1
	c.buttons = make([]*Button, 0)

	c.fillDays()

	return &c
}

func (c *DayOfMonthSelector) XExpandable() bool {
	return true
}

func (c *DayOfMonthSelector) YExpandable() bool {
	return true
}

func (c *DayOfMonthSelector) Subclass() string {
	if c.pressed && c.mousePointerInRect {
		return "pressed"
	}
	return c.Control.Subclass()
}

func (c *DayOfMonthSelector) SetYearAndMonth(year, month int) {
	c.year = year
	c.month = month
	c.fillDays()
	c.Update("DayOfMonthSelector")
}

func (c *DayOfMonthSelector) SelectedDay() int {
	return c.selectedDay
}

func (c *DayOfMonthSelector) ControlType() string {
	return "DateTimePicker"
}

func (c *DayOfMonthSelector) SetEnabled(enabled bool) {
	c.enabled_ = enabled
	c.Update("DateTimePicker")
}

func (c *DayOfMonthSelector) fillDays() {
	c.buttons = make([]*Button, 0)
	c.RemoveAllWidgets()

	loc, _ := time.LoadLocation("")
	t := time.Date(c.year, time.Month(c.month), 1, 1, 1, 1, 1, loc)

	xIndex := int(t.Weekday())
	yIndex := 0

	lastDay := 0
	for t.Month() == time.Month(c.month) {
		dayText := t.Format("02")
		btn := NewButton(c, dayText, c.onClickButton)
		btn.SetUserData("day", int(t.Day()))
		lastDay = t.Day()
		c.AddWidgetOnGrid(btn, xIndex, yIndex)
		c.buttons = append(c.buttons, btn)

		t = t.Add(24 * time.Hour)

		xIndex++
		if xIndex > 6 {
			xIndex = 0
			yIndex++
		}
	}

	if c.selectedDay > lastDay {
		c.SetDay(lastDay)
	}

	c.updateButtonColors()
}

func (c *DayOfMonthSelector) updateButtonColors() {
	for _, btn := range c.buttons {
		if btn.UserData("day").(int) == c.selectedDay {
			btn.SetForeColor(DefaultBackColor)
			btn.SetBackColor(c.ForeColor())
		} else {
			btn.SetForeColor(c.ForeColor())
			btn.SetBackColor(DefaultBackColor)
		}
	}
}

func (c *DayOfMonthSelector) Day() int {
	return c.selectedDay
}

func (c *DayOfMonthSelector) SetDay(day int) {
	c.selectedDay = day
	c.updateButtonColors()
	c.Update("DateTimePicker")

	if c.DayChanged != nil {
		c.DayChanged()
	}
}

func (c *DayOfMonthSelector) onClickButton(ev *Event) {
	c.SetDay(ev.Sender.(*Button).UserData("day").(int))
}
*/
