package example16

import "github.com/ipoluianov/goforms/ui"

type MainForm struct {
	ui.Form
}

func NewMainForm() *MainForm {
	var c MainForm
	return &c
}

func (c *MainForm) OnInit() {
	mainSplit := c.Panel().AddSplitContainerOnGrid(0, 0)
	panel1 := NewFilePanel(mainSplit.Panel1)
	panel2 := NewFilePanel(mainSplit.Panel2)
	mainSplit.Panel1.AddWidget(panel1)
	mainSplit.Panel2.AddWidget(panel2)
	mainSplit.SetPositionRelative(0.5)
}
