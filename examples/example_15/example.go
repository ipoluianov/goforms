package example15

import (
	"fmt"

	"github.com/ipoluianov/goforms/ui"
)

type MainForm struct {
	ui.Form
}

func (c *MainForm) OnInit() {
	lv := c.Panel().AddListView()
	lv.AddColumn("Name", 100)
	lv.AddColumn("Size", 100)
	lv.AddColumn("Size", 100)
	lv.AddColumn("Size", 100)
	lv.AddColumn("Size", 100)
	lv.AddColumn("Size", 100)
	lv.AddColumn("Size", 100)
	lv.AddColumn("Size", 100)
	lv.AddColumn("Size", 100)
	lv.AddColumn("Size", 100)

	for i := 0; i < 100; i++ {
		lv.AddItem2("Item "+fmt.Sprint(i), "100")
	}

	lv.Focus()
}

func newMainForm() *MainForm {
	var c MainForm
	return &c
}

func Run() {
	ui.StartMainForm(newMainForm())
}
