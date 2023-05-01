package uiinterfaces

import (
	"github.com/gazercloud/gazerui/ui"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window interface {
	LoopUI()

	Show()
	Draw() bool

	OnInit()
	OnClose() bool
	Maximize()

	Init()
	Close()

	// process events from OS
	// mouse
	ProcessMouseMove(x, y int)
	ProcessMouseWheel(delta int)
	ProcessMouseDown(button uievents.MouseButton)
	ProcessMouseUp(button uievents.MouseButton)
	ProcessClick(x, y int, button uievents.MouseButton)
	// keyboard
	ProcessCharInput(ch rune)
	ProcessKeyModifiers(shift bool, control bool, alt bool)
	ProcessKeyDown(key glfw.Key)
	ProcessKeyUp(key glfw.Key)
	// window
	ProcessWindowResize(width, height int)
	ProcessWindowMove(x, y int)
	ProcessFocus()

	KeyModifiers() uievents.KeyModifiers

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
	Position() ui.Point

	SetParent(window Window)
	Parent() Window
	Menu() Menu

	//CreatePopupForm(window Window, x int, y int)
	CreateModalForm(window Window)

	Modal() bool
	SetModal(modal bool)

	Popup() bool
	SetPopup(popup bool)

	UpdateWindow(source string)
	UpdateMenu()

	Accept()
	DialogResult() bool

	NewTimer(period int64, handler func()) *uievents.FormTimer
	RemoveTimer(timer *uievents.FormTimer)
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

	//SelectColorDialog(col color.Color, onColorChanged func(color color.Color)) (bool, color.Color)
	//SelectDateTimeDialog(dt time.Time, onDateTimeChanged func(dateTime time.Time)) (bool, time.Time)
	//MessageBox(title, text string)

	ShowMaximazed() bool

	UpdateLayout()

	SetWindow(w *glfw.Window)
	Window() *glfw.Window

	SetMouseCursor(cur ui.MouseCursor)
	CentralWidget() Widget
}
