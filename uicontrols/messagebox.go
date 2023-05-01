package uicontrols

import (
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type MessageBox struct {
	Dialog
	lblText      *TextBlock
	pButtons     *Panel
	text         string
	header       string
	typeOfDialog string

	onYes    func()
	onNo     func()
	onOK     func()
	onCancel func()
}

func ShowInformationMessage(parent uiinterfaces.Widget, text string, header string) {
	c := NewMessageBox(parent, "info")
	c.SetText(text)
	c.SetHeader(header)
	c.ShowDialog()
}

func ShowErrorMessage(parent uiinterfaces.Widget, text string, header string) {
	c := NewMessageBox(parent, "error")
	c.typeOfDialog = "error"
	c.SetText(text)
	c.SetHeader(header)
	c.ShowDialog()
}

func ShowQuestionMessageOKCancel(parent uiinterfaces.Widget, text string, header string, onOK func(), onCancel func()) {
	c := NewMessageBox(parent, "question_ok_cancel")
	c.onOK = onOK
	c.onCancel = onCancel
	c.SetText(text)
	c.SetHeader(header)
	c.ShowDialog()
}

func ShowQuestionMessageYesNoCancel(parent uiinterfaces.Widget, text string, header string, onYes func(), onNo func(), onCancel func()) {
	c := NewMessageBox(parent, "question_yes_no_cancel")
	c.onYes = onYes
	c.onNo = onNo
	c.onCancel = onCancel
	c.SetText(text)
	c.SetHeader(header)
	c.ShowDialog()
}

func NewMessageBox(parent uiinterfaces.Widget, typeOfDialog string) *MessageBox {
	var c MessageBox
	c.typeOfDialog = typeOfDialog
	c.InitControl(parent, &c)
	return &c
}

func (c *MessageBox) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle(c.header)

	if c.typeOfDialog == "info" {
		c.Resize(400, 200)
		c.lblText = c.ContentPanel().AddTextBlockOnGrid(0, 0, c.text)
		c.lblText.SetXExpandable(true)
		c.lblText.TextHAlign = canvas.HAlignCenter
	}
	if c.typeOfDialog == "error" {
		c.Resize(400, 200)
		c.lblText = c.ContentPanel().AddTextBlockOnGrid(0, 0, c.text)
		c.lblText.SetXExpandable(true)
		c.lblText.TextHAlign = canvas.HAlignCenter
	}
	if c.typeOfDialog == "question_ok_cancel" {
		c.Resize(400, 200)
		c.lblText = c.ContentPanel().AddTextBlockOnGrid(0, 0, c.text)
		c.lblText.SetXExpandable(true)
		c.lblText.SetYExpandable(true)
		c.lblText.TextHAlign = canvas.HAlignCenter
		pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)
		pButtons.AddHSpacerOnGrid(0, 0)
		btnOK := pButtons.AddButtonOnGrid(1, 0, "OK", nil)
		btnOK.SetMinWidth(70)
		btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
		btnCancel.SetMinWidth(70)
		c.SetAcceptButton(btnOK)
		c.SetRejectButton(btnCancel)
		c.TryAccept = func() bool {
			if c.onOK != nil {
				c.onOK()
			}
			return true
		}
		c.OnReject = func() {
			if c.onCancel != nil {
				c.onCancel()
			}
		}
	}

	if c.typeOfDialog == "question_yes_no_cancel" {
		c.Resize(400, 200)
		c.lblText = c.ContentPanel().AddTextBlockOnGrid(0, 0, c.text)
		c.lblText.SetXExpandable(true)
		c.lblText.SetYExpandable(true)
		c.lblText.TextHAlign = canvas.HAlignCenter
		pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)
		pButtons.AddHSpacerOnGrid(0, 0)
		btnYes := pButtons.AddButtonOnGrid(1, 0, "Yes", nil)
		btnYes.SetMinWidth(70)
		btnNo := pButtons.AddButtonOnGrid(2, 0, "No", nil)
		btnNo.SetMinWidth(70)
		btnNo.onPress = func(event *uievents.Event) {
			if c.onNo != nil {
				c.onNo()
			}
			c.Close()
		}
		btnCancel := pButtons.AddButtonOnGrid(3, 0, "Cancel", nil)
		btnCancel.SetMinWidth(70)
		c.SetAcceptButton(btnYes)
		c.SetRejectButton(btnCancel)
		c.TryAccept = func() bool {
			if c.onYes != nil {
				c.onYes()
			}
			return true
		}
		c.OnReject = func() {
			if c.onCancel != nil {
				c.onCancel()
			}
		}
	}

}

func (c *MessageBox) SetText(text string) {
	c.text = text
	if c.lblText != nil {
		c.lblText.SetText(text)
	}
}

func (c *MessageBox) SetHeader(header string) {
	c.header = header
	c.SetTitle(header)
}
