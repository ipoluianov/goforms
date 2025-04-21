package ui

import (
	"image/color"

	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/goforms/utils/uiproperties"
	"github.com/ipoluianov/nui/nuikey"
)

type ListView struct {
	Container

	items               []*ListViewItem
	displayedItems      []*displayedItem
	currentItem         *ListViewItem
	currentItemColumn   int
	lastClickedRowIndex int
	columns             []*ListViewColumn

	selectionType int

	itemHeight int
	cache      *ImageCache

	header        *ListViewHeader
	content       *ListViewContent
	headerVisible bool

	popupLineEdit *PopupLineEdit

	contentPadding int

	OnItemClicked      func(item *ListViewItem)
	OnSelectionChanged func()

	selectionBackground *uiproperties.Property
	selectionForeground *uiproperties.Property
	gridColor           *uiproperties.Property

	showing         bool
	showingProgress float64
	showingTime     *FormTimer

	AllowDeselectItems bool

	OnMouseDown func()
	OnMouseUp   func()
}

type ListViewColumn struct {
	width  int
	text   string
	hAlign canvas.HAlign
}

func NewListView(parent Widget) *ListView {
	var c ListView
	c.InitControl(parent, &c)
	c.Construct()
	return &c
}

func (c *ListView) Construct() {
	c.AllowDeselectItems = true
	c.items = make([]*ListViewItem, 0)
	c.selectionBackground = AddPropertyToWidget(c, "selectionBackground", uiproperties.PropertyTypeColor)
	c.selectionForeground = AddPropertyToWidget(c, "selectionForeground", uiproperties.PropertyTypeColor)
	c.gridColor = AddPropertyToWidget(c, "gridColor", uiproperties.PropertyTypeColor)
	c.cellPadding = 0
	c.panelPadding = 0
	c.cache = NewImageCache("ListView")
	c.itemHeight = 25
	c.contentPadding = 3
	c.headerVisible = true

	c.selectionType = 1

	c.header = NewListViewHeader(c, 0, 0, c.Width(), c.itemHeight)
	c.header.listView = c
	c.header.SetVisible(c.headerVisible)
	c.AddWidgetOnGrid(c.header, 0, 0)

	c.content = newListViewContent(c, 0, c.itemHeight, c.Width(), c.Height()-c.itemHeight)
	c.content.listView = c
	c.content.SetXExpandable(true)
	c.content.SetYExpandable(true)
	c.AddWidgetOnGrid(c.content, 0, 1)

	c.content.onScrolled = func(hScrollPos int, vScrollPos int) {
		c.header.scrollOffsetX = hScrollPos
		//c.header.scrollOffsetY = vScrollPos
	}

	c.content.horizontalScrollVisible.SetOwnValue(true)
	c.content.verticalScrollVisible.SetOwnValue(true)

	c.Window().UpdateLayout()

	c.showing = false

	c.showingTime = c.Window().NewTimer(10, func() {
		if !c.showing {
			c.showingTime.StopTimer()
			//c.showingTime = nil
		}
		c.showingProgress += 0.01
		if c.showingProgress > 1 {
			c.showing = false
		}
	})
	InitDefaultStyle(c)
	//c.showingTime.StartTimer()

}

func (c *ListView) Dispose() {
	c.Container.Dispose()
	c.items = nil
	c.displayedItems = nil
	c.currentItem = nil
	c.cache.Clear()
	c.header.listView = nil
	c.header = nil
	c.content.listView = nil
	c.content = nil
	c.OnSelectionChanged = nil
	c.selectionBackground.Dispose()
	c.selectionBackground = nil
	c.selectionForeground.Dispose()
	c.selectionForeground = nil
}

func (c *ListView) Draw(ctx DrawContext) {
	if !c.showing {
		c.Container.Draw(ctx)
		return
	}

	ctx.SetColor(c.ForeColor())
	ctx.SetStrokeWidth(1)
	ctx.DrawLine(int(c.showingProgress*float64(c.Width())), 0, int(c.showingProgress*float64(c.Width())), 100)
}

func (c *ListView) UpdateStyle() {
	c.cache.Clear()
	c.Container.UpdateStyle()
}

func (c *ListView) IsHeaderVisible() bool {
	return c.header.IsVisible()
}

func (c *ListView) SetHeaderVisible(visible bool) {
	c.header.SetVisible(visible)
	c.Window().UpdateLayout()
}

func (c *ListView) ControlType() string {
	return "ListView"
}

func (c *ListView) EnabledChanged(enabled bool) {
	c.header.SetEnabled(enabled)
	c.content.SetEnabled(enabled)
	c.cache.Clear()
	c.Update("ListView")
}

func (c *ListView) EnsureVisibleItem(index int) {
	if index >= 0 && index < len(c.items) {
		c.content.ScrollEnsureVisible(0, index*c.itemHeight)
		c.content.ScrollEnsureVisible(0, index*c.itemHeight+c.itemHeight-1)
	}
	c.Update("ListView")
}

func (c *ListView) ItemsCount() int {
	return len(c.items)
}

func (c *ListView) TabStop() bool {
	return true
}

func (c *ListView) Focus() {
	c.content.Focus()
}

func (c *ListView) AddItem(text string) *ListViewItem {
	var item ListViewItem
	item.row = len(c.items)
	item.values = make(map[int]string)
	item.values[0] = text
	item.listView = c
	item.foreColors = make(map[int]color.Color)
	c.items = append(c.items, &item)

	c.content.updateInnerSize()

	c.Update("ListView")
	return &item
}

func (c *ListView) AddItem2(col0 string, col1 string) *ListViewItem {
	item := c.AddItem(col0)
	item.SetValue(1, col1)
	return item
}

func (c *ListView) AddItem3(col0 string, col1 string, col2 string) *ListViewItem {
	item := c.AddItem(col0)
	item.SetValue(1, col1)
	item.SetValue(2, col2)
	return item
}

func (c *ListView) RemoveItems() {
	c.UnselectAllItems()

	c.items = make([]*ListViewItem, 0)
	c.cache.Clear()
	c.content.updateInnerSize()
	c.content.ScrollEnsureVisible(0, 0)
	c.Update("ListView")
}

func (c *ListView) SelectItem(rowIndex int, column int) {
	c.SetCurrentRow(rowIndex, column, false)
}

func (c *ListView) SetColumnWidth(colIndex int, width int) {
	if colIndex >= 0 && colIndex < len(c.columns) {
		c.columns[colIndex].width = width
	}
	c.Update("ListView")
}

func (c *ListView) SelectItemSelection(rowIndex int, selected bool) {
	c.items[rowIndex].selected = selected
	c.cache.Clear()
	c.Update("ListView")
}

func (c *ListView) RemoveColumns() {
	c.columns = make([]*ListViewColumn, 0)
	c.cache.Clear()
	c.content.updateInnerSize()
	c.content.ScrollEnsureVisible(0, 0)
	c.Update("ListView")
}

func (c *ListView) AddColumn(text string, width int) *ListViewColumn {
	var listViewColumn ListViewColumn
	listViewColumn.text = text
	listViewColumn.width = width
	listViewColumn.hAlign = canvas.HAlignLeft

	c.columns = append(c.columns, &listViewColumn)
	c.content.updateInnerSize()
	c.Update("ListView")

	return &listViewColumn
}

func (c *ListView) SetColumnTextAlign(columnIndex int, hAlign canvas.HAlign) {
	if columnIndex < len(c.columns) {
		c.columns[columnIndex].hAlign = hAlign
		c.cache.Clear()
		c.Update("ListView")
	}
}

func (c *ListView) Item(rowIndex int) *ListViewItem {
	if rowIndex < 0 || rowIndex >= len(c.items) {
		return nil
	}
	return c.items[rowIndex]
}

func (c *ListView) SetItemValue(rowIndex int, columnIndex int, text string) {
	if rowIndex < 0 || rowIndex >= len(c.items) {
		return
	}
	item := c.items[rowIndex]
	if item.values[columnIndex] != text {
		item.values[columnIndex] = text
		c.cache.ClearXY(columnIndex, item.row)
		c.Update("ListView")
	}
}

func (c *ListView) calcColumnXOffset(columnIndex int) int {
	columnOffset := 0
	for index, column := range c.columns {
		if index == columnIndex {
			break
		}
		columnOffset += column.width
	}
	return columnOffset
}

func (c *ListView) SelectAllItems() {
	hasChanges := false

	for _, item := range c.items {
		if !item.selected {
			item.selected = true
			hasChanges = true
		}
	}

	if hasChanges {
		c.cache.Clear()
		c.Update("ListView")
	}

	if hasChanges && c.OnSelectionChanged != nil {
		c.OnSelectionChanged()
	}
}

func (c *ListView) UnselectAllItems() {
	hasChanges := false

	for _, item := range c.items {
		if item.selected {
			item.selected = false
			hasChanges = true
		}
	}

	c.currentItem = nil

	if hasChanges {
		c.cache.Clear()
		c.Update("ListView")
	}

	if hasChanges && c.OnSelectionChanged != nil {
		c.OnSelectionChanged()
	}
}

func (c *ListView) SetCurrentRow(row int, column int, byMouse bool) {
	if row < 0 || row > len(c.items) {
		return
	}

	if column < 0 || column >= len(c.columns) {
		return
	}

	if c.currentItem != nil {
		c.removeCacheForRow(c.currentItem.row)
	}
	c.removeCacheForRow(row)
	c.currentItem = c.items[row]
	c.currentItemColumn = column

	needToClearSelection := false

	if byMouse {
		if c.Window().KeyModifiers().Shift {
			needToClearSelection = true
		} else {
			if c.Window().KeyModifiers().Ctrl {
				needToClearSelection = false
				c.lastClickedRowIndex = row
			} else {
				c.lastClickedRowIndex = row
				needToClearSelection = true
			}
		}

		if needToClearSelection {
			for _, item := range c.items {
				item.selected = false
			}
		}

		if c.Window().KeyModifiers().Ctrl {
			c.currentItem.selected = !c.currentItem.selected
		} else {
			c.currentItem.selected = true
		}

		if c.Window().KeyModifiers().Shift {
			selectionFrom := c.lastClickedRowIndex
			selectionTo := row
			if selectionTo < selectionFrom {
				selectionFrom, selectionTo = selectionTo, selectionFrom
			}
			for i := selectionFrom; i <= selectionTo; i++ {
				c.items[i].selected = true
			}
		}
	} else {

		if c.Window().KeyModifiers().Shift {
			needToClearSelection = true
		} else {
			if c.Window().KeyModifiers().Ctrl {
				needToClearSelection = true
				c.lastClickedRowIndex = row
			} else {
				needToClearSelection = true
				c.lastClickedRowIndex = row
			}
		}
		if needToClearSelection {
			for _, item := range c.items {
				item.selected = false
			}
		}

		if c.Window().KeyModifiers().Shift {
			selectionFrom := c.lastClickedRowIndex
			selectionTo := row
			if selectionTo < selectionFrom {
				selectionFrom, selectionTo = selectionTo, selectionFrom
			}
			for i := selectionFrom; i <= selectionTo; i++ {
				if i >= 0 && i < len(c.items) {
					c.items[i].selected = true
				}
			}
		} else {
			c.currentItem.selected = true
		}
	}

	c.cache.Clear()

	if c.OnSelectionChanged != nil {
		c.OnSelectionChanged()
	}
}

func (c *ListViewContent) KeyDown(event *KeyDownEvent) bool {
	selectedItemDisplayIndex := -1

	for index, dItem := range c.listView.displayedItems {
		if c.listView.currentItem == dItem.item {
			selectedItemDisplayIndex = index
			break
		}
	}

	if event.Key == nuikey.KeyA && event.Modifiers.Ctrl {
		c.listView.SelectAllItems()
		return true
	}

	if event.Key == nuikey.KeyEnter {
		c.EditCurrentCell("")
		return true
	}

	if event.Key == nuikey.KeyArrowUp {
		if selectedItemDisplayIndex > 0 {
			c.listView.SetCurrentRow(c.listView.displayedItems[selectedItemDisplayIndex-1].item.row, c.listView.currentItemColumn, false)
			c.listView.EnsureVisibleItem(c.listView.currentItem.row)
		}
		return true
	}

	if event.Key == nuikey.KeyArrowDown {
		if selectedItemDisplayIndex < len(c.listView.displayedItems)-1 {
			c.listView.SetCurrentRow(c.listView.displayedItems[selectedItemDisplayIndex+1].item.row, c.listView.currentItemColumn, false)
			c.listView.EnsureVisibleItem(c.listView.currentItem.row)
		}
		return true
	}

	if event.Key == nuikey.KeyArrowLeft {
		col := c.listView.currentItemColumn - 1
		if col < 0 {
			col = 0
		}
		c.listView.SetCurrentRow(c.listView.displayedItems[selectedItemDisplayIndex].item.row, col, false)
		c.listView.EnsureVisibleItem(c.listView.currentItem.row)
		return true
	}

	if event.Key == nuikey.KeyArrowRight {
		col := c.listView.currentItemColumn + 1
		if col >= len(c.listView.columns) {
			col = len(c.listView.columns) - 1
		}
		c.listView.SetCurrentRow(c.listView.displayedItems[selectedItemDisplayIndex].item.row, col, false)
		c.listView.EnsureVisibleItem(c.listView.currentItem.row)
		return true
	}

	if event.Key == nuikey.KeyHome {
		if len(c.listView.displayedItems) > 0 {
			c.listView.SetCurrentRow(0, c.listView.currentItemColumn, false)
			c.listView.EnsureVisibleItem(0)
		}
		return true
	}

	if event.Key == nuikey.KeyEnd {
		if len(c.listView.displayedItems) > 0 {
			c.listView.SetCurrentRow(len(c.listView.items)-1, c.listView.currentItemColumn, false)
			c.listView.EnsureVisibleItem(len(c.listView.items) - 1)
		}
		return true
	}

	if event.Key == nuikey.KeyPageUp {
		if len(c.listView.displayedItems) > 0 {
			if c.listView.currentItem != nil {
				row := c.listView.currentItem.row
				row -= c.Height() / c.listView.itemHeight
				if row < 0 {
					row = 0
				}
				c.listView.SetCurrentRow(row, c.listView.currentItemColumn, false)
				c.listView.EnsureVisibleItem(row)
			}
		}
		return true
	}

	if event.Key == nuikey.KeyPageDown {
		if len(c.listView.displayedItems) > 0 {
			if c.listView.currentItem != nil {
				row := c.listView.currentItem.row
				row += c.Height() / c.listView.itemHeight
				if row >= len(c.listView.items) {
					row = len(c.listView.items) - 1
				}
				c.listView.SetCurrentRow(row, c.listView.currentItemColumn, false)
				c.listView.EnsureVisibleItem(row)
			}
		}
		return true
	}

	return false
}

func (c *ListView) SelectedItem() *ListViewItem {
	return c.currentItem
}

func (c *ListView) SelectedItemIndex() int {
	if c.currentItem != nil {
		return c.currentItem.row
	}
	return -1
}

func (c *ListView) VisibleItems() []*ListViewItem {
	return []*ListViewItem{}
}

func (c *ListView) removeCacheForRow(row int) {
	for index := range c.columns {
		c.cache.ClearXY(index, row)
	}
}

func (c *ListView) findDisplayItemByCoordinates(x int, y int) *displayedItem {
	for _, dItem := range c.displayedItems {
		if x >= 0 && x < c.InnerWidth() && y >= dItem.currentY && y < dItem.currentY+dItem.currentHeight {
			return dItem
		}
	}
	return nil
}

func (c *ListView) findDisplayColumnByCoordinates(x int) (colIndex int) {
	for index, column := range c.columns {
		colRightBorder := c.calcColumnXOffset(index) + column.width
		if x >= c.calcColumnXOffset(index) && x < colRightBorder {
			return index
		}
	}
	return -1
}

func (c *ListView) updateItemHeight() {
	_, fontHeight, _ := canvas.MeasureText(c.FontFamily(), c.FontSize(), false, false, "1Qg", false)
	fontHeight += 2
	if c.itemHeight != fontHeight {
		c.itemHeight = fontHeight
		c.UpdateLayout()
	}
}

func (c *ListView) OnInit() {
}

func (c *ListView) ClearSelection() {
	selectedFound := false
	for _, item := range c.items {
		if item.selected {
			item.selected = false
			selectedFound = true
		}
	}

	if selectedFound {
		c.currentItem = nil
		if c.OnSelectionChanged != nil {
			c.OnSelectionChanged()
		}
		c.cache.Clear()
	}

	c.Update("ListView")
}

func (c *ListView) SelectedItems() []*ListViewItem {
	items := make([]*ListViewItem, 0)
	for _, item := range c.items {
		if item.selected {
			items = append(items, item)
		}
	}
	return items
}

func (c *ListView) SelectedItemsIndexes() []int {
	indexes := make([]int, 0)
	for index, item := range c.items {
		if item.selected {
			indexes = append(indexes, index)
		}
	}
	return indexes
}
