package ui

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/nui/nui"
)

func init() {
}

type DrawContextSW struct {
	DrawContextBase
	cnv *canvas.CanvasDirect
	cc  *gg.Context
}

func NewDrawContextSW(window nui.Window) *DrawContextSW {
	var c DrawContextSW
	c.Window = window
	c.WindowWidth, c.WindowHeight = window.Size()
	/*kx, ky := window.GetContentScale()
	c.WindowWidthScale = kx
	c.WindowHeightScale = ky*/

	c.DrawContextBase.InitBase()

	c.cnv = canvas.NewCanvasWindow(int(float32(c.WindowWidth)), int(float32(c.WindowHeight)))

	c.setViewport()
	return &c
}

func NewDrawContextSWRGBA(window nui.Window, rgba *image.RGBA) *DrawContextSW {
	var c DrawContextSW
	c.Window = window
	c.WindowWidth, c.WindowHeight = window.Size()

	c.DrawContextBase.InitBase()

	c.cnv = canvas.NewCanvasRGBA(rgba)

	c.setViewport()
	return &c
}

func NewDrawContextSWSpecial(width, height int) *DrawContextSW {
	var c DrawContextSW
	c.Window = nil
	c.WindowWidth, c.WindowHeight = width, height

	c.DrawContextBase.InitBase()

	c.cnv = canvas.NewCanvasWindow(width, height)
	c.cc = gg.NewContextForRGBA(c.GraphContextImage())

	c.setViewport()
	return &c
}

func (c *DrawContextSW) Init() {
	/*c.DrawContextBase.InitBase()
	//c.Window.MakeContextCurrent()
	c.setViewport()
	gl.Clear(gl.DEPTH_BUFFER_BIT)
	c.cc = gg.NewContextForRGBA(c.GraphContextImage())*/
	//gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (c *DrawContextSW) GG() *gg.Context {
	return c.cc
}

func (c *DrawContextSW) GraphContextImage() *image.RGBA {
	return c.cnv.Image()
}

func (c *DrawContextSW) setViewport() {
	//gl.Viewport(int32(c.CurrentClipSettings.x), int32(c.WindowHeight-c.CurrentClipSettings.y-c.CurrentClipSettings.height), int32(c.CurrentClipSettings.width)*int32(c.WindowWidthScale), int32(c.CurrentClipSettings.height*int(c.WindowHeightScale)))
}

/*func (c *DrawContextSW) xyTo01(x, y int) (float32, float32) {
	dx := 2 / float64(c.CurrentClipSettings.width)
	dy := 2 / float64(c.CurrentClipSettings.height)
	x1d := float32(dx*float64(x)) - 1
	y1d := float32(dy*float64(y)) - 1
	y1d = -y1d // OpenGL coordinates are inverted by Y
	return x1d, y1d
}*/

func (c *DrawContextSW) DrawRect(x, y, width, height int) {
	c.cnv.DrawRect(x, y, width, height, c.CurrentColor, c.StrokeWidth)
}

func (c *DrawContextSW) FillRect(x, y, width, height int) {
	c.cnv.FillRect(x, y, width, height, c.CurrentColor)
}

func (c *DrawContextSW) DrawImage(x, y, width, height int, img image.Image) {
	c.cnv.DrawImage(x, y, img)
}

func (c *DrawContextSW) DrawText(x, y, width, height int, text string) {
	c.cnv.DrawTextMultiline(x, y, width, height, c.TextHAlign, c.TextVAlign, text, c.CurrentColor, c.FontFamily, c.FontSize, c.UnderLine)
}

func (c *DrawContextSW) DrawLine(x1, y1, x2, y2 int) {
	c.cnv.DrawLine(x1, y1, x2, y2, c.StrokeWidth, c.CurrentColor)
}

func (c *DrawContextSW) DrawEllipse(x, y, width, height int) {
}

func (c *DrawContextSW) MeasureText(text string) (int, int) {
	w, h, _ := canvas.MeasureText(c.FontFamily, c.FontSize, false, false, text, true)
	return w, h
}

func (c *DrawContextSW) Save() {
	c.cnv.Save()
}

func (c *DrawContextSW) Translate(x, y int) {
	c.cnv.Translate(x, y)
}

func (c *DrawContextSW) TranslatedX() int {
	return c.State().TranslateX
}

func (c *DrawContextSW) TranslatedY() int {
	return c.State().TranslateY
}

func (c *DrawContextSW) ClippedRegion() (int, int, int, int) {
	return c.cnv.ClippedRegion()
}

func (c *DrawContextSW) ClipIn(x, y, width, height int) {
	c.cnv.ClipInTranslated(x, y, width, height)
}

/*func (c *DrawContextSW) Clip(x, y, width, height int) {
	c.cnv.ClipInTranslated(x, y, width, height)
}*/

func (c *DrawContextSW) State() canvas.CanvasDirectState {
	return c.cnv.State()
}

func (c *DrawContextSW) Load() {
	c.cnv.Load()
}

func (c *DrawContextSW) Finish(rgba *image.RGBA) {
}
