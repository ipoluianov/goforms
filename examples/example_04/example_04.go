package example04

import (
	"embed"

	"github.com/ipoluianov/goforms/ui"
	"github.com/ipoluianov/goforms/utils/canvas"
)

//go:embed button_icon.png
var button_icon []byte

func newMainForm() *ui.Form {
	_ = embed.FS.Open
	form := ui.NewForm()
	form.SetTitle("Buttons")

	vPanel := form.Panel().AddVPanel()

	row1 := vPanel.AddHPanel()
	btn11 := row1.AddButton("Text Button", func(event *ui.Event) {})
	btn11.SetMinWidth(150)
	btn11.SetTextHAlign(canvas.HAlignLeft)
	btn12 := row1.AddButton("Text Button", func(event *ui.Event) {})
	btn12.SetMinWidth(150)
	btn12.SetTextHAlign(canvas.HAlignCenter)
	btn13 := row1.AddButton("Text Button", func(event *ui.Event) {})
	btn13.SetMinWidth(150)
	btn13.SetTextHAlign(canvas.HAlignRight)
	row1.AddHSpacer()

	row2 := vPanel.AddHPanel()
	btn21 := row2.AddButton("Text Button", func(event *ui.Event) {})
	btn21.SetImage(ui.DecodeImage(button_icon))
	btn22 := row2.AddButton("Text Button", func(event *ui.Event) {})
	btn22.SetImage(ui.DecodeImage(button_icon))
	btn23 := row2.AddButton("Text Button", func(event *ui.Event) {})
	btn23.SetImage(ui.DecodeImage(button_icon))
	row2.AddHSpacer()

	vPanel.AddVSpacer()

	return form
}

func ExecMainForm() {
	ui.InitUI()
	ui.StartMainForm(newMainForm())
}
