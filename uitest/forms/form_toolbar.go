package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormToolBar struct {
	uiforms.Form
}

func (c *FormToolBar) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormToolBar")
}
