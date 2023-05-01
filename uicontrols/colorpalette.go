package uicontrols

import (
	"fmt"
	"image/color"

	"github.com/ipoluianov/goforms/ui"
	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiinterfaces"
	"golang.org/x/image/colornames"
)

type ColorPalette struct {
	Control

	color color.Color

	pressed bool

	OnColorChanged func(colorPicker *ColorPalette, color color.Color)

	loadedColors map[string]color.Color
	rectSize     int

	colCount int
	rowCount int

	xRectSize int
	yRectSize int
}

func NewColorPalette(parent uiinterfaces.Widget) *ColorPalette {
	var c ColorPalette
	c.InitControl(parent, &c)
	c.color = color.RGBA{0, 0, 0, 255}
	c.SetXExpandable(true)
	c.SetYExpandable(true)
	c.colCount = 0
	c.rowCount = 0

	c.loadedColors = make(map[string]color.Color)
	c.addColor(0, 0, colornames.White)
	c.addColor(1, 0, colornames.Black)
	c.addColor(2, 0, colornames.Red)
	c.addColor(3, 0, colornames.Green)
	c.addColor(4, 0, colornames.Blue)

	c.addColor(0, 1, colornames.Magenta)
	c.addColor(1, 1, colornames.Yellow)
	c.addColor(2, 1, colornames.Gray)

	return &c
}

func (c *ColorPalette) addColor(x int, y int, col color.Color) {
	if x > c.colCount-1 {
		c.colCount = x + 1
	}
	if y > c.rowCount-1 {
		c.rowCount = y + 1
	}
	c.loadedColors[c.posKey(x, y)] = col
}

func (c *ColorPalette) posKey(x int, y int) string {
	return fmt.Sprint(x, "_", y)
}

func (c *ColorPalette) Dispose() {
	c.OnColorChanged = nil
	c.Control.Dispose()
}

func (c *ColorPalette) Subclass() string {
	if c.pressed {
		return "pressed"
	}
	return c.Control.Subclass()
}

func (c *ColorPalette) ControlType() string {
	return "ColorPalette"
}

func (c *ColorPalette) colorByCoordinates(x int, y int) (color.Color, error) {
	cellX := x / c.xRectSize
	cellY := y / c.yRectSize

	if col, ok := c.loadedColors[c.posKey(cellX, cellY)]; ok {
		return col, nil
	}

	return nil, fmt.Errorf("no color found")
}

func (c *ColorPalette) Draw(ctx ui.DrawContext) {
	/*ctx.DrawRect(0, 0, c.InnerWidth(), c.InnerHeight(), c.rightBorderColor.Color(), 1)

	xRectSize := c.Width() / c.colCount
	yRectSize := c.Height() / c.rowCount

	for y := 0; y < c.rowCount; y++ {
		for x := 0; x < c.colCount; x++ {
			col, err := c.colorByCoordinates(x*xRectSize+2, y*yRectSize+2)
			if err == nil {
				ctx.FillRect(x*xRectSize, y*yRectSize, xRectSize, yRectSize, col)
			}
		}
	}*/
}

func (c *ColorPalette) Color() color.Color {
	return c.color
}

func (c *ColorPalette) SetColor(color color.Color) {
	c.color = color
	c.Update("ColorPicker")
}

func (c *ColorPalette) TabStop() bool {
	return false
}

func (c *ColorPalette) MouseDown(event *uievents.MouseDownEvent) {
	c.pressed = true
	c.Update("Button")
}

func (c *ColorPalette) MouseUp(event *uievents.MouseUpEvent) {
	if c.pressed {
		c.pressed = false
		c.Update("ColorPalette")
		if event.X >= 0 && event.Y >= 0 && event.X < c.InnerWidth() && event.Y < c.InnerHeight() {

		}
	}
}

func (c *ColorPalette) colorChangedInDialog(color color.Color) {
	c.color = color
	if c.OnColorChanged != nil {
		c.OnColorChanged(c, c.color)
	}
	c.Update("ColorPicker")
}

func (c *ColorPalette) selectColor() {
	/*var originalColor color.RGBA
	{
		r, g, b, a := c.color.RGBA()
		originalColor = color.RGBA{uint8(r / 256), uint8(g / 256), uint8(b / 256), uint8(a / 256)}
	}
	ok, col := c.OwnWindow.SelectColorDialog(c.color, c.colorChangedInDialog)
	r, g, b, a := col.RGBA()
	if ok {
		c.color = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	} else {
		c.color = originalColor
	}

	if c.OnColorChanged != nil {
		c.OnColorChanged(c, c.color)
	}
	c.Update("ColorPicker")*/
}
