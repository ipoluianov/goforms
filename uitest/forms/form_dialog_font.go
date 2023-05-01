package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormDialogFont struct {
	uiforms.Form
}

func (c *FormDialogFont) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormDialogFont")
}
