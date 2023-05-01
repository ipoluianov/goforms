package forms

import (
	"github.com/ipoluianov/goforms/uicontrols"
	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiforms"
	"github.com/ipoluianov/goforms/uiinterfaces"
	"github.com/ipoluianov/goforms/uiresources"
)

type FormDialog struct {
	uiforms.Form
}

type Dialog struct {
	uicontrols.Dialog
}

func NewDialog(parent uiinterfaces.Widget) *Dialog {
	var c Dialog
	c.InitControl(parent, &c)
	return &c
}

func (c *Dialog) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Dialog Title")
	c.Resize(400, 200)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	pLeft := pContent.AddPanelOnGrid(0, 0)
	pRight := pContent.AddPanelOnGrid(1, 0)
	pLeft.SetPanelPadding(0)
	//pLeft.SetBorderRight(1, c.ForeColor())
	pLeft.SetMinWidth(100)
	t1 := pRight.AddTextBoxOnGrid(0, 0)
	t1.SetTabIndex(1)
	t2 := pRight.AddTextBoxOnGrid(0, 1)
	t2.SetTabIndex(2)
	pRight.AddVSpacerOnGrid(0, 2)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)

	img := pLeft.AddImageBoxOnGrid(0, 0, uiresources.ResImgCol(uiresources.R_icons_material4_png_av_stop_materialicons_48dp_1x_baseline_stop_black_48dp_png, c.ForeColor()))
	img.SetScaling(uicontrols.ImageBoxScaleAdjustImageKeepAspectRatio)
	img.SetMinHeight(64)
	img.SetMinWidth(64)
	pLeft.AddVSpacerOnGrid(0, 1)

	pButtons.AddHSpacerOnGrid(0, 0)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)
	btnCancel.SetTabIndex(3)

	t1.Focus()

}

func (c *FormDialog) OnInit() {
	c.Resize(600, 600)
	c.SetTitle("FormDialog")

	c.Panel().AddButtonOnGrid(0, 0, "Open Dialog", func(event *uievents.Event) {
		dlg := NewDialog(c.Panel())
		dlg.ShowDialog()
	})

	c.Panel().AddVSpacerOnGrid(0, 1)
}
