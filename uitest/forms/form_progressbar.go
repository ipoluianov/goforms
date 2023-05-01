package forms

import (
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiforms"
)

type FormProgressBar struct {
	uiforms.Form

	pr *uicontrols.ProgressBar
}

func (c *FormProgressBar) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormProgressBar")

	p1 := c.Panel().AddPanelOnGrid(0, 0)
	c.pr = p1.AddProgressBarOnGrid(0, 0)
	p2 := c.Panel().AddPanelOnGrid(0, 1)
	p2.AddButtonOnGrid(0, 0, "Do It", func(event *uievents.Event) {
		c.pr.SetValue(78)
		c.pr.SetText("Hello")
	})

	c.Panel().AddVSpacerOnGrid(0, 2)
}
