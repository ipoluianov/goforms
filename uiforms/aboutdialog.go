package uiforms

import (
	"github.com/ipoluianov/goforms/ui"
	"github.com/ipoluianov/goforms/utils/canvas"
)

type AboutDialog struct {
	ui.Form
	text     string
	txtBlock *ui.TextBlock
	btnOK    *ui.Button
}

func (f *AboutDialog) onBtnOK(event *ui.Event) {
	f.Close()
}

func (f *AboutDialog) OnInit() {
	f.SetTitle("MessageBox")
	f.Resize(200, 200)

	// Text
	f.txtBlock = f.Panel().AddTextBlockOnGrid(0, 0, f.text)

	// Button OK
	f.btnOK = f.Panel().AddButtonOnGrid(0, 1, "OK", f.onBtnOK)

	f.adjustSizeToContent(f.text)

	//f.onSizeChanged = f.onFormSizeChanged
}

func (f *AboutDialog) onFormSizeChanged(event *ui.FormSizeChangedEvent) {
	f.btnOK.SetX(f.Width()/2 - f.btnOK.Width()/2)
}

func (f *AboutDialog) adjustSizeToContent(text string) {
	textWidth, textHeight, _ := canvas.MeasureText(f.txtBlock.FontFamily(), f.txtBlock.FontSize(), f.txtBlock.FontBold(), f.txtBlock.FontItalic(), text, true)
	if textWidth < 300 {
		textWidth = 300
	}

	if textHeight < 100 {
		textHeight = 100
	}
	f.Resize(textWidth+10, textHeight+10)
}

func (f *AboutDialog) SetText(text string) {
	f.text = text
	if f.txtBlock != nil {
		f.txtBlock.SetText(f.text)
		f.adjustSizeToContent(f.text)
	}
}

func ShowAboutDialog(parent ui.Window, title string, text string) {
	var form AboutDialog
	form.SetTitle(title)
	form.SetText(text)
	ui.StartModalForm(parent, &form)
}
