package example07

import "github.com/ipoluianov/goforms/ui"

func newMainForm() *ui.Form {
	form := ui.NewForm()
	column := form.Panel().AddVPanel()
	groupBox := column.AddGroupBox("GroupBox")
	groupBox.Panel().AddTextBlock("Label1")
	groupBox.Panel().AddTextBlock("Label2")
	groupBox.Panel().AddTextBlock("Label3")
	inGroupPanel := groupBox.Panel().AddHPanel()
	inGroupPanel.AddButton("Change Name", func(event *ui.Event) { groupBox.SetTitle("Changed Name") })
	inGroupPanel.AddHSpacer()
	column.AddVSpacer()
	return form
}

func ExecMainForm() {
	ui.InitUI()
	ui.StartMainForm(newMainForm())
}
