package example01

import "github.com/ipoluianov/goforms/ui"

func newMainForm() *ui.Form {
	form := ui.NewForm()

	panelTop := form.Panel().AddVPanel()
	lblText := panelTop.AddTextBlock("TextBlock value")
	txtBox := panelTop.AddTextBox()
	txtBox.OnTextChanged = func(txtBox *ui.TextBox, oldValue, newValue string) {
		lblText.SetText(newValue)
	}
	panelTop.AddVSpacer()

	// Bottom Button Box
	panelBottom := form.Panel().AddHPanel()
	panelBottom.AddHSpacer()
	panelBottom.AddButton("Change text", func(event *ui.Event) { lblText.SetText("Button is pressed") })
	panelBottom.AddButton("Cloce", func(event *ui.Event) { form.Close() })

	return form
}

func ExecMainForm() {
	ui.InitUI()
	ui.StartMainForm(newMainForm())
}
