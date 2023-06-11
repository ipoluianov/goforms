package main

import (
	"github.com/ipoluianov/goforms/ui"
)

func newMainForm() *ui.Form {
	form := ui.NewForm()

	panelTop := form.Panel().AddVPanel()
	lv := panelTop.AddListView()
	lv.AddColumn("Col1", 100)
	lv.AddColumn("Col2", 100)
	lv.AddColumn("Col3", 100)
	lv.AddItem3("1", "222", "333")
	lv.AddItem3("2", "222", "333")
	lv.AddItem3("3", "222", "333")
	panelTop.AddTextBlock("Label")
	panelTop.AddTextBox()
	//panelTop.AddVSpacer()

	panelBottom := form.Panel().AddHPanel()
	panelBottom.AddHSpacer()
	panelBottom.AddButton("OK", func(event *ui.Event) { form.Close() })
	panelBottom.AddButton("Cancel", func(event *ui.Event) { form.Close() })

	return form
}

func main() {
	ui.InitUI()
	ui.StartMainForm(newMainForm())
}
