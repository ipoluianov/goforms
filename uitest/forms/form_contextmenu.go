package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormContextMenu struct {
	uiforms.Form
}

func (c *FormContextMenu) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormContextMenu")
}
