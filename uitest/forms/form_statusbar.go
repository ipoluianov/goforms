package forms

import (
	"github.com/ipoluianov/goforms/uiforms"
)

type FormStatusBar struct {
	uiforms.Form
}

func (c *FormStatusBar) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormStatusBar")
}
