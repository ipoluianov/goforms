package forms

import (
	"time"

	"github.com/ipoluianov/goforms/uicontrols"
	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiforms"
)

type FormTabControl struct {
	uiforms.Form

	tab  *uicontrols.TabControl
	txt1 *uicontrols.TextBlock
	txt2 *uicontrols.TextBlock
}

func (c *FormTabControl) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormTabControl")

	c.tab = c.Panel().AddTabControlOnGrid(0, 0)
	p1 := c.tab.AddPage()
	p1.SetText("111")

	p2 := c.tab.AddPage()
	p2.SetText("222")

	{
		p1.AddButtonOnGrid(0, 1, "btn1", nil)
		p1.AddButtonOnGrid(0, 2, "btn2", func(event *uievents.Event) {
			c.txt1.SetText(time.Now().String())
		})
		p1.AddButtonOnGrid(0, 3, "btn3", nil)
		c.txt1 = p1.AddTextBlockOnGrid(0, 4, "---")
		p1.AddVSpacerOnGrid(0, 10)
	}

	{
		p2.AddButtonOnGrid(0, 0, "btn11", nil)
		p2.AddButtonOnGrid(0, 1, "btn22", func(event *uievents.Event) {
			c.txt2.SetText(time.Now().String() + "+++")
		})
		p2.AddButtonOnGrid(0, 2, "btn33", nil)
		c.txt2 = p2.AddTextBlockOnGrid(0, 3, "---+++")
		tv := p2.AddTreeViewOnGrid(0, 4)
		tv.AddColumn("qqq", 100)
		tv.AddColumn("www", 200)
		tv.AddNode(nil, "111")
	}
}
