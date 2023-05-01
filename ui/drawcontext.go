package ui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/ipoluianov/goforms/canvas"
)

type DrawContext interface {
	Init()

	SetColor(col color.Color)
	SetStrokeWidth(w int)
	SetFontFamily(fontFamily string)
	SetFontSize(s float64)
	SetTextAlign(h canvas.HAlign, v canvas.VAlign)
	SetUnderline(underline bool)

	DrawLine(x1, y1, x2, y2 int)
	DrawEllipse(x, y, width, height int)
	DrawRect(x, y, width, height int)
	FillRect(x, y, width, height int)
	Save()
	Load()
	Translate(x, y int)
	//Clip(x, y, width, height int)
	ClipIn(x, y, width, height int)
	DrawImage(x, y, width, height int, img image.Image)
	DrawText(x, y, width, height int, text string)
	MeasureText(text string) (int, int)

	TranslatedX() int
	TranslatedY() int
	ClippedRegion() (int, int, int, int)

	GraphContextImage() *image.RGBA
	GG() *gg.Context

	State() canvas.CanvasDirectState

	Finish()
}
