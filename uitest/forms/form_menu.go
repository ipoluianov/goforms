package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormMenu struct {
	uiforms.Form
}

func (c *FormMenu) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormMenu")
}
