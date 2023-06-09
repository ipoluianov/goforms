package ui

import (
	"image"
	"image/draw"
	"math/rand"
	"strconv"

	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/goforms/utils/uiproperties"
)

type Table struct {
	Control
	verticalHeaderWidth    int
	horizontalHeaderHeight int
	columns                []*TableColumn
	rows                   []*TableRow
	focusedCell            Widget
	lastMouseDownWidget    Widget
}

type TableColumn struct {
	Name  string
	Width int

	leftBorderWidth  *uiproperties.Property
	leftBorderColor  *uiproperties.Property
	rightBorderWidth *uiproperties.Property
	rightBorderColor *uiproperties.Property
}

type TableRow struct {
	Name   string
	Height int
	Cells  []Widget
	rows   []*TableRow

	topBorderWidth    *uiproperties.Property
	topBorderColor    *uiproperties.Property
	bottomBorderWidth *uiproperties.Property
	bottomBorderColor *uiproperties.Property
}

func NewTable(parent Widget, x int, y int, width int, height int) *Table {
	var c Table
	c.InitControl(nil, &c)

	c.columns = make([]*TableColumn, 0)
	c.rows = make([]*TableRow, 0)

	c.horizontalHeaderHeight = 26
	c.verticalHeaderWidth = 100

	for i := 0; i < 50; i++ {
		c.AddColumn("Column" + strconv.Itoa(i))
	}

	for i := 0; i < 50; i++ {
		c.AddRow("Row " + strconv.Itoa(i))
	}

	return &c
}

func (c *Table) btnClick(event *Event) {
	//c.OwnWindow.MessageBox("123", "asd")
}

func (c *Table) AddRow(name string) {
	row := TableRow{Name: name, Height: 21}
	row.topBorderWidth = uiproperties.NewProperty("TopBorderWidth", uiproperties.PropertyTypeInt)
	row.topBorderWidth.SetOwnValue(0)
	row.bottomBorderWidth = uiproperties.NewProperty("BottomBorderWidth", uiproperties.PropertyTypeInt)
	row.bottomBorderWidth.SetOwnValue(1)
	c.rows = append(c.rows, &row)
	row.Cells = make([]Widget, 0)
	for i := 0; i < len(c.columns); i++ {
		textBox := NewTextBox(c)
		textBox.SetText("text " + strconv.Itoa(rand.Int()/1000000000000))
		row.Cells = append(row.Cells, textBox)
	}
	c.updateCellSizes()
}

func (c *Table) AddColumn(name string) {
	col := TableColumn{Name: name, Width: 100}
	col.leftBorderWidth = uiproperties.NewProperty("LeftBorderWidth", uiproperties.PropertyTypeInt)
	col.leftBorderWidth.SetOwnValue(0)
	col.rightBorderWidth = uiproperties.NewProperty("RightBorderWidth", uiproperties.PropertyTypeInt)
	col.rightBorderWidth.SetOwnValue(1)
	c.columns = append(c.columns, &col)
	//c.updateCellSizes()
}

func (c *Table) Draw(ctx DrawContext) {
	/*c.updateInnerSize()

	var visibleRect image.Rectangle
	visibleRect.Min.X = c.ScrollOffsetX()
	visibleRect.Min.Y = c.ScrollOffsetY()
	visibleRect.Max.X = c.Width() + c.ScrollOffsetX()
	visibleRect.Max.Y = c.Height() + c.ScrollOffsetY()

	xOffset := c.verticalHeaderWidth
	for _, col := range c.columns {
		if xOffset >= visibleRect.Min.X && xOffset <= visibleRect.Max.X {
			ctx.DrawText(xOffset, 0, col.Name, c.fontFamily.String(), c.fontSize.Float64(), c.foregroundColor.Color())
		}

		//ctx.DrawRect(xOffset, 0, col.Width, c.height, colorForeGround)
		//ctx.DrawRect(xOffset, 0, col.leftBorderWidth.Int(), c.height.Int(), c.foregroundColor.Color())
		//ctx.DrawRect(xOffset + col.Width + col.leftBorderWidth.Int(), 0, col.rightBorderWidth.Int(), c.innerHeight(), c.foregroundColor.Color())
		xOffset += col.Width
		xOffset += col.leftBorderWidth.Int()
		xOffset += col.rightBorderWidth.Int()
	}

	yOffset := c.horizontalHeaderHeight
	for _, row := range c.rows {
		xOffset = c.verticalHeaderWidth
		if yOffset >= visibleRect.Min.Y && yOffset <= visibleRect.Max.Y {

			//ctx.DrawRect(0, yOffset, c.width.Int(), row.topBorderWidth.Int(), c.foregroundColor.Color())
			//ctx.DrawRect(0, yOffset + row.topBorderWidth.Int() + row.Height, c.innerWidth(), row.bottomBorderWidth.Int(), c.foregroundColor.Color())
			ctx.DrawText(0, yOffset, row.Name, c.fontFamily.String(), c.fontSize.Float64(), c.foregroundColor.Color())

			colIndex := 0
			for _, cell := range row.Cells {
				if xOffset >= visibleRect.Min.X && xOffset <= visibleRect.Max.X {
					cell.DrawWidgetOnCanvas(ctx)
				}
				col := c.columns[colIndex]
				xOffset += col.Width
				xOffset += col.leftBorderWidth.Int()
				xOffset += col.rightBorderWidth.Int()
				colIndex++
			}
		}
		yOffset += row.Height + row.topBorderWidth.Int() + row.bottomBorderWidth.Int()
	}*/

}

func drawWidgetOnCanvasTest(ctx *canvas.CanvasDirect, w Widget) {

	// Content of control
	cnvInner := canvas.NewCanvas(w.InnerWidth(), w.InnerHeight())

	//w.drawBackground(cnvInner)
	//w.Draw(cnvInner)

	// 100
	bInner := cnvInner.Image().Bounds()
	bInner.Min.X += w.LeftBorderWidth() - w.ScrollOffsetX()
	bInner.Min.Y += w.TopBorderWidth() - w.ScrollOffsetY()
	bInner.Max.X += w.LeftBorderWidth() - w.ScrollOffsetX()
	bInner.Max.Y += w.TopBorderWidth() - w.ScrollOffsetY()

	cnv := canvas.NewCanvas(w.Width(), w.Height())
	//w.DrawBackground(cnv)
	draw.Draw(cnv.Image(), bInner, cnvInner.Image(), image.ZP, draw.Over) // Content 20
	//w.DrawBorders(cnv)                                                    // Borders 10
	//w.DrawScrollBars(cnv)

	// Result to parent canvas
	b := cnv.Image().Bounds()
	b.Min.X += w.X()
	b.Min.Y += w.Y()
	b.Max.X += w.X()
	b.Max.Y += w.Y()
	draw.Draw(ctx.Image(), b, cnv.Image(), image.ZP, draw.Src)
}

func (c *Table) updateCellSizes() {
	yOffset := c.horizontalHeaderHeight
	for _, row := range c.rows {
		xOffset := c.verticalHeaderWidth
		for colIndex, col := range c.columns {
			row.Cells[colIndex].SetX(xOffset + col.leftBorderWidth.Int())
			row.Cells[colIndex].SetY(yOffset + row.topBorderWidth.Int())
			row.Cells[colIndex].SetWidth(col.Width)
			row.Cells[colIndex].SetHeight(row.Height)
			xOffset += col.Width
			xOffset += col.leftBorderWidth.Int()
			xOffset += col.rightBorderWidth.Int()
		}

		yOffset += row.Height
		yOffset += row.topBorderWidth.Int()
		yOffset += row.bottomBorderWidth.Int()
	}
}

func (c *Table) MouseClick(event *MouseClickEvent) {
	w := c.widgetByPoint(event.X, event.Y)
	c.clearFocusForCells()
	c.focusedCell = nil
	if w != nil {
		w.SetFocus(true)
		c.focusedCell = w
		ev := *event
		ev.X -= w.X()
		ev.Y -= w.Y()
		w.ProcessMouseClick(&ev)
	}
}

func (c *Table) MouseMove(event *MouseMoveEvent) {
	if c.lastMouseDownWidget != nil {
		c.lastMouseDownWidget.ProcessMouseMove(event.Translate(c.lastMouseDownWidget))
	} else {
		w := c.widgetByPoint(event.X, event.Y)
		if w != nil {
			ev := *event
			ev.X -= w.X()
			ev.Y -= w.Y()
			w.ProcessMouseMove(&ev)
		}
	}
}

func (c *Table) MouseDown(event *MouseDownEvent) {
	w := c.widgetByPoint(event.X, event.Y)
	if w != nil {
		ev := *event
		ev.X -= w.X()
		ev.Y -= w.Y()
		w.ProcessMouseDown(&ev)
		c.lastMouseDownWidget = w
	}
}

func (c *Table) MouseUp(event *MouseUpEvent) {
	if c.lastMouseDownWidget != nil {
		c.lastMouseDownWidget.ProcessMouseUp(event.Translate(c.lastMouseDownWidget))
		c.lastMouseDownWidget = nil
	} else {
		w := c.widgetByPoint(event.X, event.Y)
		if w != nil {
			ev := *event
			ev.X -= w.X()
			ev.Y -= w.Y()
			w.ProcessMouseUp(&ev)
		}
	}
}

func (c *Table) KeyChar(event *KeyCharEvent) {
	if c.focusedCell != nil {
		c.focusedCell.ProcessKeyChar(event)
	}
}

func (c *Table) KeyDown(event *KeyDownEvent) bool {
	if c.focusedCell != nil {
		c.focusedCell.KeyDown(event)
	}
	return true
}

func (c *Table) KeyUp(event *KeyUpEvent) {
	if c.focusedCell != nil {
		c.focusedCell.KeyUp(event)
	}
}

func (c *Table) widgetByPoint(x, y int) Widget {
	wCol, wRow := c.cellByPoint(x, y)
	if wCol < 0 || wRow < 0 {
		return nil
	}

	return c.rows[wRow].Cells[wCol]
}

func (c *Table) cellByPoint(x, y int) (int, int) {
	for rowIndex, row := range c.rows {
		for colIndex := range c.columns {
			if x > row.Cells[colIndex].X() && x < row.Cells[colIndex].X()+row.Cells[colIndex].Width() && y > row.Cells[colIndex].Y() && y < row.Cells[colIndex].Y()+row.Cells[colIndex].Height() {
				return colIndex, rowIndex
			}
		}
	}

	return -1, -1
}

func (c *Table) clearFocus() {
	c.Control.ClearFocus()
	c.clearFocusForCells()
}

func (c *Table) clearFocusForCells() {
	for _, row := range c.rows {
		for colIndex := range c.columns {
			row.Cells[colIndex].ClearFocus()
		}
	}
}

func (c *Table) setFocus(focus bool) {
	c.Control.SetFocus(focus)
	if !focus {
		c.clearFocusForCells()
	}
}

func (c *Table) updateInnerSize() {

	height := c.horizontalHeaderHeight
	for _, row := range c.rows {
		height += row.Height
		height += row.topBorderWidth.Int()
		height += row.bottomBorderWidth.Int()
	}
	width := c.verticalHeaderWidth
	for _, column := range c.columns {
		width += column.Width
		width += column.leftBorderWidth.Int()
		width += column.rightBorderWidth.Int()
	}

	c.innerHeightOverloaded = height
	c.innerWidthOverloaded = width
	c.innerSizeOverloaded = true
}
