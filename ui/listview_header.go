package ui

import (
	"math"

	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/nui/nuimouse"
)

type ListViewHeader struct {
	Control

	pressed             bool
	columnResizing      bool
	columnResizingIndex int

	listView *ListView
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
	if c.columnResizing {
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

func (c *ListViewHeader) Draw(ctx DrawContext) {
	c.listView.updateItemHeight()

	xOffset := 0
	yOffset := c.listView.contentPadding

	ctx.SetColor(c.ForeColor())
	ctx.SetStrokeWidth(1)
	ctx.SetFontSize(c.FontSize())

	for colIndex, column := range c.listView.columns {
		var cnv *canvas.CanvasDirect
		cnv = c.listView.cache.GetXY(colIndex, -1)
		if cnv == nil {
			cnv = canvas.NewCanvas(column.width, c.listView.itemHeight)
			cnv.DrawText(0, 0, column.text, c.FontFamily(), c.FontSize(), c.ForeColor(), false)
			c.listView.cache.SetXY(colIndex, -1, cnv)
		}

		clipX := xOffset + c.listView.contentPadding
		clipW := column.width - c.listView.contentPadding*2

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
