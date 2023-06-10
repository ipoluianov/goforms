package ui

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/ipoluianov/goforms/utils"
	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/goforms/utils/uiproperties"
	"golang.org/x/image/colornames"
)

// type DrawControl func(ctx *gg.Context)

type Control struct {
	utils.Obj
	UserDataContainer

	disposed bool

	name     *uiproperties.Property
	isInited bool

	x      *uiproperties.Property
	y      *uiproperties.Property
	width  *uiproperties.Property
	height *uiproperties.Property

	anchors *uiproperties.Property

	innerSizeOverloaded   bool
	innerWidthOverloaded  int
	innerHeightOverloaded int

	hover bool
	focus bool

	scrollOffsetX int
	scrollOffsetY int

	hScroll Rect
	vScroll Rect

	gridX       int
	gridY       int
	minWidth    int
	minHeight   int
	maxWidth    int
	maxHeight   int
	xExpandable bool
	yExpandable bool

	visible bool
	enabled bool

	parent    Widget
	ownWindow Window

	theme  string
	widget Widget

	toolTip *uiproperties.Property

	foregroundColor *uiproperties.Property
	accentColor     *uiproperties.Property
	backgroundColor *uiproperties.Property
	inactiveColor   *uiproperties.Property

	leftBorderWidth   *uiproperties.Property
	leftBorderColor   *uiproperties.Property
	rightBorderWidth  *uiproperties.Property
	rightBorderColor  *uiproperties.Property
	topBorderWidth    *uiproperties.Property
	topBorderColor    *uiproperties.Property
	bottomBorderWidth *uiproperties.Property
	bottomBorderColor *uiproperties.Property

	/*	leftMargin   *uiproperties.Property
		rightMargin  *uiproperties.Property
		topMargin    *uiproperties.Property
		bottomMargin *uiproperties.Property*/

	properties map[string]*uiproperties.Property

	fontFamily *uiproperties.Property
	fontSize   *uiproperties.Property
	fontBold   *uiproperties.Property
	fontItalic *uiproperties.Property

	verticalScrollVisible   *uiproperties.Property
	verticalScrollWidth     *uiproperties.Property
	horizontalScrollVisible *uiproperties.Property
	horizontalScrollWidth   *uiproperties.Property

	verticalScrollDisplayed   bool
	horizontalScrollDisplayed bool

	verticalScrollMoving           bool
	verticalScrollMovingMousePos   int
	verticalScrollMovingOriginal   int
	horizontalScrollMoving         bool
	horizontalScrollMovingMousePos int
	horizontalScrollMovingOriginal int

	contextMenu       IMenu
	OnContextMenuNeed func(x, y int) IMenu

	//innerCanvas  *CanvasDirect
	widgetCanvas *canvas.CanvasDirect
	needToRedraw bool
	alwaysRedraw bool
	isUpdating   bool

	acceptReturnKey bool

	onMouseWheel    func(event *MouseWheelEvent)
	onMouseMove     func(event *MouseMoveEvent)
	onMouseDown     func(event *MouseDownEvent)
	onMouseUp       func(event *MouseUpEvent)
	onMouseClick    func(event *MouseClickEvent)
	onMouseDblClick func(event *MouseDblClickEvent)
	onKeyChar       func(event *KeyCharEvent)
	onKeyDown       func(event *KeyDownEvent) bool
	onKeyUp         func(event *KeyUpEvent)
	onMouseEnter    func()
	onMouseLeave    func()

	onPreMouseWheel    func(event *MouseWheelEvent)
	onPreMouseMove     func(event *MouseMoveEvent)
	onPreMouseDown     func(event *MouseDownEvent)
	onPreMouseUp       func(event *MouseUpEvent)
	onPreMouseClick    func(event *MouseClickEvent)
	onPreMouseDblClick func(event *MouseDblClickEvent)
	onPreKeyChar       func(event *KeyCharEvent)
	onPreKeyDown       func(event *KeyDownEvent)
	onPreKeyUp         func(event *KeyUpEvent)
	onScrolled         func(hScrollPos int, vScrollPos int)

	OnMouseDrop func(droppedValue interface{}, x, y int)

	cursor MouseCursor

	isTabPlate bool
	tabIndex   int
}

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (c *Rect) Contains(x, y int) bool {
	if x >= c.X && x < c.X+c.Width && y > c.Y && y < c.Y+c.Height {
		return true
	}
	return false
}

const ANCHOR_LEFT = 1
const ANCHOR_TOP = 2
const ANCHOR_RIGHT = 4
const ANCHOR_BOTTOM = 8
const ANCHOR_ALL = 0xF
const ANCHOR_HOR = ANCHOR_TOP | ANCHOR_LEFT | ANCHOR_RIGHT
const MAX_WIDTH = 100000
const MAX_HEIGHT = 100000

func NewControl(parent Widget) *Control {
	var c Control
	c.InitControl(parent, &c)
	return &c
}

func (c *Control) InitControl(parent Widget, w Widget) {
	if parent != nil {
		if !parent.Initialized() {
			panic("Parent control is not initialized")
		}
	}

	if parent == nil {
		//fmt.Error("no parent for control")
	}

	c.parent = parent
	if parent != nil {
		c.SetWindow(parent.Window())
	}
	c.widget = w
	//c.properties = make(map[string]*uiproperties.Property)
	c.visible = true
	c.enabled = true
	//c.alwaysRedraw = true

	c.name = AddUnstyledPropertyToWidget(w, "name", uiproperties.PropertyTypeString)
	c.x = AddUnstyledPropertyToWidget(w, "x", uiproperties.PropertyTypeInt)
	c.y = AddUnstyledPropertyToWidget(w, "y", uiproperties.PropertyTypeInt)
	c.width = AddUnstyledPropertyToWidget(w, "width", uiproperties.PropertyTypeInt)
	c.height = AddUnstyledPropertyToWidget(w, "height", uiproperties.PropertyTypeInt)
	c.anchors = AddUnstyledPropertyToWidget(w, "anchors", uiproperties.PropertyTypeInt)
	c.toolTip = AddUnstyledPropertyToWidget(w, "toolTip", uiproperties.PropertyTypeString)
	c.toolTip.SetOwnValue("")

	/*c.SetX(x)
	c.SetY(y)
	c.SetWidth(width)
	c.SetHeight(height)

	c.SetAnchors(anchors)*/

	c.foregroundColor = AddPropertyToWidget(w, "foregroundColor", uiproperties.PropertyTypeColor)
	c.accentColor = AddPropertyToWidget(w, "accentColor", uiproperties.PropertyTypeColor)
	c.backgroundColor = AddPropertyToWidget(w, "backgroundColor", uiproperties.PropertyTypeColor)
	c.inactiveColor = AddPropertyToWidget(w, "inactiveColor", uiproperties.PropertyTypeColor)

	c.leftBorderWidth = AddPropertyToWidget(w, "leftBorderWidth", uiproperties.PropertyTypeInt)
	c.leftBorderColor = AddPropertyToWidget(w, "leftBorderColor", uiproperties.PropertyTypeColor)
	c.rightBorderWidth = AddPropertyToWidget(w, "rightBorderWidth", uiproperties.PropertyTypeInt)
	c.rightBorderColor = AddPropertyToWidget(w, "rightBorderColor", uiproperties.PropertyTypeColor)
	c.topBorderWidth = AddPropertyToWidget(w, "topBorderWidth", uiproperties.PropertyTypeInt)
	c.topBorderColor = AddPropertyToWidget(w, "topBorderColor", uiproperties.PropertyTypeColor)
	c.bottomBorderWidth = AddPropertyToWidget(w, "bottomBorderWidth", uiproperties.PropertyTypeInt)
	c.bottomBorderColor = AddPropertyToWidget(w, "bottomBorderColor", uiproperties.PropertyTypeColor)

	c.fontFamily = AddPropertyToWidget(w, "fontFamily", uiproperties.PropertyTypeString)
	c.fontSize = AddPropertyToWidget(w, "fontSize", uiproperties.PropertyTypeDouble)
	c.fontBold = AddPropertyToWidget(w, "fontBold", uiproperties.PropertyTypeBool)
	c.fontItalic = AddPropertyToWidget(w, "fontItalic", uiproperties.PropertyTypeBool)

	c.verticalScrollVisible = AddPropertyToWidget(w, "verticalScrollVisible", uiproperties.PropertyTypeBool)
	c.verticalScrollWidth = AddPropertyToWidget(w, "verticalScrollWidth", uiproperties.PropertyTypeInt)
	c.horizontalScrollVisible = AddPropertyToWidget(w, "horizontalScrollVisible", uiproperties.PropertyTypeBool)
	c.horizontalScrollWidth = AddPropertyToWidget(w, "horizontalScrollWidth", uiproperties.PropertyTypeInt)

	c.verticalScrollWidth.SetOwnValue(10)
	c.horizontalScrollWidth.SetOwnValue(10)

	c.minWidth = 0
	c.minHeight = 0
	c.maxWidth = MAX_WIDTH
	c.maxHeight = MAX_HEIGHT
	c.xExpandable = false
	c.yExpandable = false

	c.theme = "light"

	InitDefaultStyle(w)

	c.Obj.InitObj("UI_"+w.ControlType(), time.Now().Format("2006-01-02 15-04-05-999")+"_"+fmt.Sprint(rand.Int())+"_###_"+c.FullPath())
	runtime.SetFinalizer(c, finalizerControl)

	c.isInited = true
	c.SetPos(0, 0)
	c.SetSize(42, 42)

	w.OnInit()
}

func finalizerControl(c *Control) {
	c.Obj.UninitObj()
}

func (c *Control) OnInit() {

}

func (c *Control) UpdateStyle() {
	InitDefaultStyle(c)
	if c.contextMenu != nil {
		if popupMenu, ok := c.contextMenu.(*PopupMenu); ok {
			InitDefaultStyle(popupMenu)
			popupMenu.UpdateStyle()
		}
	}
	c.Update("ControlUpdateStyle")
}

func (c *Control) Widgets() []Widget {
	return []Widget{}
}

func (c *Control) SetGridPos(x int, y int) {
	c.SetGridX(x)
	c.SetGridY(y)
}

func (c *Control) Parent() Widget {
	return c.parent
}

func (c *Control) IsVisible() bool {
	return c.visible
}

func (c *Control) IsVisibleRec() bool {
	if !c.visible {
		return false
	}

	if c.parent != nil {
		if !c.parent.IsVisible() {
			return false
		}
	}
	return c.visible
}

func (c *Control) SetVisible(visible bool) {
	c.visible = visible

	c.ownWindow.UpdateLayout()
	c.Window().UpdateLayout()

	c.Update("Control")
}

func (c *Control) SetEnabled(enabled bool) {
	if c.enabled != enabled {
		c.enabled = enabled
		c.widget.EnabledChanged(c.enabled)
		c.Update("Control")
	}
}

func (c *Control) EnabledChanged(enabled bool) {
}

func (c *Control) IsTabPlate() bool {
	return c.isTabPlate
}

func (c *Control) SetIsTabPlate(isTabPlate bool) {
	c.isTabPlate = isTabPlate
}

func (c *Control) SetTabIndex(index int) {
	c.tabIndex = index
}

func (c *Control) TabIndex() int {
	return c.tabIndex
}

func (c *Control) FullPath() string {
	result := ""
	if c.parent != nil {
		result = c.parent.FullPath()
	}

	name := c.Name()
	if len(name) == 0 {
		if len(name) == 0 {
			name = "[" + c.widget.ControlType() + ":" + c.widget.Text() + "]"
		}
	}

	if c.widget == nil {
		panic("No widget")
	}

	return result + "/" + name
}

func (c *Control) Disposed() bool {
	return c.disposed
}

func (c *Control) Dispose() {
	c.disposed = true
	c.UserDataContainer.Dispose()

	for _, prop := range c.properties {
		prop.Dispose()
	}

	if c.ownWindow != nil {
		c.ownWindow.ControlRemoved()
	}

	c.ownWindow = nil
	c.parent = nil
	c.widget = nil
	if c.contextMenu != nil {
		c.contextMenu.DisposeMenu()
		c.contextMenu = nil
	}
}

func (c *Control) Init() {
}

func (c *Control) Name() string {
	return c.name.String()
}

func (c *Control) SetName(name string) {
	c.name.SetOwnValue(name)
}

func (c *Control) Update(source string) {
	if c.disposed {
		return
	}

	c.needToRedraw = true

	if c.isUpdating {
		return
	}

	if c.parent != nil {
		c.parent.Update(source)
	} else {
		if c.ownWindow != nil {
			c.ownWindow.UpdateWindow(source)
		}
	}
}

func (c *Control) ControlType() string {
	return "Control"
}

func (c *Control) SetTooltip(text string) {
	c.toolTip.SetOwnValue(text)
}

func (c *Control) Tooltip() string {
	return c.toolTip.ValueOwn().(string)
}

func (c *Control) SetTheme(theme string) {
	c.theme = theme
}

func AddPropertyToWidget(w Widget, name string, propertyType uiproperties.PropertyType) *uiproperties.Property {
	if w == nil {
		panic("No widget for property " + name)
	}
	p := uiproperties.NewProperty(name, propertyType)
	p.Name = name
	p.Init(name, w)
	w.AddProperty(name, p)
	return p
}

func AddUnstyledPropertyToWidget(w Widget, name string, propertyType uiproperties.PropertyType) *uiproperties.Property {
	if w == nil {
		panic("No widget for property " + name)
	}
	p := uiproperties.NewProperty(name, propertyType)
	p.Name = name
	p.Init(name, w)
	p.SetUnstyled(true)
	w.AddProperty(name, p)
	return p
}

func (c *Control) AddProperty(name string, prop *uiproperties.Property) {
	if c.properties == nil {
		c.properties = make(map[string]*uiproperties.Property)
	}
	c.properties[name] = prop
}

func (c *Control) SetX(x int) {
	c.x.SetOwnValue(x)
}

func (c *Control) SetY(y int) {
	c.y.SetOwnValue(y)
}

func (c *Control) SetWidth(width int) {
	delta := width - c.Width()
	c.width.SetOwnValue(width)

	// Logic correct scrolling after resize
	if c.leftBorderWidth != nil && c.rightBorderWidth != nil {
		if c.Width()-c.LeftBorderWidth()-c.RightBorderWidth()+c.scrollOffsetX > c.InnerWidth() && delta > 0 {
			c.scrollOffsetX -= delta
			if c.scrollOffsetX < 0 {
				c.scrollOffsetX = 0
			}
		}
	}
}

func (c *Control) SetHeight(height int) {
	delta := height - c.Height()
	c.height.SetOwnValue(height)

	// Logic correct scrolling after resize
	if c.topBorderWidth != nil && c.bottomBorderWidth != nil {
		if c.Height()-c.TopBorderWidth()-c.BottomBorderWidth()+c.scrollOffsetY > c.InnerHeight() && delta > 0 {
			c.scrollOffsetY -= delta
			if c.scrollOffsetY < 0 {
				c.scrollOffsetY = 0
			}
		}
	}
}

func (c *Control) SetAnchors(anchors int) {
	c.anchors.SetOwnValue(anchors)
}

func (c *Control) X() int {
	return c.x.Int()
}

func (c *Control) Y() int {
	return c.y.Int()
}

func (c *Control) Width() int {
	return c.width.Int()
}

func (c *Control) Height() int {
	return c.height.Int()
}

func (c *Control) ClientWidth() int {
	return c.Width() - c.LeftBorderWidth() - c.RightBorderWidth()
}

func (c *Control) ClientHeight() int {
	return c.Height() - c.TopBorderWidth() - c.BottomBorderWidth()
}

func (c *Control) InnerWidth() int {
	if c.innerSizeOverloaded {
		return c.innerWidthOverloaded
	}
	return c.width.Int() - c.leftBorderWidth.Int() - c.rightBorderWidth.Int()
}

func (c *Control) InnerHeight() int {
	if c.innerSizeOverloaded {
		return c.innerHeightOverloaded
	}
	return c.height.Int() - c.topBorderWidth.Int() - c.bottomBorderWidth.Int()
}

func (c *Control) LeftBorderWidth() int {
	return c.leftBorderWidth.Int()
}

func (c *Control) RightBorderWidth() int {
	return c.rightBorderWidth.Int()
}

func (c *Control) TopBorderWidth() int {
	return c.topBorderWidth.Int()
}

func (c *Control) BottomBorderWidth() int {
	return c.bottomBorderWidth.Int()
}

func (c *Control) FontFamily() string {
	return c.fontFamily.String()
}

func (c *Control) FontSize() float64 {
	return c.fontSize.Float64()
}

func (c *Control) SetFontSize(fontSize float64) {
	c.fontSize.SetOwnValue(fontSize)
	if c.Window() != nil {
		c.Window().UpdateLayout()
	}
}

func (c *Control) FontBold() bool {
	return c.fontBold.Bool()
}

func (c *Control) FontItalic() bool {
	return c.fontItalic.Bool()
}

func (c *Control) Anchors() int {
	return c.anchors.Int()
}

func (c *Control) SetHover(hover bool) {
	if hover {
		if c.cursor != MouseCursorNotDefined && hover {
			c.Window().SetMouseCursor(c.cursor)
		}
	} else {
		if c.cursor != MouseCursorNotDefined {
			c.Window().SetMouseCursor(MouseCursorArrow)
		}
	}
	c.hover = hover
}

func (c *Control) MouseCursor() MouseCursor {
	return c.cursor
}

func (c *Control) FindWidgetUnderPointer(x, y int) Widget {
	return nil
}

func (c *Control) ProcessFindWidgetUnderPointer(x, y int) Widget {
	cX := x
	cY := y

	cX = cX - c.widget.X()
	cY = cY - c.widget.Y()

	cX -= c.widget.LeftBorderWidth()
	cY -= c.widget.TopBorderWidth()
	cX += c.widget.ScrollOffsetX()
	cY += c.widget.ScrollOffsetY()

	return c.widget.FindWidgetUnderPointer(cX, cY)
}

func (c *Control) Hover() bool {
	return c.hover
}

func (c *Control) ClearHover() {
	if c.cursor != MouseCursorNotDefined {
		c.Window().SetMouseCursor(MouseCursorArrow)
	}
	c.hover = false
}

func (c *Control) Focus() {
	c.ownWindow.SetFocusForWidget(c.widget)
	c.Update("Control")
}

func (c *Control) SetFocus(focus bool) {
	c.focus = focus
}

func (c *Control) HasFocus() bool {
	return c.focus
}

func (c *Control) BackColor() color.Color {
	return c.backgroundColor.Color().(color.Color)
}

func (c *Control) ForeColor() color.Color {
	return c.foregroundColor.Color()
}

func (c *Control) SetBackColor(backColor color.Color) {
	c.backgroundColor.SetOwnValue(backColor)
}

func (c *Control) SetForeColor(foreColor color.Color) {
	c.foregroundColor.SetOwnValue(foreColor)
}

func (c *Control) AccentColor() color.Color {
	return c.accentColor.Color()
}

func (c *Control) InactiveColor() color.Color {
	return c.inactiveColor.Color()
}

func (c *Control) ClearFocus() {
	c.focus = false
	c.redraw()
}

func (c *Control) OnScroll(scrollPositionX int, scrollPositionY int) {

}

func (c *Control) CloseTopPopup() {
}

func (c *Control) ProcessMouseDown(event *MouseDownEvent) {
	if !c.enabled {
		return
	}

	me := *event

	me.X = me.X - c.widget.X()
	me.Y = me.Y - c.widget.Y()

	// ScrollBars
	processed := false

	vRect := c.vScroll

	if vRect.Contains(me.X, me.Y) {
		c.verticalScrollMoving = true
		c.verticalScrollMovingMousePos = me.Y
		c.verticalScrollMovingOriginal = c.scrollOffsetY

		//fmt.Println("VSCROLL INIT: ", me.Y, c.scrollOffsetY)

		processed = true
		me.SetUserData("processedWidget", c)
	}

	hRect := c.hScroll

	if hRect.Contains(me.X, me.Y) {
		c.horizontalScrollMoving = true
		c.horizontalScrollMovingMousePos = me.X
		c.horizontalScrollMovingOriginal = c.scrollOffsetX

		processed = true
		me.SetUserData("processedWidget", c)
	}

	me.X -= c.widget.LeftBorderWidth()
	me.Y -= c.widget.TopBorderWidth()
	me.X += c.widget.ScrollOffsetX()
	me.Y += c.widget.ScrollOffsetY()

	if !processed {
		if c.onPreMouseDown != nil {
			c.onPreMouseDown(&me)
		}
		if !me.Ignore {
			c.widget.MouseDown(&me)
			if c.onMouseDown != nil {
				c.onMouseDown(&me)
			}
		}
	}
}

func (c *Control) ProcessMouseUp(event *MouseUpEvent) {
	if !c.enabled {
		return
	}

	me := *event

	me.X = me.X - c.widget.X()
	me.Y = me.Y - c.widget.Y()

	me.X -= c.widget.LeftBorderWidth()
	me.Y -= c.widget.TopBorderWidth()
	me.X += c.widget.ScrollOffsetX()
	me.Y += c.widget.ScrollOffsetY()

	c.verticalScrollMoving = false
	c.horizontalScrollMoving = false

	if c.onPreMouseUp != nil {
		c.onPreMouseUp(&me)
	}
	if !me.Ignore {
		c.widget.MouseUp(&me)
		if c.onMouseUp != nil {
			c.onMouseUp(&me)
		}
	}
}

func (c *Control) ProcessMouseMove(event *MouseMoveEvent) {
	if !c.enabled {
		return
	}

	me := *event

	me.X = me.X - c.widget.X()
	me.Y = me.Y - c.widget.Y()

	processed := false
	if c.verticalScrollMoving {
		deltaMouse := me.Y - c.verticalScrollMovingMousePos
		heightIn := c.ClientHeight()
		barSizeK := float64(heightIn) / float64(c.InnerHeight())
		delta := int(float64(deltaMouse) / barSizeK)

		c.scrollOffsetY = c.verticalScrollMovingOriginal + delta

		if c.scrollOffsetY < 0 {
			c.scrollOffsetY = 0
		}
		if c.scrollOffsetY > c.InnerHeight()-heightIn {
			c.scrollOffsetY = c.InnerHeight() - heightIn
		}
		processed = true
		if c.onScrolled != nil {
			c.onScrolled(c.scrollOffsetX, c.scrollOffsetY)
		}
	}

	if c.horizontalScrollMoving {
		deltaMouse := event.X - c.horizontalScrollMovingMousePos

		widthIn := c.Width() - c.LeftBorderWidth() - c.RightBorderWidth()
		barSizeK := float64(widthIn) / float64(c.InnerWidth())
		delta := int(float64(deltaMouse) / barSizeK)

		c.scrollOffsetX = c.horizontalScrollMovingOriginal + delta

		if c.scrollOffsetX < 0 {
			c.scrollOffsetX = 0
		}
		if c.scrollOffsetX > c.InnerWidth()-widthIn {
			c.scrollOffsetX = c.InnerWidth() - widthIn
		}

		c.widget.OnScroll(c.scrollOffsetX, c.scrollOffsetY)

		processed = true
		if c.onScrolled != nil {
			c.onScrolled(c.scrollOffsetX, c.scrollOffsetY)
		}
	}

	me.X -= c.widget.LeftBorderWidth()
	me.Y -= c.widget.TopBorderWidth()
	me.X += c.widget.ScrollOffsetX()
	me.Y += c.widget.ScrollOffsetY()

	if !processed {
		if c.onPreMouseMove != nil {
			c.onPreMouseMove(&me)
		}
		if !me.Ignore {
			c.widget.MouseMove(&me)
			if c.onMouseMove != nil {
				c.onMouseMove(&me)
			}
		}
	}
}

func (c *Control) ProcessMouseClick(event *MouseClickEvent) {
	if !c.enabled {
		return
	}

	me := *event

	me.X = me.X - c.widget.X()
	me.Y = me.Y - c.widget.Y()

	me.X -= c.widget.LeftBorderWidth()
	me.Y -= c.widget.TopBorderWidth()
	me.X += c.widget.ScrollOffsetX()
	me.Y += c.widget.ScrollOffsetY()

	if c.onPreMouseClick != nil {
		c.onPreMouseClick(&me)
	}
	if !me.Ignore {
		contextMenuFound := false

		if event.Button == MouseButtonRight {
			wX, wY := c.RectClientAreaOnWindow()
			if c.ContextMenu() != nil {
				c.ContextMenu().ShowMenu(wX+me.X-c.ScrollOffsetX(), wY+me.Y-c.ScrollOffsetY())
				contextMenuFound = true
			} else {
				if c.OnContextMenuNeed != nil {
					m := c.OnContextMenuNeed(me.X, me.Y)
					if m != nil {
						m.ShowMenu(wX+me.X-c.ScrollOffsetX(), wY+me.Y-c.ScrollOffsetY())
						contextMenuFound = true
					}
				}
			}
		}
		if !contextMenuFound {
			c.widget.MouseClick(&me)
			if c.onMouseClick != nil {
				c.onMouseClick(&me)
			}
		}
	}
}

func (c *Control) ProcessMouseDblClick(event *MouseDblClickEvent) {
	if !c.enabled {
		return
	}

	me := *event

	me.X = me.X - c.widget.X()
	me.Y = me.Y - c.widget.Y()

	me.X -= c.widget.LeftBorderWidth()
	me.Y -= c.widget.TopBorderWidth()
	me.X += c.widget.ScrollOffsetX()
	me.Y += c.widget.ScrollOffsetY()

	if c.onPreMouseDblClick != nil {
		c.onPreMouseDblClick(&me)
	}
	if !me.Ignore {
		c.widget.MouseDblClick(&me)
		if c.onMouseDblClick != nil {
			c.onMouseDblClick(&me)
		}
	}
}

func (c *Control) ProcessMouseWheel(event *MouseWheelEvent) {
	if !c.enabled {
		return
	}

	me := *event

	me.X = me.X - c.widget.X()
	me.Y = me.Y - c.widget.Y()

	me.X -= c.widget.LeftBorderWidth()
	me.Y -= c.widget.TopBorderWidth()
	me.X += c.widget.ScrollOffsetX()
	me.Y += c.widget.ScrollOffsetY()

	if c.innerSizeOverloaded {
		widthIn := c.Width() - c.LeftBorderWidth() - c.RightBorderWidth()
		if c.InnerWidth() > widthIn {
		}
		heightIn := c.Height() - c.TopBorderWidth() - c.BottomBorderWidth()
		if c.InnerHeight() > heightIn {
			c.scrollOffsetY -= me.Delta * 20
			if c.scrollOffsetY < 0 {
				c.scrollOffsetY = 0
			}
			if c.scrollOffsetY > c.InnerHeight()-heightIn {
				c.scrollOffsetY = c.InnerHeight() - heightIn
			}
		}
		return
	}

	if c.onPreMouseWheel != nil {
		c.onPreMouseWheel(&me)
	}
	if !me.Ignore {
		c.widget.MouseWheel(&me)
		if c.onMouseWheel != nil {
			c.onMouseWheel(&me)
		}
	}
}

func (c *Control) ProcessKeyChar(event *KeyCharEvent) {
	if !c.enabled {
		return
	}

	if c.onPreKeyChar != nil {
		c.onPreKeyChar(event)
	}
	if !event.Ignore {
		c.widget.KeyChar(event)
		if c.onKeyChar != nil {
			c.onKeyChar(event)
		}
	}
}

func (c *Control) SetOnKeyDown(callback func(event *KeyDownEvent) bool) {
	c.onKeyDown = callback
}

func (c *Control) ProcessKeyDown(event *KeyDownEvent) bool {
	if !c.enabled {
		return false
	}

	processed := false
	if c.onPreKeyDown != nil {
		c.onPreKeyDown(event)
		processed = true
	}
	if !event.Ignore {
		processed = c.widget.KeyDown(event)
		if c.onKeyDown != nil {
			c.onKeyDown(event)
		}
	}
	return processed
}

func (c *Control) ProcessKeyUp(event *KeyUpEvent) {
	if !c.enabled {
		return
	}

	if c.onPreKeyUp != nil {
		c.onPreKeyUp(event)
	}
	if !event.Ignore {
		c.widget.KeyUp(event)
		if c.onKeyUp != nil {
			c.onKeyUp(event)
		}
	}
}

func (c *Control) MouseWheel(event *MouseWheelEvent) {
	if c.innerSizeOverloaded {
		widthIn := c.Width() - c.LeftBorderWidth() - c.RightBorderWidth()
		if c.InnerWidth() > widthIn {
		}
		heightIn := c.Height() - c.TopBorderWidth() - c.BottomBorderWidth()
		if c.InnerHeight() > heightIn {
			c.scrollOffsetY++
		}
	}
}

func (c *Control) MouseMove(event *MouseMoveEvent) {
}

func (c *Control) MouseDown(event *MouseDownEvent) {
}

func (c *Control) MouseUp(event *MouseUpEvent) {
}

func (c *Control) MouseDrop(event *MouseDropEvent) {
	if c.OnMouseDrop != nil {
		c.OnMouseDrop(event.DroppingObject, event.X, event.Y)
	}
}

func (c *Control) MouseValidateDrop(event *MouseValidateDropEvent) {
	event.AllowDrop = true
}

func (c *Control) MouseClick(event *MouseClickEvent) {
}

func (c *Control) MouseDblClick(event *MouseDblClickEvent) {
}

func (c *Control) KeyChar(event *KeyCharEvent) {
}

func (c *Control) KeyDown(event *KeyDownEvent) bool {
	return false
}

func (c *Control) KeyUp(event *KeyUpEvent) {
}

func (c *Control) DrawBorders(ctx DrawContext) {

	// Borders
	if c.topBorderWidth.Int() > 0 {
		ctx.SetColor(c.topBorderColor.Color())
		ctx.FillRect(0, 0, c.width.Int(), c.topBorderWidth.Int())
	}
	if c.bottomBorderWidth.Int() > 0 {
		ctx.SetColor(c.bottomBorderColor.Color())
		ctx.FillRect(0, c.height.Int()-c.bottomBorderWidth.Int(), c.width.Int(), c.bottomBorderWidth.Int())
	}
	if c.leftBorderWidth.Int() > 0 {
		ctx.SetColor(c.leftBorderColor.Color())
		ctx.FillRect(0, 0, c.leftBorderWidth.Int(), c.height.Int())
	}
	if c.rightBorderWidth.Int() > 0 {
		ctx.SetColor(c.rightBorderColor.Color())
		ctx.FillRect(c.width.Int()-c.rightBorderWidth.Int(), 0, c.rightBorderWidth.Int(), c.height.Int())
	}
}

func (c *Control) DrawScrollBars(ctx DrawContext) {
	c.verticalScrollDisplayed = false
	c.horizontalScrollDisplayed = false
	if c.innerSizeOverloaded {

		scrollBarsColor := ColorWithAlpha(c.leftBorderColor.Color(), 192)

		c.vScroll = c.verticalScrollRect()
		if c.verticalScrollDisplayed {
			ctx.SetColor(scrollBarsColor)
			ctx.FillRect(c.vScroll.X, c.vScroll.Y, c.vScroll.Width, c.vScroll.Height)
		} else {
			//c.scrollOffsetY = 0
		}

		c.hScroll = c.horizontalScrollRect()
		if c.horizontalScrollDisplayed {
			ctx.SetColor(scrollBarsColor)
			ctx.FillRect(c.hScroll.X, c.hScroll.Y, c.hScroll.Width, c.hScroll.Height)
		} else {
			//c.scrollOffsetX = 0
		}
	}
}

func (c *Control) horizontalScrollRect() Rect {
	var result Rect
	widthIn := c.Width() - c.LeftBorderWidth() - c.RightBorderWidth()
	innW := c.InnerWidth()
	if innW > widthIn {
		if c.horizontalScrollVisible.Bool() {
			w := c.horizontalScrollWidth.Int()
			barSizeK := float64(widthIn) / float64(c.InnerWidth())
			barSize := int(barSizeK * float64(widthIn))
			barOffset := int(float64(c.ScrollOffsetX()) * barSizeK)
			result.X = c.LeftBorderWidth() + barOffset
			result.Y = c.Height() - c.BottomBorderWidth() - w
			result.Width = barSize
			result.Height = w
			c.horizontalScrollDisplayed = true
		}
	}
	return result
}

func (c *Control) SetVerticalScrollVisible(verticalScrollVisible bool) {
	c.verticalScrollVisible.SetOwnValue(verticalScrollVisible)
}

func (c *Control) SetHorizontalScrollVisible(horizontalScrollVisible bool) {
	c.horizontalScrollVisible.SetOwnValue(horizontalScrollVisible)
}

func (c *Control) verticalScrollRect() Rect {
	var result Rect
	heightIn := c.Height() - c.TopBorderWidth() - c.BottomBorderWidth()
	if c.InnerHeight() > heightIn {
		if c.verticalScrollVisible.Bool() {
			w := c.verticalScrollWidth.Int()
			barSizeK := float64(heightIn) / float64(c.InnerHeight())
			barSize := int(barSizeK * float64(heightIn))
			barOffset := int(float64(c.ScrollOffsetY()) * barSizeK)
			result.X = c.Width() - c.RightBorderWidth() - w
			result.Y = c.TopBorderWidth() + barOffset
			result.Width = w
			result.Height = barSize
			c.verticalScrollDisplayed = true
		}
	}
	return result
}

func (c *Control) DrawBackground(ctx DrawContext) {
	_, _, _, a := c.BackColor().RGBA()
	if a < 1 {
		return
	}
	ctx.SetColor(c.BackColor())
	ctx.FillRect(0, 0, c.Width(), c.Height())
}

func (c *Control) ScrollOffsetX() int {
	return c.scrollOffsetX
}

func (c *Control) ScrollOffsetY() int {
	return c.scrollOffsetY
}

func (c *Control) ScrollEnsureVisible(x1, y1 int) {

	if y1 < c.scrollOffsetY {
		c.scrollOffsetY = y1
	}
	if y1 > c.scrollOffsetY+c.ClientHeight() {
		c.scrollOffsetY = y1 - c.ClientHeight()
	}

	if x1 < c.scrollOffsetX {
		c.scrollOffsetX = x1
	}
	if x1 > c.scrollOffsetX+c.ClientWidth() {
		c.scrollOffsetX = x1 - c.ClientWidth()
	}
}

func drawCopyOver(dst *image.RGBA, src *image.RGBA, x int, y int) {
	height := src.Rect.Max.Y
	for yy := y; yy < y+height; yy++ {
		srcIndex := (yy - y) * src.Stride
		srcIndexEnd := srcIndex + src.Stride
		destIndex := yy*dst.Stride + x*4
		destIndexEnd := destIndex + src.Stride
		if srcIndex > len(src.Pix) || srcIndexEnd > len(src.Pix) {
			continue
		}
		if destIndex > len(dst.Pix) || destIndexEnd > len(dst.Pix) {
			continue
		}
		copy(dst.Pix[destIndex:destIndexEnd], src.Pix[srcIndex:srcIndexEnd])
	}
}

func (c *Control) Subclass() string {
	result := "default"

	if !c.enabled {
		return "disabled"
	}

	if c.hover {
		result = "hover"
	} else {
		if c.focus {
			result = "focus"
		}
	}
	return result
}

func (c *Control) Theme() string {
	return c.theme
}

func (c *Control) MouseEnter() {
	if c.onMouseEnter != nil {
		c.onMouseEnter()
	}
}

func (c *Control) MouseLeave() {
	if c.onMouseLeave != nil {
		c.onMouseLeave()
	}
}

func (c *Control) FocusChanged(focus bool) {
	c.redraw()
}

func (c *Control) ApplyStyleLine(controlName string, controlType string, styleClass string, stylePseudoClass string, propertyName string, value string) {
	score := 0

	if len(controlType) > 0 && controlType != c.widget.ControlType() && controlType != "Control" {
		return
	}

	if len(controlName) > 0 && controlName != c.widget.Name() {
		return
	}

	if controlType == c.widget.ControlType() {
		score += 10
	}

	if controlType == "Control" {
		score += 5
	}

	/*for _, controlClass := range c.widget.Classes() {
		if controlClass == styleClass {
			score += 10
		}
	}*/

	if score >= c.widget.CurrentStyleValueScore(stylePseudoClass, propertyName) {
		c.widget.SetStyledValue(stylePseudoClass, propertyName, value, score)
	}
}

func (c *Control) CurrentStyleValueScore(subclass string, propertyName string) int {
	if val, ok := c.properties[propertyName]; ok {
		return val.ValueStyleScore(subclass)
	}
	return 0
}

func (c *Control) SetStyledValue(subclass string, propertyName string, value string, score int) {
	if val, ok := c.properties[propertyName]; ok {
		val.SetStyledValue(subclass, value, score)
	}
}

func (c *Control) StyledValue(subclass string, propertyName string) interface{} {
	if val, ok := c.properties[propertyName]; ok {
		return val.ValueStyle(subclass)
	}
	return nil
}

func (c *Control) Classes() []string {
	classes := make([]string, 0)
	return classes
}

func (c *Control) ClearRadioButtons() {
}

func (c *Control) SetParent(p Widget) {
	c.parent = p
}

func (c *Control) RectOnWindow() (int, int) {
	var x, y int

	if c.parent != nil {
		xx, yy := c.parent.RectClientAreaOnWindow()
		x += xx
		y += yy

		x -= c.parent.ScrollOffsetX()
		y -= c.parent.ScrollOffsetY()
	}

	return x, y
}

func (c *Control) RectClientAreaOnWindow() (int, int) {
	var x, y int
	x = c.X() + c.LeftBorderWidth()
	y = c.Y() + c.TopBorderWidth()

	if c.parent != nil {
		xx, yy := c.parent.RectClientAreaOnWindow()
		x += xx
		y += yy

		x -= c.parent.ScrollOffsetX()
		y -= c.parent.ScrollOffsetY()
	}

	return x, y
}

func (c *Control) Window() Window {
	if c.parent != nil {
		return c.parent.Window()
	}
	return c.ownWindow
}

func (c *Control) TranslateX(x int) int {
	return x - c.X() + c.ScrollOffsetX() - c.LeftBorderWidth()
}

func (c *Control) TranslateY(y int) int {
	return y - c.Y() + c.ScrollOffsetY() - c.TopBorderWidth()
}

func (c *Control) AcceptsReturn() bool {
	return false
}

func (c *Control) AcceptsTab() bool {
	return false
}

func (c *Control) NextFocusControl() Widget {
	return nil
}

func (c *Control) FirstFocusControl() Widget {
	return nil
}

func (c *Control) translateMouseEvent(event *MouseEvent) {
	event.X = c.TranslateX(event.X)
	event.Y = c.TranslateX(event.Y)
}

func (c *Control) SetContextMenu(menu IMenu) {
	c.contextMenu = menu
}

func (c *Control) ContextMenu() IMenu {
	return c.contextMenu
}

func (c *Control) SetWindow(w Window) {
	c.ownWindow = w
}

func (c *Control) redraw() {
	c.needToRedraw = true
}

func (c *Control) SetMouseCursor(cursor MouseCursor) {
	c.cursor = cursor
}

func (c *Control) VisibleInnerRect() Rect {
	var r Rect
	r.X = c.scrollOffsetX
	r.Y = c.scrollOffsetY
	r.Width = c.Width() + c.scrollOffsetX
	r.Height = c.Height() + c.scrollOffsetY
	return r
}

func (c *Control) BeginUpdate() {
	c.isUpdating = true
}

func (c *Control) EndUpdate() {
	c.isUpdating = false
	c.Update("Control")
}

func (c *Control) BeginDrag(draggingObject interface{}) {
	if c.ownWindow != nil {
		c.ownWindow.BeginDrag(draggingObject)
	}
}

func (c *Control) ClosePopup() {
}

func (c *Control) GridX() int {
	return c.gridX
}

func (c *Control) GridY() int {
	return c.gridY
}

func (c *Control) SetGridX(gridX int) {
	c.gridX = gridX
}

func (c *Control) SetGridY(gridY int) {
	c.gridY = gridY
}

func (c *Control) MinWidth() int {
	return c.minWidth
}

func (c *Control) MinHeight() int {
	return c.minHeight
}

func (c *Control) MaxWidth() int {
	return c.maxWidth
}

func (c *Control) MaxHeight() int {
	return c.maxHeight
}

func (c *Control) SetXExpandable(xExpandable bool) {
	c.xExpandable = xExpandable
}

func (c *Control) SetYExpandable(yExpandable bool) {
	c.yExpandable = yExpandable
}

func (c *Control) XExpandable() bool {
	return c.xExpandable
}

func (c *Control) YExpandable() bool {
	return c.yExpandable
}

func (c *Control) SetMinWidth(minWidth int) {
	c.minWidth = minWidth
}

func (c *Control) SetMinHeight(minHeight int) {
	c.minHeight = minHeight
}

func (c *Control) SetMaxWidth(maxWidth int) {
	c.maxWidth = maxWidth
}

func (c *Control) SetMaxHeight(maxHeight int) {
	c.maxHeight = maxHeight
}

func (c *Control) Initialized() bool {
	return c.isInited
}

func (c *Control) SetInnerSizeDirect(w int, h int) {
	c.innerWidthOverloaded = w
	c.innerHeightOverloaded = h
	c.innerSizeOverloaded = true
	c.verticalScrollVisible.SetOwnValue(true)
	c.horizontalScrollVisible.SetOwnValue(true)
}

func (c *Control) ResetInnerSizeDirect() {
	c.innerSizeOverloaded = false
	c.verticalScrollVisible.SetOwnValue(false)
	c.horizontalScrollVisible.SetOwnValue(false)
}

func (c *Control) SetPos(x, y int) {
	c.SetX(x)
	c.SetY(y)
}

func (c *Control) SetSize(w, h int) {
	c.SetWidth(w)
	c.SetHeight(h)
}

func (c *Control) SetFixedSize(w int, h int) {
	c.widget.SetMinWidth(w)
	c.widget.SetMaxWidth(w)
	c.widget.SetMinHeight(h)
	c.widget.SetMaxHeight(h)
}

func (c *Control) SetBorderLeft(width int, col color.Color) {
	c.leftBorderWidth.SetOwnValue(width)
	c.leftBorderColor.SetOwnValue(col)
}

func (c *Control) SetBorderRight(width int, col color.Color) {
	c.rightBorderWidth.SetOwnValue(width)
	c.rightBorderColor.SetOwnValue(col)
}

func (c *Control) SetBorderTop(width int, col color.Color) {
	c.topBorderWidth.SetOwnValue(width)
	c.topBorderColor.SetOwnValue(col)
}

func (c *Control) SetBorderBottom(width int, col color.Color) {
	c.bottomBorderWidth.SetOwnValue(width)
	c.bottomBorderColor.SetOwnValue(col)
}

func (c *Control) SetBorders(width int, col color.Color) {
	c.SetBorderLeft(width, col)
	c.SetBorderRight(width, col)
	c.SetBorderTop(width, col)
	c.SetBorderBottom(width, col)
}

func (c *Control) Draw(ctx DrawContext) {
	ctx.SetColor(colornames.Black)
	ctx.DrawRect(1, 1, c.Width()-2, c.Height()-2)
	ctx.DrawText(5, 0, c.Width()-1, c.Height()-1, c.ControlType())
}

func (c *Control) DrawControl(ctx DrawContext) {
	if !c.widget.IsVisible() {
		return
	}

	ctx.Save()
	ctx.Translate(c.widget.X(), c.widget.Y())

	c.widget.DrawBackground(ctx)

	ctx.Save()
	ctx.Translate(c.widget.LeftBorderWidth()-c.widget.ScrollOffsetX(), c.widget.TopBorderWidth()-c.widget.ScrollOffsetY())

	clipX := c.widget.ScrollOffsetX()
	clipY := c.widget.ScrollOffsetY()
	clipW := c.widget.Width() - c.widget.LeftBorderWidth() - c.widget.RightBorderWidth() + c.widget.ScrollOffsetX()
	clipH := c.widget.Height() - c.widget.TopBorderWidth() - c.widget.BottomBorderWidth() + c.widget.ScrollOffsetY()

	if clipX < 0 {
		clipW += clipX
		clipX = 0
	}

	if clipY < 0 {
		clipH += clipY
		clipY = 0
	}

	if clipW > c.widget.Width()-c.widget.LeftBorderWidth()-c.widget.RightBorderWidth() {
		clipW = c.widget.Width() - c.widget.LeftBorderWidth() - c.widget.RightBorderWidth()
	}

	if clipH > c.widget.Height()-c.widget.TopBorderWidth()-c.widget.BottomBorderWidth() {
		clipH = c.widget.Height() - c.widget.TopBorderWidth() - c.widget.BottomBorderWidth()
	}

	ctx.ClipIn(clipX, clipY, clipW, clipH)
	c.widget.Draw(ctx)
	ctx.Load()

	c.widget.DrawBorders(ctx)
	c.widget.DrawScrollBars(ctx)

	ctx.Load()
}

/*func (c *Control) DrawGL(ctx *DrawContextGL) {
	ctx.SetColor(colornames.Blue)
	ctx.FillRect(0, 0, c.Width(), c.Height())

	//UpdateResolution(c.Width(), c.Height())
	//DrawTextInRect(c.Window().Window(), 0, 0, c.Width(), c.Height(), c.ControlType(), 5)
	//DrawText(c.Window().Window(), c.Width() / 2, c.Height() / 2, c.ControlType(), 5)
}
*/

func (c *Control) String(level int) string {
	return strings.Repeat("    ", level) + fmt.Sprint("{", c.widget.ControlType(), " X:", c.X(), " Y:", c.Y(), " W:", c.Width(), " H:", c.Height(), ")", "}")
}

func (c *Control) UpdateLayout() {
}

func (c *Control) Text() string {
	return ""
}

func (c *Control) ClearLayoutCache() {
}

func (c *Control) BorderColors() (color.Color, color.Color, color.Color, color.Color) {
	return c.leftBorderColor.Color(), c.topBorderColor.Color(), c.rightBorderColor.Color(), c.bottomBorderColor.Color()
}
