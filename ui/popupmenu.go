package ui

import (
	"image"
)

const PopupMenuItemHeight = 32

type PopupMenu struct {
	Panel
	menuWidth  int
	menuHeight int
	items      []*PopupMenuItem
	CloseEvent func()
	//closingAll func()
	parentMenu *PopupMenu
}

func NewPopupMenu(parent Widget) *PopupMenu {
	var c PopupMenu
	c.SetAbsolutePositioning(true)
	c.Panel.InitControl(parent.Window().CentralWidget(), &c)
	//c.Panel.SetSize(100, 100)
	c.SetName("PopupMenuPanel")
	return &c
}

func (c *PopupMenu) ControlType() string {
	return "PopupMenu"
}

func (c *PopupMenu) DisposeMenu() {
	if c.items != nil {
		for _, item := range c.items {
			if item != nil {
				item.Dispose()
			}
		}
	}
	c.items = nil
	c.Panel.Dispose()
	c.parentMenu = nil
}

func (c *PopupMenu) Dispose() {
}

/*func (c *PopupMenu) Draw(ctx ui.DrawContext) {
	ctx.SetColor(c.backgroundColor.Color())
	ctx.FillRect(0, 0, c.InnerWidth(), c.InnerHeight())
	c.Panel.Draw(ctx)
	ctx.SetColor(c.leftBorderColor.Color())
	ctx.SetStrokeWidth(1)
	ctx.DrawRect(0, 0, c.InnerWidth(), c.InnerHeight())
}*/

func (c *PopupMenu) ShowMenu(x int, y int) {
	c.SetX(x)
	c.SetY(y)
	c.rebuildVisualElements()
	c.Window().AppendPopup(c)
}

func (c *PopupMenu) showMenu(x int, y int, parentMenu *PopupMenu) {
	c.CloseAfterPopupWidget(parentMenu)
	c.parentMenu = parentMenu
	c.SetX(x)
	c.SetY(y)
	c.rebuildVisualElements()
	c.Window().AppendPopup(c)
}

func (c *PopupMenu) ClosePopup() {
	if c.CloseEvent != nil {
		c.CloseEvent()
	}
}

func (c *PopupMenu) AddItem(text string, onClick func(event *Event), img image.Image, keyCombination string) *PopupMenuItem {
	var item PopupMenuItem
	item.InitControl(c, &item)
	item.SetAnchors(ANCHOR_LEFT)
	item.SetText(text)
	item.OnClick = onClick
	item.Image = img
	item.KeyCombination = keyCombination
	pItem := &item
	c.items = append(c.items, pItem)
	c.Panel.addWidget(&item)
	return pItem
}

func (c *PopupMenu) AddItemWithUiResImage(text string, onClick func(event *Event), img []byte, keyCombination string) *PopupMenuItem {
	var item PopupMenuItem
	item.InitControl(c, &item)
	item.SetAnchors(ANCHOR_LEFT)
	item.SetText(text)
	item.OnClick = onClick
	item.ImageResource = img
	item.KeyCombination = keyCombination
	pItem := &item
	c.items = append(c.items, pItem)
	c.Panel.addWidget(&item)
	return pItem
}

func (c *PopupMenu) AddItemWithSubmenu(text string, img image.Image, innerMenu *PopupMenu) *PopupMenuItem {
	var item PopupMenuItem
	item.InitControl(c, &item)
	item.SetAnchors(ANCHOR_LEFT)
	item.SetText(text)
	item.Image = img
	item.innerMenu = innerMenu
	pItem := &item
	c.items = append(c.items, pItem)
	c.Panel.addWidget(&item)
	return pItem
}

func (c *PopupMenu) RemoveAllItems() {
	c.Panel.RemoveAllWidgets()
	c.rebuildVisualElements()
	c.Update("PopupMenu")
}

func (c *PopupMenu) OnInit() {
	c.rebuildVisualElements()
}

func (c *PopupMenu) needToClose() {
	c.Window().CloseTopPopup()
	if c.parentMenu != nil {
		c.parentMenu.needToClose()
	}
}

func (c *PopupMenu) rebuildVisualElements() {
	//c.Panel.RemoveAllWidgets()
	yOffset := 0
	for _, item := range c.items {

		//var item PopupMenuItem
		//item.InitControl(&c.Panel, &item, 0, 0, 200, PopupMenuItemHeight, 0)
		/*item.Text = itemOrig.Text
		item.OnClick = itemOrig.OnClick
		item.Image = itemOrig.Image
		item.KeyCombination = itemOrig.KeyCombination
		item.innerMenu = itemOrig.innerMenu*/
		item.needToClosePopupMenu = c.needToClose
		item.parentMenu = c

		item.SetX(0)
		item.SetY(yOffset)
		item.SetWidth(c.Width())
		item.SetHeight(PopupMenuItemHeight)
		yOffset += PopupMenuItemHeight
	}
	c.SetWidth(300)
	c.SetHeight(yOffset)

	c.menuWidth = 300
	c.menuHeight = yOffset
}
