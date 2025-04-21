package ui

type ListViewContent struct {
	Control

	listView *ListView
}

func newListViewContent(parent Widget, x int, y int, width int, height int) *ListViewContent {
	var c ListViewContent
	c.InitControl(parent, &c)
	return &c
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
			yOffset += c.listView.items[index].draw(ctx, yOffset, index)
		}
	}

	//c.drawGrid(ctx)
	c.updateInnerSize()
}

func (c *ListViewContent) ControlType() string {
	return "ListViewContent"
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

	col := c.listView.findDisplayColumnByCoordinates(event.X)

	if event.X > dItem.currentX {
		c.listView.SetCurrentRow(dItem.item.row, col, true)
		c.ScrollEnsureVisible(dItem.currentX, dItem.currentY)
		c.ScrollEnsureVisible(dItem.currentX, dItem.currentY+c.listView.itemHeight)
		if c.listView.OnItemClicked != nil {
			c.listView.OnItemClicked(dItem.item)
		}
	}

	c.Update("ListView")
}

func (c *ListViewContent) KeyChar(event *KeyCharEvent) {
	str := string(event.Ch)
	c.EditCurrentCell(str)
}

func (c *ListViewContent) EditCurrentCell(enteredText string) {
	c.Window().IgnoreUpdates()

	initText := c.listView.currentRow.Value(c.listView.currentColumn)
	if len(enteredText) > 0 {
		initText = enteredText
	}

	posX, posY := c.RectClientAreaOnWindow()
	posX += c.listView.calcColumnXOffset(c.listView.currentColumn) - c.listView.content.scrollOffsetX
	posY += c.listView.currentRow.row*c.listView.itemHeight - c.listView.content.scrollOffsetY

	columnWidth := c.listView.columns[c.listView.currentColumn].width
	rowHeight := c.listView.itemHeight

	c.listView.popupLineEdit = NewPopupLineEdit(c, initText, len(enteredText) == 0, columnWidth, rowHeight)
	c.listView.popupLineEdit.ShowPopupLineEdit(posX, posY)
	c.listView.popupLineEdit.CloseEvent = func() {
		if c.listView.popupLineEdit.enterPressed {
			txt := c.listView.popupLineEdit.Text()
			c.listView.currentRow.SetValue(c.listView.currentColumn, txt)
		}
		c.listView.popupLineEdit = nil
		c.listView.Update("ListView")
		c.Focus()
	}
	c.Window().UnIgnoreUpdates()
	c.Update("ListView")
}

func (c *ListViewContent) MouseDblClick(event *MouseDblClickEvent) {
	if c.listView.currentRow != nil {
		c.EditCurrentCell("")
	}
	c.Update("ListView")
}
