package canvas

import (
	"flag"
	"fmt"
	"github.com/gazercloud/gazerui/grid/stats"
	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"math"
	"runtime"
	"strings"
)

type CanvasDirect struct {
	stats.Obj
	image       *image.RGBA
	state       CanvasDirectState
	savedStates []CanvasDirectState
}

type VAlign int
type HAlign int

const VAlignTop VAlign = 0
const VAlignCenter VAlign = 1
const VAlignBottom VAlign = 2

const HAlignLeft HAlign = 0
const HAlignCenter HAlign = 1
const HAlignRight HAlign = 2

type CanvasDirectState struct {
	TranslateX int
	TranslateY int

	ClipX      int
	ClipY      int
	ClipWidth  int
	ClipHeight int
}

func (c *CanvasDirectState) String() string {
	return fmt.Sprint("tr[", c.TranslateX, ":", c.TranslateY, "] clip[", c.ClipX, c.ClipY, c.ClipWidth, c.ClipHeight, "]")
}

var dpi = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")

var staticImageBuffer1 []uint8
var staticImageBuffer2 []uint8
var currentBufferIndex int

func init() {
	staticImageBuffer1 = make([]uint8, 1024*1024*10)
	staticImageBuffer2 = make([]uint8, 1024*1024*10)
}

func NewCanvas(width int, height int) *CanvasDirect {
	var canvas CanvasDirect
	canvas.image = image.NewRGBA(image.Rect(0, 0, width, height))
	canvas.state.ClipX = 0
	canvas.state.ClipY = 0
	canvas.state.ClipWidth = 10000000
	canvas.state.ClipHeight = 10000000
	canvas.Obj.InitObj("CanvasDirect", fmt.Sprint("W:", width, " H:", height))
	runtime.SetFinalizer(&canvas, finalizerCanvas)
	return &canvas
}

func NewCanvasWindow(width int, height int) *CanvasDirect {
	var canvas CanvasDirect
	//canvas.image = image.NewRGBA(image.Rect(0, 0, width, height))
	//staticImageBuffer = make([]uint8, 4*width*height)

	var currentBuffer []uint8

	if currentBufferIndex == 0 {
		currentBufferIndex = 1
		currentBuffer = staticImageBuffer1
	} else {
		currentBufferIndex = 0
		currentBuffer = staticImageBuffer2
	}

	for i := range currentBuffer {
		currentBuffer[i] = 0
	}

	canvas.image = &image.RGBA{
		Pix:    currentBuffer,
		Stride: 4 * width,
		Rect:   image.Rect(0, 0, width, height),
	}
	//canvas.image = image.NewRGBA(image.Rect(0, 0, 10, 10))
	canvas.state.ClipX = 0
	canvas.state.ClipY = 0
	canvas.state.ClipWidth = 10000000
	canvas.state.ClipHeight = 10000000
	canvas.Obj.InitObj("CanvasDirect", fmt.Sprint("W:", width, " H:", height))
	runtime.SetFinalizer(&canvas, finalizerCanvas)
	return &canvas
}

func finalizerCanvas(c *CanvasDirect) {
	c.Obj.UninitObj()
}

func (c *CanvasDirect) Image() *image.RGBA {
	return c.image
}

func (c *CanvasDirect) State() CanvasDirectState {
	return c.state
}

func (c *CanvasDirect) Save() {
	c.savedStates = append(c.savedStates, c.state)
}

func (c *CanvasDirect) Load() {
	if len(c.savedStates) > 0 {
		lenOfStack := len(c.savedStates)
		c.state = c.savedStates[lenOfStack-1]
		c.savedStates = c.savedStates[:lenOfStack-1]
	}
}

func (c *CanvasDirect) Translate(x, y int) {
	c.state.TranslateX += x
	c.state.TranslateY += y
}

func (c *CanvasDirect) TranslatedX() int {
	return c.state.TranslateX
}

func (c *CanvasDirect) TranslatedY() int {
	return c.state.TranslateY
}

/*func (c *CanvasDirect) Clip(x, y, width, height int) {
	c.state.ClipX = x
	c.state.ClipY = y
	c.state.ClipWidth = width
	c.state.ClipHeight = height
}*/

func (c *CanvasDirect) ClipIn(x, y, width, height int) {
	nClipX := x
	nClipY := y
	nClipW := width
	nClipH := height

	if nClipX < c.state.ClipX {
		nClipW -= c.state.ClipX - nClipX
		nClipX = c.state.ClipX
	}

	if nClipX+nClipW > c.state.ClipX+c.state.ClipWidth {
		nClipW -= (nClipX + nClipW) - (c.state.ClipX + c.state.ClipWidth)
	}

	if nClipY < c.state.ClipY {
		nClipH -= c.state.ClipY - nClipY
		nClipY = c.state.ClipY
	}

	if nClipY+nClipH > c.state.ClipY+c.state.ClipHeight {
		nClipH -= (nClipY + nClipH) - (c.state.ClipY + c.state.ClipHeight)
	}

	if nClipW < 0 {
		nClipW = 0
	}

	if nClipH < 0 {
		nClipH = 0
	}

	c.state.ClipX = nClipX
	c.state.ClipY = nClipY
	c.state.ClipWidth = nClipW
	c.state.ClipHeight = nClipH
}

/*func (c *CanvasDirect) ClipTranslated(x, y, width, height int) {
	c.state.ClipX = c.state.TranslateX + x
	c.state.ClipY = c.state.TranslateY + y
	c.state.ClipWidth = width
	c.state.ClipHeight = height
}*/

func (c *CanvasDirect) ClipInTranslated(x, y, width, height int) {
	c.ClipIn(c.state.TranslateX+x, c.state.TranslateY+y, width, height)
}

func (c *CanvasDirect) ClippedRegion() (int, int, int, int) {
	return c.state.ClipX, c.state.ClipY, c.state.ClipX + c.state.ClipWidth, c.state.ClipY + c.state.ClipHeight
}

func (c *CanvasDirect) DrawPoint(x int, y int, color color.Color) {
	if x < c.state.ClipX || x > c.state.ClipX+c.state.ClipWidth {
		return
	}
	if y < c.state.ClipY || y > c.state.ClipY+c.state.ClipHeight {
		return
	}

	c.image.Set(x, y, color)
}

func (c *CanvasDirect) MixPixel(x int, y int, rgba color.Color) {

	if x < c.state.ClipX || x > c.state.ClipX+c.state.ClipWidth {
		return
	}
	if y < c.state.ClipY || y > c.state.ClipY+c.state.ClipHeight {
		return
	}

	cOld := c.image.At(x, y)
	oR, oG, oB, _ := cOld.RGBA()
	cR, cG, cB, cA := rgba.RGBA()
	cR = cR >> 8
	cG = cG >> 8
	cB = cB >> 8
	cA = cA >> 8

	oR = oR >> 8
	oG = oG >> 8
	oB = oB >> 8

	alpha := uint32(cA)
	antialpha := 255 - alpha

	if cA > 0 && cA < 255 {
		cA += 1
	}

	//alpha, antialpha = antialpha, alpha
	nR := uint8(((uint32(oR) * antialpha) >> 8) + ((cR * alpha) >> 8))
	nG := uint8(((uint32(oG) * antialpha) >> 8) + ((cG * alpha) >> 8))
	nB := uint8(((uint32(oB) * antialpha) >> 8) + ((cB * alpha) >> 8))

	c.image.SetRGBA(x, y, color.RGBA{nR, nG, nB, 255})
}

func (c *CanvasDirect) DrawLine(x1 int, y1 int, x2 int, y2 int, width int, color color.Color) {
	x1 = x1 + c.state.TranslateX
	y1 = y1 + c.state.TranslateY
	x2 = x2 + c.state.TranslateX
	y2 = y2 + c.state.TranslateY

	x1, y1, x2, y2, visible := CohenSutherland(x1, y1, x2, y2, c.state.ClipX, c.state.ClipY, c.state.ClipX+c.state.ClipWidth, c.state.ClipY+c.state.ClipHeight)
	if visible {
		if (x1 == x2 || y1 == y2) && width == 1 {
			if x1 == x2 {
				if y1 > y2 {
					y1, y2 = y2, y1
				}
				for y := y1; y < y2; y++ {
					c.MixPixel(x1, y, color)
				}
			}
			if y1 == y2 {
				if x1 > x2 {
					x1, x2 = x2, x1
				}
				for x := x1; x < x2; x++ {
					c.MixPixel(x, y1, color)
				}
			}
		} else {
			script := c.MakeScriptLine(float64(x1), float64(y1), float64(x2), float64(y2), float64(width))
			script.Bounds = image.Rectangle{Min: image.Point{X: c.state.ClipX, Y: c.state.ClipY}, Max: image.Point{X: c.state.ClipX + c.state.ClipWidth, Y: c.state.ClipY + c.state.ClipHeight}}
			script.DrawToRGBA(c.image, color)
		}
	}
}

func (c *CanvasDirect) DrawRect(x int, y int, width int, height int, color color.Color, strokeWidth int) {

	if strokeWidth == 1 {
		x = x + c.state.TranslateX
		y = y + c.state.TranslateY

		if width > 0 {
			for yy := 0; yy < height; yy++ {
				c.DrawPoint(x, yy+y, color)
				c.DrawPoint(x+width-1, yy+y, color)
			}
		}

		if height > 0 {
			for xx := 0; xx < width; xx++ {
				c.DrawPoint(xx+x, y, color)
				c.DrawPoint(xx+x, y+height-1, color)
			}
		}
	} else {

		/*scr := c.MakeScriptLine(float64(x-strokeWidth/2), float64(y), float64(x+width-1+strokeWidth/2), float64(y), float64(strokeWidth))
		scr.append(c.MakeScriptLine(float64(x+width-1), float64(y), float64(x+width-1), float64(y+height-1), float64(strokeWidth)))
		scr.append(c.MakeScriptLine(float64(x+width-1+strokeWidth/2), float64(y+height-1), float64(x-strokeWidth/2), float64(y+height-1), float64(strokeWidth)))
		scr.append(c.MakeScriptLine(float64(x), float64(y+height-1), float64(x), float64(y), float64(strokeWidth)))
		scr.Bounds = image.Rectangle{Min: image.Point{X: c.state.ClipX, Y: c.state.ClipY}, Max: image.Point{X: c.state.ClipX + c.state.ClipWidth, Y: c.state.ClipY + c.state.ClipHeight}}
		scr.DrawToRGBA(c.image, color)*/

		c.FillRect(x-strokeWidth/2, y-strokeWidth/2, width+strokeWidth, strokeWidth, color)        // top
		c.FillRect(x-strokeWidth/2, y+height-strokeWidth/2, width+strokeWidth, strokeWidth, color) // bottom

		c.FillRect(x-strokeWidth/2, y-strokeWidth/2, strokeWidth, height+strokeWidth, color)       // left
		c.FillRect(x+width-strokeWidth/2, y-strokeWidth/2, strokeWidth, height+strokeWidth, color) // right
	}
}

func (c *CanvasDirect) MakeScriptLineBresenham(x1 float64, y1 float64, x2 float64, y2 float64) *DrawScript {
	x1 = math.Round(x1)
	y1 = math.Round(y1)
	x2 = math.Round(x2)
	y2 = math.Round(y2)

	result := NewDrawScript()

	stepByX := true
	if math.Abs(x2-x1) < math.Abs(y2-y1) {
		stepByX = false
	}

	if stepByX {
		if x2 < x1 { // Swap
			xo := x1
			x1 = x2
			x2 = xo
			yo := y1
			y1 = y2
			y2 = yo
		}

		beginX := int(x1)
		endX := int(x2)

		deltaX := math.Abs(x2 - x1)
		deltaY := math.Abs(y2 - y1)

		y := int(y1)
		var err float64 = 0
		errDelta := deltaY / deltaX
		var yDir = 1
		if (y2 - y1) < 0 {
			yDir = -1
		}

		for x := beginX; x <= endX; x++ {
			result.plot(x, y, 1)
			err += errDelta
			if err > 0.5 {
				y += yDir
				err -= 1.0
			}
		}
	} else {
		if y2 < y1 {
			xo := x1
			x1 = x2
			x2 = xo
			yo := y1
			y1 = y2
			y2 = yo
		}

		beginY := int(y1)
		endY := int(y2)

		deltaX := math.Abs(x2 - x1)
		deltaY := math.Abs(y2 - y1)

		x := int(x1)
		err := 0.0
		errDelta := deltaX / deltaY
		xDir := 1
		if (x2 - x1) < 0 {
			xDir = -1
		}

		for y := beginY; y <= endY; y++ {
			result.plot(x, y, 1)
			err += errDelta
			if err > 0.5 {
				x += xDir
				err -= 1.0
			}
		}
	}

	return result
}

func ipart(x float64) int {
	rn := math.Floor(x)
	return int(rn)
}

func round(x float64) int {
	iprt := ipart(x + 0.5)
	return iprt
}

func fpart(x float64) float64 {
	return x - float64(ipart(x))
}

func (c *CanvasDirect) MakeScriptLineWu(x1, y1, x2, y2 float64) *DrawScript {

	result := NewDrawScript()

	x1 = math.Round(x1)
	y1 = math.Round(y1)
	x2 = math.Round(x2)
	y2 = math.Round(y2)

	if math.Abs(x1-x2) < 0.1 && math.Abs(y1-y2) < 0.1 {
		return result
	}

	stepByX := true
	if math.Abs(x2-x1) < math.Abs(y2-y1) {
		stepByX = false
	}

	if stepByX {
		if x2 < x1 {
			xo := x1
			x1 = x2
			x2 = xo
			yo := y1
			y1 = y2
			y2 = yo
		}

		dx := x2 - x1
		dy := y2 - y1
		gradient := dy / dx
		offset := y1 + gradient*(float64(round(x1))-x1)

		workFrom := int(x1)
		workTo := int(x2)

		xgap := 1 - fpart(x1+0.5)
		result.plot(workFrom, int(math.Floor(y1)), (1-fpart(offset))*xgap)
		result.plot(workFrom, int(math.Floor(y1))+1, fpart(offset)*xgap)

		xgap = fpart(x1 + 0.5)
		result.plot(workTo, int(math.Floor(y2)), (1-fpart(offset))*xgap)
		result.plot(workTo, int(math.Floor(y2))+1, fpart(offset)*xgap)

		offset += gradient

		for workIndex := workFrom + 1; workIndex <= workTo-1; workIndex++ {
			err := fpart(offset)
			if err > 0 {
				c1 := 1 - err
				c2 := err
				result.plot(workIndex, ipart(offset), c1)
				result.plot(workIndex, ipart(offset)+1, c2)
			} else {
				c1 := 1 - math.Abs(err)
				c2 := math.Abs(err)
				result.plot(workIndex, ipart(offset), c1)
				result.plot(workIndex, ipart(offset)-1, c2)
			}
			offset = offset + gradient
		}
	} else {
		if y2 < y1 {
			xo := x1
			x1 = x2
			x2 = xo
			yo := y1
			y1 = y2
			y2 = yo
		}

		dx := x2 - x1
		dy := y2 - y1

		gradient := dx / dy

		offset := x1 + gradient*(float64(round(y1))-y1)

		workFrom := int(y1)
		workTo := int(y2)

		xgap := 1 - fpart(y1+0.5)
		result.plot(int(math.Floor(x1)), workFrom, (1-fpart(offset))*xgap)
		result.plot(int(math.Floor(x1))+1, workFrom, fpart(offset)*xgap)

		xgap = fpart(y1 + 0.5)
		result.plot(int(math.Floor(x2)), workTo, (1-fpart(offset))*xgap)
		result.plot(int(math.Floor(x2))+1, workTo, fpart(offset)*xgap)

		offset += gradient

		for workIndex := workFrom + 1; workIndex <= workTo-1; workIndex++ {
			err := fpart(offset)
			if err > 0 {
				c1 := 1 - err
				c2 := err
				result.plot(ipart(offset), workIndex, c1)
				result.plot(ipart(offset)+1, workIndex, c2)
			} else {
				c1 := 1 - math.Abs(err)
				c2 := math.Abs(err)
				result.plot(ipart(offset), workIndex, c1)
				result.plot(ipart(offset)-1, workIndex, c2)
			}
			offset = offset + gradient
		}
	}

	return result
}

func (c *CanvasDirect) DrawImage(x int, y int, img image.Image) {
	bInner := image.Rectangle{}
	bInner.Min.X = x + c.state.TranslateX
	bInner.Min.Y = y + c.state.TranslateY
	bInner.Max.X += bInner.Min.X + img.Bounds().Max.X
	bInner.Max.Y += bInner.Min.Y + img.Bounds().Max.Y

	xOffset := 0
	yOffset := 0

	if bInner.Min.X < c.state.ClipX {
		xOffset = c.state.ClipX - bInner.Min.X
		bInner.Min.X = c.state.ClipX
	}

	if bInner.Min.Y < c.state.ClipY {
		yOffset = c.state.ClipY - bInner.Min.Y
		bInner.Min.Y = c.state.ClipY
	}

	if bInner.Max.X > c.state.ClipX+c.state.ClipWidth {
		bInner.Max.X = c.state.ClipX + c.state.ClipWidth
	}

	if bInner.Max.Y > c.state.ClipY+c.state.ClipHeight {
		bInner.Max.Y = c.state.ClipY + c.state.ClipHeight
	}

	if bInner.Min.X >= bInner.Max.X {
		return
	}

	if bInner.Min.Y >= bInner.Max.Y {
		return
	}

	draw.Draw(c.image, bInner, img, image.Point{xOffset, yOffset}, draw.Over)

	/*maxY := img.Bounds().Max.Y
	maxX := img.Bounds().Max.X

	for yy := yOffset; yy < maxY; yy++ {
		for xxx := xOffset; xxx < maxX; xxx++ {
			c.image.Set(xxx + bInner.Min.X, yy + bInner.Min.Y, img.At(xxx, yy))
		}
	}*/

}

func AdjustImageForColor(mask image.Image, width int, height int, col color.Color) image.Image {
	if mask == nil {
		return nil
	}
	img := resize.Resize(uint(width), uint(height), mask, resize.Bicubic)

	// Src
	src := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	draw.Draw(src, src.Bounds(), &image.Uniform{col}, image.ZP, draw.Src)

	// Dest
	dest := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})

	draw.DrawMask(dest, dest.Bounds(), src, image.ZP, img, image.ZP, draw.Over)
	return dest
}

func (c *CanvasDirect) FillEntire(colr color.Color) {
	line := make([]uint8, c.image.Stride)
	cc := color.RGBAModel.Convert(colr).(color.RGBA)
	for xx := 0; xx < c.image.Stride; xx += 4 {
		line[xx] = cc.R
		line[xx+1] = cc.G
		line[xx+2] = cc.B
		line[xx+3] = cc.A
	}

	for yy := 0; yy < c.image.Rect.Max.Y; yy++ {
		offset := yy * c.image.Stride
		copy(c.image.Pix[offset:offset+c.image.Stride], line[:])
	}
}

func (c *CanvasDirect) FillRect(x int, y int, width int, height int, colr color.Color) {

	x = x + c.state.TranslateX
	y = y + c.state.TranslateY

	if x < 0 {
		width += x
		x = 0
	}

	if y < 0 {
		height += y
		y = 0
	}

	if x+width > c.image.Rect.Max.X {
		width = c.image.Rect.Max.X - x
	}

	if y+height > c.image.Rect.Max.Y {
		height = c.image.Rect.Max.Y - y
	}

	if width < 0 {
		return
	}

	if height < 0 {
		return
	}

	if x < c.state.ClipX {
		width -= c.state.ClipX - x
		x = c.state.ClipX
	}

	if y < c.state.ClipY {
		height -= c.state.ClipY - y
		y = c.state.ClipY
	}

	if x+width > c.state.ClipX+c.state.ClipWidth {
		width = c.state.ClipX + c.state.ClipWidth - x
	}

	if y+height > c.state.ClipY+c.state.ClipHeight {
		height = c.state.ClipY + c.state.ClipHeight - y
	}

	if width < 1 {
		return
	}

	if height < 1 {
		return
	}

	if _, _, _, a := colr.RGBA(); a == 65535 {
		line := make([]uint8, width*4)
		cc := color.RGBAModel.Convert(colr).(color.RGBA)
		for xx := 0; xx < width*4; xx += 4 {
			line[xx] = cc.R
			line[xx+1] = cc.G
			line[xx+2] = cc.B
			line[xx+3] = cc.A
		}

		for yy := y; yy < y+height; yy++ {
			offset := yy*c.image.Stride + x*4
			copy(c.image.Pix[offset:offset+width*4], line[:])
		}
	} else {
		for yy := y; yy < y+height; yy++ {
			for xx := x; xx < x+width; xx++ {
				c.MixPixel(xx, yy, colr)
			}
		}
	}
}

func (c *CanvasDirect) fillEx(points []image.Point) *DrawScript {
	result := NewDrawScript()
	if len(points) < 3 {
		return result
	}

	// Borders
	lastPoint := points[0]
	for i := 1; i < len(points); i++ {
		result.append(c.MakeScriptLineBresenham(float64(lastPoint.X), float64(lastPoint.Y), float64(points[i].X), float64(points[i].Y)))
		lastPoint = points[i]
	}
	result.append(c.MakeScriptLineBresenham(float64(lastPoint.X), float64(lastPoint.Y), float64(points[0].X), float64(points[0].Y)))

	minX := int(math.MaxInt32)
	maxX := int(math.MinInt32)
	minY := int(math.MaxInt32)
	maxY := int(math.MinInt32)

	for _, pp := range points {
		if pp.X > maxX {
			maxX = pp.X
		}
		if pp.X < minX {
			minX = pp.X
		}
		if pp.Y > maxY {
			maxY = pp.Y
		}
		if pp.Y < minY {
			minY = pp.Y
		}
	}

	beginX := 0

	for y := minY + 1; y <= maxY-1; y++ {
		countOfTriggers := 0
		flag := false

		for x := minX; x <= maxX+1; x++ {
			if result.hasPixel(x, y) && !result.hasPixel(x+1, y) {
				countOfTriggers++
			}
		}

		if countOfTriggers > 1 {
			for x := minX; x <= maxX; x++ {
				if result.hasPixel(x, y) && !result.hasPixel(x+1, y) {
					if flag {
						var line DrawScriptHorLine
						line.X1 = beginX
						line.X2 = x
						line.Y = y
						line.C = 1
						result.horLines = append(result.horLines, line)
					} else {
						beginX = x
					}

					flag = !flag
				}
			}
		}
	}

	return result
}

func (c *CanvasDirect) MakeScriptLine(x1, y1, x2, y2, width float64) *DrawScript {
	x1 = math.Round(x1)
	y1 = math.Round(y1)
	x2 = math.Round(x2)
	y2 = math.Round(y2)
	width = math.Round(width)

	result := NewDrawScript()

	{
		minX := math.Min(x1, x2)
		maxX := math.Max(x1, x2)
		minY := math.Min(y1, y2)
		maxY := math.Max(y1, y2)

		if int(minX) > c.state.ClipX+c.state.ClipWidth {
			return result
		}
		if int(maxX) < c.state.ClipX {
			return result
		}
		if int(minY) > c.state.ClipY+c.state.ClipHeight {
			return result
		}
		if int(maxY) < c.state.ClipY {
			return result
		}
	}

	if width > 1.1 {

		width -= 1
		width /= 2

		deltaX := x2 - x1
		deltaY := y2 - y1

		tnA := deltaY / deltaX

		katY := width * math.Cos(math.Atan(tnA))
		katX := math.Sqrt(math.Abs(width*width - katY*katY))

		if deltaX <= 0 && deltaY <= 0 {
			katY = -katY
		}

		if deltaX <= 0 && deltaY > 0 {
			katX = -katX
			katY = -katY
		}

		if deltaX > 0 && deltaY > 0 {
			katX = -katX
		}

		a1_x := x1 - katX
		a1_y := y1 - katY
		a2_x := x1 + katX
		a2_y := y1 + katY

		b1_x := x2 - katX
		b1_y := y2 - katY
		b2_x := x2 + katX
		b2_y := y2 + katY

		var ps []image.Point
		ps = make([]image.Point, 0)
		ps = append(ps, image.Pt(int(a1_x), int(a1_y)))
		ps = append(ps, image.Pt(int(a2_x), int(a2_y)))
		ps = append(ps, image.Pt(int(b2_x), int(b2_y)))
		ps = append(ps, image.Pt(int(b1_x), int(b1_y)))
		result.append(c.fillEx(ps))

		result.append(c.MakeScriptLineWu(a1_x, a1_y, a2_x, a2_y))
		result.append(c.MakeScriptLineWu(b1_x, b1_y, b2_x, b2_y))
		result.append(c.MakeScriptLineWu(a1_x, a1_y, b1_x, b1_y))
		result.append(c.MakeScriptLineWu(a2_x, a2_y, b2_x, b2_y))
	} else {
		result.append(c.MakeScriptLineWu(x1, y1, x2, y2))
	}

	return result
}

func (c *CanvasDirect) DrawText(x int, y int, text string, fontFamily string, fontSize float64, colr color.Color, underline bool) {
	if math.IsNaN(fontSize) {
		return
	}

	x += c.state.TranslateX
	y += c.state.TranslateY

	fontContext, err := FontContext(fontFamily, fontSize, false, false, colr)
	if err != nil {
		return
	}

	rectBounds := image.Rectangle{}
	rectBounds.Min.X = c.state.ClipX
	rectBounds.Min.Y = c.state.ClipY
	rectBounds.Max.X = c.state.ClipX + c.state.ClipWidth
	rectBounds.Max.Y = c.state.ClipY + c.state.ClipHeight

	fontContext.SetClip(rectBounds) // cnv.image.Bounds()
	fontContext.SetDst(c.image)

	_, face, err := Font(fontFamily, fontSize, false, false)
	if err != nil {
		return
	}

	//face := FontFace(f, fontSize)
	pt := freetype.Pt(x, y+face.Metrics().Ascent.Ceil())
	fontContext.DrawString(text, pt)
}

func (c *CanvasDirect) DrawTextMultiline(x int, y int, width int, height int, hAlign HAlign, vAlign VAlign, text string, colr color.Color, fontFamily string, fontSize float64, underline bool) {
	lines := strings.Split(text, "\r\n")

	yOffset := 0

	_, textHeight, err := MeasureText(fontFamily, fontSize, false, false, "Йg", true)
	if err != nil {
		return
	}

	//textHeight := 20

	fulltextHeight := textHeight * len(lines)

	switch vAlign {
	case VAlignTop:
		yOffset = 0
	case VAlignCenter:
		yOffset = height/2 - fulltextHeight/2
	case VAlignBottom:
		yOffset = height - fulltextHeight
	}

	c.Save()
	c.ClipIn(x+c.TranslatedX(), y+c.TranslatedY(), width, height)

	for _, str := range lines {
		xx := x
		textWidth, _, err := MeasureText(fontFamily, fontSize, false, false, str, false)
		if err != nil {
			return
		}

		if hAlign != HAlignLeft {

			switch hAlign {
			case HAlignLeft:
				xx = x
			case HAlignCenter:
				xx = (width / 2) - (textWidth / 2) + x
			case HAlignRight:
				xx = width - textWidth + x
			}
		}

		c.DrawText(xx, yOffset+y, str, fontFamily, fontSize, colr, underline)

		if underline {
			underLineWidth := fontSize / 20
			if underLineWidth < 1 {
				underLineWidth = 1
			}
			c.DrawLine(xx, yOffset+y+textHeight-1, xx+textWidth, yOffset+y+textHeight-1, int(underLineWidth), colr)
		}

		yOffset += textHeight
	}

	c.Load()
}

func CohenSutherland(x1, y1, x2, y2, left, top, right, bottom int) (int, int, int, int, bool) {

	invalid := false

	for {
		// Left-Right-Top-Bottom
		code1 := 0
		if x1 < left {
			code1 |= 8
		}
		if x1 > right {
			code1 |= 4
		}
		if y1 < top {
			code1 |= 2
		}
		if y1 > bottom {
			code1 |= 1
		}

		code2 := 0
		if x2 < left {
			code2 |= 8
		}
		if x2 > right {
			code2 |= 4
		}
		if y2 < top {
			code2 |= 2
		}
		if y2 > bottom {
			code2 |= 1
		}

		if code1 == 0 && code2 == 0 {
			break
		}

		if (code1 & code2) != 0 {
			invalid = true
			break // Outside
		}

		code := code1
		if code == 0 {
			code = code2
		}

		x := 0
		y := 0

		if (code & 2) != 0 {
			x = x1 + (x2-x1)*(top-y1)/(y2-y1)
			y = top
		} else {
			if (code & 1) != 0 {
				x = x1 + (x2-x1)*(bottom-y1)/(y2-y1)
				y = bottom
			} else {
				if (code & 4) != 0 {
					y = y1 + (y2-y1)*(right-x1)/(x2-x1)
					x = right
				} else {
					if (code & 8) != 0 {
						y = y1 + (y2-y1)*(left-x1)/(x2-x1)
						x = left
					}
				}
			}
		}

		if code1 != 0 {
			x1 = x
			y1 = y
		} else {
			x2 = x
			y2 = y
		}
	}

	return x1, y1, x2, y2, !invalid
}

/*func DrawText(text string, fontFamily string, fontSize float64, colr color.Color) image.Image {

	fontContext, err := FontContext(fontFamily, fontSize, false, false, colr)
	if err != nil {
		return nil
	}

	w, h, _ := MeasureText(fontFamily, fontSize, false, false, text, false)
	img := image.NewRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{w, h}})

	fontContext.SetDst(img)

	_, face, err := Font(fontFamily, fontSize, false, false)
	if err != nil {
		return nil
	}

	rectBounds := image.Rectangle{}
	rectBounds.Min.X = 0
	rectBounds.Min.Y = 0
	rectBounds.Max.X = 10000000
	rectBounds.Max.Y = 10000000
	fontContext.SetClip(rectBounds)

	pt := freetype.Pt(0, face.Metrics().Ascent.Ceil())
	fontContext.DrawString(text, pt)

	return img
}
*/
