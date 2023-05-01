package main

import (
	"github.com/ipoluianov/goforms/ui"
	"github.com/ipoluianov/goforms/uiforms"
	"github.com/ipoluianov/goforms/uitest/forms"
)

func main() {
	ui.InitUI()
	{
		var form forms.MainForm
		uiforms.StartMainForm(&form)
		Window := form.WindowToCreate
		form.Dispose()

		uiforms.StartMainForm(Window)
	}
}
