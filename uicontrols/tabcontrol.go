package uicontrols

import (
	"image"
	"image/color"

	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiinterfaces"
	"github.com/ipoluianov/goforms/uiresources"
)

type TabControl struct {
	Container

	headerButtons         []*Button
	pages                 []*TabPage
	visiblePages          []int
	currentPageIndex      int
	closeButtonIndexHover int

	hoverTabIndex int

	header        *Panel
	pagesPanel    *Panel
	btnAdd        *Button
	headerHSpacer *HSpacer

	showAddButton bool

	OnPageSelected     func(index int)
	OnNeedClose        func(index int)
	OnAddButtonPressed func()
}

type TabPage struct {
	Panel
	ShowImage       bool
	ShowText        bool
	Img             image.Image
	text            string
	ShowCloseButton bool
	tabControl      *TabControl
	headerButton    *Button
}

func NewTabControl(parent uiinterfaces.Widget) *TabControl {
	var c TabControl
	c.InitControl(parent, &c)
	c.pages = make([]*TabPage, 0)
	c.SetPanelPadding(0)
	c.SetCellPadding(0)

	c.header = NewPanel(&c)
	c.AddWidgetOnGrid(c.header, 0, 0)
	c.header.SetCellPadding(0)
	c.header.SetPanelPadding(0)
	//c.header.AddHSpacerOnGrid(100, 0)

	c.pagesPanel = NewPanel(&c)
	c.AddWidgetOnGrid(c.pagesPanel, 0, 1)
	c.pagesPanel.SetCellPadding(0)
	c.pagesPanel.SetPanelPadding(0)
	c.pagesPanel.SetBorders(1, c.ForeColor())

	c.visiblePages = make([]int, 0)

	c.xExpandable = true
	c.yExpandable = true

	return &c
}

func (c *TabPage) Dispose() {
	c.Panel.Dispose()
}

func (c *TabPage) ControlType() string {
	return "TabPage"
}

func (c *TabPage) SetText(text string) {
	c.text = text
	c.headerButton.SetText(text)
	c.Update("TabPage")
}

func (c *TabPage) SetVisible(visible bool) {
	c.headerButton.SetVisible(visible)
	c.Panel.SetVisible(visible)
}

func (c *TabPage) MouseMove(event *uievents.MouseMoveEvent) {
	c.Panel.MouseMove(event)
}

func (c *TabControl) Dispose() {
	for _, p := range c.pages {
		p.Dispose()
	}

	c.Control.Dispose()

	c.pages = nil
	c.OnNeedClose = nil
}

func (c *TabControl) SetShowAddButton(showAddButton bool) {
	c.showAddButton = showAddButton
	c.updateHeaderButtons()
}

func (c *TabControl) ControlType() string {
	return "TabControl"
}

func (c *TabControl) AddPage() *TabPage {
	var t TabPage
	t.InitControl(c, &t)
	t.SetWindow(c.OwnWindow)

	pageIndex := len(c.pages)
	c.pagesPanel.AddWidgetOnGrid(&t, pageIndex, 0)
	c.pages = append(c.pages, &t)
	t.ShowText = true
	t.ShowImage = false
	t.tabControl = c

	if len(c.pages) == 1 {
		c.SetCurrentPage(0)
	} else {
		c.SetCurrentPage(pageIndex)
	}

	t.SetYExpandable(true)

	return &t
}

func (c *TabControl) RemovePage(index int) {
	if index >= 0 && index < len(c.pages) {
		c.pagesPanel.RemoveWidget(c.pages[index].widget)
		c.pages[index].Dispose()
		c.pages = append(c.pages[:index], c.pages[index+1:]...)

		if c.currentPageIndex >= len(c.pages) {
			c.SetCurrentPage(len(c.pages) - 1)
		} else {
			c.SetCurrentPage(c.currentPageIndex)
		}

		c.updateHeaderButtons()
		c.Update("TabControl")
	}
}

func (c *TabControl) Page(index int) *TabPage {
	return c.pages[index]
}

func (c *TabControl) SetCurrentPage(index int) {
	if index >= 0 && index < len(c.pages) {
		for i, page := range c.pages {
			if index == i {
				page.Panel.SetVisible(true)
			} else {
				page.Panel.SetVisible(false)
			}
		}
		c.currentPageIndex = index
		c.Update("TabControl")
		if c.OnPageSelected != nil {
			c.OnPageSelected(index)
		}
	}
	c.updateHeaderButtons()
}

func (c *TabControl) updateHeaderButtons() {
	c.header.RemoveAllWidgets()
	c.headerButtons = make([]*Button, 0)
	for pageIndex, page := range c.pages {
		btnPanel := c.header.AddPanelOnGrid(pageIndex, 0)
		btnPanel.SetPanelPadding(0)
		btnPanel.SetCellPadding(0)
		btn := btnPanel.AddButtonOnGrid(0, 0, "TabName", func(event *uievents.Event) {
			c.SetCurrentPage(event.Sender.(*Button).UserData("index").(int))
		})
		btn.SetUserData("index", pageIndex)
		btn.SetMinWidth(150)
		c.headerButtons = append(c.headerButtons, btn)
		page.headerButton = btn
		btn.SetText(page.text)
		btn.SetBorderRight(0, color.RGBA{})

		btnClose := btnPanel.AddButtonOnGrid(1, 0, "", func(event *uievents.Event) {
			if c.OnNeedClose != nil {
				c.OnNeedClose(event.Sender.(*Button).UserData("index").(int))
			}
		})
		btnClose.SetUserData("index", pageIndex)
		btnClose.SetBorderLeft(0, color.RGBA{})
		btnClose.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_close_materialiconsoutlined_48dp_1x_outline_close_black_48dp_png, c.ForeColor()))
		btnClose.SetImageSize(16, 16)

		if pageIndex == c.currentPageIndex {
			/*btn.SetForeColor(c.BackColor())
			btn.SetBackColor(c.ForeColor())
			btnClose.SetForeColor(c.BackColor())
			btnClose.SetBackColor(c.ForeColor())*/
			btn.parent.SetBorderBottom(5, c.AccentColor())
		} else {
			/*btn.SetForeColor(nil)
			btn.SetBackColor(nil)
			btnClose.SetForeColor(nil)
			btnClose.SetBackColor(nil)*/
			btn.parent.SetBorderBottom(5, c.BackColor())
		}
	}

	pageIndex := len(c.pages)

	if c.showAddButton {
		c.btnAdd = c.header.AddButtonOnGrid(pageIndex+1, 0, "", func(event *uievents.Event) {
			if c.OnAddButtonPressed != nil {
				c.OnAddButtonPressed()
			}
		})
		c.btnAdd.SetBorderBottom(5, c.BackColor())
		c.btnAdd.SetImage(uiresources.ResImgCol(uiresources.R_icons_material4_png_content_add_materialicons_48dp_1x_baseline_add_black_48dp_png, c.ForeColor()))
		c.btnAdd.SetImageSize(32, 16)
		c.headerHSpacer = c.header.AddHSpacerOnGrid(pageIndex+2, 0)
	} else {
		c.headerHSpacer = c.header.AddHSpacerOnGrid(pageIndex+1, 0)
	}

	c.Update("TabControl")
}

func (c *TabControl) PagesCount() int {
	return len(c.pages)
}

func (c *TabControl) Tooltip() string {
	if c.hoverTabIndex > -1 {
		if c.hoverTabIndex < len(c.pages) {
			return c.pages[c.hoverTabIndex].text
		}
	}
	return ""
}
