package main

import (
	"github.com/ipoluianov/goforms/ui"
)

func main() {
	ui.InitUI()
	form := ui.NewForm()
	panelTop := form.Panel().AddPanelOnGrid(0, 0)
	panelTop.AddTextBlockOnGrid(0, 0, "Label")
	panelTop.AddVSpacerOnGrid(0, 1)
	panelBottom := form.Panel().AddPanelOnGrid(0, 1)
	panelBottom.AddHSpacerOnGrid(0, 0)
	btnOK := panelBottom.AddButtonOnGrid(1, 0, "OK", nil)
	btnOK.SetMinWidth(100)
	btnCancel := panelBottom.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(100)
	ui.StartMainForm(form)
}
