package example15

import "github.com/ipoluianov/goforms/ui"

type MainForm struct {
	ui.Form
}

func (c *MainForm) OnInit() {
	lv := c.Panel().AddListView()
	lv.AddColumn("Name", 100)
	lv.AddColumn("Size", 100)

	lv.AddItem2("Item 1", "100")
	lv.AddItem2("Item 2", "200")
	lv.AddItem2("Item 3", "300")
}

func newMainForm() *MainForm {
	var c MainForm
	return &c
}

func Run() {
	ui.StartMainForm(newMainForm())
}
