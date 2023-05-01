package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormTable struct {
	uiforms.Form
}

func (c *FormTable) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormTable")
}
