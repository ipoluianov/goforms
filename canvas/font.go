package canvas

import (
	"errors"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
)

type FontInfo struct {
	family         string
	fontRegular    *truetype.Font
	fontBold       *truetype.Font
	fontItalic     *truetype.Font
	fontBoldItalic *truetype.Font
	faces          map[float64]font.Face
}

var globalFonts map[string]*FontInfo

func init() {

	globalFonts = make(map[string]*FontInfo)
	for _, assetName := range AssetNames() {
		if strings.HasPrefix(assetName, "fonts") {
			nameParts := strings.Split(assetName, "/")
			fontName := nameParts[len(nameParts)-2]
			fileName := nameParts[len(nameParts)-1]

			if _, ok := globalFonts[fontName]; !ok {
				var fontInfo FontInfo
				fontInfo.family = strings.ToLower(fontName)
				fontInfo.faces = make(map[float64]font.Face)
				globalFonts[fontName] = &fontInfo
			}

			val := globalFonts[fontName]

			if strings.Contains(strings.ToLower(fileName), "regular") {
				f, err := readFontFromAsset(assetName)
				if err == nil {
					val.fontRegular = f
				}
			}
			if strings.Contains(strings.ToLower(fileName), "bold") && !strings.Contains(strings.ToLower(fileName), "italic") {
				f, err := readFontFromAsset(assetName)
				if err == nil {
					val.fontBold = f
				}
			}
			if strings.Contains(strings.ToLower(fileName), "italic") && !strings.Contains(strings.ToLower(fileName), "bold") {
				f, err := readFontFromAsset(assetName)
				if err == nil {
					val.fontItalic = f
				}
			}
			if strings.Contains(strings.ToLower(fileName), "italic") && strings.Contains(strings.ToLower(fileName), "bold") {
				f, err := readFontFromAsset(assetName)
				if err == nil {
					val.fontBoldItalic = f
				}
			}
		}
	}
}

func readFontFromAsset(assentName string) (*truetype.Font, error) {
	fontBytes, err := Asset(assentName)
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Font(family string, size float64, bold bool, italic bool) (result *truetype.Font, face font.Face, err error) {
	var val *FontInfo
	ok := false

	if val, ok = globalFonts[strings.ToLower(family)]; ok {
		if !bold && !italic {
			result = val.fontRegular
		}
		if bold && !italic {
			result = val.fontBold
		}
		if !bold && italic {
			result = val.fontItalic
		}
		if bold && italic {
			result = val.fontBoldItalic
		}
	}

	if result == nil {
		err = fmt.Errorf("No font found")
	}

	if _, ok := val.faces[size]; !ok {
		val.faces[size] = MakeFontFace(result, size)
	}

	return result, val.faces[size], err
}

func FontContext(family string, size float64, bold bool, italic bool, color color.Color) (*Context, error) {
	fg := image.NewUniform(color)
	f, _, err := Font(family, size, bold, italic)
	if err != nil {
		return nil, err
	}

	c := NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	//c.SetClip(cnv.image.Bounds())
	//c.SetDst(cnv.image)
	c.SetSrc(fg)
	return c, nil
}

func MakeFontFace(f *truetype.Font, size float64) font.Face {
	h := font.HintingNone
	face := truetype.NewFace(f, &truetype.Options{
		Size:              size,
		DPI:               *dpi,
		Hinting:           h,
		GlyphCacheEntries: 64,
	})
	return face
}

type FontMeasureData struct {
	Width  int
	Height int
}

var fontMeasureData map[string]*FontMeasureData

func GetFontMeasureData(family string, size float64, bold bool, italic bool, text string, multiline bool) *FontMeasureData {

	if fontMeasureData == nil {
		fontMeasureData = make(map[string]*FontMeasureData)
	}

	/*lines := strings.FieldsFunc(text, func(r rune) bool {
		return r == '\n'
	})
	linesCount := len(lines)
	if linesCount < 1 {
		linesCount = 1
	}*/

	fontHash := "font-" + family + "-size-" + strconv.FormatFloat(size, 'f', 10, 64) + "-bold-" + strconv.FormatBool(bold) + "-italic-" + strconv.FormatBool(italic) + text
	if fmd, ok := fontMeasureData[fontHash]; ok {
		return fmd
	}

	var fmdNew FontMeasureData
	fmdNew.Width, fmdNew.Height, _ = MeasureTextFreeType(family, size, bold, italic, text, multiline)
	fontMeasureData[fontHash] = &fmdNew

	return &fmdNew
}

func MeasureText(family string, size float64, bold bool, italic bool, text string, multiline bool) (int, int, error) {
	measureData := GetFontMeasureData(family, size, bold, italic, text, multiline)
	return measureData.Width, measureData.Height, nil
}

func MeasureTextFreeType(family string, size float64, bold bool, italic bool, text string, multiline bool) (int, int, error) {

	if math.IsNaN(size) {
		return 0, 0, errors.New("font size is NaN")
	}

	_, face, err := Font(family, size, bold, italic)

	if err != nil {
		return 0, 0, err
	}

	//face := FontFace(f, size)

	d := &font.Drawer{
		Dst:  nil,
		Src:  nil,
		Face: face,
	}

	var textWidth int
	var textHeight int

	var lines []string

	if multiline {
		lines = strings.Split(text, "\r\n")
	} else {
		lines = make([]string, 0)
		lines = append(lines, text)
	}

	for _, str := range lines {
		//pt := freetype.Pt(x, y+face.Metrics().Ascent.Ceil())

		w := d.MeasureString(str) / 64

		h := face.Metrics().Ascent + face.Metrics().Descent
		textHeight += int(h)
		if int(w) > textWidth {
			textWidth = int(w)
		}
	}

	return int(textWidth), int(textHeight / 64), nil
}

func CharPositions(family string, size float64, bold bool, italic bool, text string) ([]int, error) {
	runes := []rune(text)
	result := make([]int, len(runes)+1)
	for pos, _ := range runes {
		safeSubstring := string(runes[0:pos])
		w, _, err := MeasureText(family, size, bold, italic, safeSubstring, false)
		if err != nil {
			return nil, err
		}
		result[pos] = w
	}
	w, _, err := MeasureText(family, size, bold, italic, text, false)
	if err != nil {
		return nil, err
	}
	result[len(runes)] = w

	return result, err
}
