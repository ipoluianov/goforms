package ui

import (
	"github.com/fogleman/gg"
	"github.com/gazercloud/gazerui/canvas"
	"github.com/gazercloud/gazerui/opengl/gl11/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"image/draw"
	"unsafe"
)

var texturesMap map[int]*image.RGBA

func init() {
	texturesMap = make(map[int]*image.RGBA)
}

type DrawContextSW struct {
	DrawContextBase
	cnv *canvas.CanvasDirect
	cc  *gg.Context
}

func NewDrawContextSW(window *glfw.Window) *DrawContextSW {
	var c DrawContextSW
	c.Window = window
	c.WindowWidth, c.WindowHeight = window.GetSize()

	c.DrawContextBase.InitBase()

	c.cnv = canvas.NewCanvasWindow(c.WindowWidth, c.WindowHeight)

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
	c.DrawContextBase.InitBase()
	//c.Window.MakeContextCurrent()
	c.setViewport()
	gl.Clear(gl.DEPTH_BUFFER_BIT)
	c.cc = gg.NewContextForRGBA(c.GraphContextImage())
	//gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (c *DrawContextSW) GG() *gg.Context {
	return c.cc
}

func (c *DrawContextSW) GraphContextImage() *image.RGBA {
	return c.cnv.Image()
}

func (c *DrawContextSW) setViewport() {
	gl.Viewport(int32(c.CurrentClipSettings.x), int32(c.WindowHeight-c.CurrentClipSettings.y-c.CurrentClipSettings.height), int32(c.CurrentClipSettings.width), int32(c.CurrentClipSettings.height))
}

func (c *DrawContextSW) xyTo01(x, y int) (float32, float32) {
	dx := 2 / float64(c.CurrentClipSettings.width)
	dy := 2 / float64(c.CurrentClipSettings.height)
	x1d := float32(dx*float64(x)) - 1
	y1d := float32(dy*float64(y)) - 1
	y1d = -y1d // OpenGL coordinates are inverted by Y
	return x1d, y1d
}

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

func (c *DrawContextSW) drawCanvasRect(img image.Image, x, y int, textureSize int) {

	p1x, p1y := c.xyTo01(x, y)
	p2x, p2y := c.xyTo01(x+textureSize, y+textureSize)

	var rgba *image.RGBA
	var ok bool
	if rgba, ok = texturesMap[textureSize]; !ok {
		rgba = image.NewRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{textureSize, textureSize}})
		texturesMap[textureSize] = rgba
	} else {
		rgba = texturesMap[textureSize]
	}
	draw.Draw(rgba, img.Bounds(), img, image.Point{x, y}, draw.Src)

	var texid []uint32
	texid = make([]uint32, 1)
	texid[0] = 0xFFFFFFFF

	//fmt.Println("123131231", p1x, p1y, p2x, p2y)

	// Init
	gl.ShadeModel(gl.FLAT)
	gl.Enable(gl.DEPTH_TEST)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.GenTextures(1, &texid[0])
	gl.BindTexture(gl.TEXTURE_2D, texid[0])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(textureSize), int32(textureSize), 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&rgba.Pix[0]))

	// Draw
	gl.Enable(gl.TEXTURE_2D)

	//gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.DECAL)
	//gl.BindTexture(gl.TEXTURE_2D, texid[0])

	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0.0, 0.0)
	gl.Vertex3f(p1x, p1y, 0.0)
	gl.TexCoord2f(0.0, 1.0)
	gl.Vertex3f(p1x, p2y, 0.0)
	gl.TexCoord2f(1.0, 1.0)
	gl.Vertex3f(p2x, p2y, 0.0)
	gl.TexCoord2f(1.0, 0.0)
	gl.Vertex3f(p2x, p1y, 0.0)
	gl.End()

	gl.Flush()

	//gl.Disable(gl.BLEND)
	//gl.DepthMask(true)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.TEXTURE_2D)

	gl.DeleteTextures(1, &texid[0])
}

func (c *DrawContextSW) Finish() {
	img := c.cnv.Image()

	var imgScaled image.Image
	imgScaled = img
	/*winWidth, winHeight := c.Window.GetSize()
	if img.Bounds().Max.X != winWidth || img.Bounds().Max.Y != winHeight {
		imgScaled = resize.Resize(uint(winWidth), uint(winHeight), img, resize.Lanczos3)
	}*/

	var maxTextSize int32
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &maxTextSize)

	tw := imgScaled.Bounds().Max.X
	th := imgScaled.Bounds().Max.Y

	textureSize := 1024

	rectsByX := (tw / textureSize) + 1
	rectsByY := (th / textureSize) + 1

	for x := 0; x < rectsByX; x++ {
		for y := 0; y < rectsByY; y++ {
			c.drawCanvasRect(imgScaled, x*textureSize, y*textureSize, textureSize)
		}
	}
}
