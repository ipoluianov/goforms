package main

import (
	"github.com/ipoluianov/goforms/ui"
	"github.com/ipoluianov/goforms/uiforms"
	"github.com/ipoluianov/goforms/uitest/forms"
)

func main() {
	ui.InitUISystem()
	var mainForm forms.MainForm
	uiforms.StartMainForm(&mainForm)
}
