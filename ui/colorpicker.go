package ui

/*
type ColorPicker struct {
	Panel

	btn *Button

	color color.Color

	pressed bool

	OnColorChanged func(colorPicker *ColorPicker, color color.Color)
}

func NewColorPicker(parent Widget) *ColorPicker {
	var c ColorPicker
	c.InitControl(parent, &c)
	c.SetPanelPadding(0)
	c.btn = NewButton(&c, "Select", func(event *Event) {
		c.selectColor()
	})
	c.btn.SetMinWidth(70)
	c.btn.textImageVerticalOrientation = false

	c.AddWidgetOnGrid(c.btn, 0, 0)
	c.color = color.RGBA{0, 0, 0, 255}
	return &c
}

func (c *ColorPicker) Dispose() {
	c.OnColorChanged = nil
	c.Container.Dispose()
}


func (c *ColorPicker) Color() color.Color {
	return c.color
}

func (c *ColorPicker) SetColor(color color.Color) {
	c.color = color
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: 48,
			Y: 16,
		},
	})
	for x := 8; x < 32; x++ {
		for y := 0; y < 16; y++ {
			img.Set(x, y, c.color)
		}
	}
	c.btn.SetImage(img)
	c.btn.SetImageSize(48, 16)
	c.btn.SetShowImage(true)
	c.btn.SetText("Select...")
	c.Update("ColorPicker")
}

func (c *ColorPicker) colorChangedInDialog(color color.Color) {
	c.color = color
	if c.OnColorChanged != nil {
		c.OnColorChanged(c, c.color)
	}
	c.Update("ColorPicker")
}

func (c *ColorPicker) selectColor() {
	dialog := NewColorSelector(c, c.color)
	dialog.ShowDialog()
	dialog.OnColorSelected = func(col color.Color) {
		c.SetColor(col)
		if c.OnColorChanged != nil {
			c.OnColorChanged(c, c.color)
		}
	}

	dialog.OnAccept = func() {
		c.SetColor(dialog.ResColor())
		if c.OnColorChanged != nil {
			c.OnColorChanged(c, c.color)
		}
	}

	dialog.OnReject = func() {
		c.SetColor(dialog.ResColor())
		if c.OnColorChanged != nil {
			c.OnColorChanged(c, c.color)
		}
	}

	c.Update("ColorPicker")
}
*/
