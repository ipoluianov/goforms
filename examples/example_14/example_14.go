package example14

import (
	"github.com/ipoluianov/goforms/ui"
)

type MainForm struct {
	ui.Form

	// UI
	/*
		Form:
			VPanel:
				HPanel:
					LineEdit
					Button
				HPanel:
					TextBlock
				VSpacer
	*/

	// First line
	p1       *ui.Panel
	lineEdit *ui.LineEdit
	btn      *ui.Button

	// Second line
	p2    *ui.Panel
	label *ui.TextBlock
}

func newMainForm() *MainForm {
	var c MainForm
	return &c
}

func (c *MainForm) OnInit() {
	vpanel := c.Panel().AddVPanel()

	// P1
	c.p1 = vpanel.AddHPanel()
	c.lineEdit = ui.NewLineEdit(c.p1)
	c.p1.AddWidget(c.lineEdit)
	btn := ui.NewButton(c.p1, "Create", func(event *ui.Event) {
		c.label.SetText("[" + c.lineEdit.Text() + "]")
	})
	c.p1.AddWidget(btn)

	// P2
	c.p2 = vpanel.AddHPanel()
	c.label = c.p2.AddTextBlock("---")

	vpanel.AddVSpacer()

	c.lineEdit.Focus()
}

func ExecMainForm() {
	ui.StartMainForm(newMainForm())
}
