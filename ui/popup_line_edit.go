package ui

import "github.com/ipoluianov/nui/nuikey"

type PopupLineEdit struct {
	Panel
	le           *LineEdit
	enterPressed bool
	initText     string
	width        int
	height       int
	selectAll    bool
	CloseEvent   func()
}

func NewPopupLineEdit(parent Widget, initText string, selectAll bool, width int, height int) *PopupLineEdit {
	var c PopupLineEdit
	c.initText = initText
	c.width = width
	c.height = height
	c.selectAll = selectAll
	c.SetAbsolutePositioning(true)
	c.Panel.InitControl(parent.Window().CentralWidget(), &c)
	c.SetName("PopupListEditPanel")
	return &c
}

func (c *PopupLineEdit) ControlType() string {
	return "PopupLineEdit"
}

func (c *PopupLineEdit) Dispose() {
}

func (c *PopupLineEdit) ShowPopupLineEdit(x int, y int) {
	c.BeginUpdate()
	c.SetX(x)
	c.SetY(y)
	c.rebuildVisualElements()
	c.Window().AppendPopup(c)
	c.le.Focus()
	if c.selectAll {
		c.le.SelectAllText()
	} else {
		c.le.MoveCursorToEnd()
	}
	c.EndUpdate()
}

func (c *PopupLineEdit) ClosePopup() {
	if c.CloseEvent != nil {
		c.CloseEvent()
	}
}

func (c *PopupLineEdit) OnInit() {
	c.rebuildVisualElements()
}

func (c *PopupLineEdit) needToClose() {
	c.Window().CloseTopPopup()
}

func (c *PopupLineEdit) Text() string {
	if c.le != nil {
		return c.le.Text()
	}
	return ""
}

func (c *PopupLineEdit) rebuildVisualElements() {
	c.SetWidth(c.width)
	c.SetHeight(c.height)

	if c.le != nil {
		c.RemoveAllWidgets()
	}

	c.enterPressed = false

	c.absolutePositioning = true
	c.le = NewLineEdit(c)
	c.le.SetX(0)
	c.le.SetY(0)
	c.le.SetWidth(c.width)
	c.le.SetHeight(c.height)
	c.le.SetText(c.initText)
	c.le.onKeyDown = func(event *KeyDownEvent) bool {
		if event.Key == nuikey.KeyEnter {
			c.enterPressed = true
			c.needToClose()
		}

		if event.Key == nuikey.KeyEsc {
			c.enterPressed = false
			c.needToClose()
		}

		return true
	}

	c.AddWidget(c.le)
}
