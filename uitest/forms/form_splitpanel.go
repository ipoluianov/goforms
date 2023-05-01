package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormSplitPanel struct {
	uiforms.Form
}

func (c *FormSplitPanel) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormSplitPanel")
}
