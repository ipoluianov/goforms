package example14

import (
	"fmt"

	"github.com/ipoluianov/goforms/ui"
)

type CustomListView struct {
	ui.ListView
}

func NewCustomListView(parent ui.Widget) *CustomListView {
	var c CustomListView
	c.InitControl(parent, &c)
	c.Construct()
	c.AddColumn("Col1", 100)

	for i := 0; i < 100; i++ {
		c.AddItem("item_" + fmt.Sprint(i))
	}

	return &c
}

func (c *CustomListView) OnInit() {
}
