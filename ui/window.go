package ui

import (
	"image"

	"github.com/ipoluianov/nui/nui"
	"github.com/ipoluianov/nui/nuikey"
	"github.com/ipoluianov/nui/nuimouse"
)

type Window interface {
	LoopUI()

	Show()
	Draw(rgba *image.RGBA) bool

	OnInit()
	OnClose() bool
	Maximize()

	Init()
	Close()

	// process events from OS
	// mouse
	ProcessMouseMove(x, y int)
	ProcessMouseWheel(delta int)
	ProcessMouseDown(button nuimouse.MouseButton)
	ProcessMouseUp(button nuimouse.MouseButton)
	ProcessMouseDblClick(button nuimouse.MouseButton)
	ProcessClick(x, y int, button nuimouse.MouseButton)
	// keyboard
	ProcessCharInput(ch rune)
	ProcessKeyModifiers(shift bool, control bool, alt bool)
	ProcessKeyDown(key nuikey.Key)
	ProcessKeyUp(key nuikey.Key)
	// window
	ProcessWindowResize(width, height int)
	ProcessWindowMove(x, y int)
	ProcessFocus()

	KeyModifiers() nuikey.KeyModifiers

	// title
	Title() string
	SetTitle(title string)

	// size
	Width() int
	Height() int
	Resize(width, height int)

	IsMainWindow() bool
	SetIsMainWindow(isMainWindow bool)
	Id() int
	SetId(id int)
	Position() Point

	SetParent(window Window)
	Parent() Window
	Menu() Menu

	CreateModalForm(window Window)

	Modal() bool
	SetModal(modal bool)

	Popup() bool
	SetPopup(popup bool)

	UpdateWindow(source string)
	UpdateMenu()

	Accept()
	DialogResult() bool

	NewTimer(period int64, handler func()) *FormTimer
	RemoveTimer(timer *FormTimer)
	MainTimer()

	BeginDrag(object interface{})
	CurrentDraggingObject() interface{}

	ShowTooltip(x, y int, text string)
	SetFocusForWidget(c Widget)
	FocusedWidget() Widget

	AppendPopup(c Widget)
	CloseAllPopup()
	CloseAfterPopupWidget(w Widget)
	CloseTopPopup()

	ProcessTabDown()

	ControlRemoved()

	ShowMaximazed() bool

	UpdateLayout()

	SetWindow(w nui.Window)
	Window() nui.Window

	SetMouseCursor(cur nuimouse.MouseCursor)
	CentralWidget() Widget
}
