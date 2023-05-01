package forms

import (
	"github.com/ipoluianov/goforms/uiforms"
)

type FormFont struct {
	uiforms.Form
}

func (c *FormFont) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormFont")
}
