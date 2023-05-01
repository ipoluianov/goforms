package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormDialogDirectory struct {
	uiforms.Form
}

func (c *FormDialogDirectory) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormDialogDirectory")
}
