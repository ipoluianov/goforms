package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormSpinBox struct {
	uiforms.Form
}

func (c *FormSpinBox) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormSpinBox")

	c.Panel().AddSpinBoxOnGrid(0, 0)
	c.Panel().AddVSpacerOnGrid(0, 1)
	c.Panel().AddHSpacerOnGrid(1, 0)
}
