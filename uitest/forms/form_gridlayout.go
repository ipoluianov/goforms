package forms

import (
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiforms"
)

type FormGridLayout struct {
	uiforms.Form

	btn1 *uicontrols.Button
	btn2 *uicontrols.Button
	btn3 *uicontrols.Button
	btn4 *uicontrols.Button
}

func (c *FormGridLayout) OnInit() {
	c.Resize(800, 800)
	c.SetTitle("FormGridLayout")
	c.Panel().SetAbsolutePositioning(false)
	c.btn1 = c.Panel().AddButtonOnGrid(0, 0, "100-150", nil)
	c.btn1.SetXExpandable(true)
	c.btn1.SetMinWidth(100)
	c.btn1.SetMaxWidth(150)
	c.btn1.SetMinHeight(100)
	c.btn1.SetMaxHeight(150)
	c.btn2 = c.Panel().AddButtonOnGrid(1, 0, "222", nil)
	c.btn2.SetMinWidth(200)
	c.btn3 = c.Panel().AddButtonOnGrid(0, 1, "333", nil)
	c.btn4 = c.Panel().AddButtonOnGrid(1, 1, "444", nil)
}
