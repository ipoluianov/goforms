package ui

import (
	"image/color"
	"math"

	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/goforms/utils/uiproperties"
	"github.com/ipoluianov/nui/nuikey"
	"github.com/ipoluianov/nui/nuimouse"
)

type ListView struct {
	Container

	items               []*ListViewItem
	displayedItems      []*displayedItem
	currentItem         *ListViewItem
	lastClickedRowIndex int
	columns             []*ListViewColumn

	itemHeight int
	cache      *ImageCache

	header        *ListViewHeader
	content       *ListViewContent
	headerVisible bool

	contentPadding int

	OnItemClicked      func(item *ListViewItem)
	OnSelectionChanged func()

	selectionBackground *uiproperties.Property
	gridColor           *uiproperties.Property

	showing         bool
	showingProgress float64
	showingTime     *FormTimer

	AllowDeselectItems bool

	OnMouseDown func()
	OnMouseUp   func()
}

type ListViewHeader struct {
	Control

	pressed             bool
	columnResizing      bool
	columnResizingIndex int

	listView *ListView
}

type ListViewContent struct {
	Control

	listView *ListView
}

type ListViewColumn struct {
	width  int
	text   string
	hAlign canvas.HAlign
}

type displayedItem struct {
	currentX      int
	currentY      int
	currentWidth  int
	currentHeight int
	item          *ListViewItem
}

type ListViewItem struct {
	UserDataContainer
	row        int
	values     map[int]string
	selected   bool
	listview   *ListView
	foreColors map[int]color.Color
	//UserData interface{}
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
	c.gridColor = AddPropertyToWidget(c, "gridColor", uiproperties.PropertyTypeColor)
	c.cellPadding = 0
	c.panelPadding = 0
	c.cache = NewImageCache("ListView")
	c.itemHeight = 25
	c.contentPadding = 3
	c.headerVisible = true

	c.header = NewListViewHeader(c, 0, 0, c.Width(), c.itemHeight)
	c.header.listView = c
	c.header.SetVisible(c.headerVisible)
	c.AddWidgetOnGrid(c.header, 0, 0)

	c.content = NewListViewContent(c, 0, c.itemHeight, c.Width(), c.Height()-c.itemHeight)
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
	//c.showingTime.StartTimer()

}

func NewListViewContent(parent Widget, x int, y int, width int, height int) *ListViewContent {
	var c ListViewContent
	c.InitControl(parent, &c)
	return &c
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
	item.listview = c
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

func (c *ListView) SelectItem(rowIndex int) {
	c.SetCurrentRow(rowIndex, false)
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

func (c *ListView) SetCurrentRow(row int, byMouse bool) {
	if row < 0 || row > len(c.items) {
		return
	}

	if c.currentItem != nil {
		c.removeCacheForRow(c.currentItem.row)
	}
	c.removeCacheForRow(row)
	c.currentItem = c.items[row]
	c.currentItem = c.items[row]

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

	if event.Key == nuikey.KeyEsc {
		c.listView.showing = true
		c.listView.showingProgress = 0
		c.listView.showingTime.StartTimer()
		return true
	}

	if event.Key == nuikey.KeyArrowUp {
		if selectedItemDisplayIndex > 0 {
			c.listView.SetCurrentRow(c.listView.displayedItems[selectedItemDisplayIndex-1].item.row, false)
			c.listView.EnsureVisibleItem(c.listView.currentItem.row)
		}
		return true
	}

	if event.Key == nuikey.KeyArrowDown {
		if selectedItemDisplayIndex < len(c.listView.displayedItems)-1 {
			c.listView.SetCurrentRow(c.listView.displayedItems[selectedItemDisplayIndex+1].item.row, false)
			c.listView.EnsureVisibleItem(c.listView.currentItem.row)
		}
		return true
	}

	if event.Key == nuikey.KeyHome {
		if len(c.listView.displayedItems) > 0 {
			c.listView.SetCurrentRow(0, false)
			c.listView.EnsureVisibleItem(0)
		}
		return true
	}

	if event.Key == nuikey.KeyEnd {
		if len(c.listView.displayedItems) > 0 {
			c.listView.SetCurrentRow(len(c.listView.items)-1, false)
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
				c.listView.SetCurrentRow(row, false)
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
				c.listView.SetCurrentRow(row, false)
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

func NewListViewHeader(parent Widget, x int, y int, width int, height int) *ListViewHeader {
	var c ListViewHeader
	c.InitControl(parent, &c)
	return &c
}

func (c *ListViewHeader) MouseDown(event *MouseDownEvent) {
	c.pressed = true
	c.Update("ListView")

	colRightBorder := c.findColumnRightBorder(event.X, event.Y)
	if colRightBorder >= 0 {
		c.columnResizing = true
		c.columnResizingIndex = colRightBorder
	}

	if c.listView.OnMouseDown != nil {
		c.listView.OnMouseDown()
	}
}

func (c *ListViewHeader) MouseUp(event *MouseUpEvent) {
	if c.pressed {
		c.pressed = false
		c.columnResizing = false
		c.Update("ListView")
	}

	if c.listView.OnMouseUp != nil {
		c.listView.OnMouseUp()
	}
}

func (c *ListViewHeader) MouseMove(event *MouseMoveEvent) {
	//fmt.Println("ListView Move: ", event.X, " ", event.Y)
	if c.columnResizing {
		//fmt.Println("Move: event.Y:", event.X, " c.calcColumnXOffset(c.columnResizingIndex):", c.listView.calcColumnXOffset(c.columnResizingIndex), " w: ", event.Y-c.listView.calcColumnXOffset(c.columnResizingIndex))
		c.listView.columns[c.columnResizingIndex].width = event.X - c.listView.calcColumnXOffset(c.columnResizingIndex)
		if c.listView.columns[c.columnResizingIndex].width < 10 {
			c.listView.columns[c.columnResizingIndex].width = 10
		}
		c.listView.cache.Clear()
	}

	if event.Y < c.listView.itemHeight {
		if c.findColumnRightBorder(event.X, event.Y) >= 0 {
			c.Window().SetMouseCursor(nuimouse.MouseCursorResizeHor)
		} else {
			c.Window().SetMouseCursor(nuimouse.MouseCursorArrow)
		}
	}
}

func (c *ListViewHeader) findColumnRightBorder(x, y int) int {
	for index, column := range c.listView.columns {
		colRightBorder := c.listView.calcColumnXOffset(index) + column.width
		if math.Abs(float64(x-colRightBorder)) < 5 {
			return index
		}
	}
	return -1
}

func (c *ListViewHeader) MouseLeave() {
	c.Window().SetMouseCursor(nuimouse.MouseCursorArrow)
}

func (c *ListView) updateItemHeight() {
	_, fontHeight, _ := canvas.MeasureText(c.FontFamily(), c.FontSize(), false, false, "1Qg", false)
	fontHeight += 2
	if c.itemHeight != fontHeight {
		c.itemHeight = fontHeight
		c.UpdateLayout()
	}
}

func (c *ListViewHeader) Draw(ctx DrawContext) {
	c.listView.updateItemHeight()

	xOffset := 0
	yOffset := c.listView.contentPadding

	ctx.SetColor(c.ForeColor())
	ctx.SetStrokeWidth(1)
	ctx.SetFontSize(c.FontSize())

	for colIndex, column := range c.listView.columns {
		/*if index > 0 {
			ctx.DrawLine(xOffset, 0, xOffset, c.listView.itemHeight + c.listView.contentPadding * 2)
		}*/

		var cnv *canvas.CanvasDirect
		cnv = c.listView.cache.GetXY(colIndex, -1)
		if cnv == nil {
			cnv = canvas.NewCanvas(column.width, c.listView.itemHeight)
			cnv.DrawText(0, 0, column.text, c.FontFamily(), c.FontSize(), c.ForeColor(), false)
			c.listView.cache.SetXY(colIndex, -1, cnv)
		}

		clipX := xOffset + c.listView.contentPadding
		//clipY := yOffset
		clipW := column.width - c.listView.contentPadding*2
		//clipH := c.listView.itemHeight

		effectiveWidth := c.Width()

		if clipX+clipW > effectiveWidth {
			clipW = effectiveWidth - clipX
		}
		if clipW < 0 {
			clipW = 0
		}
		if clipX < 0 {
			clipW += clipX
			clipX = 0
		}

		//ctx.Save()
		//ctx.Clip(clipX, clipY, clipW, clipH)
		ctx.DrawImage(xOffset+c.listView.contentPadding, yOffset, column.width, c.listView.itemHeight, cnv.Image())
		//ctx.Load()
		ctx.SetColor(c.listView.leftBorderColor.Color())
		ctx.DrawLine(xOffset+column.width, 0, xOffset+column.width, c.listView.itemHeight+c.listView.contentPadding*2)
		xOffset += column.width

		/*ctx.Save()
		ctx.Clip(xOffset+c.listView.contentPadding, yOffset, column.width-c.listView.contentPadding*2, c.listView.itemHeight)
		ctx.DrawText(xOffset+c.listView.contentPadding, yOffset, column.width, c.listView.itemHeight, column.text)
		ctx.Load()
		ctx.DrawLine(xOffset+column.width, 0, xOffset+column.width, c.listView.itemHeight+c.listView.contentPadding*2)
		xOffset += column.width*/
	}

	//ctx.DrawLine(0, c.listView.itemHeight-1, c.InnerWidth(), c.listView.itemHeight-1)*/
}

func (c *ListViewHeader) MinWidth() int {
	return c.Control.MinWidth()
}

func (c *ListViewHeader) MinHeight() int {
	return c.listView.itemHeight + c.listView.contentPadding*2
}

func (c *ListViewHeader) ControlType() string {
	return "ListViewHeader"
}

func (c *ListViewContent) Draw(ctx DrawContext) {
	yOffset := 0
	c.listView.displayedItems = make([]*displayedItem, 0)

	visRect := c.VisibleInnerRect()
	beginIndex := visRect.Y/c.listView.itemHeight - 1
	if beginIndex < 0 {
		beginIndex = 0
	}
	if beginIndex >= len(c.listView.items) {
		beginIndex = len(c.listView.items) - 1
	}
	endIndex := (visRect.Y+visRect.Height)/c.listView.itemHeight + 1
	if endIndex < 0 {
		endIndex = 0
	}
	if endIndex >= len(c.listView.items) {
		endIndex = len(c.listView.items) - 1
	}

	yOffset += beginIndex * c.listView.itemHeight

	if beginIndex >= 0 && beginIndex <= endIndex {
		for index := beginIndex; index <= endIndex; index++ {
			yOffset += c.drawItem(ctx, c.listView.items[index], yOffset, index)
		}
	}

	c.drawGrid(ctx)
	c.updateInnerSize()
}

func (c *ListViewContent) ControlType() string {
	return "ListViewContent"
}

func (c *ListViewContent) drawItem(ctx DrawContext, item *ListViewItem, y int, itemIndex int) int {
	yOffset := 0

	//rowX := 0
	//rowWidth := c.InnerWidth()
	//rowY := y
	rowHeight := c.listView.itemHeight

	visRect := c.VisibleInnerRect()
	/*if rowX+rowWidth < visRect.X {
		yOffset += rowHeight
		return yOffset
	}
	if rowY+rowHeight < visRect.Y {
		yOffset += rowHeight
		return yOffset
	}
	if rowX > visRect.X+visRect.Width {
		yOffset += rowHeight
		return yOffset
	}
	if rowY > visRect.Y+visRect.Height {
		yOffset += rowHeight
		return yOffset
	}*/

	var dItem displayedItem
	dItem.currentX = 0
	dItem.currentY = y
	dItem.currentWidth = 100
	dItem.currentHeight = c.listView.itemHeight
	dItem.item = item
	c.listView.displayedItems = append(c.listView.displayedItems, &dItem)

	xOffset := 0
	for index, column := range c.listView.columns {
		cellX := xOffset
		cellWidth := column.width
		cellY := y
		cellHeight := rowHeight

		if cellX+cellWidth < visRect.X {
			xOffset += column.width
			continue
		}
		if cellY+cellHeight < visRect.Y {
			xOffset += column.width
			continue
		}
		if cellX > visRect.X+visRect.Width {
			xOffset += column.width
			continue
		}
		if cellY > visRect.Y+visRect.Height {
			xOffset += column.width
			continue
		}

		var cnv *canvas.CanvasDirect
		cnv = c.listView.cache.GetXY(index, itemIndex)
		if cnv == nil {
			cnv = canvas.NewCanvas(column.width, rowHeight)

			value := item.values[index]
			backColor := c.BackColor()
			foreColor := c.ForeColor()
			if item.selected {
				backColor = c.listView.selectionBackground.Color()
				//foreColor = c.BackColor()
			}

			cnv.FillRect(0, 0, column.width, rowHeight, backColor)

			/*if c.listView.currentItem == item {
				cnv.DrawRect(0, 0, column.width, rowHeight, c.ForeColor(), 1)
			}*/

			col := foreColor
			if colorForCell, ok := item.foreColors[index]; ok {
				if colorForCell != nil {
					col = colorForCell
				}
			}

			//cnv.DrawText(c.listView.contentPadding, 0, value, c.fontFamily.String(), c.fontSize.Float64(), col)
			cnv.DrawTextMultiline(c.listView.contentPadding, 0, column.width-c.listView.contentPadding*2, rowHeight, column.hAlign, canvas.VAlignCenter, value, col, c.fontFamily.String(), c.fontSize.Float64(), false)
		}

		ctx.DrawImage(cellX, cellY, c.Width(), c.Height(), cnv.Image())

		c.listView.cache.SetXY(index, itemIndex, cnv)

		xOffset += column.width
	}

	yOffset += rowHeight

	return yOffset
}

func (c *ListViewContent) drawGrid(ctx DrawContext) {
	xOffset := 0

	ctx.SetColor(c.listView.gridColor.Color())
	ctx.SetStrokeWidth(1)

	for _, column := range c.listView.columns {
		xOffset += column.width
		ctx.DrawLine(xOffset, c.scrollOffsetY, xOffset, c.Height()+c.scrollOffsetY)
	}
}

func (c *ListViewContent) updateInnerSize() {
	c.InnerSizeOverloaded = true
	_, c.InnerHeightOverloaded = c.calcTreeColumnSize()
	c.updateInnerWidth()

}

func (c *ListViewContent) updateInnerWidth() {
	c.InnerWidthOverloaded = 0
	for _, column := range c.listView.columns {
		c.InnerWidthOverloaded += column.width
	}
}

func (c *ListViewContent) calcTreeColumnSize() (int, int) {
	width := 0
	height := 0

	for _, col := range c.listView.columns {
		width += col.width
	}

	height = len(c.listView.items) * c.listView.itemHeight

	return width, height
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

func (c *ListViewContent) MouseClick(event *MouseClickEvent) {
	dItem := c.listView.findDisplayItemByCoordinates(event.X, event.Y)
	if c.listView.OnMouseDown != nil {
		c.listView.OnMouseDown()
	}

	if dItem == nil {
		if c.listView.AllowDeselectItems {
			c.listView.ClearSelection()
		}
		return
	}

	if event.X > dItem.currentX {
		c.listView.SetCurrentRow(dItem.item.row, true)
		c.ScrollEnsureVisible(dItem.currentX, dItem.currentY)
		c.ScrollEnsureVisible(dItem.currentX, dItem.currentY+c.listView.itemHeight)
		if c.listView.OnItemClicked != nil {
			c.listView.OnItemClicked(dItem.item)
		}
	}

	c.Update("ListView")
}

func (c *ListViewItem) SetValue(column int, text string) {
	for index, item := range c.listview.items {
		if item == c {
			c.listview.SetItemValue(index, column, text)
		}
	}
}

func (c *ListViewItem) SetForeColorForRow(color color.Color) {
	for colIndex := 0; colIndex < len(c.listview.columns); colIndex++ {
		if colorForCell, ok := c.foreColors[colIndex]; ok {
			if colorForCell != color {
				c.foreColors[colIndex] = color
				c.listview.cache.Clear()
				c.listview.Update("ListViewItem")
			}
		} else {
			c.foreColors[colIndex] = color
			c.listview.cache.Clear()
			c.listview.Update("ListViewItem")
		}
	}
}

func (c *ListViewItem) SetForeColorForCell(colIndex int, color color.Color) {
	if colorForCell, ok := c.foreColors[colIndex]; ok {
		if colorForCell != color {
			c.foreColors[colIndex] = color
			c.listview.cache.Clear()
			c.listview.Update("ListViewItem")
		}
	} else {
		c.foreColors[colIndex] = color
		c.listview.cache.Clear()
		c.listview.Update("ListViewItem")
	}
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

func (c *ListViewItem) Value(column int) string {
	if v, ok := c.values[column]; ok {
		return v
	}
	return ""
}
