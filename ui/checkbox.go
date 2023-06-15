package ui

import (
	"github.com/ipoluianov/goforms/uiresources"
	"github.com/ipoluianov/goforms/utils/canvas"
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
	var c CheckBox
	c.text = text
	c.InitControl(parent, &c)
	return &c
}

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
		c.imgBox.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_toggle_check_box_materialiconsoutlined_48dp_1x_outline_check_box_black_48dp_png, c.ForeColor())) /////////////////////////////////
		c.imgBox.SetScaling(ImageBoxScaleAdjustImageKeepAspectRatio)
	} else {
		c.imgBox.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_toggle_check_box_outline_blank_materialiconsoutlined_48dp_1x_outline_check_box_outline_blank_black_48dp_png, c.ForeColor()))
		c.imgBox.SetScaling(ImageBoxScaleAdjustImageKeepAspectRatio)
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
