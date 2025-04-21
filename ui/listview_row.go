package ui

import (
	"image/color"

	"github.com/ipoluianov/goforms/utils/canvas"
)

type displayedItem struct {
	currentX      int
	currentY      int
	currentWidth  int
	currentHeight int
	item          *ListViewRow
}

type ListViewRow struct {
	UserDataContainer
	row        int
	unitedRows int
	unitedCols int
	//values     map[int]string
	cells      map[int]*listViewCell
	selected   bool
	listView   *ListView
	foreColors map[int]color.Color
}

func (c *ListViewRow) SetValue(column int, text string) {
	for index, item := range c.listView.items {
		if item == c {
			c.listView.SetCellText(index, column, text)
		}
	}
}

func (c *ListViewRow) SetForeColorForRow(color color.Color) {
	for colIndex := 0; colIndex < len(c.listView.columns); colIndex++ {
		if colorForCell, ok := c.foreColors[colIndex]; ok {
			if colorForCell != color {
				c.foreColors[colIndex] = color
				c.listView.cache.Clear()
				c.listView.Update("ListViewItem")
			}
		} else {
			c.foreColors[colIndex] = color
			c.listView.cache.Clear()
			c.listView.Update("ListViewItem")
		}
	}
}

func (c *ListViewRow) SetForeColorForCell(colIndex int, color color.Color) {
	if colorForCell, ok := c.foreColors[colIndex]; ok {
		if colorForCell != color {
			c.foreColors[colIndex] = color
			c.listView.cache.Clear()
			c.listView.Update("ListViewItem")
		}
	} else {
		c.foreColors[colIndex] = color
		c.listView.cache.Clear()
		c.listView.Update("ListViewItem")
	}
}

func (c *ListViewRow) CellText(column int) string {
	if v, ok := c.cells[column]; ok {
		return v.text
	}
	return ""
}

func (c *ListViewRow) draw(ctx DrawContext, y int, itemIndex int) int {
	yOffset := 0

	rowHeight := c.listView.itemHeight

	visRect := c.listView.content.VisibleInnerRect()

	var dItem displayedItem
	dItem.currentX = 0
	dItem.currentY = y
	dItem.currentWidth = 100
	dItem.currentHeight = c.listView.itemHeight
	dItem.item = c
	c.listView.displayedItems = append(c.listView.displayedItems, &dItem)

	xOffset := 0
	for columnIndex, column := range c.listView.columns {
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
		cnv = c.listView.cache.GetXY(columnIndex, itemIndex)
		if cnv == nil {
			cnv = canvas.NewCanvas(column.width, rowHeight)

			//value := c.values[columnIndex]
			var value string
			if v, ok := c.cells[columnIndex]; ok {
				value = v.text
			}

			backColor := c.listView.content.BackColor()
			foreColor := c.listView.content.ForeColor()

			if c.listView.selectionType == 0 {
				if c.selected {
					backColor = c.listView.selectionBackground.Color()
					foreColor = c.listView.selectionForeground.Color()
				}
			}

			if c.listView.selectionType == 1 {
				if c.selected && c.listView.currentColumn == columnIndex {
					backColor = c.listView.selectionBackground.Color()
					foreColor = c.listView.selectionForeground.Color()
				}
			}

			cnv.FillRect(0, 0, column.width, rowHeight, backColor)

			col := foreColor
			if colorForCell, ok := c.foreColors[columnIndex]; ok {
				if colorForCell != nil {
					col = colorForCell
				}
			}

			cnv.DrawTextMultiline(c.listView.contentPadding, 0, column.width-c.listView.contentPadding*2, rowHeight, column.hAlign, canvas.VAlignCenter, value, col, c.listView.content.fontFamily.String(), c.listView.content.fontSize.Float64(), false)

			// Draw borders
			if itemIndex > 0 {
				// Top border
				cnv.DrawLine(0, 0, column.width, 0, 1, c.listView.gridColor.Color())
			}
			if itemIndex == len(c.listView.items)-1 {
				// Bottom border
				cnv.DrawLine(0, rowHeight-1, column.width, rowHeight-1, 1, c.listView.gridColor.Color())
			}
			if columnIndex > 0 {
				// Left border
				cnv.DrawLine(0, 0, 0, rowHeight, 1, c.listView.gridColor.Color())
			}
			if columnIndex == len(c.listView.columns)-1 {
				// Right border
				cnv.DrawLine(column.width-1, 0, column.width-1, rowHeight, 1, c.listView.gridColor.Color())
			}
		}

		ctx.DrawImage(cellX, cellY, c.listView.content.Width(), c.listView.content.Height(), cnv.Image())

		c.listView.cache.SetXY(columnIndex, itemIndex, cnv)

		xOffset += column.width
	}

	yOffset += rowHeight

	return yOffset
}
