package example03

import (
	"embed"

	"github.com/ipoluianov/goforms/ui"
)

//go:embed example_image.png
var exampleImageBS []byte

func newMainForm() *ui.Form {
	_ = embed.FS.Open
	form := ui.NewForm()
	panel := form.Panel().AddVPanel()
	panel.AddImageBox(ui.DecodeImage(exampleImageBS))
	panel.AddVSpacer()
	return form
}

func ExecMainForm() {
	ui.InitUI()
	ui.StartMainForm(newMainForm())
}
