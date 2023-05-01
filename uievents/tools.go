package uievents

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"image/color"
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

func (c *FormTimer) StopTimer() {
	c.Enabled = false
}

type UserDataContainer struct {
	userDataContainer map[string]interface{}
	TempData          string
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

type MouseButton int

const MouseButtonLeft MouseButton = 0x01
const MouseButtonMiddle MouseButton = 0x02
const MouseButtonRight MouseButton = 0x04

type KeyModifiers struct {
	Shift   bool
	Control bool
	Alt     bool
}

type Event struct {
	UserDataContainer
	Sender    interface{}
	Modifiers KeyModifiers
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
	Button MouseButton
}

func NewEvent(sender interface{}) *Event {
	var event Event
	event.InitDataContainer()
	event.Sender = sender
	return &event
}

func NewMouseDownEvent(x, y int, button MouseButton, modifiers KeyModifiers) *MouseDownEvent {
	var event MouseDownEvent
	event.InitDataContainer()
	event.Modifiers = modifiers
	event.X = x
	event.Y = y
	event.Button = button
	return &event
}

func NewMouseDropEvent(x, y int, button MouseButton, modifiers KeyModifiers, droppingObject interface{}) *MouseDropEvent {
	var event MouseDropEvent
	event.InitDataContainer()
	event.Modifiers = modifiers
	event.X = x
	event.Y = y
	event.Button = button
	event.DroppingObject = droppingObject
	return &event
}

func NewMouseValidateDropEvent(x, y int, button MouseButton, modifiers KeyModifiers, droppingObject interface{}) *MouseValidateDropEvent {
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
	Button MouseButton
}

type MouseDropEvent struct {
	MouseEvent
	Button         MouseButton
	DroppingObject interface{}
}

type MouseValidateDropEvent struct {
	MouseEvent
	Button         MouseButton
	DroppingObject interface{}
	AllowDrop      bool
}

type MouseClickEvent struct {
	MouseEvent
	Button MouseButton
}

type MouseDblClickEvent struct {
	MouseEvent
	Button MouseButton
}

type KeyDownEvent struct {
	Event
	Key glfw.Key
}

type KeyUpEvent struct {
	Event
	Key glfw.Key
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
	var result MouseWheelEvent
	result = *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseMoveEvent) Translate(w CoordinateTranslator) *MouseMoveEvent {
	var result MouseMoveEvent
	result = *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseDropEvent) Translate(w CoordinateTranslator) *MouseDropEvent {
	var result MouseDropEvent
	result = *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseValidateDropEvent) Translate(w CoordinateTranslator) *MouseValidateDropEvent {
	var result MouseValidateDropEvent
	result = *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseDownEvent) Translate(w CoordinateTranslator) *MouseDownEvent {
	var result MouseDownEvent
	result = *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseUpEvent) Translate(w CoordinateTranslator) *MouseUpEvent {
	var result MouseUpEvent
	result = *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseClickEvent) Translate(w CoordinateTranslator) *MouseClickEvent {
	var result MouseClickEvent
	result = *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func (event *MouseDblClickEvent) Translate(w CoordinateTranslator) *MouseDblClickEvent {
	var result MouseDblClickEvent
	result = *event
	result.X = w.TranslateX(event.X)
	result.Y = w.TranslateY(event.Y)
	return &result
}

func ColorWithAlpha(col color.Color, alpha uint8) color.Color {
	r, b, g, _ := col.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), alpha}
}
