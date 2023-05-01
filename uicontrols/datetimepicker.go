package uicontrols

import (
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"image"
	"time"
)

type DateTimePicker struct {
	Control

	dateTime time.Time

	pressed            bool
	mousePointerInRect bool

	enabled_ bool
	checked_ bool
	Image    image.Image

	DateTimeChanged func(event *uievents.Event)
}

func NewDateTimePicker(parent uiinterfaces.Widget) *DateTimePicker {
	var c DateTimePicker
	c.InitControl(parent, &c)
	c.enabled_ = true
	c.dateTime = time.Now()
	c.SetMinWidth(200)
	return &c
}

func (c *DateTimePicker) Subclass() string {
	if c.pressed && c.mousePointerInRect {
		return "pressed"
	}
	return c.Control.Subclass()
}

func (c *DateTimePicker) ControlType() string {
	return "DateTimePicker"
}

func (c *DateTimePicker) SetChecked(checked bool) {
	c.checked_ = checked
	c.Update("DateTimePicker")
}

func (c *DateTimePicker) SetEnabled(enabled bool) {
	c.enabled_ = enabled
	c.Update("DateTimePicker")
}

func (c *DateTimePicker) DateTimeToString() string {
	return c.dateTime.Format("2006-01-02 15:04:05.999")
}

func (c *DateTimePicker) Draw(ctx ui.DrawContext) {
	if c.checked_ {
		ctx.SetStrokeWidth(2)
		ctx.SetColor(c.rightBorderColor.Color())
		ctx.DrawRect(0, 0, c.InnerWidth(), c.InnerHeight())
	} else {
		ctx.SetStrokeWidth(1)
		ctx.SetColor(c.rightBorderColor.Color())
		ctx.DrawRect(0, 0, c.InnerWidth(), c.InnerHeight())
	}

	text := c.DateTimeToString()

	if c.enabled_ {
		ctx.SetColor(c.ForeColor())
		ctx.SetTextAlign(canvas.HAlignLeft, canvas.VAlignCenter)
		ctx.SetFontFamily(c.FontFamily())
		ctx.SetFontSize(c.FontSize())
		ctx.DrawText(3, 1, c.InnerWidth()-6, c.InnerHeight()-2, text)
	} else {
		ctx.SetColor(c.InactiveColor())
		ctx.SetTextAlign(canvas.HAlignLeft, canvas.VAlignCenter)
		ctx.SetFontFamily(c.FontFamily())
		ctx.SetFontSize(c.FontSize())
		ctx.DrawText(3, 1, c.InnerWidth()-6, c.InnerHeight()-2, text)
	}

	if c.Image != nil {
		ctx.DrawImage(1, 1, c.InnerWidth()-2, c.InnerHeight()-2, c.Image)
	}
}

func (c *DateTimePicker) DateTime() time.Time {
	return c.dateTime
}

func (c *DateTimePicker) SetDateTime(dateTime time.Time) {
	c.dateTime = dateTime
	c.Update("DateTimePicker")
}

func (c *DateTimePicker) TabStop() bool {
	return true
}

func (c *DateTimePicker) KeyDown(event *uievents.KeyDownEvent) bool {
	return false
}

func (c *DateTimePicker) MouseEnter() {
	c.mousePointerInRect = true
}

func (c *DateTimePicker) MouseLeave() {
	c.mousePointerInRect = false
}

func (c *DateTimePicker) MouseDown(event *uievents.MouseDownEvent) {
	c.pressed = true
	c.Update("DateTimePicker")
}

func (c *DateTimePicker) MouseUp(event *uievents.MouseUpEvent) {
	if c.pressed {
		c.pressed = false
		c.Update("DateTimePicker")
		if event.X >= 0 && event.Y >= 0 && event.X < c.InnerWidth() && event.Y < c.InnerHeight() {
			c.Press()
		}
	}
}

func (c *DateTimePicker) MouseClick(event *uievents.MouseClickEvent) {
}

func (c *DateTimePicker) Press() {
	c.selectDateTime()
}

func (c *DateTimePicker) dateTimeChangedInDialog(dateTime time.Time) {
	c.dateTime = dateTime
	if c.DateTimeChanged != nil {
		c.DateTimeChanged(uievents.NewEvent(c))
	}
	c.Update("DateTimePicker")
}

func (c *DateTimePicker) selectDateTime() {
	if !c.enabled_ {
		return
	}

	dialog := NewDateTimePickerDialog(c)
	dialog.SetDateTime(c.dateTime)
	dialog.OnAccept = func() {
		c.dateTimeChangedInDialog(dialog.DateTime())
	}
	x, y := c.RectClientAreaOnWindow()
	dialog.ShowDialogAtPos(x, y+c.Height())
}
