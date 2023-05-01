package uicontrols

import (
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiresources"

	"github.com/nfnt/resize"
	"image"
	"time"
)

type PopupMenuItem struct {
	Control
	text                 string
	OnClick              func(event *uievents.Event)
	Image                image.Image
	ImageResource        []byte
	KeyCombination       string
	parentMenu           *PopupMenu
	needToClosePopupMenu func()

	timerShowInnerMenu  *uievents.FormTimer
	AdjustColorForImage bool

	innerMenu *PopupMenu
}

func (c *PopupMenuItem) OnInit() {
	c.timerShowInnerMenu = c.OwnWindow.NewTimer(200, c.timerShowInnerMenuHandler)
	c.AdjustColorForImage = true
}

func (c *PopupMenuItem) SetText(text string) {
	c.text = text
	c.Update("PopupMenuItem")
}

func (c *PopupMenuItem) Dispose() {
	if c.OwnWindow != nil {
		c.OwnWindow.RemoveTimer(c.timerShowInnerMenu)
	}

	if c.innerMenu != nil {
		c.innerMenu.Dispose()
		c.innerMenu = nil
	}

	c.parentMenu = nil
	c.OnClick = nil
	c.needToClosePopupMenu = nil
	c.timerShowInnerMenu = nil

	c.Control.Dispose()
}

func (c *PopupMenuItem) ControlType() string {
	return "PopupMenuItem"
}

func (c *PopupMenuItem) Draw(ctx ui.DrawContext) {
	ctx.SetColor(c.backgroundColor.Color())
	ctx.FillRect(0, 0, c.InnerWidth(), c.InnerHeight())

	xOffset := 0
	if c.Image != nil || c.ImageResource != nil {
		imageSource := c.Image
		if c.ImageResource != nil {
			imageSource = uiresources.ResImgCol(c.ImageResource, c.ForeColor())
		}

		img := resize.Resize(24, 24, imageSource, resize.Bicubic)
		if c.AdjustColorForImage {
			img = canvas.AdjustImageForColor(img, img.Bounds().Max.X, img.Bounds().Max.Y, c.foregroundColor.Color())
		}

		height := img.Bounds().Max.Y
		yOffset := (c.Height() - height) / 2

		ctx.DrawImage(3, yOffset, img.Bounds().Max.X, height, img)
		xOffset += 32
	}

	ctx.SetColor(c.foregroundColor.Color())
	ctx.SetFontSize(c.fontSize.Float64())
	textWidth := c.InnerWidth()
	if c.innerMenu != nil {
		textWidth -= c.InnerHeight()
	}
	ctx.DrawText(xOffset+5, 0, textWidth, c.InnerHeight(), c.text)
	if c.innerMenu != nil {
		imgArrow := uiresources.ResImgCol(uiresources.R_icons_material4_png_av_play_arrow_materialicons_48dp_1x_baseline_play_arrow_black_48dp_png, c.ForeColor())
		ctx.DrawImage(c.InnerWidth()-c.InnerHeight(), 0, imgArrow.Bounds().Max.X, imgArrow.Bounds().Max.Y, resize.Resize(uint(c.InnerHeight()), uint(c.InnerHeight()), imgArrow, resize.Bicubic))
	}
}

func (c *PopupMenuItem) MouseClick(event *uievents.MouseClickEvent) {
	c.timerShowInnerMenu.Enabled = false

	if c.innerMenu != nil {
		x, y := c.RectClientAreaOnWindow()
		w := c.Width()
		c.innerMenu.showMenu(x+w, y, c.parentMenu)
		return
	}

	if c.needToClosePopupMenu != nil {
		c.needToClosePopupMenu()
	}

	if c.OnClick != nil {
		var ev uievents.Event
		ev.Sender = c
		c.OnClick(&ev)
	}
}

func (c *PopupMenuItem) MouseEnter() {
	c.OwnWindow.CloseAfterPopupWidget(c.parentMenu)

	if c.innerMenu != nil {
		c.timerShowInnerMenu.Enabled = true
		c.timerShowInnerMenu.LastElapsedDTMSec = time.Now().UnixNano() / 1000000
		return
	}

}

func (c *PopupMenuItem) MouseLeave() {
	c.timerShowInnerMenu.Enabled = false
}

func (c *PopupMenuItem) MouseMove(event *uievents.MouseMoveEvent) {
	c.Update("PopupMenuItem")
}

func (c *PopupMenuItem) SetInnerMenu(menu *PopupMenu) {
	c.innerMenu = menu
}

func (c *PopupMenuItem) timerShowInnerMenuHandler() {
	if c.timerShowInnerMenu.Enabled {
		c.timerShowInnerMenu.Enabled = false

		x, y := c.RectClientAreaOnWindow()
		w := c.Width()
		c.innerMenu.showMenu(x+w, y, c.parentMenu)
	}
}
