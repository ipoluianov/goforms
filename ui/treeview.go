package ui

import (
	"image"
	"image/color"
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ipoluianov/goforms/canvas"
	"github.com/ipoluianov/goforms/uiresources"
	"github.com/nfnt/resize"
	"golang.org/x/image/colornames"
)

type TreeView struct {
	Container

	nodes          []*TreeNode
	displayedNodes []*displayedNode
	selectedNode   *TreeNode
	columns        []*TreeColumn

	header  *TreeViewHeader
	content *TreeViewContent

	cache    *ImageCache
	imgPlus  image.Image
	imgMinus image.Image

	nextNodeId int

	itemHeight     int
	contentPadding int
	headerVisible  bool

	OnExpand       func(treeView *TreeView, node *TreeNode)
	OnSelectedNode func(treeView *TreeView, node *TreeNode)
	OnBeginDrag    func(treeView *TreeView, node *TreeNode) interface{}
	OnDropOnNode   func(treeView *TreeView, node *TreeNode, parameter interface{})
}

type TreeViewHeader struct {
	Control
	pressed             bool
	columnResizing      bool
	columnResizingIndex int
	treeView            *TreeView
}

type TreeViewContent struct {
	Control

	treeView *TreeView

	lastMouseDownNode  *TreeNode
	lastMouseHoverNode *displayedNode
	lastMouseDownPoint Point
}

type TreeColumn struct {
	width int
	text  string
}

type displayedNode struct {
	currentX      int
	currentY      int
	currentWidth  int
	currentHeight int
	node          *TreeNode
}

type TreeNode struct {
	id         int
	text       string
	ParentNode *TreeNode
	nodes      []*TreeNode
	expanded   bool
	values     map[int]string
	Icon       image.Image
	UserData   interface{}
	TextColor  color.Color
	ToolTip    string
}

func (c *TreeNode) Text() string {
	return c.text
}

func NewTreeView(parent Widget) *TreeView {
	var c TreeView
	c.nodes = make([]*TreeNode, 0)
	c.InitControl(parent, &c)
	c.cellPadding = 0
	c.panelPadding = 0
	c.cache = NewImageCache("TreeView")
	c.itemHeight = 20
	c.contentPadding = 3
	c.headerVisible = true

	c.contentPadding = 3

	c.header = NewTreeViewHeader(&c)
	c.header.treeView = &c
	c.header.SetVisible(c.headerVisible)
	c.AddWidgetOnGrid(c.header, 0, 0)

	c.content = NewTreeViewContent(&c)
	c.content.treeView = &c
	c.content.SetXExpandable(true)
	c.content.SetYExpandable(true)
	c.AddWidgetOnGrid(c.content, 0, 1)

	c.imgPlus = canvas.AdjustImageForColor(uiresources.ResImage(uiresources.R_icons_material4_png_content_add_materialicons_48dp_1x_baseline_add_black_48dp_png), c.itemHeight, c.itemHeight, c.ForeColor())
	c.imgMinus = canvas.AdjustImageForColor(uiresources.ResImage(uiresources.R_icons_material4_png_content_remove_materialicons_48dp_1x_baseline_remove_black_48dp_png), c.itemHeight, c.itemHeight, c.ForeColor())

	c.AddColumn("Items", 200)

	return &c
}

func NewTreeViewContent(parent Widget) *TreeViewContent {
	var c TreeViewContent
	c.InitControl(parent, &c)
	c.horizontalScrollVisible.SetOwnValue(true)
	c.verticalScrollVisible.SetOwnValue(true)
	return &c
}

func NewTreeViewHeader(parent Widget) *TreeViewHeader {
	var c TreeViewHeader
	c.InitControl(parent, &c)
	return &c
}

func (c *TreeView) Dispose() {
	c.Container.Dispose()
	c.cache.Clear()
	c.content.treeView = nil
	c.content = nil
	c.header.treeView = nil
	c.header = nil
}

func (c *TreeView) ControlType() string {
	return "TreeView"
}

func (c *TreeView) IsHeaderVisible() bool {
	return c.header.IsVisible()
}

func (c *TreeView) SetHeaderVisible(visible bool) {
	c.header.SetVisible(visible)
	c.Window().UpdateLayout()
}

func (c *TreeViewContent) Draw(ctx DrawContext) {

	yOffset := 0
	c.treeView.displayedNodes = make([]*displayedNode, 0)

	visRect := c.VisibleInnerRect()
	beginIndex := visRect.Y / c.treeView.itemHeight
	if beginIndex >= len(c.treeView.nodes) {
		beginIndex = len(c.treeView.nodes) - 1
	}
	endIndex := (visRect.Y + visRect.Height) / c.treeView.itemHeight
	if endIndex >= len(c.treeView.nodes) {
		endIndex = len(c.treeView.nodes) - 1
	}

	yOffset += beginIndex * c.treeView.itemHeight

	if beginIndex >= 0 && beginIndex <= endIndex {
		for index := beginIndex; index <= endIndex; index++ {
			yOffset += c.drawNode(ctx, c.treeView.nodes[index], yOffset, index)
		}
	}

	c.drawGrid(ctx)
	c.updateInnerSize()

}

func (c *TreeView) updateItemHeight() {
	_, fontHeight, _ := canvas.MeasureText(c.FontFamily(), c.FontSize(), false, false, "1Qg", false)
	fontHeight += 2
	if c.itemHeight != fontHeight {
		c.itemHeight = fontHeight
		c.UpdateLayout()
	}
}

func (c *TreeViewHeader) Draw(ctx DrawContext) {
	c.treeView.updateItemHeight()

	xOffset := 0
	yOffset := c.treeView.contentPadding

	ctx.SetColor(c.ForeColor())
	ctx.SetStrokeWidth(1)
	ctx.SetFontSize(c.FontSize())

	for colIndex, column := range c.treeView.columns {
		var cnv *canvas.CanvasDirect
		cnv = c.treeView.cache.GetXY(colIndex, -1)
		if cnv == nil {
			cnv = canvas.NewCanvas(column.width, c.treeView.itemHeight)
			cnv.DrawText(0, 0, column.text, c.FontFamily(), c.FontSize(), c.ForeColor(), false)
			c.treeView.cache.SetXY(colIndex, -1, cnv)
		}

		clipX := xOffset + c.treeView.contentPadding
		clipW := column.width - c.treeView.contentPadding*2

		effectiveWidth := c.Width()

		if clipX+clipW > effectiveWidth {
			clipW = effectiveWidth - clipX
		}
		if clipW < 0 {
			clipW = 0
		}
		if clipX < 0 {
			clipW += clipX
			clipX = 0
		}

		ctx.DrawImage(xOffset+c.treeView.contentPadding, yOffset, column.width, c.treeView.itemHeight, cnv.Image())
		ctx.SetColor(c.treeView.leftBorderColor.Color())
		ctx.DrawLine(xOffset+column.width, 0, xOffset+column.width, c.treeView.itemHeight+c.treeView.contentPadding*2)
		xOffset += column.width
	}

}

func (c *TreeViewContent) MouseDrop(event *MouseDropEvent) {
	node := c.findDisplayNodeByCoordinates(event.X, event.Y)
	if node != nil {
		if c.treeView.OnDropOnNode != nil {
			c.treeView.OnDropOnNode(c.treeView, node.node, event.DroppingObject)
		}
	}
}

func (c *TreeViewContent) drawNode(ctx DrawContext, node *TreeNode, x int, y int) int {
	yOffset := 0
	rowHeight := c.treeView.itemHeight

	var dNode displayedNode
	dNode.currentX = x
	dNode.currentY = y
	dNode.currentWidth = 100
	dNode.currentHeight = c.treeView.itemHeight
	dNode.node = node
	c.treeView.displayedNodes = append(c.treeView.displayedNodes, &dNode)

	if len(node.nodes) > 0 {
		if node.expanded {
			c.drawMinus(ctx, x, y)
		} else {
			c.drawPlus(ctx, x, y)
		}
	}

	c.drawIcon(ctx, x, y, node.Icon)

	xOffset := 0
	for index, column := range c.treeView.columns {
		textXOffset := 0

		value := node.values[index]
		if index == 0 {
			textXOffset = c.treeView.itemHeight + x
			if node.Icon != nil {
				textXOffset += c.treeView.itemHeight
			}
			value = node.text
		}

		cellX := xOffset + textXOffset
		cellY := y
		ctx.SetColor(c.BackColor())

		if c.lastMouseHoverNode != nil {
			if node == c.lastMouseHoverNode.node {
				ctx.SetColor(colornames.Yellow)
			}
		}

		//ctx.FillRect(xOffset+textXOffset, y, column.width, 20)
		ctx.SetColor(node.TextColor)
		ctx.SetFontSize(c.FontSize())

		var cnv *canvas.CanvasDirect
		cnv = c.treeView.cache.GetXY(index, node.id)
		if cnv == nil {
			cnv = canvas.NewCanvas(column.width, rowHeight)

			foreColor := c.foregroundColor.Color()
			backColor := c.BackColor()
			if c.treeView.selectedNode == node {
				backColorTemp := backColor
				backColor = foreColor
				foreColor = backColorTemp
			}

			cnv.FillRect(2, 2, column.width, c.treeView.itemHeight-4, backColor)
			cnv.DrawText(3, 0, value, c.fontFamily.String(), c.fontSize.Float64(), foreColor, false)
		}

		ctx.DrawImage(cellX, cellY, c.Width(), c.Height(), cnv.Image())

		c.treeView.cache.SetXY(index, node.id, cnv)

		xOffset += column.width
	}

	yOffset += c.treeView.itemHeight

	if node.expanded {
		for _, node := range node.nodes {
			yOffset += c.drawNode(ctx, node, x+c.treeView.itemHeight, y+yOffset)
		}
	}

	return yOffset
}

func (c *TreeViewContent) drawGrid(ctx DrawContext) {
	xOffset := 0

	ctx.SetColor(c.foregroundColor.Color())
	ctx.SetStrokeWidth(1)

	for _, column := range c.treeView.columns {
		xOffset += column.width
		ctx.DrawLine(xOffset, 0, xOffset, c.Height()+c.scrollOffsetY)
	}
}

func (c *TreeView) drawHeader(ctx *canvas.CanvasDirect) int {
	xOffset := 0
	yOffset := 0

	for _, column := range c.columns {
		//ctx.DrawRect(xOffset, yOffset, column.width, 20, c.foregroundColor.Color(), 1)
		ctx.DrawLine(xOffset, 0, xOffset, c.itemHeight, 1, c.foregroundColor.Color())
		ctx.DrawText(xOffset+5, yOffset, column.text, c.fontFamily.String(), c.fontSize.Float64(), c.foregroundColor.Color(), false)
		xOffset += column.width
	}
	ctx.DrawLine(0, c.itemHeight-1, c.InnerWidth(), 19, 1, c.foregroundColor.Color())

	return yOffset + c.itemHeight
}

func (c *TreeView) TabStop() bool {
	return true
}

func (c *TreeView) AddNode(parentNode *TreeNode, text string) *TreeNode {
	var node TreeNode
	node.text = text
	node.ParentNode = parentNode
	node.values = make(map[int]string)
	node.TextColor = c.foregroundColor.Color()
	node.id = c.nextNodeId
	c.nextNodeId++

	if parentNode == nil {
		c.nodes = append(c.nodes, &node)
	} else {
		parentNode.nodes = append(parentNode.nodes, &node)
	}

	c.content.updateInnerSize()
	c.Update("TreeView")
	return &node
}

func (c *TreeView) RemoveNode(node *TreeNode) {
	node.nodes = make([]*TreeNode, 0)

	for index, n := range c.nodes {
		if n == node {
			c.nodes = append(c.nodes[:index], c.nodes[index+1:]...)
			return
		}
	}

	c.cache.Clear()
	c.Update("TreeView")
}

func (c *TreeView) RemoveAllNodes() {
	c.selectedNode = nil
	c.nodes = make([]*TreeNode, 0)
	c.cache.Clear()
	c.content.updateInnerSize()
	c.content.ScrollEnsureVisible(0, 0)
	c.Update("TreeView")
}

func (c *TreeView) EnsureVisibleDisplayedNode(node *displayedNode) {

	if node != nil {
		c.content.ScrollEnsureVisible(node.currentX+node.currentWidth, node.currentY+node.currentHeight)
		c.content.ScrollEnsureVisible(node.currentX, node.currentY)
	}
	c.Update("TreeView")
}

func (c *TreeView) EnsureVisibleNode(node *TreeNode) {
	for _, dn := range c.displayedNodes {
		if dn.node == node {
			c.EnsureVisibleDisplayedNode(dn)
			break
		}
	}
}

func (c *TreeView) RemoveNodes(node *TreeNode) {
	node.nodes = make([]*TreeNode, 0)
	c.cache.Clear()
	c.Update("TreeView")
}

func (c *TreeView) AddColumn(text string, width int) *TreeColumn {
	var treeColumn TreeColumn
	treeColumn.text = text
	treeColumn.width = width

	c.columns = append(c.columns, &treeColumn)
	c.content.updateInnerSize()
	c.Update("TreeView")

	return &treeColumn
}

func (c *TreeView) SetColumnWidth(colIndex int, width int) {
	if colIndex >= 0 && colIndex < len(c.columns) {
		c.columns[colIndex].width = width
	}
	c.content.updateInnerSize()
	c.Update("TreeView")
}

func (c *TreeView) SetNodeValue(node *TreeNode, columnIndex int, text string) {
	if columnIndex == 0 {
		node.text = text
	}
	node.values[columnIndex] = text
	c.cache.ClearXY(columnIndex, node.id)
	c.Update("TreeView")
}

func (c *TreeView) GetNodeIndexInParent(node *TreeNode) int {
	index := -1
	for i, n := range node.ParentNode.nodes {
		if n == node {
			index = i
			break
		}
	}
	return index
}

func (c *TreeViewContent) updateInnerSize() {
	c.innerSizeOverloaded = true
	_, c.innerHeightOverloaded = c.calcTreeColumnSize()
	c.updateInnerWidth()
}

func (c *TreeViewContent) updateInnerWidth() {
	c.innerWidthOverloaded = 0
	for _, column := range c.treeView.columns {
		c.innerWidthOverloaded += column.width
	}
}

func (c *TreeView) ExpandNode(node *TreeNode) {
	cNode := node
	for cNode != nil {
		cNode.expanded = true
		if c.OnExpand != nil {
			c.OnExpand(c, node)
		}
		cNode = cNode.ParentNode
	}
	c.Update("TreeView")
}

func (c *TreeView) SelectNode(node *TreeNode) {
	if c.selectedNode != nil {
		c.removeCacheForRow(c.selectedNode.id)
	}
	c.selectedNode = node
	c.removeCacheForRow(node.id)
	c.Update("TreeView")
	c.EnsureVisibleNode(node)

	if c.OnSelectedNode != nil {
		c.OnSelectedNode(c, node)
	}

}

func (c *TreeView) removeCacheForRow(row int) {
	for index := range c.columns {
		c.cache.ClearXY(index, row)
	}
}

func (c *TreeView) CollapseNode(node *TreeNode) {
	node.expanded = false
	c.SelectNode(node)
	c.Update("TreeView")
}

func (c *TreeViewContent) drawPlus(ctx DrawContext, x int, y int) {
	ctx.DrawImage(x, y, 5, 5, c.treeView.imgPlus)
}

func (c *TreeViewContent) drawMinus(ctx DrawContext, x int, y int) {
	ctx.DrawImage(x, y, 5, 5, c.treeView.imgMinus)
}

func (c *TreeViewContent) drawIcon(ctx DrawContext, x int, y int, img image.Image) {
	if img != nil {
		img = resize.Resize(uint(c.treeView.itemHeight), uint(c.treeView.itemHeight), img, resize.Bicubic)
		ctx.DrawImage(x+c.treeView.itemHeight, y, c.treeView.itemHeight, c.treeView.itemHeight, img)
	}
}

func (c *TreeViewContent) findDisplayNodeByCoordinates(x int, y int) *displayedNode {
	for _, dNode := range c.treeView.displayedNodes {
		if x >= 0 && x < c.InnerWidth() && y >= dNode.currentY && y < dNode.currentY+dNode.currentHeight {
			return dNode
		}
	}
	return nil
}

func (c *TreeViewContent) calcTreeColumnSize() (int, int) {
	width := 0
	height := 0
	for _, dNode := range c.treeView.displayedNodes {
		maxX := dNode.currentX + dNode.currentWidth
		maxY := dNode.currentY + dNode.currentHeight
		if maxX > width {
			width = maxX
		}
		if maxY > height {
			height = maxY
		}
	}
	return width, height
}

func (c *TreeViewContent) MouseClick(event *MouseClickEvent) {
	dNode := c.findDisplayNodeByCoordinates(event.X, event.Y)
	if dNode == nil {
		return
	}

	if event.X > dNode.currentX && event.X < dNode.currentX+c.treeView.itemHeight {
		if dNode.node.expanded {
			c.treeView.CollapseNode(dNode.node)
		} else {
			c.treeView.ExpandNode(dNode.node)
		}
	}

	c.Update("TreeView")
}

func (c *TreeViewHeader) MouseDown(event *MouseDownEvent) {
	c.pressed = true
	c.Update("TreeView")

	if event.Y < c.treeView.itemHeight {
		colRightBorder := c.treeView.findColumnRightBorder(event.X, event.Y)
		if colRightBorder >= 0 {
			c.columnResizing = true
			c.columnResizingIndex = colRightBorder
		}
	}

}

func (c *TreeViewContent) MouseDown(event *MouseDownEvent) {
	dNode := c.findDisplayNodeByCoordinates(event.X, event.Y)
	if dNode != nil {
		c.lastMouseDownNode = dNode.node
		c.lastMouseDownPoint = Point{event.X, event.Y}

		if event.X > dNode.currentX+c.treeView.itemHeight {
			c.treeView.SelectNode(dNode.node)
		}
	}
}

func (c *TreeViewHeader) MouseUp(event *MouseUpEvent) {
	if c.pressed {
		c.pressed = false
		c.columnResizing = false
		c.Update("TreeView")
	}
}

func (c *TreeViewContent) MouseUp(event *MouseUpEvent) {
	c.lastMouseDownNode = nil
}

func (c *TreeViewHeader) MouseMove(event *MouseMoveEvent) {
	if c.columnResizing {
		c.treeView.columns[c.columnResizingIndex].width = event.X - c.treeView.calcColumnXOffset(c.columnResizingIndex)
		if c.treeView.columns[c.columnResizingIndex].width < 10 {
			c.treeView.columns[c.columnResizingIndex].width = 10
		}

		c.treeView.cache.Clear()
	}

	if event.Y < c.treeView.itemHeight {
		if c.treeView.findColumnRightBorder(event.X, event.Y) >= 0 {
			c.Window().SetMouseCursor(MouseCursorResizeHor)
		} else {
			c.Window().SetMouseCursor(MouseCursorArrow)
		}
	}
}

func (c *TreeViewContent) MouseMove(event *MouseMoveEvent) {
	if c.lastMouseDownNode != nil {
		if math.Abs(float64(c.lastMouseDownPoint.X-event.X)) > 5 || math.Abs(float64(c.lastMouseDownPoint.Y-event.Y)) > 5 {
			var dragParameter interface{}
			dragParameter = c.lastMouseDownNode
			if c.treeView.OnBeginDrag != nil {
				dragParameter = c.treeView.OnBeginDrag(c.treeView, c.lastMouseDownNode)
			}
			c.BeginDrag(dragParameter)
		}
	}
	c.lastMouseHoverNode = c.findDisplayNodeByCoordinates(event.X, event.Y)
}

func (c *TreeViewContent) Tooltip() string {
	if c.lastMouseHoverNode != nil {
		return c.lastMouseHoverNode.node.ToolTip
	}
	return ""
}

func (c *TreeView) calcColumnXOffset(columnIndex int) int {
	columnOffset := 0
	for index, column := range c.columns {
		if index == columnIndex {
			break
		}
		columnOffset += column.width
	}
	return columnOffset
}

func (c *TreeView) findColumnRightBorder(x, y int) int {
	if y > c.itemHeight {
		return -1
	}
	for index, column := range c.columns {
		colRightBorder := c.calcColumnXOffset(index) + column.width
		if math.Abs(float64(x-colRightBorder)) < 5 {
			return index
		}
	}
	return -1
}

func (c *TreeView) KeyDown(event *KeyDownEvent) bool {

	selectedNodeDisplayIndex := -1

	for index, dNode := range c.displayedNodes {
		if c.selectedNode == dNode.node {
			selectedNodeDisplayIndex = index
			break
		}
	}

	if event.Key == glfw.KeyUp {
		if selectedNodeDisplayIndex > 0 {
			c.SelectNode(c.displayedNodes[selectedNodeDisplayIndex-1].node)
		}
		return true
	}

	if event.Key == glfw.KeyDown {
		if selectedNodeDisplayIndex < len(c.displayedNodes)-1 {
			c.SelectNode(c.displayedNodes[selectedNodeDisplayIndex+1].node)
		}
		return true
	}

	if event.Key == glfw.KeyLeft {
		if c.selectedNode != nil {
			if c.selectedNode.expanded {
				c.CollapseNode(c.selectedNode)
			} else {
				if c.selectedNode.ParentNode != nil {
					c.SelectNode(c.selectedNode.ParentNode)
				}
			}
		}
		return true
	}

	if event.Key == glfw.KeyRight {
		if c.selectedNode != nil {
			c.ExpandNode(c.selectedNode)
		}
		return true
	}

	if event.Key == glfw.KeyHome {
		if len(c.displayedNodes) > 0 {
			c.SelectNode(c.displayedNodes[0].node)
			c.EnsureVisibleDisplayedNode(c.displayedNodes[0])
		}
		return true
	}

	if event.Key == glfw.KeyEnd {
		if len(c.displayedNodes) > 0 {
			c.SelectNode(c.displayedNodes[len(c.displayedNodes)-1].node)
			c.EnsureVisibleDisplayedNode(c.displayedNodes[len(c.displayedNodes)-1])
		}
		return true
	}
	return false
}

func (c *TreeView) SelectedNode() *TreeNode {
	return c.selectedNode
}

func (c *TreeView) VisibleNodes() []*TreeNode {
	return c.visibleNodes(nil)
}

func (c *TreeView) Nodes() []*TreeNode {
	return c.nodes
}

func (c *TreeView) Children(node *TreeNode) []*TreeNode {
	if node == nil {
		return c.nodes
	}
	return node.nodes
}

func (c *TreeView) visibleNodes(node *TreeNode) []*TreeNode {
	result := make([]*TreeNode, 0)

	nodes := c.nodes

	if node != nil {
		nodes = node.nodes
	}

	for _, n := range nodes {
		result = append(result, n)
		if n.expanded {
			result = append(result, c.visibleNodes(n)...)
		}
	}
	return result
}

func (c *TreeViewContent) OnScroll(scrollPositionX int, scrollPositionY int) {
	c.treeView.header.scrollOffsetX = scrollPositionX
}

func (c *TreeViewHeader) MinWidth() int {
	return c.Control.MinWidth()
}

func (c *TreeViewHeader) MinHeight() int {
	return c.treeView.itemHeight + c.treeView.contentPadding*2
}

func (c *TreeViewHeader) ControlType() string {
	return "TreeViewHeader"
}

func (c *TreeViewHeader) MouseLeave() {
	c.Window().SetMouseCursor(MouseCursorArrow)
}
