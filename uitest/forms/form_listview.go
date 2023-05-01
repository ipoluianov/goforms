package forms

import (
	"fmt"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiforms"
)

type FormListView struct {
	uiforms.Form

	lvItems *uicontrols.ListView
}

func (c *FormListView) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormListView")
	c.lvItems = c.Panel().AddListViewOnGrid(0, 0)
	c.lvItems.AddColumn("Column 1", 200)
	//c.lvItems.AddColumn("Column 2", 50)
	//c.lvItems.AddColumn("Column 3", 100)

	for i := 0; i < 100; i++ {
		item := c.lvItems.AddItem("Элемент " + fmt.Sprint(i))
		item.SetValue(1, "111")
		item.SetValue(2, "222")
	}

	p := c.Panel().AddPanelOnGrid(1, 0)
	p.AddButtonOnGrid(0, 0, "Show header", func(event *uievents.Event) {
		c.lvItems.SetHeaderVisible(!c.lvItems.IsHeaderVisible())
	})
	p.AddVSpacerOnGrid(0, 100)

}
