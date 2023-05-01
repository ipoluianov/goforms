package forms

import (
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiforms"
	"github.com/gazercloud/gazerui/uiresources"
)

type FormButton struct {
	uiforms.Form

	btn *uicontrols.Button

	btnSetHAlignLeft   *uicontrols.Button
	btnSetHAlignCenter *uicontrols.Button
	btnSetHAlignRight  *uicontrols.Button

	btnSetVAlignTop    *uicontrols.Button
	btnSetVAlignCenter *uicontrols.Button
	btnSetVAlignBottom *uicontrols.Button

	btnSetText *uicontrols.Button
	txtText    *uicontrols.TextBox

	data []int
}

func (c *FormButton) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormButton")

	c.btn = c.Panel().AddButtonOnGrid(0, 0, "Button", c.onBtn)

	grHAlign := c.Panel().AddGroupBoxOnGrid(0, 1, "HAlign")
	grHAlign.Panel().AddButtonOnGrid(0, 0, "HAlign Left", c.onBtnSetHAlignLeft)
	grHAlign.Panel().AddButtonOnGrid(1, 0, "HAlign Center", c.onBtnSetHAlignCenter)
	grHAlign.Panel().AddButtonOnGrid(2, 0, "HAlign Right", c.onBtnSetHAlignRight)

	grVAlign := c.Panel().AddGroupBoxOnGrid(0, 2, "VAlign")
	grVAlign.Panel().AddButtonOnGrid(0, 0, "VAlign Top", c.onBtnSetVAlignTop)
	grVAlign.Panel().AddButtonOnGrid(1, 0, "VAlign Center", c.onBtnSetVAlignCenter)
	grVAlign.Panel().AddButtonOnGrid(2, 0, "VAlign Bottom", c.onBtnSetVAlignBottom)

	c.txtText = c.Panel().AddTextBoxOnGrid(0, 3)
	c.txtText.SetMultiline(true)
	c.Panel().AddButtonOnGrid(0, 4, "Set Text", c.onBtnSetText)

	c.Panel().AddButtonOnGrid(0, 5, "Set picture", func(event *uievents.Event) {
		c.btn.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_refresh_materialicons_48dp_1x_baseline_refresh_black_48dp_png, c.Panel().ForeColor()))
		c.btn.SetImageSize(100, 100)
	})
}

func (c *FormButton) onBtn(event *uievents.Event) {
}

func (c *FormButton) onBtnSetHAlignLeft(event *uievents.Event) {
	c.btn.SetTextHAlign(canvas.HAlignLeft)
}

func (c *FormButton) onBtnSetHAlignCenter(event *uievents.Event) {
	c.btn.SetTextHAlign(canvas.HAlignCenter)
}

func (c *FormButton) onBtnSetHAlignRight(event *uievents.Event) {
	c.btn.SetTextHAlign(canvas.HAlignRight)
}

func (c *FormButton) onBtnSetVAlignTop(event *uievents.Event) {
	c.btn.SetTextVAlign(canvas.VAlignTop)
}

func (c *FormButton) onBtnSetVAlignCenter(event *uievents.Event) {
	c.btn.SetTextVAlign(canvas.VAlignCenter)
}

func (c *FormButton) onBtnSetVAlignBottom(event *uievents.Event) {
	c.btn.SetTextVAlign(canvas.VAlignBottom)
}

func (c *FormButton) onBtnSetText(event *uievents.Event) {
	c.btn.SetText(c.txtText.Text())
}
