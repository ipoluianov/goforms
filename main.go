package main

import (
	"github.com/ipoluianov/goforms/ui"
)

func main() {
	ui.InitUI()
	form := ui.NewForm()
	ui.StartMainForm(form)
}
