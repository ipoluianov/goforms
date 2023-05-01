package uiinterfaces

import (
	"image/color"

	"github.com/ipoluianov/goforms/ui"
	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiproperties"
)

type Widget interface {
	Draw(ctx ui.DrawContext)
	DrawControl(ctx ui.DrawContext)

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

	ProcessMouseWheel(event *uievents.MouseWheelEvent)
	ProcessMouseMove(event *uievents.MouseMoveEvent)
	ProcessMouseDown(event *uievents.MouseDownEvent)
	ProcessMouseUp(event *uievents.MouseUpEvent)
	ProcessMouseClick(event *uievents.MouseClickEvent)
	ProcessMouseDblClick(event *uievents.MouseDblClickEvent)
	ProcessKeyChar(event *uievents.KeyCharEvent)
	ProcessKeyDown(event *uievents.KeyDownEvent) bool
	ProcessKeyUp(event *uievents.KeyUpEvent)

	MouseWheel(event *uievents.MouseWheelEvent)
	MouseMove(event *uievents.MouseMoveEvent)
	MouseDown(event *uievents.MouseDownEvent)
	MouseUp(event *uievents.MouseUpEvent)
	MouseDrop(event *uievents.MouseDropEvent)
	MouseValidateDrop(event *uievents.MouseValidateDropEvent)
	MouseClick(event *uievents.MouseClickEvent)
	MouseDblClick(event *uievents.MouseDblClickEvent)
	KeyChar(event *uievents.KeyCharEvent)
	KeyDown(event *uievents.KeyDownEvent) bool
	KeyUp(event *uievents.KeyUpEvent)

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

	DrawBorders(ctx ui.DrawContext)
	DrawBackground(ctx ui.DrawContext)
	DrawScrollBars(ctx ui.DrawContext)

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

	SetContextMenu(menu Menu)
	ContextMenu() Menu
	SetWindow(window Window)

	MouseCursor() ui.MouseCursor

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
