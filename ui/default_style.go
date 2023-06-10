package ui

import (
	"image/color"

	"github.com/ipoluianov/goforms/utils/uiproperties"
	"gopkg.in/go-playground/colors.v1"
)

func InitDefaultStyle(w Widget) {
	if CurrentStyle == StyleLight {
		InitLight(w)
	}
	if CurrentStyle == StyleDarkBlue {
		InitDarkBlue(w)
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

var CurrentStyle = StyleDarkWhite

var DefaultBackColor color.Color

func InitLight(w Widget) {
	var c CSS
	c.Parse(`
Control 
{
	fontFamily: roboto;
	fontSize: 16;
	backgroundColor:#FFFFFFFF;
	foregroundColor:#555555;
	inactiveColor:#888888;
	accentColor:#4169e1;

	leftBorderColor: #777;
	rightBorderColor: #777;
	topBorderColor: #777;
	bottomBorderColor: #777;

	selectionBackground: #ADD8E6;

	verticalScrollVisible: true;

	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 0;
}

Panel, HSpacer, VSpacer, TextBlock, ImageBox, Container, CheckBox
{
	backgroundColor:#FFFFFF00;
}

Dialog {
	backgroundColor:#F8F8F8;
}

Button, TextBox, ListView, TreeView, ProgressBar, ComboBox, TimeChart
{
	leftBorderWidth: 1;
	rightBorderWidth: 1;
	topBorderWidth: 1;
	bottomBorderWidth: 1;
	barColor: #800;
}

Button:hover 
{
	backgroundColor:#CCCCCCFF;
}

Button:focus 
{
	backgroundColor:#CCCCCCFF;
}

Button:disabled
{
	backgroundColor:#FFFFFFFF;
}

Button:clicked 
{
	backgroundColor:#AAAAAAFF;
}

ListView {
	gridColor: #DDDDDD;
}

ListViewHeader
{
	leftBorderColor: #777;
	rightBorderColor: #777;
	topBorderColor: #777;
	bottomBorderColor: #777;
	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 1;
}

TreeViewHeader
{
	leftBorderColor: #777;
	rightBorderColor: #777;
	topBorderColor: #777;
	bottomBorderColor: #777;
	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 1;
}

PopupMenu {
	leftBorderColor: #777;
	rightBorderColor: #777;
	topBorderColor: #777;
	bottomBorderColor: #777;
	leftBorderWidth: 1;
	rightBorderWidth: 1;
	topBorderWidth: 1;
	bottomBorderWidth: 1;
}

PopupMenuItem:hover {
	backgroundColor:#CCCCCCFF;
}

ProgressBar
{
	barColor: #800;
}

`)

	for _, line := range c.lines {

		if line.elementType == "Control" {
			if line.propertyName == "backgroundColor" {
				v, _ := uiproperties.ParseCSSProperty(uiproperties.PropertyTypeColor, line.value)
				DefaultBackColor = v.(color.Color)
			}
		}

		w.ApplyStyleLine(line.elementName, line.elementType, "", line.subclass, line.propertyName, line.value)
	}

	return

}

func InitDarkGreen(w Widget) {
	var c CSS
	c.Parse(`
Control 
{
	fontFamily: Roboto;
	fontSize: 16;
	backgroundColor:#222222FF;
	foregroundColor:#148014;
	inactiveColor:#225522;
	accentColor:#555500;

	leftBorderColor: #444444;
	rightBorderColor: #444444;
	topBorderColor: #444444;
	bottomBorderColor: #444444;

	selectionBackground: #144014;

	verticalScrollVisible: true;

	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 0;
}

Panel, HSpacer, VSpacer, TextBlock, ImageBox, Container, CheckBox
{
	backgroundColor:#FFFFFF00;
}

Dialog {
	backgroundColor:#272727FF;
}

Button, TextBox, ListView, TreeView, ProgressBar, ComboBox, TimeChart
{
	leftBorderColor: #777;
	rightBorderColor: #777;
	topBorderColor: #777;
	bottomBorderColor: #777;
	leftBorderWidth: 1;
	rightBorderWidth: 1;
	topBorderWidth: 1;
	bottomBorderWidth: 1;
	barColor: #800;
}

Button:hover 
{
	backgroundColor:#444444FF;
}

Button:clicked 
{
	backgroundColor:#777777FF;
}

Button:disabled
{
	backgroundColor:#222222FF;
}

ListView {
	gridColor: #00FF00;
}

ListViewHeader
{
	leftBorderColor: #777;
	rightBorderColor: #777;
	topBorderColor: #777;
	bottomBorderColor: #777;
	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 1;
}

TreeViewHeader
{
	leftBorderColor: #777;
	rightBorderColor: #777;
	topBorderColor: #777;
	bottomBorderColor: #777;
	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 1;
}

ProgressBar
{
	barColor: #800;
}

`)

	for _, line := range c.lines {

		if line.elementType == "Control" {
			if line.propertyName == "backgroundColor" {
				v, _ := uiproperties.ParseCSSProperty(uiproperties.PropertyTypeColor, line.value)
				DefaultBackColor = v.(color.Color)
			}
		}

		w.ApplyStyleLine(line.elementName, line.elementType, "", line.subclass, line.propertyName, line.value)
	}

	return

}

func InitDarkBlue(w Widget) {
	var c CSS
	c.Parse(`
Control 
{
	fontFamily: Roboto;
	fontSize: 16;
	backgroundColor:#303030FF;
	foregroundColor:#2298EB;
	inactiveColor:#444;
	accentColor:#ff8c00;

	leftBorderColor: #444;
	rightBorderColor: #444;
	topBorderColor: #444;
	bottomBorderColor: #444;

	selectionBackground: #00489B;

	verticalScrollVisible: true;

	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 0;
}

Control:disabled
{
	foregroundColor:#004090;
}

Panel, HSpacer, VSpacer, TextBlock, ImageBox, Container, CheckBox
{
	backgroundColor:#FFFFFF00;
}

Dialog {
	backgroundColor:#101010FF;
}

TimeChart {
	color0: #ff8c00;
	color1: #ff2000;
	color2: #8cff00;
	color3: #005566;
	color4: #440000;
}

Button, TextBox, ListView, TreeView, ProgressBar, ComboBox, TimeChart
{
	leftBorderWidth: 1;
	rightBorderWidth: 1;
	topBorderWidth: 1;
	bottomBorderWidth: 1;
	barColor: #800;
}

Button:hover 
{
	backgroundColor:#004070;
}

Button:focus 
{
	backgroundColor:#202040FF;
}

TextBlock:disabled
{
	foregroundColor:#004455FF;
}

Button:disabled
{
	backgroundColor:#202020FF;
	foregroundColor:#004455FF;
}

Button:clicked 
{
	backgroundColor:#777777FF;
}

ListViewHeader
{
	foregroundColor:#777;
	backgroundColor:#242424;
	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 1;
	bottomBorderColor: #444;
}

ListView {
	gridColor: #333;
}


TreeViewHeader
{
	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 1;
}

PopupMenu {
	leftBorderColor: #777;
	rightBorderColor: #777;
	topBorderColor: #777;
	bottomBorderColor: #777;
	leftBorderWidth: 1;
	rightBorderWidth: 1;
	topBorderWidth: 1;
	bottomBorderWidth: 1;
}

PopupMenuItem:hover {
	backgroundColor:#444444FF;
}

ProgressBar
{
	barColor: #800;
}

`)

	for _, line := range c.lines {

		if line.elementType == "Control" {
			if line.propertyName == "backgroundColor" {
				v, _ := uiproperties.ParseCSSProperty(uiproperties.PropertyTypeColor, line.value)
				DefaultBackColor = v.(color.Color)
			}
		}

		w.ApplyStyleLine(line.elementName, line.elementType, "", line.subclass, line.propertyName, line.value)
	}

	return

}

func InitDarkWhite(w Widget) {
	var c CSS
	c.Parse(`
Control 
{
	fontFamily: Roboto;
	fontSize: 16;
	backgroundColor:#303030;
	foregroundColor:#AAAAAA;
	inactiveColor:#444;
	accentColor:#ff8c00;

	leftBorderColor: #444;
	rightBorderColor: #444;
	topBorderColor: #444;
	bottomBorderColor: #444;

	selectionBackground: #00489B;

	verticalScrollVisible: true;

	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 0;
}

Control:disabled
{
	foregroundColor:#004090;
}

Panel, HSpacer, VSpacer, TextBlock, ImageBox, Container, CheckBox
{
	backgroundColor:#00FFFFFF;
}

Dialog {
	backgroundColor:#FF101010;
}

TimeChart {
	color0: #ff8c00;
	color1: #ff2000;
	color2: #8cff00;
	color3: #005566;
	color4: #440000;
}

Button, TextBox, ListView, TreeView, ProgressBar, ComboBox, TimeChart
{
	leftBorderWidth: 1;
	rightBorderWidth: 1;
	topBorderWidth: 1;
	bottomBorderWidth: 1;
	barColor: #800;
}

Button:hover 
{
	backgroundColor:#004070;
}

Button:focus 
{
	backgroundColor:#202040;
}

TextBlock:disabled
{
	foregroundColor:#004455;
}

Button:disabled
{
	backgroundColor:#202020;
	foregroundColor:#004455;
}

Button:clicked 
{
	backgroundColor:#777777;
}

ListViewHeader
{
	foregroundColor:#777;
	backgroundColor:#242424;
	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 1;
	bottomBorderColor: #444;
}

ListView {
	gridColor: #333;
}


TreeViewHeader
{
	leftBorderWidth: 0;
	rightBorderWidth: 0;
	topBorderWidth: 0;
	bottomBorderWidth: 1;
}

PopupMenu {
	leftBorderColor: #777;
	rightBorderColor: #777;
	topBorderColor: #777;
	bottomBorderColor: #777;
	leftBorderWidth: 1;
	rightBorderWidth: 1;
	topBorderWidth: 1;
	bottomBorderWidth: 1;
}

PopupMenuItem:hover {
	backgroundColor:#444444;
}

ProgressBar
{
	barColor: #800;
}

`)

	for _, line := range c.lines {

		if line.elementType == "Control" {
			if line.propertyName == "backgroundColor" {
				v, _ := uiproperties.ParseCSSProperty(uiproperties.PropertyTypeColor, line.value)
				DefaultBackColor = v.(color.Color)
			}
		}

		w.ApplyStyleLine(line.elementName, line.elementType, "", line.subclass, line.propertyName, line.value)
	}

	return

}

func parseHexColor(str string) color.Color {
	hex, _ := colors.ParseHEX(str)
	rgb := hex.ToRGB()
	var col color.RGBA
	col.R = rgb.R
	col.G = rgb.G
	col.B = rgb.B
	col.A = 255
	return col
}
