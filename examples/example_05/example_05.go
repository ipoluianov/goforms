package example05

import "github.com/ipoluianov/goforms/ui"

func newMainForm() *ui.Form {
	form := ui.NewForm()
	column := form.Panel().AddVPanel()
	chkBox := column.AddCheckBox("CheckBox")
	lbl := column.AddTextBlock("")
	chkBox.OnCheckedChanged = func(checkBox *ui.CheckBox, checked bool) {
		if checked {
			lbl.SetText("Checked")
		} else {
			lbl.SetText("Unchecked")
		}
	}
	chkBox.SetChecked(true)

	row := column.AddHPanel()
	row.AddButton("Check", func(event *ui.Event) {
		chkBox.SetChecked(true)
	})
	row.AddButton("Uncheck", func(event *ui.Event) {
		chkBox.SetChecked(false)
	})
	row.AddHSpacer()
	column.AddVSpacer()
	return form
}

func ExecMainForm() {
	ui.StartMainForm(newMainForm())
}
