package uicontrols

import (
	"image"
	"image/color"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ipoluianov/goforms/canvas"
	"github.com/ipoluianov/goforms/ui"
	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiinterfaces"
)

type Button struct {
	Container

	text        string
	img         image.Image
	imgDisabled image.Image

	imageBeforeText              bool
	textImageVerticalOrientation bool

	showText   bool
	txtBlock   *TextBlock
	textHAlign canvas.HAlign
	textVAlign canvas.VAlign

	imgBox    *ImageBox
	showImage bool

	pressed            bool
	mousePointerInRect bool

	imageForeColor bool

	checked_ bool

	imageWidth     int
	imageHeight    int
	padding        int
	drawBackground bool

	onPress func(event *uievents.Event)
}

func NewButton(parent uiinterfaces.Widget, text string, onPress func(event *uievents.Event)) *Button {
	var c Button
	c.InitControl(parent, &c)
	c.SetStatistics("btn_" + text)

	c.text = text
	c.img = nil
	c.onPress = onPress
	c.showImage = false
	c.showText = false
	c.imageBeforeText = true
	c.imageWidth = 24
	c.imageHeight = 24
	c.textImageVerticalOrientation = true
	c.textHAlign = canvas.HAlignCenter
	c.textVAlign = canvas.VAlignCenter
	c.padding = 3

	c.cellPadding = 0

	if len(c.text) > 0 {
		c.showText = true
	} else {
		c.showText = false
	}

	c.rebuildContent()
	return &c
}

func (c *Button) SetOnPress(onPress func(ev *uievents.Event)) {
	c.onPress = onPress
}

func (c *Button) Subclass() string {
	if c.pressed && c.mousePointerInRect {
		return "pressed"
	}
	return c.Control.Subclass()
}

func (c *Button) ControlType() string {
	return "Button"
}

func (c *Button) SetTextImageVerticalOrientation(textImageVerticalOrientation bool) {
	c.textImageVerticalOrientation = textImageVerticalOrientation
	c.rebuildContent()
}

func (c *Button) SetForeColor(foreColor color.Color) {
	c.Container.SetForeColor(foreColor)
	if c.txtBlock != nil {
		c.txtBlock.SetForeColor(foreColor)
	}
}

func (c *Button) SetPadding(padding int) {
	c.padding = padding
	c.rebuildContent()
	c.Window().UpdateLayout()
}

func (c *Button) SetChecked(checked bool) {
	c.checked_ = checked
	c.Update("Button")
}

func (c *Button) SetEnabled(enabled bool) {
	c.Container.SetEnabled(enabled)
}

func (c *Button) Text() string {
	if c.txtBlock != nil {
		return c.txtBlock.Text()
	}
	return ""
}

func (c *Button) Draw(ctx ui.DrawContext) {

	if c.imgBox != nil {
		if c.imageForeColor {
			img := c.imgBox.image
			c.imgBox.SetImage(canvas.AdjustImageForColor(img, img.Bounds().Max.X, img.Bounds().Max.Y, c.ForeColor()))
		}
	}

	c.Container.Draw(ctx)
}

func (c *Button) DrawBackground(ctx ui.DrawContext) {
	c.drawBackground = true // ??
	if c.hover {
		c.drawBackground = true
	}

	if c.drawBackground {
		c.Container.DrawBackground(ctx)
	}
}

func (c *Button) SetDrawBackground(drawBackground bool) {
	c.drawBackground = drawBackground
}

func (c *Button) rebuildContent() {
	c.RemoveAllWidgets()
	c.txtBlock = nil
	c.imgBox = nil
	c.panelPadding = c.padding

	if c.showText {
		c.txtBlock = NewTextBlock(c, c.text)
		c.txtBlock.TextHAlign = c.textHAlign
		c.txtBlock.TextVAlign = c.textVAlign
		if c.foregroundColor.ValueOwn() != nil {
			c.txtBlock.SetForeColor(c.foregroundColor.Color())
		}
		c.txtBlock.SetEnabled(c.enabled)
	}
	if c.showImage {
		if c.enabled {
			c.imgBox = NewImageBox(c, c.img)
		} else {
			if c.imgDisabled != nil {
				c.imgBox = NewImageBox(c, c.imgDisabled)
			} else {
				c.imgBox = NewImageBox(c, c.img)
			}
		}

		c.imgBox.SetMinWidth(c.imageWidth)
		c.imgBox.SetMinHeight(c.imageHeight)
		c.imgBox.SetScaling(ImageBoxScaleAdjustImageKeepAspectRatio)
	}

	if c.imgBox == nil && c.txtBlock == nil {
		c.txtBlock = NewTextBlock(c, " ")
		c.txtBlock.TextHAlign = c.textHAlign
		c.txtBlock.TextVAlign = c.textVAlign
	}

	txtGridX := 0
	txtGridY := 0
	imgGridX := 0
	imgGridY := 0

	if c.textImageVerticalOrientation {
		if c.imageBeforeText {
			txtGridX = 0
			txtGridY = 1
			imgGridX = 0
			imgGridY = 0
		} else {
			txtGridX = 0
			txtGridY = 0
			imgGridX = 0
			imgGridY = 1
		}
	} else {
		if c.imageBeforeText {
			txtGridX = 1
			txtGridY = 0
			imgGridX = 0
			imgGridY = 0
		} else {
			txtGridX = 0
			txtGridY = 0
			imgGridX = 1
			imgGridY = 0
		}
		vs := NewHSpacer(c)
		c.AddWidgetOnGrid(vs, 2, 0)
	}

	if c.txtBlock != nil {
		c.AddWidgetOnGrid(c.txtBlock, txtGridX, txtGridY)
	}

	if c.imgBox != nil {
		c.AddWidgetOnGrid(c.imgBox, imgGridX, imgGridY)
	}

	c.Window().UpdateLayout()
}

func (c *Button) SetText(text string) {
	c.text = text
	if len(c.text) > 0 {
		c.showText = true
	} else {
		c.showText = false
	}

	c.rebuildContent()
}

func (c *Button) SetShowText(showText bool) {
	c.showText = showText
	c.rebuildContent()
}

func (c *Button) SetShowImage(showImage bool) {
	c.showImage = showImage
	c.rebuildContent()
}

func (c *Button) SetImageSize(width, height int) {
	c.imageWidth = width
	c.imageHeight = height
	c.rebuildContent()
}

func (c *Button) SetImage(img image.Image) {
	c.img = img
	if c.img == nil {
		c.showImage = false
	} else {
		c.showImage = true
	}
	c.rebuildContent()
	c.Window().UpdateLayout()
	c.Update("Button")
}

func (c *Button) SetImageDisabled(img image.Image) {
	c.imgDisabled = img
	if c.img == nil {
		c.showImage = false
	} else {
		c.showImage = true
	}
	c.rebuildContent()
	c.Window().UpdateLayout()
	c.Update("Button")
}

func (c *Button) TextHAlign() canvas.HAlign {
	return c.txtBlock.TextHAlign
}

func (c *Button) SetTextHAlign(textHAlign canvas.HAlign) {
	c.textHAlign = textHAlign
	c.rebuildContent()
}

func (c *Button) TextVAlign() canvas.VAlign {
	return c.txtBlock.TextVAlign
}

func (c *Button) SetTextVAlign(textVAlign canvas.VAlign) {
	c.textVAlign = textVAlign
	c.rebuildContent()
	c.Update("Button")
}

func (c *Button) KeyDown(event *uievents.KeyDownEvent) bool {
	if event.Key == glfw.KeySpace {
		if c.onPress != nil {
			var ev uievents.Event
			ev.Sender = c
			c.onPress(&ev)
			return true
		}
	}
	return false
}

func (c *Button) MouseEnter() {
	c.mousePointerInRect = true
	c.Container.MouseEnter()
}

func (c *Button) MouseLeave() {
	c.mousePointerInRect = false
	c.Container.MouseLeave()
}

func (c *Button) MouseDown(event *uievents.MouseDownEvent) {
	c.pressed = true
	c.Update("Button")
}

func (c *Button) MouseUp(event *uievents.MouseUpEvent) {
	if c.enabled && c.pressed {
		c.pressed = false
		c.Update("Button")
		c.Press()
		event.Ignore = true
	}
}

func (c *Button) Press() {
	if c.onPress != nil {
		var ev uievents.Event
		ev.Sender = c
		c.onPress(&ev)
	}
}

func (c *Button) FindWidgetUnderPointer(x, y int) uiinterfaces.Widget {
	return nil
}

func (c *Button) EnabledChanged(enabled bool) {
	c.rebuildContent()
}
