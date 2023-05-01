package forms

import (
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiforms"
)

type FormRadioButton struct {
	uiforms.Form

	rb11 *uicontrols.RadioButton
	rb12 *uicontrols.RadioButton
	rb13 *uicontrols.RadioButton

	rb21 *uicontrols.RadioButton
	rb22 *uicontrols.RadioButton
	rb23 *uicontrols.RadioButton

	txt *uicontrols.TextBlock
}

func (c *FormRadioButton) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormRadioButton")

	p1 := c.Panel().AddPanelOnGrid(0, 0)
	c.rb11 = p1.AddRadioButtonOnGrid(0, 0, "111", nil)
	c.rb12 = p1.AddRadioButtonOnGrid(0, 1, "222", nil)
	c.rb13 = p1.AddRadioButtonOnGrid(0, 2, "333", nil)
	p1.AddVSpacerOnGrid(0, 3)

	p2 := c.Panel().AddPanelOnGrid(1, 0)
	c.rb21 = p2.AddRadioButtonOnGrid(0, 0, "111", nil)
	c.rb22 = p2.AddRadioButtonOnGrid(1, 0, "222", nil)
	c.rb23 = p2.AddRadioButtonOnGrid(2, 0, "333", nil)
	p2.AddVSpacerOnGrid(0, 1)
	p2.AddHSpacerOnGrid(3, 0)

	p3 := c.Panel().AddPanelOnGrid(0, 1)
	p3.AddButtonOnGrid(0, 0, "bbb", nil)

	p4 := c.Panel().AddPanelOnGrid(1, 1)
	c.txt = p4.AddTextBlockOnGrid(0, 0, "---")
}
