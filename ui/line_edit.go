package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/atotto/clipboard"
	"github.com/ipoluianov/goforms/utils/canvas"
	"github.com/ipoluianov/goforms/utils/uiproperties"
	"github.com/ipoluianov/nui/nuikey"
	"github.com/ipoluianov/nui/nuimouse"
	"golang.org/x/image/colornames"
)

type LineEdit struct {
	Control
	lines []string

	cursorPosX          int
	cursorPosY          int
	selectionLeftX      int
	selectionRightX     int
	mouseButtonPressed  bool
	cursorWidth         int
	leftAndRightPadding int

	dragingCursor bool
	readonly      bool
	isPassword    bool

	blockUpdate bool
	emptyText   string

	selectionBackground *uiproperties.Property

	timerBlink    *FormTimer
	cursorVisible bool

	OnTextChanged    func(txtBox *LineEdit, oldValue string, newValue string)
	onValidateNeeded func(oldValue string, newValue string) bool
}

/*type lineeditModifyCommand int

const lineeditModifyCommandInsertChar textboxModifyCommand = 0
const lineeditModifyCommandInsertString textboxModifyCommand = 1
const lineeditModifyCommandInsertReturn textboxModifyCommand = 2
const lineeditModifyCommandBackspace textboxModifyCommand = 3
const lineeditModifyCommandDelete textboxModifyCommand = 4
const lineeditModifyCommandSetText textboxModifyCommand = 5*/

type LineEditSelection struct {
	X1, Y1, X2, Y2 int
	Text           string
}

func (c *LineEdit) ControlType() string {
	return "TextBox"
}

func NewLineEdit(parent Widget) *LineEdit {
	var b LineEdit
	b.lines = make([]string, 1)
	b.selectionBackground = AddPropertyToWidget(&b, "selectionBackground", uiproperties.PropertyTypeColor)
	b.InitControl(parent, &b)
	b.SetXExpandable(true)
	b.SetText("")
	b.cursorWidth = 1
	b.leftAndRightPadding = 3
	b.alwaysRedraw = false
	b.cursorVisible = true
	b.timerBlink = b.Window().NewTimer(250, b.timerCursorBlinking)
	b.timerBlink.StartTimer()
	b.ScrollToBegin()
	b.updateInnerSize()

	menu := NewPopupMenu(&b)
	menu.AddItem("Cut", func(event *Event) {
		b.cutSelected()
	}, nil, "")
	menu.AddItem("Copy", func(event *Event) {
		b.copySelected()
	}, nil, "")
	menu.AddItem("Paste", func(event *Event) {
		b.paste()
	}, nil, "")
	menu.AddItem("Delete", func(event *Event) {
		if !b.readonly {
			b.modifyText(textboxModifyCommandDelete, nuikey.KeyModifiers{}, nil)
		}
	}, nil, "")
	menu.AddItem("Select All", func(event *Event) {
		b.SelectAllText()
	}, nil, "")
	b.SetContextMenu(menu)

	return &b
}

func (c *LineEdit) Dispose() {
	c.selectionBackground.Dispose()
	c.selectionBackground = nil

	c.OnTextChanged = nil
	c.onValidateNeeded = nil

	c.Window().RemoveTimer(c.timerBlink)
	c.timerBlink = nil

	c.Control.Dispose()
}

func (c *LineEdit) OnInit() {
}

func (c *LineEdit) SetReadOnly(readonly bool) {
	c.readonly = readonly
}

func (c *LineEdit) SetIsPassword(isPassword bool) {
	c.isPassword = isPassword
}

func (c *LineEdit) timerCursorBlinking() {
	if c.focus {
		c.cursorVisible = !c.cursorVisible
		c.Update("TextBox cursor")
	}
}

func (c *LineEdit) SetText(text string) {
	c.redraw()
	var modifiers nuikey.KeyModifiers
	c.modifyText(textboxModifyCommandSetText, modifiers, text)
	c.updateInnerSize()
	c.ScrollToBegin()
	c.Update("TextBox")
}

func (c *LineEdit) SetEmptyText(text string) {
	c.redraw()
	c.emptyText = text
	c.updateInnerSize()
	c.ScrollToBegin()
	c.Update("TextBox")
}

func (c *LineEdit) Text() string {
	return c.AssemblyText(c.lines)
}

func (c *LineEdit) AssemblyText(lines []string) string {
	result := ""
	for pos, line := range lines {
		result += line
		if pos < len(lines)-1 {
			result += "\r\n"
		}
	}
	return result
}

func (c *LineEdit) updateInnerSize() {
	_, textHeight, err := canvas.MeasureText(c.FontFamily(), c.FontSize(), c.FontBold(), c.FontItalic(), "0", false)
	if err != nil {
		return
	}
	c.InnerHeightOverloaded = textHeight * len(c.lines)
	var maxTextWidth int
	for _, line := range c.lines {
		textWidth, _, err := canvas.MeasureText(c.FontFamily(), c.FontSize(), c.FontBold(), c.FontItalic(), line, false)
		if err != nil {
			return
		}
		if textWidth > maxTextWidth {
			maxTextWidth = textWidth
		}
	}
	c.InnerWidthOverloaded = maxTextWidth + c.leftAndRightPadding*3
	c.InnerSizeOverloaded = true

	c.InnerHeightOverloaded = textHeight
}

func (c *LineEdit) lineToPasswordChars(line string) string {
	if c.isPassword {
		lenOfLine := utf8.RuneCountInString(line)
		line = ""
		for i := 0; i < lenOfLine; i++ {
			line += "*"
		}
	}
	return line
}

func (c *LineEdit) Draw(ctx DrawContext) {

	oneLineHeight := c.OneLineHeight()

	yStaticOffset := c.ClientHeight()/2 - oneLineHeight/2

	// Selection
	if len(c.selectedLines()) > 0 {
		selection := c.selectionRange()
		for selY := selection.Y1; selY <= selection.Y2; selY++ {
			lineCharPos, err := canvas.CharPositions(c.FontFamily(), c.FontSize(), c.FontBold(), c.FontItalic(), c.lines[selY])
			if err != nil {
				return
			}
			for i := 0; i < len(lineCharPos); i++ {
				lineCharPos[i] = lineCharPos[i] + c.leftAndRightPadding
			}

			selXBegin := 0
			selXWidth := lineCharPos[len(lineCharPos)-1]
			if selY == selection.Y1 {
				selXBegin = lineCharPos[selection.X1]
				selXWidth = lineCharPos[len(lineCharPos)-1] - selXBegin
			}
			if selY == selection.Y2 {
				if selection.X2 < len(lineCharPos) {
					selXWidth = lineCharPos[selection.X2] - selXBegin
				}
			}

			//rectY := selY * oneLineHeight

			rectY := yStaticOffset

			ctx.SetColor(c.selectionBackground.Color())
			ctx.FillRect(selXBegin, rectY, selXWidth, oneLineHeight)
		}
	}

	// Text
	yOffset := 0

	for _, line := range c.lines {
		line = c.lineToPasswordChars(line)

		ctx.SetColor(c.foregroundColor.Color())
		ctx.SetFontSize(c.fontSize.Float64())
		ctx.SetTextAlign(canvas.HAlignLeft, canvas.VAlignCenter)
		_, textHeightInLine, err := canvas.MeasureText(c.FontFamily(), c.FontSize(), c.FontBold(), c.FontItalic(), line, false)
		ctx.DrawText(c.leftAndRightPadding, yStaticOffset+yOffset, c.InnerWidth(), textHeightInLine, line)

		if err != nil {
			return
		}
		yOffset += oneLineHeight
	}

	// Cursor
	if c.focus && c.cursorVisible {
		charPos, err := canvas.CharPositions(c.FontFamily(), c.fontSize.Float64(), c.FontBold(), c.FontItalic(), c.lineToPasswordChars(c.lines[c.cursorPosY]))
		for i := 0; i < len(charPos); i++ {
			charPos[i] = charPos[i] + c.leftAndRightPadding
		}
		if err != nil {
			return
		}
		cursorPosInPixels := charPos[c.cursorPosX]
		curX := cursorPosInPixels - (c.cursorWidth / 2)
		curY := yStaticOffset + c.cursorPosY*oneLineHeight
		ctx.SetColor(c.foregroundColor.Color())
		ctx.FillRect(curX, curY, c.cursorWidth, oneLineHeight)
	}

	if ServiceDrawBorders {
		info := ""
		info += fmt.Sprint("InnerWidth: ", c.InnerWidth(), "\r\n")
		info += fmt.Sprint("InnerHeight: ", c.InnerHeight(), "\r\n")
		ctx.SetColor(colornames.Red)
		ctx.DrawText(100, 0, 200, 200, info)
	}

	if c.Text() == "" && c.emptyText != "" && !c.focus {
		ctx.SetColor(c.inactiveColor.Color())
		ctx.SetFontSize(c.fontSize.Float64())
		ctx.SetTextAlign(canvas.HAlignLeft, canvas.VAlignCenter)
		//_, textHeightInLine, _ := canvas.MeasureText(c.FontFamily(), c.FontSize(), c.FontBold(), c.FontItalic(), c.emptyText, false)
		ctx.DrawText(c.leftAndRightPadding, 0, c.Width(), c.Height(), c.emptyText)
	}
}

/*func (c *TextBox) onFocusChanged(focus bool) {
}*/

func (c *LineEdit) KeyChar(event *KeyCharEvent) {

	if c.readonly {
		return
	}

	c.redraw()

	ch := event.Ch
	if ch < 32 {
		return
	}

	c.modifyText(textboxModifyCommandInsertChar, event.Modifiers, ch)
}

func (c *LineEdit) cutSelected() {
	clipboard.WriteAll(c.SelectedText())
	if !c.readonly {
		c.modifyText(textboxModifyCommandDelete, nuikey.KeyModifiers{}, nil)
	}
}

func (c *LineEdit) copySelected() {
	clipboard.WriteAll(c.SelectedText())
}

func (c *LineEdit) paste() {
	str, err := clipboard.ReadAll()
	if err != nil {
		return
	}

	c.modifyText(textboxModifyCommandInsertString, nuikey.KeyModifiers{
		Shift: false,
		Ctrl:  true,
		Alt:   false,
	}, str)
}

func (c *LineEdit) KeyDown(event *KeyDownEvent) bool {
	c.redraw()

	if event.Modifiers.Ctrl && event.Key == nuikey.KeyA {
		c.SelectAllText()
		return true
	}

	if event.Modifiers.Ctrl && event.Key == nuikey.KeyX {
		c.cutSelected()
		return true
	}

	if event.Modifiers.Ctrl && event.Key == nuikey.KeyV {
		c.paste()
		/*str := glfw.GetClipboardString()
		c.modifyText(textboxModifyCommandInsertString, event.Modifiers, str)*/
		return true
	}

	if event.Modifiers.Ctrl && event.Key == nuikey.KeyC {
		c.copySelected()
		return true
	}

	if event.Key == nuikey.KeyArrowLeft {
		c.moveCursor(c.cursorPosX-1, c.cursorPosY, event.Modifiers, true)
		return true
	}

	if event.Key == nuikey.KeyArrowRight {
		c.moveCursor(c.cursorPosX+1, c.cursorPosY, event.Modifiers, true)
		return true
	}

	if event.Key == nuikey.KeyArrowUp {
		c.moveCursor(c.cursorPosX, c.cursorPosY-1, event.Modifiers, true)
		return true
	}

	if event.Key == nuikey.KeyArrowDown {
		c.moveCursor(c.cursorPosX, c.cursorPosY+1, event.Modifiers, true)
		return true
	}

	if event.Key == nuikey.KeyHome {
		c.moveCursor(0, c.cursorPosY, event.Modifiers, true)
		return true
	}

	if event.Key == nuikey.KeyEnter {
		return !c.readonly
	}

	if event.Key == nuikey.KeyEnd {
		runes := []rune(c.lines[c.cursorPosY])
		c.moveCursor(len(runes), c.cursorPosY, event.Modifiers, true)
		return true
	}

	if event.Key == nuikey.KeyBackspace {
		if !c.readonly {
			c.modifyText(textboxModifyCommandBackspace, event.Modifiers, nil)
		}
		return true
	}

	if event.Key == nuikey.KeyDelete {
		if !c.readonly {
			c.modifyText(textboxModifyCommandDelete, event.Modifiers, nil)
		}
		return true
	}

	return false
}

func (c *LineEdit) KeyUp(event *KeyUpEvent) {
}

func (c *LineEdit) MouseDown(event *MouseDownEvent) {
	if event.Button == nuimouse.MouseButtonLeft {
		c.redraw()
		fmt.Println("LineEdit MouseDown")
		c.mouseButtonPressed = true
		c.moveCursorNearPoint(event.X, event.Y, event.Modifiers)
		c.selectionLeftX = c.cursorPosX
		c.selectionRightX = c.cursorPosX
		c.dragingCursor = true
		c.Update("TextBox")
	}
}

func (c *LineEdit) MouseMove(event *MouseMoveEvent) {
	c.redraw()
	if c.mouseButtonPressed {
		c.moveCursorNearPoint(event.X, event.Y, event.Modifiers)
	}
	c.Update("TextBox")
}

func (c *LineEdit) moveCursorNearPoint(x, y int, modifiers nuikey.KeyModifiers) {

	_, textHeight, err := canvas.MeasureText(c.FontFamily(), c.FontSize(), c.FontBold(), c.FontItalic(), "0", false)
	if err != nil {
		return
	}
	lineNumber := y / textHeight

	if lineNumber >= len(c.lines) {
		lineNumber = len(c.lines) - 1
	}

	if lineNumber < 0 {
		lineNumber = 0
	}

	charPos, _ := canvas.CharPositions(c.FontFamily(), c.FontSize(), c.FontBold(), c.FontItalic(), c.lines[lineNumber])
	for i := 0; i < len(charPos); i++ {
		charPos[i] = charPos[i] + c.leftAndRightPadding
	}

	if len(charPos) == 1 {
		c.moveCursor(0, lineNumber, modifiers, true)
		return
	}

	if x < charPos[1]-(charPos[1]-charPos[0])/2 {
		c.moveCursor(0, lineNumber, modifiers, true)
	}

	for pos := 1; pos < len(charPos)-1; pos++ {
		left := charPos[pos] - (charPos[pos]-charPos[pos-1])/2
		right := charPos[pos] + (charPos[pos+1]-charPos[pos])/2
		if x >= left && x < right {
			c.moveCursor(pos, lineNumber, modifiers, true)
			break
		}
	}

	if x > charPos[len(charPos)-1] {
		c.moveCursor(len(charPos)-1, lineNumber, modifiers, true)
	}
}

func (c *LineEdit) MouseUp(event *MouseUpEvent) {
	c.dragingCursor = false
	c.redraw()
	c.mouseButtonPressed = false
	fmt.Println("LineEdit MouseUp")
	c.Update("TextBox")
}

func (c *LineEdit) selectionRange() TextBoxSelection {
	var result TextBoxSelection

	if c.selectionLeftX > c.selectionRightX {
		result.X1 = c.selectionRightX
		result.X2 = c.selectionLeftX
	} else {
		result.X2 = c.selectionRightX
		result.X1 = c.selectionLeftX
	}

	return result
}

func (c *LineEdit) selectedLines() []int {
	var result []int
	result = make([]int, 0)
	selection := c.selectionRange()
	if selection.Y2 != selection.Y1 {
		for i := selection.Y1; i <= selection.Y2; i++ {
			result = append(result, i)
		}
	} else {
		if selection.X1 != selection.X2 {
			result = append(result, selection.Y1)
		}
	}
	return result
}

func (c *LineEdit) moveCursor(posX int, posY int, modifiers nuikey.KeyModifiers, allowSelection bool) {

	if posY < 0 {
		return
	}

	if posY >= len(c.lines) {
		return
	}

	runes := []rune(c.lines[posY])

	if posX < 0 {
		return
	}

	if posX > len(runes) {
		posX = len(runes)
	}

	c.cursorPosX = posX
	c.cursorPosY = posY

	if !modifiers.Shift && !c.mouseButtonPressed {
		c.clearSelection()
	}

	if allowSelection {
		if modifiers.Shift || c.dragingCursor {
			c.selectionRightX = c.cursorPosX
		}
	}

	if !c.blockUpdate {
		c.ensureVisibleCursor()
	}

	c.cursorVisible = true
	if c.timerBlink != nil {
		c.timerBlink.RestartTimer()
	}
	c.Update("LineEdit")
}

func (c *LineEdit) SelectedText() string {
	result := ""

	//lines := make([]string, 0)
	selection := c.selectionRange()

	if selection.Y1 == selection.Y2 {
		runes1 := []rune(c.lines[selection.Y1])
		result += string(runes1[selection.X1:selection.X2])
	} else {
		runes1 := []rune(c.lines[selection.Y1])
		result += string(runes1[selection.X1:])
		result += "\r\n"

		if selection.Y2-selection.Y1 > 1 {
			for row := selection.Y1 + 1; row < selection.Y2; row++ {
				result += c.lines[row]
				result += "\r\n"
			}
		}

		runes2 := []rune(c.lines[selection.Y2])
		result += string(runes2[0:selection.X2])
	}

	return result
}

func (c *LineEdit) removeSelectedText(modifiers nuikey.KeyModifiers) (bool, []string, int, int) {
	lines := make([]string, 0)
	modified := false
	selection := c.selectionRange()
	curPosX := c.cursorPosX
	curPosY := c.cursorPosY
	if len(c.selectedLines()) > 0 {
		lines = append(lines, c.lines[0:selection.Y1]...)
		runes1 := []rune(c.lines[selection.Y1])
		runes2 := []rune(c.lines[selection.Y2])
		lines = append(lines, string(runes1[0:selection.X1])+string(runes2[selection.X2:]))
		lines = append(lines, c.lines[selection.Y2+1:]...)
		modified = true
		curPosX = selection.X1
		curPosY = selection.Y1
	} else {
		lines = append(lines, c.lines...)
	}

	return modified, lines, curPosX, curPosY
}

func (c *LineEdit) ensureVisibleCursor() {
	_, oneLineHeight, _ := canvas.MeasureText(c.FontFamily(), c.FontSize(), c.FontBold(), c.FontItalic(), "Q", false)
	charPos, err := canvas.CharPositions(c.FontFamily(), c.fontSize.Float64(), c.FontBold(), c.FontItalic(), c.lines[c.cursorPosY])
	for i := 0; i < len(charPos); i++ {
		charPos[i] = charPos[i] + c.leftAndRightPadding
	}
	if err != nil {
		return
	}
	cursorPosInPixels := charPos[c.cursorPosX]
	curX := cursorPosInPixels - (c.cursorWidth / 2)
	curY := c.cursorPosY * oneLineHeight
	// ctx.FillRect(curX, curY, c.cursorWidth, oneLineHeight)
	c.ScrollEnsureVisible(curX, curY)
	c.ScrollEnsureVisible(curX+c.cursorWidth, curY+oneLineHeight)
}

func (c *LineEdit) clearSelection() {
	c.selectionLeftX = c.cursorPosX
	c.selectionRightX = c.cursorPosX
}

func (c *LineEdit) modifyText(cmd textboxModifyCommand, modifiers nuikey.KeyModifiers, data interface{}) {
	c.redraw()

	valid := true
	selectedTextRemoved, lines, curPosX, curPosY := c.removeSelectedText(modifiers)
	allowSelection := true

	switch cmd {
	case textboxModifyCommandInsertChar:
		{
			out := []rune(lines[curPosY])
			left := string(out[0:curPosX])
			right := string(out[curPosX:])
			lines[curPosY] = left + string(data.(rune)) + right
			curPosX += 1
			allowSelection = false
		}
	case textboxModifyCommandInsertReturn:
		{
			runes := []rune(lines[curPosY])
			left := string(runes[0:curPosX])
			right := string(runes[curPosX:])
			linesBefore := lines[0:curPosY]
			linesAfter := lines[curPosY:]
			lines = append(linesBefore, right)
			lines = append(lines, linesAfter...)
			lines[curPosY] = left
			curPosX = 0
			curPosY++
		}
	case textboxModifyCommandBackspace:
		{
			runes := []rune(lines[curPosY])
			if !selectedTextRemoved {
				if curPosX > 0 {
					left := string(runes[0 : curPosX-1])
					right := string(runes[curPosX:])
					lines[curPosY] = left + right
					curPosX = curPosX - 1
				} else {
					if curPosY > 0 {
						runes := []rune(lines[curPosY-1])
						newCursorPosX := len(runes)
						linesTemp := make([]string, 0)
						linesTemp = append(linesTemp, lines[0:curPosY]...)
						linesTemp[curPosY-1] += lines[curPosY]
						linesTemp = append(linesTemp, lines[curPosY+1:]...)
						lines = linesTemp
						curPosX = newCursorPosX
						curPosY = curPosY - 1
					}
				}
			}
		}
	case textboxModifyCommandDelete:
		{
			runes := []rune(lines[curPosY])
			if !selectedTextRemoved {
				if curPosX < len(runes) {
					left := string(runes[0:curPosX])
					right := string(runes[curPosX+1:])
					lines[curPosY] = left + right
				} else {
					if curPosY < len(lines)-1 {
						linesTemp := make([]string, 0)
						linesTemp = append(linesTemp, lines[0:curPosY+1]...)
						linesTemp[curPosY] += lines[curPosY+1]
						linesTemp = append(linesTemp, lines[curPosY+2:]...)
						lines = linesTemp
					}
				}
			}
		}
	case textboxModifyCommandSetText:
		{
			lines = strings.Split(strings.Replace(data.(string), "\r", "", -1), "\n")
			//curPosX = 0
			//curPosY = 0
		}
	case textboxModifyCommandInsertString:
		{
			c.blockUpdate = true
			runes := string(data.(string))
			for _, ch := range runes {
				if ch < 32 {
					if ch == 10 {
						//c.insertReturn(modifiers)
					}
				}

				var ev KeyCharEvent
				ev.Modifiers = modifiers
				ev.Ch = ch
				c.KeyChar(&ev)
			}
			lines = c.lines
			curPosX = c.cursorPosX
			curPosY = c.cursorPosY
			c.blockUpdate = false
		}
	}

	if c.onValidateNeeded != nil {
		oldValue := c.Text()
		newValue := c.AssemblyText(lines)
		valid = c.onValidateNeeded(oldValue, newValue)
	}

	if valid {
		c.lines = lines
		c.moveCursor(curPosX, curPosY, modifiers, allowSelection)

		if !c.blockUpdate {
			c.clearSelection()
			c.updateInnerSize()

			if c.OnTextChanged != nil {
				oldValue := c.Text()
				newValue := c.AssemblyText(lines)
				c.OnTextChanged(c, oldValue, newValue)
			}
		}

	}

	//c.Update("TextBox")
}

func (c *LineEdit) SelectAllText() {
	runesLast := []rune(c.lines[len(c.lines)-1])
	c.selectionLeftX = 0
	c.selectionRightX = len(runesLast)
}

func (c *LineEdit) ScrollToBegin() {
	c.ScrollEnsureVisible(0, 0)
	c.ScrollEnsureVisible(0, 1)
}

func (c *LineEdit) OneLineHeight() int {
	_, fontHeight, _ := canvas.MeasureText(c.FontFamily(), c.FontSize(), false, false, "1Qg", false)
	return fontHeight
}

func (c *LineEdit) MinHeight() int {
	return c.OneLineHeight() + 4 + c.TopBorderWidth() + c.BottomBorderWidth()
}

/*func (c *TextBox) MaxHeight() int {
	if c.multiline {
		return MAX_HEIGHT
	}
	return c.OneLineHeight()
}*/

func (c *LineEdit) AcceptsReturn() bool {
	return false
}
