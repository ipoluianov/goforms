package forms

import (
	"github.com/ipoluianov/goforms/uicontrols"
	"github.com/ipoluianov/goforms/uiforms"
)

type FormAbsLayout struct {
	uiforms.Form

	btn1 *uicontrols.Button
	btn2 *uicontrols.Button
	btn3 *uicontrols.Button
	btn4 *uicontrols.Button
}

func (c *FormAbsLayout) OnInit() {
	c.Resize(400, 400)
	c.SetTitle("FormGridLayout")

	c.btn1 = uicontrols.NewButton(c.Panel(), "100x100", nil)
	c.Panel().AddWidget(c.btn1)

	c.btn2 = uicontrols.NewButton(c.Panel(), "200x200", nil)

	c.Panel().AddWidget(c.btn2)
}
