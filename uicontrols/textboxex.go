package uicontrols

import (
	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiinterfaces"
)

type TextBoxExt struct {
	Container
	CurrentItemIndex int
	btnOpen          *Button
	txtBlock         *TextBox
	OnSelect         func(textBoxExt *TextBoxExt)
	OnTextChanged    func(txtBox *TextBoxExt, oldValue string, newValue string)
}

func NewTextBoxExt(parent uiinterfaces.Widget, text string, onSelect func(textBoxExt *TextBoxExt)) *TextBoxExt {
	var c TextBoxExt
	c.InitControl(parent, &c)
	c.OnSelect = onSelect

	c.txtBlock = NewTextBox(&c)
	c.txtBlock.OnTextChanged = func(txtBox *TextBox, oldValue string, newValue string) {
		if c.OnTextChanged != nil {
			c.OnTextChanged(&c, oldValue, newValue)
		}
	}
	c.AddWidgetOnGrid(c.txtBlock, 0, 0)
	c.btnOpen = NewButton(&c, "...", func(event *uievents.Event) {
		if c.OnSelect != nil {
			c.OnSelect(&c)
		}
	})
	c.AddWidgetOnGrid(c.btnOpen, 1, 0)

	c.cellPadding = 0
	c.panelPadding = 0

	return &c
}

func (c *TextBoxExt) SetText(text string) {
	if c.txtBlock != nil {
		c.txtBlock.SetText(text)
	}
}

func (c *TextBoxExt) Text() string {
	if c.txtBlock != nil {
		return c.txtBlock.Text()
	}
	return ""
}

func (c *TextBoxExt) ControlType() string {
	return "TextBoxExt"
}
