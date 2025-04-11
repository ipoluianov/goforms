package example15

import "github.com/ipoluianov/goforms/ui"

type MainForm struct {
	ui.Form
}

func newMainForm() *MainForm {
	var c MainForm
	return &c
}

func Run() {
	ui.StartMainForm(newMainForm())
}
