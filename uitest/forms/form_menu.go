package forms

import (
	"github.com/ipoluianov/goforms/uiforms"
)

type FormMenu struct {
	uiforms.Form
}

func (c *FormMenu) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormMenu")
}
