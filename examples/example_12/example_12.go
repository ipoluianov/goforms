package example12

import "github.com/ipoluianov/goforms/ui"

func newMainForm() *ui.Form {
	form := ui.NewForm()
	tabControl := form.Panel().AddTabControl()
	page1 := tabControl.AddPage()
	page1.SetText("Page1")
	page1.AddTextBlock("123")
	page2 := tabControl.AddPage()
	page2.SetText("Page2")
	page3 := tabControl.AddPage()
	page3.SetText("Page3")

	return form
}

func ExecMainForm() {
	ui.StartMainForm(newMainForm())
}
