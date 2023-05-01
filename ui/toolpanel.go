package ui

import (
	"image"

	"github.com/ipoluianov/goforms/canvas"
	"github.com/nfnt/resize"
)

type ToolPanel struct {
	Panel

	horizontalHeaderLine *Panel
	panelContent         *Panel
	lblHeader            *TextBlock

	toolPanelButtons []*toolPanelButton
	img              image.Image

	title string
}

type toolPanelButton struct {
	isSpace bool
	text    string
	img     image.Image
	onClick func()
	button  *Button
}

func NewToolPanel(parent Widget, title string) *ToolPanel {
	var c ToolPanel
	c.InitControl(parent, &c)
	c.title = title
	return &c
}

func (c *ToolPanel) SetImage(img image.Image) {
	c.img = resize.Resize(32, 32, img, resize.Bicubic)
}

func (c *ToolPanel) AddButton(text string, img image.Image, onClick func()) {
	var b toolPanelButton
	b.text = text
	b.img = img
	b.onClick = onClick

	c.toolPanelButtons = append(c.toolPanelButtons, &b)

	c.makeContent()
}

func (c *ToolPanel) AddSpace() {
	var b toolPanelButton
	b.isSpace = true

	c.toolPanelButtons = append(c.toolPanelButtons, &b)

	c.makeContent()
}

func (c *ToolPanel) Dispose() {
	c.horizontalHeaderLine = nil
	c.panelContent = nil
	c.lblHeader = nil
	c.toolPanelButtons = nil

	c.Panel.Dispose()
}

func (c *ToolPanel) makeContent() {
	c.RemoveAllWidgets()
	c.horizontalHeaderLine = c.AddPanelOnGrid(0, 0)
	c.horizontalHeaderLine.SetBackColor(c.ForeColor())
	c.horizontalHeaderLine.SetMinHeight(2)
	c.panelContent = c.AddPanelOnGrid(0, 1)
	c.panelContent.AddImageBoxOnGrid(0, 0, c.img)
	c.lblHeader = c.panelContent.AddTextBlockOnGrid(1, 0, c.title)
	c.lblHeader.SetFontSize(16)
	c.lblHeader.TextHAlign = canvas.HAlignLeft

	for index, b := range c.toolPanelButtons {
		if b.isSpace {
			btn := NewControl(c.panelContent)
			c.panelContent.AddWidgetOnGrid(btn, index+3, 0)
		} else {
			btn := NewButton(c.panelContent, b.text, c.onBtnClick)
			btn.SetImage(b.img)
			btn.imageWidth = 24
			btn.imageHeight = 24
			btn.showText = false
			btn.SetFixedSize(32, 32)
			btn.SetUserData("onClick", b.onClick)
			c.panelContent.AddWidgetOnGrid(btn, index+3, 0)
		}
	}

	c.Update("ToolPanel")
}

func (c *ToolPanel) onBtnClick(ev *Event) {
	btn := ev.Sender.(*Button)
	btnOnClick := btn.UserData("onClick").(func())
	if btnOnClick != nil {
		btnOnClick()
	}
}

func (c *ToolPanel) XExpandable() bool {
	return true
}

func (c *ToolPanel) YExpandable() bool {
	return false
}
