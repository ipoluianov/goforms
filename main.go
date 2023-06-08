package main

import (
	"github.com/ipoluianov/goforms/ui"
)

func main() {
	ui.InitUI()
	form := ui.NewForm()
	panelTop := form.Panel().AddVPanel()
	panelTop.AddTextBlock("Label")
	panelBottom := form.Panel().AddHPanel()
	btnOK := panelBottom.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(100)
	btnCancel := panelBottom.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(100)
	ui.StartMainForm(form)
}
