package uicontrols

import (
	"image"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ipoluianov/goforms/canvas"
	"github.com/ipoluianov/goforms/ui"
	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiinterfaces"
	"github.com/ipoluianov/goforms/uiresources"
	"golang.org/x/image/colornames"
)

type ComboBoxItem struct {
	uievents.UserDataContainer
	Text string
}

type ComboBoxEvent struct {
	uievents.Event
	CurrentIndex int
	Item         *ComboBoxItem
}

type ComboBox struct {
	Container
	CurrentItemIndex int
	Items            []*ComboBoxItem
	txtBlock         *TextBlock
	popupWidget      uiinterfaces.Widget
	img              image.Image
	enabled_         bool
	pressed          bool

	OnCurrentIndexChanged func(event *ComboBoxEvent)
}

func NewComboBox(parent uiinterfaces.Widget) *ComboBox {
	var c ComboBox
	c.InitControl(parent, &c)
	c.Items = make([]*ComboBoxItem, 0)
	c.enabled_ = true

	c.txtBlock = NewTextBlock(&c, "")
	c.txtBlock.SetBorderRight(0, colornames.White)
	c.AddWidgetOnGrid(c.txtBlock, 0, 0)

	c.cellPadding = 0
	c.panelPadding = 0

	c.UpdateStyle()

	return &c
}

func (c *ComboBox) UpdateStyle() {
	c.Container.UpdateStyle()
	c.img = uiresources.ResImgCol(uiresources.R_icons_material4_png_navigation_arrow_drop_down_materialicons_48dp_1x_baseline_arrow_drop_down_black_48dp_png, c.ForeColor())
}

func (c *ComboBox) ControlType() string {
	return "ComboBox"
}

func (c *ComboBox) Draw(ctx ui.DrawContext) {
	c.Container.Draw(ctx)

	if c.img != nil {
		img := canvas.AdjustImageForColor(c.img, c.Height(), c.Height(), c.ForeColor())
		ctx.DrawImage(c.Width()-c.Height(), 0, 10, 10, img)
	}
}

func (c *ComboBox) MouseDown(event *uievents.MouseDownEvent) {
	c.pressed = true
	c.Update("ComboBox")
}

func (c *ComboBox) MouseUp(event *uievents.MouseUpEvent) {
	if c.enabled_ && c.pressed {
		c.pressed = false
		c.Update("ComboBox")
		c.Press()
		event.Ignore = true
	}
}

func (c *ComboBox) Press() {
	c.ShowPopupForm()
}

func (c *ComboBox) AddItem(text string, userData interface{}) {
	var item ComboBoxItem
	item.Text = text
	item.SetUserData("key", userData)
	c.Items = append(c.Items, &item)
}

func (c *ComboBox) SetCurrentItemIndex(index int) {
	if index >= 0 && index < len(c.Items) {
		c.CurrentItemIndex = index
		c.txtBlock.SetText(c.Items[index].Text)
	}
	if c.OnCurrentIndexChanged != nil {
		var ev ComboBoxEvent
		ev.InitDataContainer()
		ev.Sender = c
		ev.CurrentIndex = c.CurrentItemIndex
		ev.Item = c.Items[c.CurrentItemIndex]
		c.OnCurrentIndexChanged(&ev)
	}
}

func (c *ComboBox) SetCurrentItemKey(key string) {
	for index, item := range c.Items {
		if item.UserData("key") == key {
			c.SetCurrentItemIndex(index)
			break
		}
	}
}

func (c *ComboBox) CurrentItemKey() interface{} {
	if c.CurrentItemIndex >= 0 && c.CurrentItemIndex < len(c.Items) {
		return c.Items[c.CurrentItemIndex].UserData("key")
	}
	return nil
}

func (c *ComboBox) TabStop() bool {
	return true
}

func (c *ComboBox) ClosePopup() {
}

func (c *ComboBox) KeyDown(event *uievents.KeyDownEvent) bool {
	//fmt.Println(event)
	if event.Key == glfw.KeySpace {
		c.SetCurrentItemIndex(0)
		return true
	}
	if event.Key == glfw.KeyDown && event.Modifiers.Alt {
		c.ShowPopupForm()
		return true
	}
	if event.Key == glfw.KeyDown && !event.Modifiers.Alt {
		c.SetCurrentItemIndex(c.CurrentItemIndex + 1)
		return true
	}
	if event.Key == glfw.KeyUp && !event.Modifiers.Alt {
		c.SetCurrentItemIndex(c.CurrentItemIndex - 1)
		return true
	}
	if event.Key == glfw.KeyEscape {
		c.ClosePopup()
		return true
	}
	return false
}

func (c *ComboBox) ShowPopupForm() {
	container := NewPanel(nil)
	container.SetWindow(c.Window())
	container.cellPadding = 0
	container.panelPadding = 0
	wX, wY := c.RectClientAreaOnWindow()
	container.SetPos(wX, wY+c.Height())

	lvItems := NewListView(container)
	container.AddWidgetOnGrid(lvItems, 0, 0)
	container.horizontalScrollVisible.SetOwnValue(false)

	height := lvItems.itemHeight*len(c.Items) + 5
	if height > 300 {
		height = 300
	}

	container.SetSize(c.ClientWidth(), height)
	c.Window().AppendPopup(container)
	lvItems.AddColumn("", 130)
	lvItems.content.horizontalScrollVisible.SetOwnValue(false)

	for _, item := range c.Items {
		lvItem := lvItems.AddItem(item.Text)
		lvItem.TempData = item.TempData
	}

	c.popupWidget = container

	lvItems.OnItemClicked = func(item *ListViewItem) {
		c.SetCurrentItemIndex(item.row)
		c.Window().CloseTopPopup()
	}

	lvItems.SetHeaderVisible(false)
}

func (c *ComboBox) FindWidgetUnderPointer(x, y int) uiinterfaces.Widget {
	return nil
}
