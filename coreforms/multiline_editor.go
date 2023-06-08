package coreforms

/*
import (
	"github.com/ipoluianov/goforms/uicontrols"
	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiinterfaces"
)

type MultilineEditor struct {
	uicontrols.Dialog
	resValue      string
	txtUnitName   *uicontrols.TextBox
	btnOK         *uicontrols.Button
	OnTextChanged func(txtMultiline *MultilineEditor, oldText string, newText string)
}

func NewMultilineEditor(parent uiinterfaces.Widget, value string) *MultilineEditor {
	var c MultilineEditor
	c.resValue = value
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	pRight.AddTextBlockOnGrid(0, 0, "Text:")
	c.txtUnitName = pRight.AddTextBoxOnGrid(1, 0)
	c.txtUnitName.SetText(value)
	c.txtUnitName.SetMultiline(true)
	c.txtUnitName.OnTextChanged = func(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
		if c.OnTextChanged != nil {
			c.OnTextChanged(&c, oldValue, newValue)
		}
	}

	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", func(event *uievents.Event) {
		c.Accept()
	})
	c.TryAccept = func() bool {
		c.btnOK.SetEnabled(false)
		c.resValue = c.txtUnitName.Text()
		c.TryAccept = nil
		c.Accept()
		return false
	}

	c.btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetRejectButton(btnCancel)

	return &c
}

func (c *MultilineEditor) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Edit value")
	c.Resize(600, 400)
}

func (c *MultilineEditor) Text() string {
	if c.txtUnitName == nil {
		return ""
	}
	return c.txtUnitName.Text()
}

func (c *MultilineEditor) ResultText() string {
	return c.resValue
}
*/
