package coreforms

/*
type LuaEditor struct {
	uiforms.Form

	txtCode *uicontrols.TextBox

	code string

	btnOK     *uicontrols.Button
	btnCancel *uicontrols.Button
}

func NewLuaEditor() *LuaEditor {
	var c LuaEditor
	return &c
}

func (c *LuaEditor) OnInit() {
	c.Resize(1000, 500)
	c.Panel().SetAbsolutePositioning(true)
	c.SetTitle("Value editor")

	c.txtCode = uicontrols.NewTextBox(c.Panel())
	c.txtCode.SetAnchors(uicontrols.ANCHOR_ALL)
	c.txtCode.SetText(c.code)
	c.txtCode.SetMultiline(true)
	c.AddWidget(c.txtCode)

	c.btnOK = uicontrols.NewButton(c.Panel(), "OK", c.onBtnOK)
	c.btnOK.SetAnchors(uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_BOTTOM)
	c.AddWidget(c.btnOK)
	c.btnCancel = uicontrols.NewButton(c.Panel(), "Cancel", c.RejectButton)
	c.btnCancel.SetAnchors(uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_BOTTOM)
	c.AddWidget(c.btnCancel)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(c.btnCancel)
}

func (c *LuaEditor) SetCode(code string) {
	c.code = code
	if c.txtCode != nil {
		c.txtCode.SetText(c.code)
	}
}

func (c *LuaEditor) Code() string {
	return c.code
}

func (c *LuaEditor) onBtnOK(ev *uievents.Event) {
	c.code = c.txtCode.Text()
	c.Accept()
}
*/
