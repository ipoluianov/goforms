package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormTextBlock struct {
	uiforms.Form
}

func (c *FormTextBlock) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormTextBlock")
}
