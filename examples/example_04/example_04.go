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
	btn22.SetImageBeforeText(false)
	row2.AddHSpacer()

	row3 := vPanel.AddHPanel()
	btn31 := row3.AddButton("Text Button", func(event *ui.Event) {})
	btn31.SetImage(ui.DecodeImage(button_icon))
	btn31.SetTextImageVerticalOrientation(false)
	btn32 := row3.AddButton("Text Button", func(event *ui.Event) {})
	btn32.SetImage(ui.DecodeImage(button_icon))
	btn32.SetTextImageVerticalOrientation(false)
	btn32.SetImageBeforeText(false)
	btn32.SetTextHAlign(canvas.HAlignRight)
	//btn32.SetMinWidth(300)
	row3.AddHSpacer()

	vPanel.AddVSpacer()

	return form
}

func ExecMainForm() {
	ui.InitUI()
	ui.StartMainForm(newMainForm())
}
