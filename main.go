package main

import (
	"github.com/ipoluianov/goforms/ui"
)

func makeMainForm() *ui.Form {
	form := ui.NewForm()

	panelTop := form.Panel().AddVPanel()
	panelTop.AddTextBlock("Label")
	txtBox := panelTop.AddTextBox()
	txtBox.SetText("Hello!")
	panelTop.AddVSpacer()

	panelBottom := form.Panel().AddHPanel()
	panelBottom.AddHSpacer()
	panelBottom.AddButton("OK", func(event *ui.Event) { form.Close() })
	panelBottom.AddButton("Cancel", func(event *ui.Event) { form.Close() })

	return form
}

func main() {
	ui.InitUI()
	ui.StartMainForm(makeMainForm())
}
