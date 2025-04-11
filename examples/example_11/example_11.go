package example11

import "github.com/ipoluianov/goforms/ui"

func newMainForm() *ui.Form {
	form := ui.NewForm()
	grid := form.Panel().AddPanelOnGrid(0, 0)
	grid.AddTextBlockOnGrid(0, 0, "111")
	grid.AddTextBlockOnGrid(0, 1, "222")
	grid.AddTextBlockOnGrid(1, 0, "333")
	grid.AddTextBlockOnGrid(1, 1, "444")
	return form
}

func ExecMainForm() {
	ui.StartMainForm(newMainForm())
}
