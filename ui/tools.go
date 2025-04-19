package ui

import (
	"bytes"
	"image"
	"image/color"
	"time"

	"github.com/ipoluianov/nui/nuikey"
	"github.com/ipoluianov/nui/nuimouse"
)

type FormTimer struct {
	Enabled           bool
	LastElapsedDTMSec int64
	Period            int64
	Handler           func()
}

type CoordinateTranslator interface {
	TranslateX(x int) int
	TranslateY(y int) int
}

func (c *FormTimer) StartTimer() {
	c.Enabled = true
}

func (c *FormTimer) RestartTimer() {
	c.Enabled = true
	nowMSec := time.Now().UnixNano() / 1000000
	c.LastElapsedDTMSec = nowMSec
}

func (c *FormTimer) StopTimer() {
	c.Enabled = false
}

type UserDataContainer struct {
	userDataContainer map[string]interface{}
	//TempData          string
}

func (c *UserDataContainer) SetUserData(key string, data interface{}) {
	c.InitDataContainer()
	c.userDataContainer[key] = data
}

func (c *UserDataContainer) UserData(key string) interface{} {
	c.InitDataContainer()
	return c.userDataContainer[key]
}

func (c *UserDataContainer) AllUserData() interface{} {
	return c.userDataContainer
}

func (c *UserDataContainer) InitDataContainer() {
	if c.userDataContainer == nil {
		c.userDataContainer = make(map[string]interface{})
	}
}

func (c *UserDataContainer) Dispose() {
	c.userDataContainer = nil
}

type Event struct {
	UserDataContainer
	Sender    interface{}
	Modifiers nuikey.KeyModifiers
	Ignore    bool
}

func (c *Event) SetPosX(x int) {
}

func (c *Event) SetPosY(y int) {
}

func (c *Event) PosX() int {
	return 0
}

func (c *Event) PosY() int {
	return 0
}

type MouseEvent struct {
	Event
	X, Y int
}

func (c *MouseEvent) SetPosX(x int) {
	c.X = x
}

func (c *MouseEvent) SetPosY(y int) {
	c.Y = y
}

func (c *MouseEvent) PosX() int {
	return c.X
}

func (c *MouseEvent) PosY() int {
	return c.Y
}

type MouseWheelEvent struct {
	MouseEvent
	Delta int
}

type MouseMoveEvent struct {
	MouseEvent
}

type MouseDownEvent struct {
	MouseEvent
	Button nuimouse.MouseButton
}

func NewEvent(sender interface{}) *Event {
	var event Event
	event.InitDataContainer()
	event.Sender = sender
	return &event
}

func NewMouseDownEvent(x, y int, button nuimouse.MouseButton, modifiers nuikey.KeyModifiers) *MouseDownEvent {
	var event MouseDownEvent
	event.InitDataContainer()
	event.Modifiers = modifiers
	event.X = x
	event.Y = y
	event.Button = button
	return &event
}

func NewMouseDblClickEvent(x, y int, button nuimouse.MouseButton, modifiers nuikey.KeyModifiers) *MouseDblClickEvent {
	var event MouseDblClickEvent
	event.InitDataContainer()
	event.Modifiers = modifiers
	event.X = x
	event.Y = y
	event.Button = button
	return &event
}

func NewMouseDropEvent(x, y int, button nuimouse.MouseButton, modifiers nuikey.KeyModifiers, droppingObject interface{}) *MouseDropEvent {
	var event MouseDropEvent
	event.InitDataContainer()
	event.Modifiers = modifiers
	event.X = x
	event.Y = y
	event.Button = button
	event.DroppingObject = droppingObject
	return &event
}

func NewMouseValidateDropEvent(x, y int, button nuimouse.MouseButton, modifiers nuikey.KeyModifiers, droppingObject interface{}) *MouseValidateDropEvent {
	var event MouseValidateDropEvent
	event.InitDataContainer()
	event.Modifiers = modifiers
	event.X = x
	event.Y = y
	event.Button = button
	event.DroppingObject = droppingObject
	return &event
}

type MouseUpEvent struct {
	MouseEvent
	Button nuimouse.MouseButton
}

type MouseDropEvent struct {
	MouseEvent
	Button         nuimouse.MouseButton
	DroppingObject interface{}
}

type MouseValidateDropEvent struct {
	MouseEvent
	Button         nuimouse.MouseButton
	DroppingObject interface{}
	AllowDrop      bool
}

type MouseClickEvent struct {
	MouseEvent
	Button nuimouse.MouseButton
}

type MouseDblClickEvent struct {
	MouseEvent
	Button nuimouse.MouseButton
}

type KeyDownEvent struct {
	Event
	Key nuikey.Key
}

type KeyUpEvent struct {
	Event
	Key nuikey.Key
}

type KeyCharEvent struct {
	Event
	Ch rune
}

type FormSizeChangedEvent struct {
	Event
	Width  int
	Height int
}

func (event *MouseWheelEvent) Translate(w CoordinateTranslator) *MouseWheelEvent {
	//var result MouseWheelEvent
	result := *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseMoveEvent) Translate(w CoordinateTranslator) *MouseMoveEvent {
	//var result MouseMoveEvent
	result := *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseDropEvent) Translate(w CoordinateTranslator) *MouseDropEvent {
	//var result MouseDropEvent
	result := *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseValidateDropEvent) Translate(w CoordinateTranslator) *MouseValidateDropEvent {
	//var result MouseValidateDropEvent
	result := *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseDownEvent) Translate(w CoordinateTranslator) *MouseDownEvent {
	//var result MouseDownEvent
	result := *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseUpEvent) Translate(w CoordinateTranslator) *MouseUpEvent {
	//var result MouseUpEvent
	result := *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseClickEvent) Translate(w CoordinateTranslator) *MouseClickEvent {
	//var result MouseClickEvent
	result := *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseDblClickEvent) Translate(w CoordinateTranslator) *MouseDblClickEvent {
	//var result MouseDblClickEvent
	result := *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func ColorWithAlpha(col color.Color, alpha uint8) color.Color {
	r, b, g, _ := col.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), alpha}
}

func DecodeImage(bs []byte) image.Image {
	reader := bytes.NewReader(bs)
	image, _, _ := image.Decode(reader)
	return image
}
