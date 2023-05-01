package uiforms

import (
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type MessageBox struct {
	Form
	text     string
	txtBlock *uicontrols.TextBlock
	btnOK    *uicontrols.Button
}

func MessageBoxError(parent uiinterfaces.Widget, err error) {
	ShowMessageBox(parent.Window(), "Error", err.Error())
}

func (f *MessageBox) onBtnOK(event *uievents.Event) {
	f.Close()
}

func (f *MessageBox) OnInit() {
	f.SetTitle("MessageBox")
	f.Resize(200, 200)

	// Text
	f.txtBlock = f.Panel().AddTextBlockOnGrid(0, 0, f.text)

	// Button OK
	f.btnOK = f.Panel().AddButtonOnGrid(0, 1, "OK", f.onBtnOK)

	f.adjustSizeToContent(f.text)

	f.onSizeChanged = f.onFormSizeChanged
}

func (f *MessageBox) onFormSizeChanged(event *uievents.FormSizeChangedEvent) {
	f.btnOK.SetX(f.Width()/2 - f.btnOK.Width()/2)
}

func (f *MessageBox) adjustSizeToContent(text string) {
	textWidth, textHeight, _ := canvas.MeasureText(f.txtBlock.FontFamily(), f.txtBlock.FontSize(), f.txtBlock.FontBold(), f.txtBlock.FontItalic(), text, true)
	if textWidth < 300 {
		textWidth = 300
	}

	if textHeight < 100 {
		textHeight = 100
	}
	f.Resize(textWidth+10, textHeight+10)
}

func (f *MessageBox) SetText(text string) {
	f.text = text
	if f.txtBlock != nil {
		f.txtBlock.SetText(f.text)
		f.adjustSizeToContent(f.text)
	}
}

func ShowMessageBox(parent uiinterfaces.Window, title string, text string) {
	var form MessageBox
	form.SetTitle(title)
	form.SetText(text)
	StartModalForm(parent, &form)
}
