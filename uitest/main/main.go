package main

import (
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiforms"
	"github.com/gazercloud/gazerui/uitest/forms"
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
