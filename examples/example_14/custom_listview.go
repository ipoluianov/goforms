package example14

import "github.com/ipoluianov/goforms/ui"

type CustomListView struct {
	ui.ListView
}

func NewCustomListView(parent ui.Widget) *CustomListView {
	var c CustomListView
	c.InitControl(parent, &c)
	return &c
}
