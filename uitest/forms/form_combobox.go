package forms

import (
	"github.com/gazercloud/gazerui/uicontrols"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiforms"
	"strconv"
)

type FormComboBox struct {
	uiforms.Form

	cmbBox          *uicontrols.ComboBox
	cmbBoxManyItems *uicontrols.ComboBox
	txt             *uicontrols.TextBlock
}

func (c *FormComboBox) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormComboBox")
	p1 := c.Panel().AddPanelOnGrid(0, 0)
	c.cmbBox = p1.AddComboBoxOnGrid(0, 0)
	c.cmbBox.AddItem("Первый элемент", nil)
	c.cmbBox.AddItem("Второй элемент", nil)
	c.cmbBox.AddItem("Третий элемент", nil)

	c.cmbBoxManyItems = p1.AddComboBoxOnGrid(0, 1)
	for i := 0; i < 40; i++ {
		c.cmbBoxManyItems.AddItem("Элемент "+strconv.Itoa(i), nil)
	}

	p2 := c.Panel().AddPanelOnGrid(1, 0)
	p2.AddButtonOnGrid(0, 0, "Set second item", func(event *uievents.Event) {
		c.cmbBox.SetCurrentItemIndex(1)
		c.cmbBoxManyItems.SetCurrentItemIndex(20)
	})

	c.txt = p2.AddTextBlockOnGrid(0, 1, "---")

	c.cmbBoxManyItems.OnCurrentIndexChanged = func(event *uicontrols.ComboBoxEvent) {
		c.txt.SetText(event.Item.Text)
	}

	c.Panel().AddVSpacerOnGrid(0, 10)
}
