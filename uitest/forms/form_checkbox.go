package forms

import (
	"github.com/ipoluianov/goforms/uicontrols"
	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiforms"
)

type FormCheckBox struct {
	uiforms.Form

	num    *uicontrols.SpinBox
	chkBox *uicontrols.CheckBox
}

func (c *FormCheckBox) onChecked() {
	c.num.SetValue(c.num.Value() + 1)
}

func (c *FormCheckBox) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormCheckBox")

	p1 := c.Panel().AddPanelOnGrid(0, 0)
	c.chkBox = p1.AddCheckBoxOnGrid(0, 0, "CheckBoxText")
	p1.AddVSpacerOnGrid(0, 1)

	p2 := c.Panel().AddPanelOnGrid(1, 0)
	c.num = p2.AddSpinBoxOnGrid(0, 0)
	p2.AddButtonOnGrid(0, 1, "Add listener", func(event *uievents.Event) {
	})
	p2.AddButtonOnGrid(0, 2, "Remove listener", func(event *uievents.Event) {
	})
	p2.AddVSpacerOnGrid(0, 10)
}
