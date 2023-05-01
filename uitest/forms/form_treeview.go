package forms

import (
	"fmt"
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uiforms"
)

type FormTreeView struct {
	uiforms.Form

	treeView *uicontrols.TreeView
}

func (c *FormTreeView) FillNode(node *uicontrols.TreeNode) {
	for i := 0; i < 20; i++ {
		n := c.treeView.AddNode(node, "node_"+fmt.Sprint(i))
		c.treeView.SetNodeValue(n, 1, "val1")
		c.treeView.SetNodeValue(n, 2, "val2")
	}
}

func (c *FormTreeView) OnInit() {
	c.Resize(600, 800)
	c.SetTitle("FormTreeView")

	c.treeView = uicontrols.NewTreeView(c.Panel())
	c.Panel().AddWidget(c.treeView)

	c.treeView.AddColumn("Column 1", 100)
	c.treeView.AddColumn("Column 2", 100)

	node := c.treeView.AddNode(nil, "Root")

	for i := 0; i < 10; i++ {
		n := c.treeView.AddNode(node, "rootNode_"+fmt.Sprint(i))
		c.FillNode(n)
	}
}
