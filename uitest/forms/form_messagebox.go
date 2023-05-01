package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormMessageBox struct {
	uiforms.Form
}

func (c *FormMessageBox) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormMessageBox")
}
