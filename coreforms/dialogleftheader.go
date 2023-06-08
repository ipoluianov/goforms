package coreforms

/*
import (
	"image"

	"github.com/ipoluianov/goforms/canvas"
	"github.com/ipoluianov/goforms/uicontrols"
	"github.com/ipoluianov/goforms/uiinterfaces"
)

type DialogLeftHeader struct {
	uicontrols.Panel
}

func NewDialogLeftHeader(parent uiinterfaces.Widget, img image.Image, header string, text string) *DialogLeftHeader {
	var c DialogLeftHeader
	c.InitControl(parent, &c)

	mainPanel := c.AddPanelOnGrid(0, 0)

	mainPanel.AddImageBoxOnGrid(0, 0, img)

	txtHeader := mainPanel.AddTextBlockOnGrid(0, 1, header)
	txtHeader.TextHAlign = canvas.HAlignCenter
	txtHeader.SetFontSize(16)

	txtText := mainPanel.AddTextBlockOnGrid(0, 2, text)
	txtText.TextHAlign = canvas.HAlignCenter
	txtText.SetFontSize(12)

	mainPanel.AddVSpacerOnGrid(0, 10)

	mainPanel.SetBorderRight(1, c.ForeColor())

	c.SetMinWidth(150)

	return &c
}

func (c *DialogLeftHeader) ControlType() string {
	return "DialogLeftHeader"
}

func (c *DialogLeftHeader) Dispose() {
	// local dispose
	c.Panel.Dispose()
}
*/
