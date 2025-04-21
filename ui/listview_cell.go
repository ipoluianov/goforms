package ui

import "github.com/ipoluianov/goforms/utils/canvas"

type listViewCell struct {
	text       string
	unitedCols int
}

func newListViewCell(text string) *listViewCell {
	var c listViewCell
	c.text = text
	return &c
}

func (c *listViewCell) draw(cnv *canvas.CanvasDirect, listView *ListView, itemIndex int, columnIndex int, cellWidth int, column *ListViewColumn, row *ListViewRow, rowHeight int) {
	var value string
	if c != nil {
		value = c.text
	}

	backColor := listView.content.BackColor()
	foreColor := listView.content.ForeColor()

	if listView.selectionType == 0 {
		if row.selected {
			backColor = listView.selectionBackground.Color()
			foreColor = listView.selectionForeground.Color()
		}
	}

	if listView.selectionType == 1 {
		if row.selected && listView.currentColumn == columnIndex {
			backColor = listView.selectionBackground.Color()
			foreColor = listView.selectionForeground.Color()
		}
	}

	cnv.FillRect(0, 0, cellWidth, rowHeight, backColor)

	col := foreColor

	cnv.DrawTextMultiline(listView.contentPadding, 0, cellWidth-listView.contentPadding*2, rowHeight, column.hAlign, canvas.VAlignCenter, value, col, listView.content.fontFamily.String(), listView.content.fontSize.Float64(), false)

	// Draw borders
	if itemIndex > 0 {
		// Top border
		cnv.DrawLine(0, 0, cellWidth, 0, 1, listView.gridColor.Color())
	}
	if itemIndex == len(listView.items)-1 {
		// Bottom border
		cnv.DrawLine(0, rowHeight-1, cellWidth, rowHeight-1, 1, listView.gridColor.Color())
	}
	if columnIndex > 0 {
		// Left border
		cnv.DrawLine(0, 0, 0, rowHeight, 1, listView.gridColor.Color())
	}
	if columnIndex == len(listView.columns)-1 {
		// Right border
		cnv.DrawLine(cellWidth-1, 0, cellWidth-1, rowHeight, 1, listView.gridColor.Color())
	}

}
