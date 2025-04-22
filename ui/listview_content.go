package ui

import (
	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/nui/nuikey"
)

type ListViewContent struct {
	Control

	listView *ListView
}

type ListViewCustomDraw struct {
	Row    int
	Column int
	OnDraw func(cnv *canvas.CanvasDirect)
}

func newListViewContent(parent Widget, x int, y int, width int, height int) *ListViewContent {
	var c ListViewContent
	c.InitControl(parent, &c)
	return &c
}

func (c *ListViewContent) Draw(ctx DrawContext) {
	yOffset := 0
	//c.listView.displayedItems = make([]*displayedItem, 0)

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

	c.updateInnerSize()
}

func (c *ListViewContent) ControlType() string {
	return "ListViewContent"
}

func (c *ListViewContent) updateInnerSize() {
	c.InnerSizeOverloaded = true
	c.listView.header.InnerSizeOverloaded = true
	_, c.InnerHeightOverloaded = c.calcTreeColumnSize()
	c.updateInnerWidth()
}

func (c *ListViewContent) updateInnerWidth() {
	c.InnerWidthOverloaded = 0
	for _, column := range c.listView.columns {
		c.InnerWidthOverloaded += column.width
	}
	c.listView.header.InnerWidthOverloaded = c.InnerWidthOverloaded
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

func (c *ListViewContent) findColumnByX(xCoordinates int) int {
	xOffset := 0
	for colIndex, column := range c.listView.columns {
		if xCoordinates >= xOffset && xCoordinates <= xOffset+column.width {
			return colIndex
		}
		xOffset += column.width
	}
	return -1
}

func (c *ListViewContent) findRowByY(yCoordinates int) int {
	yOffset := 0
	for rowIndex, _ := range c.listView.items {
		if yCoordinates >= yOffset && yCoordinates <= yOffset+c.listView.itemHeight {
			return rowIndex
		}
		yOffset += c.listView.itemHeight
	}
	return -1
}

func (c *ListViewContent) getRowByIndex(index int) *ListViewRow {
	if index < 0 || index >= len(c.listView.items) {
		return nil
	}
	return c.listView.items[index]
}

func (c *ListViewContent) getRowYOffset(index int) int {
	if index < 0 || index >= len(c.listView.items) {
		return -1
	}
	return c.listView.items[index].row * c.listView.itemHeight
}

func (c *ListViewContent) getCell(rowIndex int, colIndex int) *listViewCell {
	if rowIndex < 0 || rowIndex >= len(c.listView.items) {
		return nil
	}
	if colIndex < 0 || colIndex >= len(c.listView.columns) {
		return nil
	}
	row := c.listView.items[rowIndex]
	if cell, ok := row.cells[colIndex]; ok {
		return cell
	}
	return nil
}

func (c *ListViewContent) KeyDown(event *KeyDownEvent) bool {
	if event.Key == nuikey.KeyA && event.Modifiers.Ctrl {
		c.listView.SelectAllItems()
		return true
	}

	if event.Key == nuikey.KeyEnter {
		c.EditCurrentCell("")
		return true
	}

	if event.Key == nuikey.KeyArrowUp {
		if c.listView.currentRow != nil && c.listView.currentRow.row > 0 {
			c.listView.SetCurrentRow(c.listView.currentRow.row-1, c.listView.currentColumn, false)
			c.listView.EnsureVisibleCell(c.listView.currentRow.row, c.listView.currentColumn)
		}
		return true
	}

	if event.Key == nuikey.KeyArrowDown {
		if c.listView.currentRow != nil && c.listView.currentRow.row < len(c.listView.items)-1 {
			c.listView.SetCurrentRow(c.listView.currentRow.row+1, c.listView.currentColumn, false)
			c.listView.EnsureVisibleCell(c.listView.currentRow.row, c.listView.currentColumn)
		}
		return true
	}

	if event.Key == nuikey.KeyArrowLeft {
		if c.listView.currentRow != nil && c.listView.currentColumn > 0 {
			col := c.listView.currentColumn - 1
			if col < 0 {
				col = 0
			}
			c.listView.SetCurrentRow(c.listView.currentRow.row, col, false)
			c.listView.EnsureVisibleCell(c.listView.currentRow.row, c.listView.currentColumn)
		}
		return true
	}

	if event.Key == nuikey.KeyArrowRight {
		if c.listView.currentRow != nil && c.listView.currentColumn < len(c.listView.columns)-1 {
			col := c.listView.currentColumn + 1
			if col >= len(c.listView.columns) {
				col = len(c.listView.columns) - 1
			}
			c.listView.SetCurrentRow(c.listView.currentRow.row, col, false)
			c.listView.EnsureVisibleCell(c.listView.currentRow.row, c.listView.currentColumn)
		}
		return true
	}

	if event.Key == nuikey.KeyHome {
		if len(c.listView.items) > 0 {
			c.listView.SetCurrentRow(0, c.listView.currentColumn, false)
			c.listView.EnsureVisibleCell(0, c.listView.currentColumn)
		}
		return true
	}

	if event.Key == nuikey.KeyEnd {
		if len(c.listView.items) > 0 {
			c.listView.SetCurrentRow(len(c.listView.items)-1, c.listView.currentColumn, false)
			c.listView.EnsureVisibleCell(len(c.listView.items)-1, c.listView.currentColumn)
		}
		return true
	}

	if event.Key == nuikey.KeyPageUp {
		if c.listView.currentRow != nil {
			row := c.listView.currentRow.row
			row += c.Height() / c.listView.itemHeight
			if row >= len(c.listView.items) {
				row = len(c.listView.items) - 1
			}
			c.listView.SetCurrentRow(row, c.listView.currentColumn, false)
			c.listView.EnsureVisibleCell(row, c.listView.currentColumn)
		}
		return true
	}

	if event.Key == nuikey.KeyPageDown {
		if c.listView.currentRow != nil {
			row := c.listView.currentRow.row
			row -= c.Height() / c.listView.itemHeight
			if row < 0 {
				row = 0
			}
			c.listView.SetCurrentRow(row, c.listView.currentColumn, false)
			c.listView.EnsureVisibleCell(row, c.listView.currentColumn)
		}
		return true
	}

	return false
}

func (c *ListViewContent) MouseClick(event *MouseClickEvent) {
	if c.listView.OnMouseDown != nil {
		c.listView.OnMouseDown()
	}

	clickedRow := c.findRowByY(event.Y)
	clickedCol := c.findColumnByX(event.X)

	if clickedCol == -1 || clickedRow == -1 {
		if c.listView.AllowDeselectItems {
			c.listView.ClearSelection()
		}
		return
	}

	colOffset := c.listView.calcColumnXOffset(clickedCol)
	row := c.getRowByIndex(clickedRow)

	rowYOffset := c.getRowYOffset(clickedRow)

	c.listView.SetCurrentRow(clickedRow, clickedCol, true)
	c.ScrollEnsureVisible(colOffset, rowYOffset)
	c.ScrollEnsureVisible(colOffset, rowYOffset+c.listView.itemHeight)
	if c.listView.OnItemClicked != nil {
		c.listView.OnItemClicked(row)
	}

	c.Update("ListView")
}

func (c *ListViewContent) KeyChar(event *KeyCharEvent) {
	str := string(event.Ch)
	c.EditCurrentCell(str)
}

func (c *ListViewContent) EditCurrentCell(enteredText string) {
	c.Window().IgnoreUpdates()

	initText := c.listView.currentRow.CellText(c.listView.currentColumn)
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
