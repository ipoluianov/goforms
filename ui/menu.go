package ui

import (
	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/goforms/utils/uiproperties"
)

type Menu struct {
	Control
	items           []*MenuItem
	visible         bool
	ownerWindow     Window
	menuOpened      bool
	openedItem      *MenuItem
	hoverItem       *MenuItem
	openedPopupMenu *PopupMenu

	hoverBackgroundColor  *uiproperties.Property
	openedBackgroundColor *uiproperties.Property
}

type MenuItem struct {
	Text    string
	OnClick func(*Event)
	items   []*MenuItem
	posX    int
	posY    int
	width   int
	height  int
	menu    *Menu
}

func NewMenu(parent Widget) *Menu {
	var c Menu
	c.ownerWindow = parent.Window()
	c.InitControl(parent, &c)
	c.hoverBackgroundColor = AddPropertyToWidget(&c, "hoverBackgroundColor", uiproperties.PropertyTypeColor)
	c.openedBackgroundColor = AddPropertyToWidget(&c, "openedBackgroundColor", uiproperties.PropertyTypeColor)
	InitDefaultStyle(&c)
	return &c
}

const MenuBarItemMargin = 10
const MenuBarItemMinWidth = 30

func (c *Menu) ControlType() string {
	return "Menu"
}

func (c *Menu) Draw(ctx DrawContext) {
	//ctx.DrawRect(2, 2, c.Width() - 4, c.Height() - 4, color.RGBA{255,255, 0, 255})
	for _, item := range c.items {
		item.Draw(ctx)
	}
}

func (c *Menu) preferredHeight() int {
	return 22
}

func (c *Menu) IsVisible() bool {
	return c.visible
}

func (c *Menu) SetVisible(visible bool) {
	c.visible = visible
	c.ownerWindow.UpdateMenu()
}

func (c *Menu) AddItem(text string) *MenuItem {
	var item MenuItem
	item.Text = text
	item.menu = c
	c.items = append(c.items, &item)
	c.updateItemsSize()
	return &item
}

func (c *Menu) updateItemsSize() {
	xOffset := 0
	for _, item := range c.items {
		width, _, _ := canvas.MeasureText(c.fontFamily.String(), c.fontSize.Float64(), false, false, item.Text, false)
		width += MenuBarItemMargin * 2
		if width < MenuBarItemMinWidth {
			width = MenuBarItemMinWidth
		}
		item.posX = xOffset
		item.posY = 0
		item.width = width
		item.height = PopupMenuItemHeight
		xOffset += width
	}

}

func (c *Menu) MouseClick(event *MouseClickEvent) {
	item := c.findMenuItemByXOffset(event.X)
	if item == nil {
		return
	}
	c.menuOpened = true
	c.ShowSubmenu(item)
}

func (c *Menu) FocusChanged(focus bool) {
	if !focus {
		c.menuOpened = false
	}
}

func (c *Menu) MouseLeave() {
	c.hoverItem = nil
}

func (c *Menu) MouseMove(event *MouseMoveEvent) {
	/*item := c.findMenuItemByXOffset(event.X)
	if item == nil {
		return
	}
	c.hoverItem = item

	if c.menuOpened {
		if c.openedItem != item {
			if c.openedPopupMenu != nil {
				c.openedItem = item
				wX, wY := c.rectOnWindow()
				wnd := c.window()
				c.remakePopupMenu()
				//c.openedPopupMenu.Move(wnd.Position().X+wX+item.posX, wnd.Position().Y+wY+item.height-2)
			}
		}
	}*/
}

func (c *Menu) onClickItem(ev *Event) {
	//c.OwnWindow.MessageBox("Title", "OnClick")
}

func (c *Menu) ShowSubmenu(item *MenuItem) {
	/*ctxMenu := NewPopupMenu(c.window())

	wX, wY := c.rectOnWindow()
	wnd := c.window()
	ctxMenu.CloseEvent = c.submenuClosed
	c.openedPopupMenu = ctxMenu
	c.openedItem = item
	c.remakePopupMenu()

	c.menuOpened = true
	ctxMenu.ShowMenu(c, wnd.Position().X+wX+item.posX, wnd.Position().Y+wY+item.height-2)*/
}

func (c *Menu) remakePopupMenu() {
	if c.openedPopupMenu != nil {
		c.openedPopupMenu.RemoveAllItems()

		for _, item := range c.openedItem.items {
			c.openedPopupMenu.AddItem(item.Text, nil, nil, "")
		}
	}
}

func (c *Menu) submenuClosed() {
	c.openedPopupMenu = nil
	c.openedItem = nil
	c.menuOpened = false
}

func (c *Menu) findMenuItemByXOffset(xOffset int) *MenuItem {
	for _, item := range c.items {
		if xOffset > item.posX && xOffset < item.posX+item.width {
			return item
		}
	}
	return nil
}

func (c *MenuItem) AddItem(text string, menu *Menu) {
	var item MenuItem
	item.Text = text
	item.menu = menu
	c.items = append(c.items, &item)
}

func (c *MenuItem) Draw(ctx DrawContext) {
	/*if c == c.menu.hoverItem {
		ctx.FillRect(c.posX, c.posY, c.width, c.height, c.menu.hoverBackgroundColor.Color())
	}
	if c.menu.menuOpened {
		if c == c.menu.openedItem {
			ctx.FillRect(c.posX, c.posY, c.width, c.height, c.menu.openedBackgroundColor.Color())
		}
	}
	ctx.DrawTextMultiline(c.posX+MenuBarItemMargin, c.posY, c.width-MenuBarItemMargin, c.height, canvas.HAlignLeft, canvas.VAlignCenter, c.Text, c.menu.foregroundColor.Color(), c.menu.fontFamily.String(), c.menu.fontSize.Float64())
	*/
}
