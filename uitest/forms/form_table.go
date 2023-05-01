package forms

import (
	"github.com/ipoluianov/goforms/uiforms"
)

type FormTable struct {
	uiforms.Form
}

func (c *FormTable) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormTable")
}
