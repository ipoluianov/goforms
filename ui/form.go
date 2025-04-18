package ui

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"runtime"
	"runtime/debug"
	"sort"
	"time"

	"github.com/ipoluianov/goforms/utils"
	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/nui/nui"
	"github.com/ipoluianov/nui/nuikey"
	"github.com/ipoluianov/nui/nuimouse"
	"golang.org/x/image/colornames"
)

var nextFormId int
var windows []Window

var windowByGLFWWindow map[nui.Window]Window

const (
	DefaultWindowWidth  = 640
	DefaultWindowHeight = 480
)

func (c *Form) LoopUI() {
	c.LoopUI_OpenGL()
}

func NewForm() *Form {
	var form Form
	form.Init()
	return &form
}

func UnInitUI() {
	//glfw.Terminate()
}

func init() {
	nextFormId = 1
	windows = make([]Window, 0)
	windowByGLFWWindow = make(map[nui.Window]Window)
}

type Form struct {
	utils.Obj
	id       int
	disposed bool
	inited   bool

	window nui.Window

	userPanel     *Panel
	menu          Menu
	width         int
	height        int
	title         string
	isMainWindow  bool
	parent        Window
	modal         bool
	popup         bool
	showMaximazed bool

	position Point

	mouseDownLastPosX   int
	mouseDownLastPosY   int
	lastMouseMoveTime   time.Time
	lastMouseMovePos    Point
	focusWidget         Widget
	hoverWidget         Widget
	lastMouseDownWidget Widget

	lastDrawTime time.Time
	needToUpdate bool

	keyModifiers nuikey.KeyModifiers

	//drawTime     []float64
	dialogResult bool

	formTimers []*FormTimer

	lastUpdateSource string

	toolTipControlProcessed bool

	onSizeChanged func(event *FormSizeChangedEvent)

	draggingObject interface{}

	acceptButton *Button
	rejectButton *Button

	currentTooltipX    int
	currentTooltipY    int
	currentTooltipText string

	activatedServiceMenu bool

	//currentCanvas    DrawContext
	//currentCanvasKey string

	popupWindow Widget

	childModal Window

	needUpdateLayout bool

	drawTimes      []time.Duration
	drawTimesIndex int
	drawTimesCount int

	ignoreUpdatesCounter int

	OnKeyDown func(event *KeyDownEvent) bool
}

func init() {
	//uiplatforms.InitUI()
}

func (c *Form) Init() {
	if c.inited {
		return
	}
	c.inited = true
	c.ProcessWindowResize(DefaultWindowWidth, DefaultWindowHeight)
	c.drawTimesCount = 5
	c.drawTimes = make([]time.Duration, c.drawTimesCount)

	c.userPanel = NewRootPanel(c)
	c.userPanel.SetBackColor(DefaultBackColor)
	c.userPanel.SetPos(0, 0)
	c.userPanel.SetSize(c.Width(), c.Height())
	c.userPanel.ownWindow = c
	c.userPanel.SetName("MainPanelOfForm")
	c.userPanel.SetIsTabPlate(true)
	//f.userPanel.SetName("MainPanel")
	//f.userPanel.SetAnchors(ANCHOR_ALL)

	//f.menu = NewMenu(f.Panel())
	//f.menu.SetAnchors(ANCHOR_LEFT | ANCHOR_RIGHT | ANCHOR_TOP)
	//f.windowPanel.AddWidget(f.menu)
	c.formTimers = make([]*FormTimer, 0)

	timer := c.NewTimer(1000, c.timerMemoryDump)
	timer.StartTimer()

	c.Obj.InitObj("Form", "Form_"+time.Now().String()+"_"+fmt.Sprint(rand.Int()))
	runtime.SetFinalizer(c, finalizerForm)

	/*vertShader, _ = gfx.NewShaderFromFile("shaders/basic.vert", gl.VERTEX_SHADER)

	fragShader, _ = gfx.NewShaderFromFile("shaders/basic.frag", gl.FRAGMENT_SHADER)

	shaderProgram, _ = gfx.NewProgram(vertShader, fragShader)*/
}

func finalizerForm(c *Form) {
	c.Obj.UninitObj()
}

func (c *Form) NewTimer(period int64, handler func()) *FormTimer {
	var timer FormTimer
	timer.Enabled = false
	timer.Period = period
	timer.LastElapsedDTMSec = 0
	timer.Handler = handler
	c.formTimers = append(c.formTimers, &timer)
	return &timer
}

func (c *Form) MakeTimerAndStart(period int64, handler func(timer *FormTimer)) *FormTimer {
	timer := c.NewTimer(period, nil)
	timer.Handler = func() {
		if handler != nil {
			handler(timer)
		}
	}
	timer.StartTimer()
	return timer
}

func (c *Form) RemoveTimer(timer *FormTimer) {
	for index, t := range c.formTimers {
		if t == timer {
			c.formTimers = append(c.formTimers[:index], c.formTimers[index+1:]...)
			break
		}
	}
}

func (c *Form) UpdateStyle() {
	c.Panel().UpdateStyle()
	c.Panel().SetBackColor(DefaultBackColor)
	c.UpdateLayout()
}

func (c *Form) SetShowMaximazed(maximazed bool) {
	c.showMaximazed = maximazed
}

func (c *Form) String() string {
	result := ""
	result += c.Panel().String(0)
	return result
}

func (c *Form) ShowMaximazed() bool {
	return c.showMaximazed
}

func (c *Form) IgnoreUpdates() {
	c.ignoreUpdatesCounter++
}

func (c *Form) UnIgnoreUpdates() {
	c.ignoreUpdatesCounter--
	if c.ignoreUpdatesCounter < 0 {
		c.ignoreUpdatesCounter = 0
	}
	if c.ignoreUpdatesCounter == 0 {
		c.UpdateWindow("UnIgnoreUpdates")
	}
}

func (c *Form) UpdateWindow(source string) {
	if c.ignoreUpdatesCounter > 0 {
		return
	}

	dt := time.Now().Format("2006-01-02 15:04:05.999999999")
	fmt.Println(dt, "UpdateWindow", source)
	c.needToUpdate = true
	c.lastUpdateSource = source
	if c.window != nil {
		c.window.Update()
	}
}

func (c *Form) Accept() {
	c.dialogResult = true
	c.Close()
}

func (c *Form) Reject() {
	c.dialogResult = false
	c.Close()
}

func (c *Form) SetIcon(img image.Image) {
	//images := make([]image.Image, 0)
	//images = append(images, img)
	//c.window.SetAppIcon(images[0])
}

func (c *Form) AcceptButton(ev *Event) {
	c.Accept()
}

func (c *Form) RejectButton(ev *Event) {
	c.Reject()
}

func (c *Form) DialogResult() bool {
	return c.dialogResult
}

func init() {
}

func (c *Form) Window() nui.Window {
	return c.window
}

func (c *Form) Draw(rgba *image.RGBA) bool {

	if c.window == nil {
		return false
	}
	if c.disposed {
		return false
	}

	c.lastDrawTime = time.Now()

	t1 := time.Now()
	c.realUpdateLayout()

	ctx := NewDrawContextSWRGBA(c.window, rgba)

	ctx.Init()

	ctx.Translate(0, 0)

	if c.userPanel != nil {
		c.userPanel.DrawControl(ctx)
	}

	c.drawDraggingObject(ctx)
	c.drawTooltip(ctx)
	if c.popupWindow != nil {
		c.popupWindow.Draw(ctx)
	}

	if ServiceDrawBorders {
		ctx.SetColor(colornames.Blue)
		ctx.SetTextAlign(canvas.HAlignLeft, canvas.VAlignTop)
		//ctx.DrawText(0, 0, c.Width(), c.Height(), c.String())
	}

	ctx.Finish(rgba)

	t2 := time.Now()

	drawTime := t2.Sub(t1)
	c.drawTimes[c.drawTimesIndex] = drawTime
	c.drawTimesIndex++
	if c.drawTimesIndex >= c.drawTimesCount {
		c.drawTimesIndex = 0
	}

	c.needToUpdate = false
	return true
}

func (c *Form) drawDraggingObject(ctx DrawContext) {
	if c.draggingObject != nil {
		ctx.SetColor(colornames.Crimson)
		ctx.SetStrokeWidth(1)
		ctx.DrawRect(c.lastMouseMovePos.X, c.lastMouseMovePos.Y, 10, 10)
	}
}

func (c *Form) drawTooltip(ctx DrawContext) {

	toolTipText := c.currentTooltipText
	toolTipX := c.currentTooltipX + 10
	toolTipY := c.currentTooltipY + 10

	if time.Since(c.lastMouseMoveTime) > 200*time.Millisecond {
		if c.hoverWidget != nil {
			if c.hoverWidget.Tooltip() != "" {
				toolTipX = c.lastMouseMovePos.X + 10
				toolTipY = c.lastMouseMovePos.Y + 10
				toolTipText = c.hoverWidget.Tooltip()
			}
		}
	}

	if toolTipText != "" {
		w, h, _ := canvas.MeasureText(c.Panel().FontFamily(), c.Panel().FontSize(), false, false, toolTipText, true)
		if w > 300 {
			w = 300
		}
		if h > 300 {
			h = 300
		}

		if toolTipX+w > c.Width() {
			toolTipX = c.Width() - w - 12
		}

		if toolTipY+h > c.Height() {
			toolTipY = c.Height() - h - 12
		}

		ctx.SetColor(c.Panel().BackColor())
		ctx.FillRect(toolTipX, toolTipY, w+10, h+6)
		ctx.SetColor(c.Panel().AccentColor())
		ctx.SetStrokeWidth(1)
		ctx.DrawRect(toolTipX, toolTipY, w+10, h+6)
		ctx.SetFontFamily(c.Panel().FontFamily())
		ctx.SetFontSize(c.Panel().FontSize())
		ctx.DrawText(toolTipX+5, toolTipY+3, w, h, toolTipText)
	}
}

func (c *Form) Dispose() {
	c.disposed = true
	if c.userPanel != nil {
		c.userPanel.Dispose()
	}
	c.userPanel = nil

	for _, t := range c.formTimers {
		t.Handler = nil
	}

	c.formTimers = nil
}

func (c *Form) OnClose() bool {
	return true
}

func (c *Form) Panel() *Panel {
	return c.userPanel
}

func (c *Form) AddWidget(w Widget) {
	c.userPanel.addWidget(w)
}

func (c *Form) SetTheme(theme string) {
	c.userPanel.SetTheme(theme)
}

func (c *Form) updatePanelInnerSize() {
	if c.userPanel == nil {
		return
	}

	width, height := c.Window().Size()

	minWidth := c.userPanel.MinWidth()
	minHeight := c.userPanel.MinHeight()
	if minWidth > width || minHeight > height {
		if minWidth < width {
			minWidth = width
		}
		if minHeight < height {
			minHeight = height
		}
		c.userPanel.SetInnerSizeDirect(minWidth, minHeight)
	} else {
		c.userPanel.ResetInnerSizeDirect()
		c.userPanel.ScrollEnsureVisible(0, 0)
	}
}

func (c *Form) ProcessWindowResize(width, height int) {

	//ClearFont()

	if c.userPanel != nil {
		c.userPanel.SetWidth(width)
		c.userPanel.SetHeight(height)
		c.userPanel.SetMaxWidth(width)
		c.userPanel.SetMaxHeight(height)
		c.userPanel.SetMaxHeight(height)
		c.userPanel.SetVerticalScrollVisible(true)
		c.userPanel.SetHorizontalScrollVisible(true)
		c.updatePanelInnerSize()
	}
	c.width = width
	c.height = height

	if c.onSizeChanged != nil {
		var event FormSizeChangedEvent
		event.Width = width
		event.Height = height
		c.onSizeChanged(&event)
	}

	c.UpdateWindow("ProcessWindowResize")

}

func (c *Form) Move(x, y int) {
	c.window.Move(x, y)
}

func (c *Form) Resize(width, height int) {
	c.window.Resize(width, height)
	//if f.id == 0 {
	c.ProcessWindowResize(width, height)
	//}
}

func (c *Form) Width() int {
	return c.width
}

func (c *Form) Height() int {
	return c.height
}

func (c *Form) OnInit() {
}

func (c *Form) SetTitle(title string) {
	c.title = title
	if c.window != nil {
		c.window.SetTitle(c.title)
	}
}

func (c *Form) IsMainWindow() bool {
	return c.isMainWindow
}

func (c *Form) SetIsMainWindow(isMainWindow bool) {
	c.isMainWindow = isMainWindow
}

func (c *Form) Popup() bool {
	return c.popup
}

func (c *Form) SetPopup(popup bool) {
	c.popup = popup
}

func StartMainForm(window Window) {
	window.SetIsMainWindow(true)
	CreateForm(window)
	window.SetTitle(window.Title())
	window.LoopUI()
}

func (c *Form) SetWindow(w nui.Window) {
	c.window = w
}

var MainForm Window

func getWindowByGLFWWindow(w nui.Window) Window {
	if window, ok := windowByGLFWWindow[w]; ok {
		return window
	}
	return nil
}

/*func OnMouseButtonCallback(w nui.Window, button nuimouse.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	window := getWindowByGLFWWindow(w)
	if window == nil {
		return
	}

	switch action {
	case glfw.Release:
		if button == nuimouse.MouseButtonLeft {
			window.ProcessMouseUp(nuimouse.MouseButtonLeft)
		}
		if button == nuimouse.MouseButtonRight {
			window.ProcessMouseUp(nuimouse.MouseButtonRight)
		}
	case glfw.Press:
		if button == nuimouse.MouseButtonLeft {
			window.ProcessMouseDown(nuimouse.MouseButtonLeft)
		}
		if button == nuimouse.MouseButtonRight {
			window.ProcessMouseDown(nuimouse.MouseButtonRight)
		}
	}
}*/

func OnWindowSizeCallback(w nui.Window, width int, height int) {
	window := getWindowByGLFWWindow(w)
	if window == nil {
		return
	}
	window.ProcessWindowResize(width, height)
	//window.Draw()
}

func OnWindowClose(w nui.Window) {
	window := getWindowByGLFWWindow(w)
	if window == nil {
		return
	}
	if window.OnClose() {
		window.Close()
	}
}

func OnWindowCursorPosCallback(w nui.Window, xpos float64, ypos float64) {
	window := getWindowByGLFWWindow(w)
	if window == nil {
		return
	}
	window.ProcessMouseMove(int(xpos), int(ypos))
	//window.Draw()
}

func OnWindowCharCallback(w nui.Window, char rune) {
	window := getWindowByGLFWWindow(w)
	if window == nil {
		return
	}
	window.ProcessCharInput(char)
	//window.Draw()
}

func OnWindowFocusCallback(w nui.Window, focused bool) {
	window := getWindowByGLFWWindow(w)
	if window == nil {
		return
	}
	window.ProcessFocus()
}

func (c *Form) ProcessFocus() {
	if c.childModal != nil {
		//c.childModal.Window().Focus()
	}
}

func OnWindowScrollCallback(w nui.Window, xoff float64, yoff float64) {
	window := getWindowByGLFWWindow(w)
	if window == nil {
		return
	}
	window.ProcessMouseWheel(int(yoff))
}

func CreateForm(form Window) {
	nuiWindow := nui.CreateWindow("UI", 800, 600, true)

	form.SetWindow(nuiWindow)
	form.SetId(nextFormId)
	nextFormId++

	nuiWindow.OnMouseEnter(func() {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}
	})

	nuiWindow.OnMouseLeave(func() {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}
	})

	//window.OnResize(OnWindowSizeCallback)
	//window.OnCloseRequest(OnWindowClose)

	nuiWindow.OnResize(func(width, height int) {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}
		window.ProcessWindowResize(width, height)
	})

	nuiWindow.OnPaint(func(rgba *image.RGBA) {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}

		window.Draw(rgba)
	})

	nuiWindow.OnTimer(func() {
		form.MainTimer()
		//form.UpdateWindow("")
	})

	nuiWindow.OnKeyDown(func(keyCode nuikey.Key, mods nuikey.KeyModifiers) {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}
		window.ProcessKeyModifiers(mods.Shift, mods.Ctrl, mods.Alt)
		window.ProcessKeyDown(keyCode)
	})
	nuiWindow.OnKeyUp(func(keyCode nuikey.Key, mods nuikey.KeyModifiers) {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}
		window.ProcessKeyModifiers(mods.Shift, mods.Ctrl, mods.Alt)
		window.ProcessKeyUp(keyCode)
	})

	nuiWindow.OnChar(func(ch rune) {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}
		window.ProcessCharInput(ch)
	})

	nuiWindow.OnMouseMove(func(xpos int, ypos int) {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}
		window.ProcessMouseMove(int(xpos), int(ypos))
		//window.Draw()
	})

	nuiWindow.OnMouseButtonDown(func(button nuimouse.MouseButton, x int, y int) {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}
		window.ProcessMouseDown(button)
	})

	nuiWindow.OnMouseButtonUp(func(button nuimouse.MouseButton, x int, y int) {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}
		window.ProcessMouseUp(button)
	})

	nuiWindow.OnMouseButtonDblClick(func(button nuimouse.MouseButton, x int, y int) {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}
		window.ProcessMouseDblClick(button)
	})

	nuiWindow.OnMouseWheel(func(x int, y int) {
		window := getWindowByGLFWWindow(nuiWindow)
		if window == nil {
			return
		}
		window.ProcessMouseWheel(y)
	})

	/*window.SetSizeCallback(OnWindowSizeCallback)
	window.SetCloseCallback(OnWindowClose)
	window.SetCursorPosCallback(OnWindowCursorPosCallback)
	window.SetMouseButtonCallback(OnMouseButtonCallback)
	window.SetKeyCallback(OnWindowKeyCallback)
	window.SetCharCallback(OnWindowCharCallback)
	window.SetFocusCallback(OnWindowFocusCallback)
	window.SetScrollCallback(OnWindowScrollCallback)*/

	windows = append(windows, form)
	windowByGLFWWindow[nuiWindow] = form

	form.Init()

	form.Resize(form.Width(), form.Height())

	if len(form.Title()) == 0 {
		form.SetTitle("Form")
	}

	form.OnInit()

	/*monitor := glfw.GetPrimaryMonitor()
	_, _, screenW, screenH := monitor.GetWorkarea()
	wW, wH := window.Size()
	window.SetPos((screenW-wW)/2, (screenH-wH)/2)*/
	nuiWindow.Show()
	form.UpdateLayout()
}

func (c *Form) CreateModalForm(window Window) {
	window.SetParent(c)
	window.SetModal(true)
	c.childModal = window
}

/*func (f *Form) CreatePopupForm(window Window, x int, y int) {
	window.init()
	window.OnInit()
	window.SetParent(f)
	window.SetPopup(true)
	wndId := createPopupWindow(f, window, x, y)
	forms[wndId] = window
}*/

func (c *Form) Maximize() {
	c.window.MaximizeWindow()
}

func StartModalForm(parent Window, window Window) {
	parent.CreateModalForm(window)
	CreateForm(window)
}

func (c *Form) Id() int {
	return c.id
}

func (c *Form) SetId(id int) {
	c.id = id
}

func (c *Form) Show() {
	c.window.Show()
	c.UpdateWindow("Form")
	c.UpdateMenu()
}

func (c *Form) Title() string {
	return c.title
}

func (c *Form) KeyModifiers() nuikey.KeyModifiers {
	return c.keyModifiers
}

func (c *Form) ProcessClick(x, y int, button nuimouse.MouseButton) {

	var event MouseClickEvent
	event.X = x
	event.Y = y
	event.Modifiers = c.keyModifiers
	event.Button = button
	c.userPanel.ProcessMouseClick(&event)
}

func (c *Form) ProcessCharInput(ch rune) {
	if c.focusWidget != nil {

		var event KeyCharEvent
		event.Modifiers = c.keyModifiers
		event.Ch = ch
		c.focusWidget.ProcessKeyChar(&event)
		//f.UpdateWindow("Form")
	}
}

func (c *Form) ProcessKeyDown(key nuikey.Key) {

	if key == nuikey.KeyF12 {
		c.activatedServiceMenu = true
		c.UpdateWindow("Service")
		return
	}

	if c.activatedServiceMenu {
		if key == nuikey.KeyF1 {
			ServiceDrawBorders = !ServiceDrawBorders
		}
		c.activatedServiceMenu = false
		c.UpdateWindow("Service")
		return
	}

	if c.focusWidget != nil {
		processed := false

		var event KeyDownEvent
		event.Modifiers = c.keyModifiers
		event.Key = key

		if c.OnKeyDown != nil {
			processed = c.OnKeyDown(&event)
		}

		if !processed {
			if key == nuikey.KeyEnter || key == nuikey.KeyTab {
				if (key == nuikey.KeyEnter) && !c.focusWidget.AcceptsReturn() {
					if c.ProcessReturnDown() {
						processed = true
					}
				}
				if key == nuikey.KeyTab && !c.focusWidget.AcceptsTab() {
					c.ProcessTabDown()
					processed = true
				}
			}
		}

		if !processed {
			ctrl := c.focusWidget
			for ctrl != nil {
				if ctrl.ProcessKeyDown(&event) {
					break
				}
				ctrl = ctrl.Parent()
			}
		}
		//f.UpdateWindow("Form")
	} else {
		if key == nuikey.KeyEnter {
			c.ProcessReturnDown()
		}
		if key == nuikey.KeyTab {
			c.ProcessTabDown()
		}
	}

	if key == nuikey.KeyEsc {
		if c.rejectButton != nil {
			c.rejectButton.Press()
		}
	}

	if key == nuikey.KeyF11 {
		//application.DumpMemoryMap()
		runtime.GC()
		debug.FreeOSMemory()
	}
}

func (c *Form) ProcessKeyUp(key nuikey.Key) {
	if c.focusWidget != nil {
		processed := false
		if key == nuikey.KeyEnter || key == nuikey.KeyTab {
			if key == nuikey.KeyEnter && !c.focusWidget.AcceptsReturn() {
				c.ProcessReturnUp()
				processed = true
			}
			if key == nuikey.KeyTab && !c.focusWidget.AcceptsTab() {
				c.ProcessTabUp()
				processed = true
			}
		}

		if !processed {
			var event KeyUpEvent
			event.Modifiers = c.keyModifiers
			event.Key = key
			c.focusWidget.ProcessKeyUp(&event)
		}
		//f.UpdateWindow("Form")
	}
}

func (c *Form) ProcessMouseWheel(delta int) {
	x := c.lastMouseMovePos.X
	y := c.lastMouseMovePos.Y
	var event MouseWheelEvent
	event.Modifiers = c.keyModifiers
	event.X = x
	event.Y = y
	event.Delta = delta
	c.userPanel.ProcessMouseWheel(&event)
	c.UpdateWindow("Form - Mouse wheel")
}

func (c *Form) updateHoverWidget(x, y int) {
	wHover := c.userPanel.ProcessFindWidgetUnderPointer(x, y)
	if wHover != c.hoverWidget {
		if c.hoverWidget != nil {
			if c.hoverWidget.Hover() {
				c.hoverWidget.SetHover(false)
				c.hoverWidget.MouseLeave()
			}
		}
	}

	if c.draggingObject != nil && wHover != nil && c.lastMouseDownWidget != nil {

		widgetUpderPoint := c.userPanel.FindWidgetUnderPointer(x, y)
		fX, fY := widgetUpderPoint.RectClientAreaOnWindow()
		ev := NewMouseValidateDropEvent(x-fX, y-fY, 0, c.keyModifiers, c.draggingObject)
		widgetUpderPoint.MouseValidateDrop(ev)
		//c.draggingObject = nil
	}

	if wHover != nil {
		if !wHover.Hover() {
			wHover.SetHover(true)
			wHover.MouseEnter()
		}
	}

	c.hoverWidget = wHover
}

func (c *Form) ProcessMouseMove(x, y int) {
	if math.Abs(float64(c.lastMouseMovePos.X-x)) > 0 || math.Abs(float64(c.lastMouseMovePos.Y-y)) > 0 {
		c.lastMouseMoveTime = time.Now()
		c.toolTipControlProcessed = false
	}
	c.lastMouseMovePos = Point{x, y}

	// HOVER
	c.updateHoverWidget(x, y)

	if c.lastMouseDownWidget != nil {
		// Event to mouseDowned Widget
		fX, fY := c.lastMouseDownWidget.RectOnWindow()
		var event MouseMoveEvent
		event.Modifiers = c.keyModifiers
		event.X = x - fX
		event.Y = y - fY
		c.lastMouseDownWidget.ProcessMouseMove(&event)
	} else {
		// Event to Widget under point
		var event MouseMoveEvent
		event.Modifiers = c.keyModifiers
		event.X = x
		event.Y = y
		c.userPanel.ProcessMouseMove(&event)
	}

	if c.hoverWidget != nil {
		if c.hoverWidget.MouseCursor() != nuimouse.MouseCursorNotDefined {
			c.SetMouseCursor(c.hoverWidget.MouseCursor())
		}
	}

	c.UpdateWindow("Form - mouse move")
}

func (c *Form) ProcessMouseDown(button nuimouse.MouseButton) {

	x := c.lastMouseMovePos.X
	y := c.lastMouseMovePos.Y

	w := c.userPanel.ProcessFindWidgetUnderPointer(x, y)
	c.SetFocusForWidget(w)

	c.mouseDownLastPosX = x
	c.mouseDownLastPosY = y

	event := NewMouseDownEvent(x, y, button, c.keyModifiers)
	c.userPanel.ProcessMouseDown(event)
	if event.UserData("processedWidget") == nil {
		c.lastMouseDownWidget = c.userPanel.ProcessFindWidgetUnderPointer(x, y)
		_ = c.lastMouseDownWidget
	} else {
		c.lastMouseDownWidget = event.UserData("processedWidget").(Widget)
	}

	c.lastMouseMoveTime = time.Now()
	c.toolTipControlProcessed = false
	c.updateHoverWidget(x, y)
}

func (c *Form) ProcessMouseUp(button nuimouse.MouseButton) {
	fmt.Println("ProcessMouseUp")
	x := c.lastMouseMovePos.X
	y := c.lastMouseMovePos.Y

	if math.Abs(float64(c.mouseDownLastPosX-x)) < 5 && math.Abs(float64(c.mouseDownLastPosY-y)) < 5 {
		c.ProcessClick(x, y, button)
	}

	if c.draggingObject != nil && c.lastMouseDownWidget != nil {

		widgetUnderPoint := c.userPanel.ProcessFindWidgetUnderPointer(x, y)
		fX, fY := widgetUnderPoint.RectClientAreaOnWindow()

		ev := NewMouseDropEvent(x-fX, y-fY, button, c.keyModifiers, c.draggingObject)
		widgetUnderPoint.MouseDrop(ev)
	}

	c.draggingObject = nil

	if c.lastMouseDownWidget != nil {
		// Event to mouseDowned Widget
		fX, fY := c.lastMouseDownWidget.RectOnWindow()
		var event MouseUpEvent
		event.Modifiers = c.keyModifiers
		event.X = x - fX
		event.Y = y - fY
		//f.lastMouseDownWidget.processMouseUp(event.Translate(f.lastMouseDownWidget))
		c.lastMouseDownWidget.ProcessMouseUp(&event)
		_ = c.lastMouseDownWidget
	} else {
		// Event to Widget under point
		var event MouseUpEvent
		event.Modifiers = c.keyModifiers
		event.X = x
		event.Y = y
		event.Button = button
		c.userPanel.ProcessMouseUp(&event)
	}

	c.lastMouseDownWidget = nil
}

func (c *Form) ProcessMouseDblClick(button nuimouse.MouseButton) {
	x := c.lastMouseMovePos.X
	y := c.lastMouseMovePos.Y

	w := c.userPanel.ProcessFindWidgetUnderPointer(x, y)
	c.SetFocusForWidget(w)

	c.mouseDownLastPosX = x
	c.mouseDownLastPosY = y

	event := NewMouseDblClickEvent(x, y, button, c.keyModifiers)
	c.userPanel.ProcessMouseDblClick(event)
	if event.UserData("processedWidget") == nil {
		c.lastMouseDownWidget = c.userPanel.ProcessFindWidgetUnderPointer(x, y)
		_ = c.lastMouseDownWidget
	} else {
		c.lastMouseDownWidget = event.UserData("processedWidget").(Widget)
	}

	c.lastMouseMoveTime = time.Now()
	c.toolTipControlProcessed = false
	c.updateHoverWidget(x, y)
}

func (c *Form) findWidgetsUnderTabPlate(parentWidget Widget) []Widget {
	if parentWidget == nil {
		return []Widget{}
	}

	result := make([]Widget, 0)
	for _, w := range parentWidget.Widgets() {
		if w.IsTabPlate() {
			continue
		}
		result = append(result, w)
		result = append(result, c.findWidgetsUnderTabPlate(w)...)
	}
	return result
}

func (c *Form) ProcessTabDown() {
	// find tab plate
	var tabPlateWidget Widget
	currentWidget := c.focusWidget
	if currentWidget != nil && currentWidget.IsTabPlate() {
		tabPlateWidget = currentWidget
	} else {
		for currentWidget != nil {
			if currentWidget.IsTabPlate() {
				tabPlateWidget = currentWidget
				break
			}
			currentWidget = currentWidget.Parent()
		}
	}

	tabPlateWidgets := c.findWidgetsUnderTabPlate(tabPlateWidget)
	currentTabIndex := -1
	if c.focusWidget != nil {
		currentTabIndex = c.focusWidget.TabIndex()
	}

	widgetsWithTabIndex := make([]Widget, 0)

	for _, w := range tabPlateWidgets {
		if w.TabIndex() > 0 {
			widgetsWithTabIndex = append(widgetsWithTabIndex, w)
		}
	}

	sort.Slice(widgetsWithTabIndex, func(i, j int) bool {
		return widgetsWithTabIndex[i].TabIndex() < widgetsWithTabIndex[j].TabIndex()
	})

	nextTabWidget := c.focusWidget

	if c.focusWidget == nil || c.focusWidget.TabIndex() < 1 {
		if len(widgetsWithTabIndex) != 0 {
			nextTabWidget = widgetsWithTabIndex[0]
		}
	} else {
		for index, w := range widgetsWithTabIndex {
			if w.TabIndex() == currentTabIndex {
				if index == len(widgetsWithTabIndex)-1 {
					nextTabWidget = widgetsWithTabIndex[0]
				} else {
					nextTabWidget = widgetsWithTabIndex[index+1]
				}
				break
			}
		}
	}

	if nextTabWidget != nil {
		nextTabWidget.Focus()
	}
}

func (c *Form) FocusedWidget() Widget {
	return c.focusWidget
}

func (c *Form) ProcessWindowMove(x, y int) {
	c.position.X = x
	c.position.Y = y
}

func (c *Form) ProcessReturnDown() bool {
	if c.acceptButton != nil {
		c.acceptButton.Press()
		return true
	}
	return false
}

func (c *Form) ProcessTabUp() {
}

func (c *Form) ProcessReturnUp() {
}

func (c *Form) ProcessKeyModifiers(shift bool, control bool, alt bool) {
	c.keyModifiers = nuikey.KeyModifiers{Shift: shift, Ctrl: control, Alt: alt}
}

/*func (f * Form) onMouseDown(x, y int) {
	f.panel.onMouseDown(x, y, f.keyModifiers)
}

func (f * Form) onMouseUp(x, y int) {
	f.panel.onMouseUp(x, y, f.keyModifiers)
}

func (f * Form) onClick(x, y int) {
}*/

func (c *Form) Close() {
	c.window.Close()

	delete(windowByGLFWWindow, c.window)

	windowIndex := -1
	for i := 0; i < len(windows); i++ {
		if windows[i].Id() == c.Id() {
			windowIndex = i
			break
		}
	}

	if windowIndex > -1 {
		windows = append(windows[:windowIndex], windows[windowIndex+1:]...)
	}
}

func (c *Form) Position() Point {
	return c.position
}

func (c *Form) SetParent(window Window) {
	c.parent = window
}

func (c *Form) Parent() Window {
	return c.parent
}

func (c *Form) Modal() bool {
	return c.modal
}

func (c *Form) SetModal(modal bool) {
	c.modal = modal
}

func (c *Form) Menu() Menu {
	return c.menu
}

func (c *Form) UpdateMenu() {
	c.ProcessWindowResize(c.width, c.height)
}

func (c *Form) MainTimer() {
	nowMSec := time.Now().UnixNano() / 1000000
	for _, timer := range c.formTimers {
		if timer.Enabled {
			if nowMSec-timer.LastElapsedDTMSec > timer.Period {
				if timer.Handler != nil {
					timer.Handler()
				}
				timer.LastElapsedDTMSec = nowMSec
			}
		}
	}

	// Tooltip
	if time.Since(c.lastMouseMoveTime) > 1*time.Second {
		if !c.toolTipControlProcessed {
			c.UpdateWindow("tool tip control")
			c.toolTipControlProcessed = true
		}
	}
}

func (c *Form) BeginDrag(draggingObject interface{}) {
	c.draggingObject = draggingObject
	c.UpdateWindow("Form/BeginDrag")
}

func (c *Form) CurrentDraggingObject() interface{} {
	return c.draggingObject
}

func (c *Form) SetAcceptButton(acceptButton *Button) {
	c.acceptButton = acceptButton
}

func (c *Form) SetRejectButton(rejectButton *Button) {
	c.rejectButton = rejectButton
}

func (c *Form) OpenFileDialog() string {
	return "tempOpenFileDialog.temp"
}

func (c *Form) ShowTooltip(x, y int, text string) {
	c.currentTooltipX = x
	c.currentTooltipY = y
	c.currentTooltipText = text
}

func (c *Form) SetFocusForWidget(w Widget) {
	fmt.Println("SetFocusForWidget", w.Name())
	if w != c.focusWidget {
		if c.focusWidget != nil {
			if c.focusWidget.HasFocus() {
				c.focusWidget.SetFocus(false)
				c.focusWidget.FocusChanged(false)
			}
		}
	}

	if w != nil {
		if !w.HasFocus() {
			w.SetFocus(true)
			w.FocusChanged(true)
		}
	}
	c.focusWidget = w
}

func (c *Form) AppendPopup(w Widget) {
	c.Panel().AppendPopupWidget(w)
	c.UpdateWindow("Form")
}

func (c *Form) CloseTopPopup() {
	c.Panel().CloseTopPopup()
	c.UpdateWindow("Form")
}

func (c *Form) CloseAllPopup() {
	c.Panel().CloseAllPopup()
	c.UpdateWindow("Form")
}

func (c *Form) CloseAfterPopupWidget(w Widget) {
	c.Panel().CloseAfterPopupWidget(w)
	c.UpdateWindow("Form")
}

func (c *Form) ControlRemoved() {
	if c.focusWidget != nil {
		if c.focusWidget.Disposed() {
			c.focusWidget = nil
		}
	}

	if c.hoverWidget != nil {
		if c.hoverWidget.Disposed() {
			c.hoverWidget = nil
		}
	}

	if c.lastMouseDownWidget != nil {
		if c.lastMouseDownWidget.Disposed() {
			c.lastMouseDownWidget = nil
		}
	}
}

func (c *Form) SelectColorDialog(col color.Color, onColorChanged func(color color.Color)) (bool, color.Color) {
	/*var dialog ColorPickerDialog
	dialog.SetColor(col)
	dialog.onColorChanged = onColorChanged //c.colorChangedInDialog
	StartModalForm(c, &dialog)
	return dialog.DialogResult(), dialog.color*/
	return false, colornames.Aliceblue
}

/*func (c *Form) SelectDateTimeDialog(dt time.Time, onDateTimeChanged func(dateTime time.Time)) (bool, time.Time) {
	var dialog DateTimePickerDialog
	dialog.SetDateTime(dt)
	dialog.onDateTimeChanged = onDateTimeChanged
	StartModalForm(c, &dialog)
	return dialog.DialogResult(), dialog.DateTime()
}*/

/*func (c *Form) MessageBox(title, text string) {
	ShowMessageBox(c, title, text)
}*/

func (c *Form) UpdateLayout() {
	c.needUpdateLayout = true
}

func (c *Form) realUpdateLayout() {
	if c.needUpdateLayout {
		c.Panel().SetSize(c.Width(), c.Height())
		c.Panel().ClearLayoutCache()
		c.Panel().UpdateLayout()
		c.needUpdateLayout = false

		x := c.lastMouseMovePos.X
		y := c.lastMouseMovePos.Y
		c.updateHoverWidget(x, y)
	}

	c.updatePanelInnerSize()
}

func (c *Form) timerMemoryDump() {
	//application.DumpMemoryMap()
}

func (c *Form) LoopUI_OpenGL() {
	c.window.EventLoop()

	/*for {
		if len(windows) < 1 {
			break
		}

		glfw.WaitEventsTimeout(0.001)
		glfw.PollEvents()

		for _, window := range windows {
			window.MainTimer()
			window.Draw()
		}
	}*/
}

func (c *Form) SetMouseCursor(cur nuimouse.MouseCursor) {
	c.window.SetMouseCursor(cur)
}

func (c *Form) CentralWidget() Widget {
	return c.Panel()
}
