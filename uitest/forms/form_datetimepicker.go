package forms

import (
	"github.com/ipoluianov/goforms/uiforms"
)

type FormDateTimePicker struct {
	uiforms.Form
}

func (c *FormDateTimePicker) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormDateTimePicker")
}
