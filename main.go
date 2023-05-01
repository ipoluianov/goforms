package main

import (
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiforms"
	"github.com/gazercloud/gazerui/uitest/forms"
)

//
func main() {
	ui.InitUISystem()
	var mainForm forms.MainForm
	uiforms.StartMainForm(&mainForm)
}
