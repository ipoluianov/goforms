package uicontrols

import (
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"github.com/gazercloud/gazerui/uistyles"
	"image"
)

type Panel struct {
	Container
}

func NewPanel(parent uiinterfaces.Widget) *Panel {
	var b Panel
	b.InitControl(parent, &b)
	b.PopupWidgets = make([]uiinterfaces.Widget, 0)
	return &b
}

func NewRootPanel(parentWindow uiinterfaces.Window) *Panel {
	var b Panel
	//b.SetName("RootPanel")
	b.InitControl(nil, &b)
	b.SetName("RootPanel")
	b.PopupWidgets = make([]uiinterfaces.Widget, 0)
	return &b
}

func (c *Panel) InitControl(parent uiinterfaces.Widget, w uiinterfaces.Widget) {
	c.Container.InitControl(parent, w)
}

func (c *Panel) Dispose() {
	c.RemoveAllWidgets()
	c.Control.Dispose()
}

func (c *Panel) ControlType() string {
	return "Panel"
}

func (c *Panel) SetEnabled(enabled bool) {
	c.Container.SetEnabled(enabled)

	for _, w := range c.Controls {
		w.SetEnabled(enabled)
	}
}

func (c *Panel) UpdateStyle() {
	uistyles.InitDefaultStyle(c)
	c.Container.UpdateStyle()
}

func (c *Panel) AddButtonOnGrid(gridX int, gridY int, text string, onPress func(event *uievents.Event)) *Button {
	control := NewButton(c, text, onPress)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddCheckBoxOnGrid(gridX int, gridY int, text string) *CheckBox {
	control := NewCheckBox(c, text)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddColorPickerOnGrid(gridX int, gridY int) *ColorPicker {
	control := NewColorPicker(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddComboBoxOnGrid(gridX int, gridY int) *ComboBox {
	control := NewComboBox(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddDateTimePickerOnGrid(gridX int, gridY int) *DateTimePicker {
	control := NewDateTimePicker(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddImageBoxOnGrid(gridX int, gridY int, img image.Image) *ImageBox {
	control := NewImageBox(c, img)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddListViewOnGrid(gridX int, gridY int) *ListView {
	control := NewListView(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddPanelOnGrid(gridX int, gridY int) *Panel {
	control := NewPanel(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddSplitContainerOnGrid(gridX int, gridY int) *SplitContainer {
	control := NewSplitContainer(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddGroupBoxOnGrid(gridX int, gridY int, title string) *GroupBox {
	control := NewGroupBox(c, title)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddProgressBarOnGrid(gridX int, gridY int) *ProgressBar {
	control := NewProgressBar(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddRadioButtonOnGrid(gridX int, gridY int, text string, onCheckChanged func(checkBox *RadioButton, checked bool)) *RadioButton {
	control := NewRadioButton(c, text)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	control.OnCheckedChanged = onCheckChanged
	c.AddWidget(control)
	return control
}

func (c *Panel) AddSpinBoxOnGrid(gridX int, gridY int) *SpinBox {
	control := NewSpinBox(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddTabControlOnGrid(gridX int, gridY int) *TabControl {
	control := NewTabControl(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddTableOnGrid(gridX int, gridY int) *Table {
	control := NewTable(c, 0, 0, 0, 0)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddTextBoxOnGrid(gridX int, gridY int) *TextBox {
	control := NewTextBox(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddTextBoxExtOnGrid(gridX int, gridY int, text string, onSelect func(textBoxExt *TextBoxExt)) *TextBoxExt {
	control := NewTextBoxExt(c, text, onSelect)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddTextBlockOnGrid(gridX int, gridY int, text string) *TextBlock {
	control := NewTextBlock(c, text)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddTreeViewOnGrid(gridX int, gridY int) *TreeView {
	control := NewTreeView(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddHSpacerOnGrid(gridX int, gridY int) *HSpacer {
	control := NewHSpacer(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}

func (c *Panel) AddVSpacerOnGrid(gridX int, gridY int) *VSpacer {
	control := NewVSpacer(c)
	control.SetGridX(gridX)
	control.SetGridY(gridY)
	c.AddWidget(control)
	return control
}
