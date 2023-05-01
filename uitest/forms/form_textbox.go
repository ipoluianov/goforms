package forms

import (
	"github.com/ipoluianov/goforms/uicontrols"
	"github.com/ipoluianov/goforms/uiforms"
)

type FormTextBox struct {
	uiforms.Form

	txt1 *uicontrols.TextBox
	txt2 *uicontrols.TextBox
}

func (c *FormTextBox) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormTextBox")

	//c.txt1 = c.Panel().AddTextBoxOnGrid(0, 0)
	c.txt2 = c.Panel().AddTextBoxOnGrid(0, 1)
	c.txt2.SetText("1 123123123\r\n2 123123123\r\n3 123123123\r\n4 123123123\r\n5 123123123\r\n6 123123123\r\n7 123123123\r\n8 123123123\r\n9 123123123\r\n10 123123123\r\n11 123123123\r\n12 123123123\r\n13 123123123\r\n14 123123123\r\n15 123123123\r\n")
	c.txt2.SetMultiline(true)
	c.Panel().AddVSpacerOnGrid(0, 2)
}
