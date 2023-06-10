package ui

import (
	"image/color"
	"math"
	"reflect"

	"golang.org/x/image/colornames"
)

const MaxUint = ^uint(0)
const MinUint = 0

const MaxInt = int(^uint(0) >> 1)
const MinInt = -MaxInt - 1

type ContainerGridColumnInfo struct {
	minWidth   int
	maxWidth   int
	expandable bool
	width      int
	collapsed  bool
}

type ContainerGridRowInfo struct {
	minHeight  int
	maxHeight  int
	expandable bool
	height     int
	collapsed  bool
}

type Container struct {
	Control
	Controls            []Widget
	PopupWidgets        []Widget
	absolutePositioning bool
	cellPadding         int
	panelPadding        int

	layoutCacheXExpandableValid bool
	layoutCacheYExpandableValid bool
	layoutCacheMinWidthValid    bool
	layoutCacheMinHeightValid   bool

	layoutCacheXExpandable bool
	layoutCacheYExpandable bool
	layoutCacheMinWidth    int
	layoutCacheMinHeight   int
}

func (c *Container) InitControl(parent Widget, w Widget) {
	c.cellPadding = 5
	c.panelPadding = 5
	c.Controls = make([]Widget, 0)
	c.PopupWidgets = make([]Widget, 0)

	c.Control.InitControl(parent, w)
}

func (c *Container) SetAbsolutePositioning(absolutePositioning bool) {
	c.absolutePositioning = absolutePositioning
}

func (c *Container) UpdateStyle() {
	InitDefaultStyle(c)

	c.Control.UpdateStyle()

	for _, w := range c.Widgets() {
		w.UpdateStyle()
	}

	for _, popupWidget := range c.PopupWidgets {
		popupWidget.UpdateStyle()
	}
}

func (c *Container) Dispose() {
	c.RemoveAllWidgets()
	c.Control.Dispose()
}

func (c *Container) ControlType() string {
	return "Panel"
}

func (c *Container) AddWidget(w Widget) {
	c.Controls = append(c.Controls, w)
	w.SetWindow(c.ownWindow)
	w.SetParent(c)
	c.Window().UpdateLayout()
	c.SetStatistics("panel_" + c.Name())
}

func (c *Container) AddWidgetOnGrid(w Widget, x int, y int) {
	w.SetGridPos(x, y)
	c.Controls = append(c.Controls, w)
	w.SetWindow(c.ownWindow)
	w.SetParent(c)
	//p.Window().UpdateLayout()
	c.SetStatistics("panel_" + c.Name())
}

func (c *Container) RemoveWidget(w Widget) {
	foundIndex := -1
	for i := 0; i < len(c.Controls); i++ {
		p1 := reflect.ValueOf(w).Pointer()
		p2 := reflect.ValueOf(c.Controls[i]).Pointer()
		if p1 == p2 {
			foundIndex = i
		}
	}

	if foundIndex > -1 {
		c.Controls = append(c.Controls[:foundIndex], c.Controls[foundIndex+1:]...)
	}
}

func (c *Container) RemoveAllWidgets() {
	if c.Controls != nil {
		for _, w := range c.Controls {
			w.Dispose()
		}
	}

	c.Controls = make([]Widget, 0)
}

func (c *Container) Draw(ctx DrawContext) {
	for _, w := range c.Controls {
		w.DrawControl(ctx)
		if ServiceDrawBorders {
			ctx.SetStrokeWidth(1)
			ctx.SetColor(colornames.Red)
			ctx.DrawRect(0, 0, w.Width(), w.Height())
			//stat := ""
			//stat += fmt.Sprint("MinW: ", w.MinWidth(), "  MinH: ", w.MinHeight(), "\r\n")
			//stat += fmt.Sprint("MaxW: ", w.MaxWidth(), "  MaxH: ", w.MaxHeight(), "\r\n")
			//stat += fmt.Sprint("W: ", w.Width(), " H: ", w.Height(), "\r\n")
			//ctx.DrawText(0, 0, 100, 100, stat)
		}
	}

	if len(c.PopupWidgets) > 0 {
		ctx.SetColor(color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 128,
		})
		//ctx.FillRect(0, 0, c.ClientWidth(), c.ClientHeight())
	}

	for index, popupWidget := range c.PopupWidgets {
		ctx.Save()
		ctx.Translate(popupWidget.X(), popupWidget.Y())
		popupWidget.DrawBackground(ctx)

		ctx.Save()
		ctx.Translate(popupWidget.LeftBorderWidth(), popupWidget.TopBorderWidth())
		ctx.ClipIn(popupWidget.ScrollOffsetX(), popupWidget.ScrollOffsetY(), popupWidget.Width()-popupWidget.LeftBorderWidth()-popupWidget.RightBorderWidth(), popupWidget.Height()-popupWidget.TopBorderWidth()-popupWidget.BottomBorderWidth())
		popupWidget.Draw(ctx)
		ctx.Load()

		popupWidget.DrawBorders(ctx)

		// Shadow parent windows
		if index < len(c.PopupWidgets)-1 {
			ctx.SetColor(color.RGBA{
				R: 0,
				G: 0,
				B: 0,
				A: 128,
			})
			//ctx.FillRect(0, 0, popupWidget.Width(), popupWidget.Height())
		}

		ctx.Load()
	}
}

func (c *Container) Widgets() []Widget {
	return c.Controls
}

func (c *Container) ClearHover() {
	c.hover = false
	for _, w := range c.Controls {
		w.ClearHover()
	}

	for _, popupWidget := range c.PopupWidgets {
		popupWidget.ClearHover()
	}
}

func (c *Container) ClearFocus() {
	c.focus = false
	for _, w := range c.Controls {
		w.ClearFocus()
	}

	for _, popupWidget := range c.PopupWidgets {
		popupWidget.ClearFocus()
	}
}

func (c *Container) FindWidgetUnderPointer(x, y int) Widget {

	for i := len(c.PopupWidgets) - 1; i >= 0; i-- {
		popupWidget := c.PopupWidgets[i]
		if x > popupWidget.X() && x < popupWidget.X()+popupWidget.Width() && y > popupWidget.Y() && y < popupWidget.Y()+popupWidget.Height() && popupWidget.IsVisible() {
			innerW := popupWidget.FindWidgetUnderPointer(x-popupWidget.X(), y-popupWidget.Y())
			if innerW != nil {
				return innerW
			} else {
				return popupWidget
			}
		}
	}

	for _, w := range c.Controls {
		if x > w.X() && x < w.X()+w.Width() && y > w.Y() && y < w.Y()+w.Height() && w.IsVisible() {
			innerW := w.ProcessFindWidgetUnderPointer(x, y)
			if innerW != nil {
				return innerW
			} else {
				return w
			}
		}
	}
	return nil
}
func (c *Container) UpdateLayout() {
	c.updateLayout(c.ClientWidth(), c.ClientHeight())
	for _, popupWidget := range c.PopupWidgets {
		popupWidget.UpdateLayout()
	}
}

func (c *Container) updateLayout(fullWidth int, fullHeight int) {

	if c.widget.ControlType() == "SplitContainer" {
		fullWidth = fullWidth - 1 + 1
	}

	if c.Name() == "#42" {
		fullWidth = fullWidth - 1 + 1
	}

	if c.absolutePositioning {
		deltaWidth := fullWidth - c.Width()
		deltaHeight := fullHeight - c.Height()

		for _, w := range c.Controls {
			a := w.Anchors()
			newX := w.X()
			newY := w.Y()
			newWidth := w.Width()
			newHeight := w.Height()

			left := a&ANCHOR_LEFT != 0
			top := a&ANCHOR_TOP != 0
			right := a&ANCHOR_RIGHT != 0
			bottom := a&ANCHOR_BOTTOM != 0

			if left && right {
				newWidth += deltaWidth
			}
			if !left && right {
				newX += deltaWidth
			}

			if top && bottom {
				newHeight += deltaHeight
			}

			if !top && bottom {
				newY += deltaHeight
			}

			w.SetX(newX)
			w.SetY(newY)
			w.SetWidth(newWidth)
			w.SetHeight(newHeight)
		}
	} else {
		/*fullWidth -= c.widget.LeftBorderWidth()
		fullWidth -= c.widget.RightBorderWidth()
		fullHeight -= c.widget.TopBorderWidth()
		fullHeight -= c.widget.BottomBorderWidth()*/

		columnsInfo, minX, maxX, allCellPaddingX := c.makeColumnsInfo(fullWidth)
		columnsInfo, _, _, _ = c.makeColumnsInfo(fullWidth - (c.panelPadding + allCellPaddingX + c.panelPadding))

		rowsInfo, minY, maxY, allCellPaddingY := c.makeRowsInfo(fullHeight)
		rowsInfo, _, _, _ = c.makeRowsInfo(fullHeight - (c.panelPadding + allCellPaddingY + c.panelPadding))

		xOffset := c.panelPadding //+ c.LeftBorderWidth()
		for x := minX; x <= maxX; x++ {
			if colInfo, ok := columnsInfo[x]; ok {
				yOffset := c.panelPadding // + c.TopBorderWidth()
				for y := minY; y <= maxY; y++ {
					if rowInfo, ok := rowsInfo[y]; ok {
						w := c.getWidgetInGridCell(x, y)

						if w != nil {

							cX := xOffset
							cY := yOffset

							wWidth := colInfo.width
							if wWidth > w.MaxWidth() {
								wWidth = w.MaxWidth()
							}
							wHeight := rowInfo.height
							if wHeight > w.MaxHeight() {
								wHeight = w.MaxHeight()
							}

							cX += (colInfo.width - wWidth) / 2
							cY += (rowInfo.height - wHeight) / 2

							w.SetX(cX)
							w.SetY(cY)
							if w.IsVisible() {
								w.SetWidth(wWidth)
								w.SetHeight(wHeight)
							} else {
								w.SetWidth(0)
								w.SetHeight(0)
							}
						}

						yOffset += rowInfo.height
						if rowInfo.height > 0 && y < maxY {
							yOffset += c.cellPadding
						}
					}
				}

				xOffset += colInfo.width
				if colInfo.width > 0 && x < maxX {
					xOffset += c.cellPadding
				}
			}
		}

		for _, w := range c.Controls {
			if !w.IsVisible() {
				w.SetWidth(0)
				w.SetHeight(0)
			}
		}
	}

	for _, child := range c.Widgets() {
		child.UpdateLayout()
	}
}

func (c *Container) getWidgetInGridCell(x, y int) Widget {
	for _, w := range c.Controls {
		if w.GridX() == x && w.GridY() == y {
			if w.IsVisible() {
				return w
			}
		}
	}
	return nil
}

func (c *Container) SetWidth(width int) {
	c.width.SetOwnValue(width)
	if c.Window() != nil {
		c.Window().UpdateLayout()
	}
}

func (c *Container) SetHeight(height int) {
	c.height.SetOwnValue(height)
	if c.Window() != nil {
		c.Window().UpdateLayout()
	}
}

func (c *Container) SetSize(width int, height int) {
	c.width.SetOwnValue(width)
	c.height.SetOwnValue(height)
	if c.Window() != nil {
		c.Window().UpdateLayout()
	}
}

func (c *Container) MouseDown(event *MouseDownEvent) {
	popupWidgetsBefore := len(c.PopupWidgets)

	for len(c.PopupWidgets) > 0 {
		topWidget := c.PopupWidgets[len(c.PopupWidgets)-1]
		if event.X > topWidget.X() && event.X < topWidget.X()+topWidget.Width() && event.Y > topWidget.Y() && event.Y < topWidget.Y()+topWidget.Height() && topWidget.IsVisible() {
			topWidget.ProcessMouseDown(event)
			return
		} else {
			c.CloseTopPopup()
			//topWidget.ClosePopup()
			//c.PopupWidgets = c.PopupWidgets[:len(c.PopupWidgets)-1]
		}
	}

	if popupWidgetsBefore != len(c.PopupWidgets) {
		return
	}

	//fmt.Println("MouseDown ", c.widget.FullPath())

	for _, w := range c.Widgets() {
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.ProcessMouseDown(event)
		}
	}
}

func (c *Container) MouseUp(event *MouseUpEvent) {
	for i := len(c.PopupWidgets) - 1; i >= 0; i-- {
		popupWidget := c.PopupWidgets[i]
		w := popupWidget
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.ProcessMouseUp(event)
			return
		}
	}

	for _, w := range c.Controls {
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.ProcessMouseUp(event)
		}
	}
}

func (c *Container) MouseClick(event *MouseClickEvent) {
	for i := len(c.PopupWidgets) - 1; i >= 0; i-- {
		popupWidget := c.PopupWidgets[i]
		w := popupWidget
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.ProcessMouseClick(event)
			return
		}
	}

	for _, w := range c.Controls {
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.ProcessMouseClick(event)
			break
		}
	}
}

func (c *Container) MouseMove(event *MouseMoveEvent) {
	for i := len(c.PopupWidgets) - 1; i >= 0; i-- {
		popupWidget := c.PopupWidgets[i]
		w := popupWidget
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.ProcessMouseMove(event)
			return
		}
	}

	for _, w := range c.Controls {
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.ProcessMouseMove(event)
		}
	}
}

func (c *Container) MouseWheel(event *MouseWheelEvent) {
	for i := len(c.PopupWidgets) - 1; i >= 0; i-- {
		popupWidget := c.PopupWidgets[i]
		w := popupWidget
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.ProcessMouseWheel(event)
			return
		}
	}

	for _, w := range c.Controls {
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.ProcessMouseWheel(event)
		}
	}
}

func (c *Container) MouseDrop(event *MouseDropEvent) {
	for i := len(c.PopupWidgets) - 1; i >= 0; i-- {
		popupWidget := c.PopupWidgets[i]
		w := popupWidget
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.MouseDrop(event)
			return
		}
	}

	for _, w := range c.Controls {
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.MouseDrop(event)
		}
	}
}

func (c *Container) MouseValidateDrop(event *MouseValidateDropEvent) {
	for i := len(c.PopupWidgets) - 1; i >= 0; i-- {
		popupWidget := c.PopupWidgets[i]
		w := popupWidget
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.MouseValidateDrop(event)
			return
		}
	}

	for _, w := range c.Controls {
		if event.X > w.X() && event.X < w.X()+w.Width() && event.Y > w.Y() && event.Y < w.Y()+w.Height() && w.IsVisible() {
			w.MouseValidateDrop(event)
		}
	}
}

func (c *Container) ClearRadioButtons() {
	for _, w := range c.Controls {
		w.ClearRadioButtons()
	}
}

func (c *Container) AppendPopupWidget(w Widget) {
	if w != nil {
		c.PopupWidgets = append(c.PopupWidgets, w)
	}
	c.Update("Form")
}

func (c *Container) CloseAfterPopupWidget(w Widget) {
	foundIndex := -1
	for index, popupWidget := range c.PopupWidgets {
		if popupWidget == w {
			foundIndex = index
			break
		}
	}

	if foundIndex > -1 {
		foundIndex++

		for i := foundIndex; i < len(c.PopupWidgets); i++ {
			popupWidget := c.PopupWidgets[i]
			popupWidget.ClosePopup()
		}

		if foundIndex < len(c.PopupWidgets) {
			c.PopupWidgets = append(c.PopupWidgets[:foundIndex], c.PopupWidgets[foundIndex+1:]...)
		}
		c.Update("Form")
	}
}

func (c *Container) CloseAllPopup() {
	for _, popupWidget := range c.PopupWidgets {
		popupWidget.ClosePopup()
	}

	c.PopupWidgets = make([]Widget, 0)
	c.Update("Form")
}

func (c *Container) CloseTopPopup() {
	if len(c.PopupWidgets) == 0 {
		return
	}
	//fmt.Println("Close popup")
	c.PopupWidgets[len(c.PopupWidgets)-1].ClosePopup()
	c.PopupWidgets[len(c.PopupWidgets)-1].Dispose()
	c.PopupWidgets = c.PopupWidgets[:len(c.PopupWidgets)-1]
}

func (c *Container) MinWidth() int {

	if c.layoutCacheMinWidthValid {
		return c.layoutCacheMinWidth
	}

	result := 0
	columnsInfo, _, _, allCellPadding := c.makeColumnsInfo(c.Width())
	columnsInfo, _, _, _ = c.makeColumnsInfo(c.Width() - (c.panelPadding + allCellPadding + c.panelPadding))

	for _, columnInfo := range columnsInfo {
		result += columnInfo.minWidth
	}

	result = result + c.panelPadding + allCellPadding + c.panelPadding + c.LeftBorderWidth() + c.RightBorderWidth()

	c.layoutCacheMinWidthValid = true
	if c.minWidth > result {
		c.layoutCacheMinWidth = c.minWidth
		return c.minWidth
	}
	c.layoutCacheMinWidth = result
	return result
}

func (c *Container) MinHeight() int {

	if c.layoutCacheMinHeightValid {
		return c.layoutCacheMinHeight
	}

	result := 0

	rowsInfo, _, _, allCellPadding := c.makeRowsInfo(c.Height())
	rowsInfo, _, _, _ = c.makeRowsInfo(c.Height() - (c.panelPadding + allCellPadding + c.panelPadding))
	for _, rowInfo := range rowsInfo {
		result += rowInfo.minHeight
	}

	result += c.panelPadding + allCellPadding + c.panelPadding + c.TopBorderWidth() + c.BottomBorderWidth()

	c.layoutCacheMinHeightValid = true
	if c.minHeight > result {
		c.layoutCacheMinHeight = c.minHeight
		return c.minHeight
	}
	c.layoutCacheMinHeight = result
	return result
}

func (c *Container) MaxWidth() int {
	return c.maxWidth
}

func (c *Container) MaxHeight() int {
	return c.maxHeight
}

func (c *Container) SetCellPadding(cellPadding int) {
	c.cellPadding = cellPadding
	c.Update("Panel")
}

func (c *Container) SetPanelPadding(panelPadding int) {
	c.panelPadding = panelPadding
	c.Update("Panel")
}

func (c *Container) makeColumnsInfo(fullWidth int) (map[int]*ContainerGridColumnInfo, int, int, int) {

	minX := MaxInt
	minY := MaxInt

	maxX := MinInt
	maxY := MinInt

	// Detect range of grid coordinates
	for _, w := range c.Controls {
		if w.GridX() < minX {
			minX = w.GridX()
		}
		if w.GridX() > maxX {
			maxX = w.GridX()
		}
		if w.GridY() < minY {
			minY = w.GridY()
		}
		if w.GridY() > maxY {
			maxY = w.GridY()
		}
	}

	columnsInfo := make(map[int]*ContainerGridColumnInfo)
	hasExpandableColumns := false

	// Fill columnsInfo
	for x := minX; x <= maxX; x++ {
		var colInfo ContainerGridColumnInfo
		colInfo.minWidth = MinInt
		colInfo.maxWidth = MaxInt
		colInfo.expandable = false
		found := false

		for y := minY; y <= maxY; y++ {
			w := c.getWidgetInGridCell(x, y)
			if w != nil {
				if w.XExpandable() {
					colInfo.expandable = true // Found expandable by X
					hasExpandableColumns = true
				}
				found = true
			}
		}

		if colInfo.expandable {
			colInfo.minWidth = MinInt
			colInfo.maxWidth = MinInt

			for y := minY; y <= maxY; y++ {
				w := c.getWidgetInGridCell(x, y)
				if w != nil {
					wMinWidth := w.MinWidth()
					if wMinWidth > colInfo.minWidth {
						colInfo.minWidth = wMinWidth
					}
					wMaxWidth := w.MaxWidth()
					if wMaxWidth > colInfo.maxWidth {
						colInfo.maxWidth = wMaxWidth
					}
				}
			}

		} else {
			colInfo.minWidth = MinInt
			colInfo.maxWidth = MinInt

			for y := minY; y <= maxY; y++ {
				w := c.getWidgetInGridCell(x, y)
				if w != nil {
					wMinWidth := w.MinWidth()
					if wMinWidth > colInfo.minWidth {
						colInfo.minWidth = w.MinWidth()
					}
					if wMinWidth > colInfo.maxWidth {
						colInfo.maxWidth = w.MaxWidth()
					}
					/*if w.MaxWidth() < colInfo.maxWidth {
						colInfo.maxWidth = w.MaxWidth()
					}*/
				}
			}
		}

		if found {
			columnsInfo[x] = &colInfo
		}
	}

	if hasExpandableColumns {
		hasNonExpandable := false
		for _, colInfo := range columnsInfo {
			if !colInfo.expandable {
				hasNonExpandable = true
				break
			}
		}
		if hasNonExpandable {
			for _, colInfo := range columnsInfo {
				if !colInfo.expandable {
					colInfo.width = colInfo.minWidth
					colInfo.collapsed = true
				}
			}
		}
	}

	width := fullWidth

	for {
		readyWidth := 0
		for _, colInfo := range columnsInfo {
			readyWidth += colInfo.width
		}
		deltaWidth := width - readyWidth
		countOfColumnCanChange := 0
		for _, colInfo := range columnsInfo {
			if deltaWidth > 0 {
				if colInfo.width < colInfo.maxWidth {
					if !colInfo.collapsed {
						countOfColumnCanChange++
					}
				}
			} else {
				if deltaWidth < 0 {
					if colInfo.width > colInfo.minWidth {
						if !colInfo.collapsed {
							countOfColumnCanChange++
						}
					}
				}
			}
		}

		if countOfColumnCanChange > 0 && deltaWidth != 0 {
			pixForOne := deltaWidth / countOfColumnCanChange
			if math.Abs(float64(pixForOne)) < 1 {
				break
			}
			for _, colInfo := range columnsInfo {
				if !colInfo.collapsed {
					colInfo.width += pixForOne
				}
			}
		} else {
			break
		}

		for _, colInfo := range columnsInfo {
			if colInfo.width > colInfo.maxWidth {
				colInfo.width = colInfo.maxWidth
			}
			if colInfo.width < colInfo.minWidth {
				colInfo.width = colInfo.minWidth
			}
		}
	}

	allCellPadding := 0
	for _, colInfo := range columnsInfo {
		if colInfo.width > 0 {
			allCellPadding++
		}
	}
	allCellPadding--
	allCellPadding *= c.cellPadding
	if allCellPadding < 0 {
		allCellPadding = 0
	}

	return columnsInfo, minX, maxX, allCellPadding

}

func (c *Container) makeRowsInfo(fullHeight int) (map[int]*ContainerGridRowInfo, int, int, int) {

	// Определяем минимальный и максимальный индекс строк
	minX := MaxInt
	minY := MaxInt
	maxX := MinInt
	maxY := MinInt
	for _, w := range c.Controls {
		if w.GridX() < minX {
			minX = w.GridX()
		}
		if w.GridX() > maxX {
			maxX = w.GridX()
		}
		if w.GridY() < minY {
			minY = w.GridY()
		}
		if w.GridY() > maxY {
			maxY = w.GridY()
		}
	}

	// Подготовка
	rowsInfo := make(map[int]*ContainerGridRowInfo)
	hasExpandableRows := false

	// Главный цикл по строкам
	for y := minY; y <= maxY; y++ {
		var rowInfo ContainerGridRowInfo
		rowInfo.minHeight = MinInt // Минимальная высота строки пока 0
		rowInfo.maxHeight = MaxInt // Максимальная высота строки пока ... максимум
		rowInfo.expandable = false // Пока думаем, что строка не мажорная
		found := false             // Признак того, что вообще есть в строке контролы

		for x := minX; x <= maxX; x++ {
			w := c.getWidgetInGridCell(x, y)
			if w != nil {
				if w.YExpandable() {
					rowInfo.expandable = true // Found expandable by Y
					hasExpandableRows = true
				}
				found = true
			}
		}

		if rowInfo.expandable {
			rowInfo.minHeight = MinInt
			rowInfo.maxHeight = MinInt

			for x := minX; x <= maxX; x++ {
				w := c.getWidgetInGridCell(x, y)
				if w != nil {
					wMinHeight := w.MinHeight()
					if wMinHeight > rowInfo.minHeight {
						rowInfo.minHeight = wMinHeight
					}
					wMaxHeight := w.MaxHeight()
					if wMaxHeight > rowInfo.maxHeight {
						rowInfo.maxHeight = wMaxHeight
					}
				}
			}

		} else {
			rowInfo.minHeight = MinInt
			rowInfo.maxHeight = MinInt

			for x := minX; x <= maxX; x++ {
				w := c.getWidgetInGridCell(x, y)
				if w != nil {
					wMinHeight := w.MinHeight()
					if wMinHeight > rowInfo.minHeight {
						rowInfo.minHeight = wMinHeight
					}
					if wMinHeight > rowInfo.maxHeight {
						rowInfo.maxHeight = w.MaxHeight()
					}
					/*if w.MaxWidth() < colInfo.maxWidth {
						colInfo.maxWidth = w.MaxWidth()
					}*/
				}
			}
		}

		if found {
			rowsInfo[y] = &rowInfo
		}
	}

	if hasExpandableRows {
		hasNonExpandable := false
		for _, rowInfo := range rowsInfo {
			if !rowInfo.expandable {
				hasNonExpandable = true
				break
			}
		}
		if hasNonExpandable {
			for _, rowsInfo := range rowsInfo {
				if !rowsInfo.expandable {
					rowsInfo.height = rowsInfo.minHeight
					rowsInfo.collapsed = true
				}
			}
		}
	}

	height := fullHeight

	for {
		readyHeight := 0
		for _, rowInfo := range rowsInfo {
			readyHeight += rowInfo.height
		}
		deltaHeight := height - readyHeight
		countOfRowCanChange := 0
		for _, rowInfo := range rowsInfo {
			if deltaHeight > 0 {
				if rowInfo.height < rowInfo.maxHeight {
					if !rowInfo.collapsed {
						countOfRowCanChange++
					}
				}
			} else {
				if deltaHeight < 0 {
					if rowInfo.height > rowInfo.minHeight {
						if !rowInfo.collapsed {
							countOfRowCanChange++
						}
					}
				}
			}
		}

		if countOfRowCanChange > 0 && deltaHeight != 0 {
			pixForOne := deltaHeight / countOfRowCanChange
			if math.Abs(float64(pixForOne)) < 1 {
				break
			}
			for _, rowInfo := range rowsInfo {
				if !rowInfo.collapsed {
					rowInfo.height += pixForOne
				}
			}
		} else {
			break
		}

		for _, rowInfo := range rowsInfo {
			if rowInfo.height > rowInfo.maxHeight {
				rowInfo.height = rowInfo.maxHeight
			}
			if rowInfo.height < rowInfo.minHeight {
				rowInfo.height = rowInfo.minHeight
			}
		}
	}

	allCellPadding := 0
	for _, rowInfo := range rowsInfo {
		if rowInfo.height > 0 {
			allCellPadding++
		}
	}
	allCellPadding--
	allCellPadding *= c.cellPadding
	if allCellPadding < 0 {
		allCellPadding = 0
	}

	return rowsInfo, minY, maxY, allCellPadding
}

func (c *Container) XExpandable() bool {
	if len(c.Widgets()) == 0 {
		return true // Panel is expandable by default
	}

	if c.layoutCacheXExpandableValid {
		return c.layoutCacheXExpandable
	}

	colsInfo, _, _, _ := c.makeColumnsInfo(1000)
	for _, ci := range colsInfo {
		if ci.expandable {
			c.layoutCacheXExpandableValid = true
			c.layoutCacheXExpandable = true
			return true
		}
	}

	c.layoutCacheXExpandableValid = true
	c.layoutCacheXExpandable = false

	return false
}

func (c *Container) YExpandable() bool {
	if len(c.Widgets()) == 0 {
		return true // Panel is expandable by default
	}

	rowsInfo, _, _, _ := c.makeRowsInfo(1000)
	for _, ri := range rowsInfo {
		if ri.expandable {
			return true
		}
	}

	return false
}

func (c *Container) String(level int) string {
	result := c.Control.String(level)
	result += "\r\n"
	for _, w := range c.Widgets() {
		result += w.String(level + 1)
		result += "\r\n"
	}
	return result
}

func (c *Container) ClearLayoutCache() {
	c.layoutCacheXExpandableValid = false
	c.layoutCacheYExpandableValid = false
	c.layoutCacheMinWidthValid = false
	c.layoutCacheMinHeightValid = false

	for _, w := range c.Widgets() {
		w.ClearLayoutCache()
	}
}
