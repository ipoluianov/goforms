package uicontrols

import (
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
)

type RadioButton struct {
	Container
	text    string
	checked bool

	imgBox  *ImageBox
	txtText *TextBlock

	OnCheckedChanged func(checkBox *RadioButton, checked bool)
}

func NewRadioButton(parent uiinterfaces.Widget, text string) *RadioButton {
	var c RadioButton
	c.text = text
	c.InitControl(parent, &c)
	return &c
}

func (c *RadioButton) InitControl(parent uiinterfaces.Widget, w uiinterfaces.Widget) {
	c.Container.InitControl(parent, w)
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

func (c *RadioButton) ControlType() string {
	return "CheckBox"
}

func (c *RadioButton) updateImage() {
	if c.checked {
		//c.imgBox.SetImage(uiresources.ResImageAdjusted("icons/material/toggle/drawable-hdpi/ic_radio_button_checked_black_24dp.png", c.ForeColor()))
		c.imgBox.SetScaling(ImageBoxScaleAdjustImageKeepAspectRatio)
	} else {
		//c.imgBox.SetImage(uiresources.ResImageAdjusted("icons/material/toggle/drawable-hdpi/ic_radio_button_unchecked_black_24dp.png", c.ForeColor()))
		c.imgBox.SetScaling(ImageBoxScaleAdjustImageKeepAspectRatio)
	}
}

func (c *RadioButton) ClearRadioButtons() {
	if c.checked {
		c.checked = false
		if c.OnCheckedChanged != nil {
			c.OnCheckedChanged(c, c.checked)
		}
	}

	c.updateImage()
	c.Update("CheckBox")
}

func (c *RadioButton) MouseClick(event *uievents.MouseClickEvent) {
	if c.parent != nil {
		c.parent.ClearRadioButtons()
	}

	c.checked = !c.checked

	c.updateImage()

	if c.OnCheckedChanged != nil {
		c.OnCheckedChanged(c, c.checked)
	}

	c.Update("CheckBox")
}

func (c *RadioButton) IsChecked() bool {
	return c.checked
}

func (c *RadioButton) SetChecked(isChecked bool) {
	if c.parent != nil {
		c.parent.ClearRadioButtons()
	}

	c.checked = isChecked

	c.updateImage()
	if c.OnCheckedChanged != nil {
		c.OnCheckedChanged(c, c.checked)
	}
	c.Update("CheckBox")
}
