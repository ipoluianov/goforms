package ui

import (
	"image/color"

	"github.com/ipoluianov/goforms/utils/canvas"
)

type ListViewRow struct {
	UserDataContainer
	row        int
	unitedRows int
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

	xOffset := 0
	skipColumnsCounter := 0
	for columnIndex, column := range c.listView.columns {
		var cell *listViewCell
		if v, ok := c.cells[columnIndex]; ok {
			cell = v
		}

		if skipColumnsCounter > 0 {
			skipColumnsCounter--
			xOffset += column.width
			continue
		}

		if cell != nil {
			if cell.unitedCols > 1 {
				skipColumnsCounter = cell.unitedCols - 1
			}
		}

		cellX := xOffset
		cellWidth := column.width
		cellY := y
		cellHeight := rowHeight

		if cell != nil {
			if cell.unitedCols > 1 {
				cellWidth = 0
				for i := 0; i < cell.unitedCols; i++ {
					cellWidth += c.listView.columns[columnIndex+i].width
				}
			}
		}

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
			cnv = canvas.NewCanvas(cellWidth, rowHeight)
			cell.draw(cnv, c.listView, itemIndex, columnIndex, cellWidth, column, c, rowHeight)
		}

		ctx.DrawImage(cellX, cellY, c.listView.content.Width(), c.listView.content.Height(), cnv.Image())

		c.listView.cache.SetXY(columnIndex, itemIndex, cnv)

		xOffset += column.width
	}

	yOffset += rowHeight

	return yOffset
}
