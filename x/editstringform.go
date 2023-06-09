package x

/*
import (
	"github.com/ipoluianov/goforms/uicontrols"
	"github.com/ipoluianov/goforms/uiinterfaces"
)

type EditStringForm struct {
	uicontrols.Dialog

	OnOK func(newValue string, oldValue string)

	initialText string
	txtText     *uicontrols.TextBox
	btnOK       *uicontrols.Button
	btnCancel   *uicontrols.Button
}

func ShowEditStringForm(parent uiinterfaces.Widget, initialValue string, OnOK func(newValue string, oldValue string)) {
	dialog := NewEditStringForm(parent, initialValue, OnOK)
	dialog.ShowDialog()
}

func NewEditStringForm(parent uiinterfaces.Widget, text string, OnOK func(newValue string, oldValue string)) *EditStringForm {
	var c EditStringForm
	c.OnOK = OnOK
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	c.txtText = pContent.AddTextBoxOnGrid(0, 0)
	c.txtText.SetText(text)
	c.initialText = text

	pContent.AddVSpacerOnGrid(0, 5)

	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)
	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.TryAccept = func() bool {
		if c.OnOK != nil {
			c.OnOK(c.txtText.Text(), c.initialText)
		}
		return true
	}
	c.btnOK.SetMinWidth(70)

	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	btnCancel.SetMinWidth(70)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(btnCancel)

	c.Resize(500, 300)
	c.SetTitle("Edit")

	return &c
}
*/
