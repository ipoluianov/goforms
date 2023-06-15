package example09

import "github.com/ipoluianov/goforms/ui"

func newMainForm() *ui.Form {
	form := ui.NewForm()
	column := form.Panel().AddVPanel()
	txtBlock := column.AddTextBox()
	menu := ui.NewPopupMenu(txtBlock)
	menu.AddItem("Item1", func(event *ui.Event) {}, nil, "")
	menu.AddItem("Item2", func(event *ui.Event) {}, nil, "")
	menu.AddItem("Item3", func(event *ui.Event) {}, nil, "")
	txtBlock.SetContextMenu(menu)
	column.AddVSpacer()
	return form
}

func ExecMainForm() {
	ui.InitUI()
	ui.StartMainForm(newMainForm())
}
