package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormColorPicker struct {
	uiforms.Form
}

func (c *FormColorPicker) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormColorPicker")
}
