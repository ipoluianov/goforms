package uicontrols

import (
	"bytes"
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/nfnt/resize"
	"image"
)

type ImageBoxScale int

const (
	ImageBoxScaleNoScaleAdjustBox           ImageBoxScale = 0
	ImageBoxScaleNoScaleImageInLeftTop      ImageBoxScale = 1
	ImageBoxScaleNoScaleImageInCenter       ImageBoxScale = 2
	ImageBoxScaleStretchImage               ImageBoxScale = 3
	ImageBoxScaleAdjustImageKeepAspectRatio ImageBoxScale = 4
)

type ImageBox struct {
	Control
	image   image.Image
	scaling ImageBoxScale
}

func NewImageBox(parent uiinterfaces.Widget, img image.Image) *ImageBox {
	var c ImageBox
	c.InitControl(parent, &c)
	c.image = img
	c.scaling = ImageBoxScaleNoScaleAdjustBox

	if c.image != nil {
		//c.SetMinWidth(c.image.Bounds().Max.X)
		//c.SetMinHeight(c.image.Bounds().Max.Y)
		/*c.SetMaxWidth(c.Image.Bounds().Max.X)
		c.SetMaxHeight(c.Image.Bounds().Max.Y)*/
	}
	return &c
}

func NewImageBoxBytes(parent uiinterfaces.Widget, data []byte) *ImageBox {
	img, _, _ := image.Decode(bytes.NewReader(data))

	var c ImageBox
	c.InitControl(parent, &c)
	c.image = img

	if c.image != nil {
	}
	return &c
}

func (c *ImageBox) SetImage(img image.Image) {
	c.image = img
	c.Window().UpdateLayout()
	c.Update("ImageBox")
}

func (c *ImageBox) SetScaling(scaling ImageBoxScale) {
	c.scaling = scaling
	c.Window().UpdateLayout()
	c.Update("ImageBox")
}

func (c *ImageBox) ControlType() string {
	return "ImageBox"
}

func (c *ImageBox) Draw(ctx ui.DrawContext) {
	if c.image == nil {
		return
	}

	if c.scaling == ImageBoxScaleNoScaleImageInLeftTop || c.scaling == ImageBoxScaleNoScaleAdjustBox {
		b := c.image.Bounds()
		ctx.DrawImage(0, 0, b.Max.X, b.Max.Y, c.image)
	}

	if c.scaling == ImageBoxScaleNoScaleImageInCenter {
		b := c.image.Bounds()
		offsetX := (c.ClientWidth() - b.Max.X) / 2
		offsetY := (c.ClientHeight() - b.Max.Y) / 2
		ctx.DrawImage(offsetX, offsetY, b.Max.X, b.Max.Y, c.image)
	}

	if c.scaling == ImageBoxScaleStretchImage {
		img := resize.Resize(uint(c.ClientWidth()), uint(c.ClientHeight()), c.image, resize.Bicubic)
		ctx.DrawImage(0, 0, c.ClientWidth(), c.ClientHeight(), img)
	}

	if c.scaling == ImageBoxScaleAdjustImageKeepAspectRatio {
		b := c.image.Bounds()
		aspRatioImg := float64(b.Max.X) / float64(b.Max.Y)
		aspRationWidget := float64(c.ClientWidth()) / float64(c.ClientHeight())
		if aspRatioImg > aspRationWidget {
			img := resize.Resize(uint(c.ClientWidth()), 0, c.image, resize.Bicubic)
			b := img.Bounds()
			offsetX := (c.ClientWidth() - b.Max.X) / 2
			offsetY := (c.ClientHeight() - b.Max.Y) / 2
			ctx.DrawImage(offsetX, offsetY, c.ClientWidth(), c.ClientHeight(), img)
		} else {
			img := resize.Resize(0, uint(c.ClientHeight()), c.image, resize.Bicubic)
			b := img.Bounds()
			offsetX := (c.ClientWidth() - b.Max.X) / 2
			offsetY := (c.ClientHeight() - b.Max.Y) / 2
			ctx.DrawImage(offsetX, offsetY, c.ClientWidth(), c.ClientHeight(), img)
		}
	}
}

func (c *ImageBox) MinWidth() int {
	if c.scaling == ImageBoxScaleNoScaleAdjustBox {
		if c.image == nil {
			return c.LeftBorderWidth() + c.RightBorderWidth()
		}
		return c.image.Bounds().Max.X + c.LeftBorderWidth() + c.RightBorderWidth()
	}
	return c.Control.MinWidth()
}

func (c *ImageBox) MinHeight() int {
	if c.scaling == ImageBoxScaleNoScaleAdjustBox {
		if c.image == nil {
			return c.TopBorderWidth() + c.BottomBorderWidth()
		}
		return c.image.Bounds().Max.Y + c.TopBorderWidth() + c.BottomBorderWidth()
	}
	return c.Control.MinHeight()
}

func (c *ImageBox) MaxWidth() int {
	if c.scaling == ImageBoxScaleNoScaleAdjustBox {
		if c.image == nil {
			return c.LeftBorderWidth() + c.RightBorderWidth()
		}
		return c.image.Bounds().Max.X + c.LeftBorderWidth() + c.RightBorderWidth()
	}
	return c.Control.MaxWidth()
}

func (c *ImageBox) MaxHeight() int {
	if c.scaling == ImageBoxScaleNoScaleAdjustBox {
		if c.image == nil {
			return c.TopBorderWidth() + c.BottomBorderWidth()
		}
		return c.image.Bounds().Max.Y + c.TopBorderWidth() + c.BottomBorderWidth()
	}
	return c.Control.MaxHeight()
}
