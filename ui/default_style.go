package ui

import (
	"embed"
	"fmt"
	"image/color"

	"github.com/ipoluianov/goforms/utils/uiproperties"
)

func InitDefaultStyle(w Widget) {
	_ = embed.FS.Open
	if CurrentStyle == StyleLight {
		InitLight(w)
	}
	if CurrentStyle == StyleDarkWhite {
		InitDarkWhite(w)
	}
}

type Style int

const (
	StyleLight     = 0
	StyleDarkBlue  = 1
	StyleDarkWhite = 2
)

var CurrentStyle = StyleLight

var DefaultBackColor color.Color

//go:embed style_dark_white.style
var styleDarkWhiteBS []byte

//go:embed style_light.style
var styleLightBS []byte

func InitDarkWhite(w Widget) {
	var c CSS
	c.Parse(string(styleDarkWhiteBS))

	for _, line := range c.lines {

		if line.elementType == "Control" {
			if line.propertyName == "backgroundColor" {
				v, _ := uiproperties.ParseCSSProperty(uiproperties.PropertyTypeColor, line.value)
				DefaultBackColor = v.(color.Color)
			}
		}

		w.ApplyStyleLine(line.elementName, line.elementType, "", line.subclass, line.propertyName, line.value)
	}
}

func InitLight(w Widget) {
	var c CSS
	c.Parse(string(styleLightBS))

	for _, line := range c.lines {

		if line.elementType == "Control" {
			if line.propertyName == "backgroundColor" {
				v, _ := uiproperties.ParseCSSProperty(uiproperties.PropertyTypeColor, line.value)
				DefaultBackColor = v.(color.Color)
			}
		}

		if w.ControlType() == "ListView" {

			if line.propertyName == "selectionBackground" {
				fmt.Println("selectionBackground", line.value)
			}
		}

		w.ApplyStyleLine(line.elementName, line.elementType, "", line.subclass, line.propertyName, line.value)
	}
}
