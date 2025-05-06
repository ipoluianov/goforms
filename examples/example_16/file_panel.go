package example16

import "github.com/ipoluianov/goforms/ui"

type FilePanel struct {
	ui.Panel

	lvItem *ui.ListView
}

func NewFilePanel(parent ui.Widget) *FilePanel {
	var c FilePanel
	c.InitControl(parent, &c)
	return &c
}

func (c *FilePanel) OnInit() {
	c.AddTextBlock("Top Text")
	c.lvItem = c.AddListView()
	c.lvItem.AddColumn("Name", 150)
	c.lvItem.AddColumn("Ext", 50)
	c.lvItem.AddColumn("Size", 100)
	c.lvItem.AddColumn("Date", 100)
	c.AddTextBlock("Bottom Text")
}
