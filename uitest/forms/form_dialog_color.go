package forms

import (
	"github.com/ipoluianov/goforms/uiforms"
)

type FormDialogColor struct {
	uiforms.Form
}

func (c *FormDialogColor) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormDialogColor")
}
