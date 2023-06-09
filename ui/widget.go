package ui

import (
	"image/color"

	"github.com/ipoluianov/goforms/utils/uiproperties"
)

type Widget interface {
	Draw(ctx DrawContext)
	DrawControl(ctx DrawContext)

	Init()
	InitControl(parent Widget, w Widget)

	SetX(x int)
	SetY(y int)
	SetWidth(width int)
	SetHeight(height int)
	SetAnchors(anchors int)
	X() int
	Y() int
	Width() int
	Height() int
	Anchors() int
	SetHover(hover bool)
	Hover() bool

	Name() string
	SetName(name string)

	Focus()
	SetFocus(focus bool)
	HasFocus() bool

	Update(source string)

	ProcessMouseWheel(event *MouseWheelEvent)
	ProcessMouseMove(event *MouseMoveEvent)
	ProcessMouseDown(event *MouseDownEvent)
	ProcessMouseUp(event *MouseUpEvent)
	ProcessMouseClick(event *MouseClickEvent)
	ProcessMouseDblClick(event *MouseDblClickEvent)
	ProcessKeyChar(event *KeyCharEvent)
	ProcessKeyDown(event *KeyDownEvent) bool
	ProcessKeyUp(event *KeyUpEvent)

	MouseWheel(event *MouseWheelEvent)
	MouseMove(event *MouseMoveEvent)
	MouseDown(event *MouseDownEvent)
	MouseUp(event *MouseUpEvent)
	MouseDrop(event *MouseDropEvent)
	MouseValidateDrop(event *MouseValidateDropEvent)
	MouseClick(event *MouseClickEvent)
	MouseDblClick(event *MouseDblClickEvent)
	KeyChar(event *KeyCharEvent)
	KeyDown(event *KeyDownEvent) bool
	KeyUp(event *KeyUpEvent)

	MouseEnter()
	MouseLeave()
	FocusChanged(focus bool)

	InnerWidth() int
	InnerHeight() int

	LeftBorderWidth() int
	RightBorderWidth() int
	TopBorderWidth() int
	BottomBorderWidth() int

	ScrollOffsetX() int
	ScrollOffsetY() int
	ScrollEnsureVisible(x, y int)

	BackColor() color.Color
	ForeColor() color.Color
	AccentColor() color.Color
	InactiveColor() color.Color

	DrawBorders(ctx DrawContext)
	DrawBackground(ctx DrawContext)
	DrawScrollBars(ctx DrawContext)

	ProcessFindWidgetUnderPointer(x, y int) Widget
	FindWidgetUnderPointer(x, y int) Widget
	ClearHover()
	ClearFocus()
	AddProperty(name string, prop *uiproperties.Property)

	Classes() []string
	Subclass() string
	ControlType() string

	CurrentStyleValueScore(subclass string, propertyName string) int
	SetStyledValue(subclass string, propertyName string, value string, score int)
	StyledValue(subclass string, propertyName string) interface{}
	ApplyStyleLine(controlName string, controlType string, styleClass string, stylePseudoClass string, propertyName string, value string)

	SetBorderLeft(width int, col color.Color)
	SetBorderRight(width int, col color.Color)
	SetBorderTop(width int, col color.Color)
	SetBorderBottom(width int, col color.Color)
	SetBorders(width int, col color.Color)

	OnInit()

	ClearRadioButtons()
	SetParent(p Widget)
	Parent() Widget
	RectOnWindow() (int, int)
	RectClientAreaOnWindow() (int, int)
	Window() Window

	TranslateX(x int) int
	TranslateY(y int) int

	IsTabPlate() bool
	TabIndex() int

	AcceptsReturn() bool
	AcceptsTab() bool

	NextFocusControl() Widget
	FirstFocusControl() Widget

	FontFamily() string
	FontSize() float64
	FontBold() bool
	FontItalic() bool

	SetContextMenu(menu IMenu)
	ContextMenu() IMenu
	SetWindow(window Window)

	MouseCursor() MouseCursor

	SetUserData(key string, data interface{})
	UserData(key string) interface{}

	BeginUpdate()
	EndUpdate()

	Dispose()

	SetTooltip(text string)
	Tooltip() string

	ClosePopup()

	GridX() int
	GridY() int

	SetGridX(x int)
	SetGridY(y int)

	SetGridPos(x int, y int)

	MinWidth() int
	MinHeight() int
	MaxWidth() int
	MaxHeight() int
	XExpandable() bool
	YExpandable() bool

	SetXExpandable(xExpandable bool)
	SetYExpandable(yExpandable bool)

	SetMinWidth(minWidth int)
	SetMinHeight(minHeight int)
	SetMaxWidth(maxWidth int)
	SetMaxHeight(maxHeight int)

	Disposed() bool
	FullPath() string

	OnScroll(scrollPositionX int, scrollPositionY int)

	IsVisible() bool
	IsVisibleRec() bool
	SetVisible(visible bool)

	SetEnabled(enabled bool)
	EnabledChanged(enabled bool)

	Initialized() bool
	SetFixedSize(w int, h int)

	SetSize(w, h int)
	SetPos(x, y int)
	String(level int) string
	UpdateLayout()

	ClientWidth() int
	ClientHeight() int
	Text() string

	ClearLayoutCache()

	UpdateStyle()
	Widgets() []Widget
}
