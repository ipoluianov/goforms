package forms

import (
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiforms"
	"github.com/gazercloud/gazerui/uiresources"
)

type FormImageBox struct {
	uiforms.Form

	imgBox *uicontrols.ImageBox
}

func (c *FormImageBox) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormImageBox")

	panelImg := c.Panel().AddPanelOnGrid(0, 0)
	c.imgBox = panelImg.AddImageBoxOnGrid(0, 0, nil)
	c.imgBox.SetFixedSize(500, 300)
	c.Panel().AddVSpacerOnGrid(0, 1)
	panelControl := c.Panel().AddPanelOnGrid(0, 2)

	panelControl.AddButtonOnGrid(0, 0, "Set Image default", func(event *uievents.Event) {
		img := uiresources.ResImg(uiresources.R_icons_material4_png_action_account_box_materialicons_48dp_1x_baseline_account_box_black_48dp_png)
		c.imgBox.SetImage(img)
	})
	panelControl.AddButtonOnGrid(0, 1, "Set Image adjusted", func(event *uievents.Event) {
		img := uiresources.ResImgCol(uiresources.R_icons_material4_png_av_play_arrow_materialicons_48dp_1x_baseline_play_arrow_black_48dp_png, c.Panel().ForeColor())
		c.imgBox.SetImage(img)
	})

	panelControl.AddButtonOnGrid(0, 2, "Set NoScaleAdjustBox", func(event *uievents.Event) {
		c.imgBox.SetScaling(uicontrols.ImageBoxScaleNoScaleAdjustBox)
	})
	panelControl.AddButtonOnGrid(0, 3, "Set NoScaleImageInLeftTop", func(event *uievents.Event) {
		c.imgBox.SetScaling(uicontrols.ImageBoxScaleNoScaleImageInLeftTop)
	})
	panelControl.AddButtonOnGrid(0, 4, "Set NoScaleImageInCenter", func(event *uievents.Event) {
		c.imgBox.SetScaling(uicontrols.ImageBoxScaleNoScaleImageInCenter)
	})
	panelControl.AddButtonOnGrid(0, 5, "Set StretchImage", func(event *uievents.Event) {
		c.imgBox.SetScaling(uicontrols.ImageBoxScaleStretchImage)
	})
	panelControl.AddButtonOnGrid(0, 6, "Set AdjustImageKeepAspectRatio", func(event *uievents.Event) {
		c.imgBox.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	})
}
