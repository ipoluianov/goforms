package uiforms

import (
	"image/color"

	"github.com/ipoluianov/goforms/uicontrols"
)

type ColorPickerDialog struct {
	Form
	numR      *uicontrols.SpinBox
	numG      *uicontrols.SpinBox
	numB      *uicontrols.SpinBox
	numA      *uicontrols.SpinBox
	btnOK     *uicontrols.Button
	btnCancel *uicontrols.Button

	color color.Color

	loading        bool
	onColorChanged func(color color.Color)

	ctrlSample *uicontrols.Panel
}

func (c *ColorPickerDialog) OnInit() {
	c.Resize(500, 500)
	//c.Panel().SetAbsolutePositioning(true)
	c.SetTitle("Select color ...")

	c.loading = true

	panelTop := c.Panel().AddPanelOnGrid(0, 0)

	c.ctrlSample = panelTop.AddPanelOnGrid(0, 0)
	c.ctrlSample.SetName("Sample control")

	panelValues := panelTop.AddPanelOnGrid(1, 0)
	panelValues.SetName("panelValues")

	panelValues.AddTextBlockOnGrid(0, 0, "R")
	c.numR = panelValues.AddSpinBoxOnGrid(1, 0)
	c.adjustColorSpinBox(c.numR)

	panelValues.AddTextBlockOnGrid(0, 1, "G")
	c.numG = panelValues.AddSpinBoxOnGrid(1, 1)
	c.adjustColorSpinBox(c.numG)

	panelValues.AddTextBlockOnGrid(0, 2, "B")
	c.numB = panelValues.AddSpinBoxOnGrid(1, 2)
	c.adjustColorSpinBox(c.numB)

	panelValues.AddTextBlockOnGrid(0, 3, "A")
	c.numA = panelValues.AddSpinBoxOnGrid(1, 3)
	c.adjustColorSpinBox(c.numA)

	panelBottom := c.Panel().AddPanelOnGrid(0, 1)
	panelBottom.AddListViewOnGrid(0, 0)

	panelButtons := c.Panel().AddPanelOnGrid(0, 2)
	panelButtons.AddHSpacerOnGrid(0, 0)
	panelButtons.AddButtonOnGrid(1, 0, "OK", c.AcceptButton)
	panelButtons.AddButtonOnGrid(2, 0, "Cancel", c.RejectButton)

	r, g, b, a := c.color.RGBA()
	c.numR.SetValue(float64(r >> 8))
	c.numG.SetValue(float64(g >> 8))
	c.numB.SetValue(float64(b >> 8))
	c.numA.SetValue(float64(a >> 8))
	c.loading = false

	c.SetColor(c.color)
}

func (c *ColorPickerDialog) adjustColorSpinBox(spinBox *uicontrols.SpinBox) {
	spinBox.OnValueChanged = c.numsChanged
	//spinBox.SetMaxWidth(40)
	spinBox.SetMinValue(0)
	spinBox.SetMaxValue(255)
	spinBox.SetPrecision(0)
	spinBox.SetXExpandable(false)
}

func (c *ColorPickerDialog) OnClose() bool {
	c.numR = nil
	c.numG = nil
	c.numB = nil
	c.numA = nil
	c.btnOK = nil
	c.btnCancel = nil
	c.onColorChanged = nil
	c.Form.Dispose()
	return true
}

func (c *ColorPickerDialog) SetColor(color color.Color) {
	c.color = color

	c.loading = true
	r, g, b, a := c.color.RGBA()
	if c.numR != nil {
		c.numR.SetValue(float64(r >> 8))
	}
	if c.numG != nil {
		c.numG.SetValue(float64(g >> 8))
	}
	if c.numB != nil {
		c.numB.SetValue(float64(b >> 8))
	}
	if c.numA != nil {
		c.numA.SetValue(float64(a >> 8))
	}

	if c.ctrlSample != nil {
		c.ctrlSample.SetBackColor(color)
	}

	c.loading = false
}

func (c *ColorPickerDialog) numsChanged(spinBox *uicontrols.SpinBox, value float64) {
	if c.loading {
		return
	}

	c.color = color.RGBA{uint8(c.numR.Value()), uint8(c.numG.Value()), uint8(c.numB.Value()), uint8(c.numA.Value())}
	c.ctrlSample.SetBackColor(c.color)

	if c.onColorChanged != nil {
		c.onColorChanged(c.color)
	}
}
