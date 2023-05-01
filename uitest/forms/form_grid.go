package forms

import (
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiforms"
	"github.com/gazercloud/gazerui/uiresources"
	"golang.org/x/image/colornames"
)

type FormGrid struct {
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

func (c *FormGrid) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormButton")

	p := c.Panel().AddPanelOnGrid(0, 0)
	p.SetName("#42")
	p.SetBorders(1, colornames.Gray)

	//tt := c.Panel().AddImageBoxOnGrid(0, 0, uiresources.ResImageAdjusted("icons/material/image/drawable-hdpi/ic_blur_on_black_48dp.png", c.Panel().ForeColor()))
	t1 := p.AddTextBlockOnGrid(0, 0, "111")
	t1.SetBorders(1, colornames.Green)

	t2 := c.Panel().AddImageBoxOnGrid(1, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_action_account_box_materialicons_48dp_1x_baseline_account_box_black_48dp_png, c.Panel().ForeColor()))
	//t2 := c.Panel().AddTextBlockOnGrid(1, 0, "222")
	t2.SetBorders(1, colornames.Yellow)

}
