package example02

import "github.com/ipoluianov/goforms/ui"

/*
	Minimal UI application
*/

func ExecMainForm() {
	ui.InitUI()
	ui.StartMainForm(ui.NewForm())
}
