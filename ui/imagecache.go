package ui

import (
	"fmt"
	"runtime"

	"github.com/ipoluianov/goforms/canvas"
	"github.com/ipoluianov/goforms/utils"
)

type ImageCache struct {
	utils.Obj
	images map[string]*canvas.CanvasDirect
}

func NewImageCache(name string) *ImageCache {
	var cObj ImageCache
	c := &cObj
	c.Obj.InitObj("ImageCache", name)
	runtime.SetFinalizer(c, finalizerImageCache)
	return c
}

func finalizerImageCache(c *ImageCache) {
	c.Obj.UninitObj()
}

func (c *ImageCache) GetXY(x, y int) *canvas.CanvasDirect {
	key := fmt.Sprint("XY", x, "-", y)
	if canvas, ok := c.images[key]; ok {
		return canvas
	}
	return nil
}

func (c *ImageCache) SetXY(x, y int, cnv *canvas.CanvasDirect) {
	if c.images == nil {
		c.images = make(map[string]*canvas.CanvasDirect)
	}
	key := fmt.Sprint("XY", x, "-", y)
	c.images[key] = cnv
}

func (c *ImageCache) ClearXY(x, y int) {
	key := fmt.Sprint("XY", x, "-", y)
	c.Obj.SetStatistics(fmt.Sprint("Count ", len(c.images)))
	delete(c.images, key)
}

func (c *ImageCache) Clear() {
	c.images = make(map[string]*canvas.CanvasDirect)
}
