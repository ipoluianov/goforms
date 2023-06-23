package example14

import "github.com/ipoluianov/goforms/ui"

func newMainForm() *ui.Form {
	form := ui.NewForm()
	column := form.Panel().AddVPanel()
	column.AddWidget(NewCustomListView(column))
	return form
}

func ExecMainForm() {
	ui.InitUI()
	ui.StartMainForm(newMainForm())
}
