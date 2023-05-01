package ui

import (
	"math"
)

type SplitContainer struct {
	Panel
	Panel1         *Panel
	Panel2         *Panel
	position       int
	leftCollapsed  bool
	rightCollapsed bool

	splitResizing bool
}

func NewSplitContainer(parent Widget) *SplitContainer {
	var c SplitContainer
	c.InitControl(parent, &c)
	c.SetAbsolutePositioning(true)

	c.Panel1 = NewPanel(&c)
	c.Panel1.SetSize(c.Width()/2-3, c.Height())
	c.Panel1.SetAnchors(ANCHOR_LEFT | ANCHOR_TOP | ANCHOR_BOTTOM)
	c.Panel1.SetName("SplitPanel1")
	c.AddWidgetOnGrid(c.Panel1, 0, 0)

	c.Panel2 = NewPanel(&c)
	c.Panel2.SetPos(c.Width()/2+3, 0)
	c.Panel2.SetSize(c.Width()/2, c.Height())
	c.Panel2.SetAnchors(ANCHOR_LEFT | ANCHOR_TOP | ANCHOR_BOTTOM | ANCHOR_RIGHT)
	c.Panel2.SetName("SplitPanel2")
	c.AddWidgetOnGrid(c.Panel2, 2, 0)

	c.SetPosition(c.Width() / 2)

	return &c
}

func (c *SplitContainer) Dispose() {
	c.Panel1 = nil
	c.Panel2 = nil
	c.Panel.Dispose()
}

func (c *SplitContainer) XExpandable() bool {
	return c.xExpandable
}

func (c *SplitContainer) YExpandable() bool {
	return c.yExpandable
}

func (c *SplitContainer) ControlType() string {
	return "SplitContainer"
}

func (c *SplitContainer) SetWidth(height int) {
	c.Panel.SetWidth(height)
	c.updateInternalPanels()
}

func (c *SplitContainer) SetHeight(height int) {
	c.Panel.SetHeight(height)
	c.updateInternalPanels()
}

func (c *SplitContainer) updateInternalPanels() {
	c.Panel1.SetHeight(c.ClientHeight())
	c.Panel2.SetHeight(c.ClientHeight())

	if c.leftCollapsed {
		c.Panel1.SetVisible(false)
		c.Panel2.SetVisible(true)
		c.Panel2.SetX(0)
		c.Panel2.SetWidth(c.ClientWidth())
		return
	}

	if c.rightCollapsed {
		c.Panel1.SetVisible(true)
		c.Panel2.SetVisible(false)
		c.Panel1.SetWidth(c.ClientWidth() - 6)
		c.Panel2.SetX(c.position + 3)
		return
	}

	if c.position == 0 {
		c.Panel1.SetVisible(false)
		c.Panel2.SetX(0)
		c.Panel2.SetWidth(c.ClientWidth())
	} else {
		c.Panel1.SetWidth(c.position - 3)
		c.Panel1.SetVisible(true)
		c.Panel2.SetVisible(true)
		c.Panel2.SetX(c.position + 3)
		c.Panel2.SetWidth(c.ClientWidth() - c.position - 6)
	}

}

func (c *SplitContainer) SetLeftCollapsed(collapsed bool) {
	c.leftCollapsed = collapsed
	c.UpdateLayout()
}

func (c *SplitContainer) SetRightCollapsed(collapsed bool) {
	c.rightCollapsed = collapsed
	c.UpdateLayout()
}

func (c *SplitContainer) SetWindow(w Window) {
	c.OwnWindow = w
	c.Panel1.SetWindow(w)
	c.Panel2.SetWindow(w)
}

func (c *SplitContainer) Draw(ctx DrawContext) {
	c.Panel.Draw(ctx)

	if c.position > 0 && !c.leftCollapsed && !c.rightCollapsed {
		ctx.SetColor(c.InactiveColor())
		ctx.SetStrokeWidth(1)

		ctx.FillRect(c.position-1, 0, 2, c.Height()/2-10)
		ctx.DrawRect(c.position-1, c.Height()/2-5, 2, 10)
		ctx.FillRect(c.position-1, c.Height()/2+10, 2, c.Height()/2-10)
	}
}

func (c *SplitContainer) MouseDown(event *MouseDownEvent) {
	c.Update("SplitContainer")
	//fmt.Println("MouseDown split ", c.widget.FullPath())

	diff := math.Abs(float64(event.X - c.position))

	if diff <= 3 {
		c.splitResizing = true
	} else {
		c.Panel.MouseDown(event)
	}
}

func (c *SplitContainer) MouseUp(event *MouseUpEvent) {
	c.splitResizing = false
	//c.Panel.MouseUp(event)
}

func (c *SplitContainer) MouseMove(event *MouseMoveEvent) {
	diff := math.Abs(float64(event.X - c.position))
	if diff <= 3 {
		c.Window().SetMouseCursor(MouseCursorResizeHor)
	} else {
		c.Window().SetMouseCursor(MouseCursorArrow)
	}

	if c.splitResizing {
		width := event.X
		c.SetPosition(width)
	} else {
		c.Panel.MouseMove(event)
	}
}

func (c *SplitContainer) FindWidgetUnderPointer(x, y int) Widget {
	diff := math.Abs(float64(x - c.position))
	if diff <= 3 {
		return c
	} else {
		return c.Panel.FindWidgetUnderPointer(x, y)
	}
}

func (c *SplitContainer) SetPosition(pos int) {
	c.position = pos
	if c.position < 0 {
		c.position = 0
	}

	c.updateInternalPanels()

	c.Update("SplitContainer")
}

func (c *SplitContainer) SetPositionRelative(pos float64) {
	if pos < 0.01 {
		pos = 0.01
	}
	if pos > 0.99 {
		pos = 0.99
	}
	c.SetPosition(int(pos * float64(c.Width())))
}
