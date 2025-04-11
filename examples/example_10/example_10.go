package example10

import "github.com/ipoluianov/goforms/ui"

func newMainForm() *ui.Form {
	form := ui.NewForm()
	column := form.Panel().AddVPanel()
	column.AddButton("Show MessageBox", func(event *ui.Event) {
		ui.ShowInformationMessage(column, "This is the message", "Title")
	})
	column.AddVSpacer()
	return form
}

func ExecMainForm() {
	ui.StartMainForm(newMainForm())
}
