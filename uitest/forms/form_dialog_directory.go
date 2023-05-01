package forms

import (
	"github.com/ipoluianov/goforms/uiforms"
)

type FormDialogDirectory struct {
	uiforms.Form
}

func (c *FormDialogDirectory) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormDialogDirectory")
}
