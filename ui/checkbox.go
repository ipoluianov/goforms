package ui

import (
	"embed"

	"github.com/ipoluianov/goforms/uiresources"
	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/nui/nuimouse"
)

type CheckBox struct {
	Container
	text    string
	checked bool

	imgBox  *ImageBox
	txtText *TextBlock

	OnCheckedChanged func(checkBox *CheckBox, checked bool)
}

func NewCheckBox(parent Widget, text string) *CheckBox {
	_ = embed.FS.Open
	var c CheckBox
	c.text = text
	c.InitControl(parent, &c)
	return &c
}

//go:embed "res/checkbox_checked.png"
var checkbox_checked []byte

//go:embed "res/checkbox_unchecked.png"
var checkbox_unchecked []byte

func (c *CheckBox) InitControl(parent Widget, w Widget) {
	c.Container.InitControl(parent, w)
	c.SetPanelPadding(0)
	c.imgBox = NewImageBox(c, nil)
	c.imgBox.SetMinWidth(22)
	c.imgBox.SetMinHeight(22)
	c.imgBox.SetMaxWidth(22)
	c.imgBox.SetMaxHeight(22)
	c.AddWidgetOnGrid(c.imgBox, 0, 0)
	c.txtText = NewTextBlock(c, c.text)
	c.txtText.TextHAlign = canvas.HAlignLeft
	c.txtText.SetMouseCursor(nuimouse.MouseCursorPointer)
	c.AddWidgetOnGrid(c.txtText, 1, 0)
	c.updateImage()
}

func (c *CheckBox) Dispose() {
	c.imgBox = nil
	c.txtText = nil
	c.Container.Dispose()
}

func (c *CheckBox) ControlType() string {
	return "CheckBox"
}

func (c *CheckBox) SetEnabled(enabled bool) {
	c.Container.SetEnabled(enabled)
	if c.txtText != nil {
		c.txtText.SetEnabled(enabled)
	}
	if c.imgBox != nil {
		c.imgBox.SetEnabled(enabled)
	}
	c.updateImage()
}

func (c *CheckBox) EnabledChanged(enabled bool) {
	c.updateImage()
	c.Update("CheckBox")
}

func (c *CheckBox) updateImage() {
	if c.checked {
		c.imgBox.SetImage(uiresources.ResImgCol(checkbox_checked, c.ForeColor()))
		c.imgBox.SetScaling(ImageBoxScaleAdjustImageKeepAspectRatio)
		c.imgBox.SetMouseCursor(nuimouse.MouseCursorPointer)
	} else {
		c.imgBox.SetImage(uiresources.ResImgCol(checkbox_unchecked, c.ForeColor()))
		c.imgBox.SetScaling(ImageBoxScaleAdjustImageKeepAspectRatio)
		c.imgBox.SetMouseCursor(nuimouse.MouseCursorPointer)
	}
}

func (c *CheckBox) MouseClick(event *MouseClickEvent) {
	c.checked = !c.checked
	c.updateImage()

	if c.OnCheckedChanged != nil {
		c.OnCheckedChanged(c, c.checked)
	}

	c.Update("CheckBox")
}

func (c *CheckBox) IsChecked() bool {
	return c.checked
}

func (c *CheckBox) SetChecked(isChecked bool) {
	if c.checked != isChecked {
		c.checked = isChecked
		if c.OnCheckedChanged != nil {
			c.OnCheckedChanged(c, c.checked)
		}
	}

	c.updateImage()
	c.Update("CheckBox")
}

func (c *CheckBox) SetText(text string) {
	c.txtText.SetText(text)
}

func (c *CheckBox) Text() string {
	if c.txtText != nil {
		return c.txtText.Text()
	}
	return ""
}
