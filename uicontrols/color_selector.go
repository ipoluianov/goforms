package uicontrols

import (
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"image/color"
)

type ColorSelector struct {
	Dialog
	col               color.Color
	resColor          color.Color
	txtUnitName       *TextBox
	OnColorSelected   func(col color.Color)
	btnOK             *Button
	panelPalette      *Panel
	panelPaletteCount int
	numR              *SpinBox
	numG              *SpinBox
	numB              *SpinBox
	numA              *SpinBox

	loadingColor bool

	panelDetails    *Panel
	txtCurrentColor *TextBlock
}

func NewColorSelector(parent uiinterfaces.Widget, col color.Color) *ColorSelector {
	var c ColorSelector
	c.resColor = col
	c.InitControl(parent, &c)

	pContent := c.ContentPanel().AddPanelOnGrid(0, 0)
	c.panelPalette = pContent.AddPanelOnGrid(0, 0)

	c.panelDetails = pContent.AddPanelOnGrid(1, 0)
	c.txtCurrentColor = c.panelDetails.AddTextBlockOnGrid(0, 0, "")
	c.txtCurrentColor.SetMinHeight(48)
	c.txtCurrentColor.SetMinWidth(48)
	c.txtCurrentColor.SetBackColor(c.col)

	panelDigits := c.panelDetails.AddPanelOnGrid(0, 1)
	panelDigits.SetPanelPadding(0)

	panelDigits.AddTextBlockOnGrid(0, 1, "R:")
	c.numR = panelDigits.AddSpinBoxOnGrid(1, 1)
	c.numR.SetMinValue(0)
	c.numR.SetMaxValue(255)
	c.numR.SetPrecision(0)
	c.numR.OnValueChanged = func(spinBox *SpinBox, value float64) {
		c.updateColorFromDigits()
	}
	panelDigits.AddTextBlockOnGrid(0, 2, "G:")
	c.numG = panelDigits.AddSpinBoxOnGrid(1, 2)
	c.numG.SetMinValue(0)
	c.numG.SetMaxValue(255)
	c.numG.SetPrecision(0)
	c.numG.OnValueChanged = func(spinBox *SpinBox, value float64) {
		c.updateColorFromDigits()
	}
	panelDigits.AddTextBlockOnGrid(0, 3, "B:")
	c.numB = panelDigits.AddSpinBoxOnGrid(1, 3)
	c.numB.SetMinValue(0)
	c.numB.SetMaxValue(255)
	c.numB.SetPrecision(0)
	c.numB.OnValueChanged = func(spinBox *SpinBox, value float64) {
		c.updateColorFromDigits()
	}
	panelDigits.AddTextBlockOnGrid(0, 4, "A:")
	c.numA = panelDigits.AddSpinBoxOnGrid(1, 4)
	c.numA.SetMinValue(0)
	c.numA.SetMaxValue(255)
	c.numA.SetPrecision(0)
	c.numA.OnValueChanged = func(spinBox *SpinBox, value float64) {
		c.updateColorFromDigits()
	}

	c.panelDetails.AddVSpacerOnGrid(0, 10)

	for i := 0; i < 8; i++ {
		v := i * 32
		if v > 255 {
			v = 255
		}
		vb := byte(v)
		c.addColorToPalette(color.RGBA{R: vb, G: vb, B: vb, A: 255})
	}

	for i := 0; i < 8; i++ {
		v := i * 32
		if v > 255 {
			v = 255
		}
		vb := byte(v)
		c.addColorToPalette(color.RGBA{R: 0, G: 0, B: vb, A: 255})
	}

	for i := 0; i < 8; i++ {
		v := i * 32
		if v > 255 {
			v = 255
		}
		vb := byte(v)
		c.addColorToPalette(color.RGBA{R: 0, G: vb, B: 0, A: 255})
	}

	for i := 0; i < 8; i++ {
		v := i * 32
		if v > 255 {
			v = 255
		}
		vb := byte(v)
		c.addColorToPalette(color.RGBA{R: vb, G: 0, B: 0, A: 255})
	}

	c.addColorToPalette(color.RGBA{R: 255, G: 64, B: 0, A: 255})
	c.addColorToPalette(color.RGBA{R: 255, G: 128, B: 0, A: 255})
	c.addColorToPalette(color.RGBA{R: 255, G: 192, B: 0, A: 255})
	c.addColorToPalette(color.RGBA{R: 255, G: 255, B: 0, A: 255})

	c.addColorToPalette(color.RGBA{R: 64, G: 255, B: 0, A: 255})
	c.addColorToPalette(color.RGBA{R: 128, G: 255, B: 0, A: 255})
	c.addColorToPalette(color.RGBA{R: 192, G: 255, B: 0, A: 255})
	c.addColorToPalette(color.RGBA{R: 255, G: 255, B: 0, A: 255})

	c.addColorToPalette(color.RGBA{R: 64, G: 0, B: 255, A: 255})
	c.addColorToPalette(color.RGBA{R: 128, G: 0, B: 255, A: 255})
	c.addColorToPalette(color.RGBA{R: 192, G: 0, B: 255, A: 255})
	c.addColorToPalette(color.RGBA{R: 255, G: 0, B: 255, A: 255})

	c.addColorToPalette(color.RGBA{R: 255, G: 0, B: 64, A: 255})
	c.addColorToPalette(color.RGBA{R: 255, G: 0, B: 128, A: 255})
	c.addColorToPalette(color.RGBA{R: 255, G: 0, B: 192, A: 255})
	c.addColorToPalette(color.RGBA{R: 255, G: 0, B: 255, A: 255})

	c.addColorToPalette(color.RGBA{R: 0, G: 255, B: 64, A: 255})
	c.addColorToPalette(color.RGBA{R: 0, G: 255, B: 128, A: 255})
	c.addColorToPalette(color.RGBA{R: 0, G: 255, B: 192, A: 255})
	c.addColorToPalette(color.RGBA{R: 0, G: 255, B: 255, A: 255})

	c.addColorToPalette(color.RGBA{R: 0, G: 64, B: 255, A: 255})
	c.addColorToPalette(color.RGBA{R: 0, G: 128, B: 255, A: 255})
	c.addColorToPalette(color.RGBA{R: 0, G: 192, B: 255, A: 255})
	c.addColorToPalette(color.RGBA{R: 0, G: 255, B: 255, A: 255})
	c.addColorToPalette(color.RGBA{R: 0, G: 0, B: 0, A: 0})

	pContent.AddVSpacerOnGrid(0, 5)

	pButtons := c.ContentPanel().AddPanelOnGrid(0, 1)
	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.TryAccept = func() bool {
		c.btnOK.SetEnabled(false)
		c.resColor = c.col
		c.TryAccept = nil
		c.Accept()
		return false
	}

	c.btnOK.SetMinWidth(70)
	btnCancel := pButtons.AddButtonOnGrid(2, 0, "Cancel", func(event *uievents.Event) {
		c.Reject()
	})
	btnCancel.SetMinWidth(70)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(btnCancel)

	c.SetColor(col)

	return &c
}

func (c *ColorSelector) updateColorFromDigits() {
	if c.loadingColor {
		return
	}
	c.SetColor(color.RGBA{R: uint8(c.numR.Value()), G: uint8(c.numG.Value()), B: uint8(c.numB.Value()), A: uint8(c.numA.Value())})
}

func (c *ColorSelector) SetColor(col color.Color) {
	c.loadingColor = true
	changed := false
	if c.col != col {
		changed = true
	}
	c.col = col
	c.txtCurrentColor.SetBackColor(col)

	r, g, b, a := c.col.RGBA()
	r /= 256
	g /= 256
	b /= 256
	a /= 256
	c.numR.SetValue(float64(r))
	c.numG.SetValue(float64(g))
	c.numB.SetValue(float64(b))
	c.numA.SetValue(float64(a))

	if changed {
		c.NotifyColorSelection()
	}
	c.loadingColor = false
}

func (c *ColorSelector) addColorToPalette(col color.Color) {
	text := ""
	r, g, b, a := col.RGBA()
	if r == 0 && g == 0 && b == 0 && a == 0 {
		text = "TR"
	}
	btn := c.panelPalette.AddButtonOnGrid(c.panelPaletteCount%8, c.panelPaletteCount/8, text, func(event *uievents.Event) {
		c.SetColor(col)
	})
	btn.SetBackColor(col)
	btn.SetMaxWidth(24)
	btn.SetMaxHeight(24)
	btn.SetMinWidth(24)
	btn.SetMinHeight(24)
	c.panelPaletteCount++
}

func (c *ColorSelector) NotifyColorSelection() {
	if c.OnColorSelected != nil {
		c.OnColorSelected(c.col)
	}
}

func (c *ColorSelector) OnInit() {
	c.Dialog.OnInit()
	c.SetTitle("Select color")
	c.Resize(400, 400)
}

func (c *ColorSelector) Color() color.Color {
	return c.col
}

func (c *ColorSelector) ResColor() color.Color {
	return c.resColor
}
