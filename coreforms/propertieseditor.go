package coreforms

/*
import (
	"fmt"
	"image/color"

	"github.com/ipoluianov/goforms/canvas"
	"github.com/ipoluianov/goforms/uicontrols"
	"github.com/ipoluianov/goforms/uievents"
	"github.com/ipoluianov/goforms/uiinterfaces"
	"github.com/ipoluianov/goforms/uiproperties"
)

type PropertiesEditor struct {
	uicontrols.Panel
	propertiesContainer uiproperties.IPropertiesContainer
	propControls        map[string]uiinterfaces.Widget
	propsMap            map[string]*uiproperties.Property

	loading bool
}

func NewPropertiesEditor(parent uiinterfaces.Widget) *PropertiesEditor {
	var c PropertiesEditor
	c.InitControl(parent, &c)
	return &c
}

func (c *PropertiesEditor) ControlType() string {
	return "PropertiesEditor"
}

func (c *PropertiesEditor) Dispose() {
	if c.propertiesContainer != nil {
		c.propertiesContainer.SetPropertyChangeNotifier(nil)
	}
	c.propControls = nil
	c.propsMap = nil
	c.propertiesContainer = nil
	c.Panel.Dispose()
}

func (c *PropertiesEditor) SetPropertiesContainer(propertiesContainer uiproperties.IPropertiesContainer) {

	c.propControls = make(map[string]uiinterfaces.Widget)
	c.propsMap = make(map[string]*uiproperties.Property)

	if c.propertiesContainer != nil {
		c.propertiesContainer.SetPropertyChangeNotifier(nil)
	}

	c.propertiesContainer = propertiesContainer
	if c.propertiesContainer != nil {
		c.propertiesContainer.SetPropertyChangeNotifier(c.OnPropertyChanged)
	}

	c.RebuildInterface()
}

func (c *PropertiesEditor) RebuildInterface() {

	c.BeginUpdate()
	c.loading = true
	c.RemoveAllWidgets()
	if c.propertiesContainer == nil {
		c.loading = false
		return
	}

	index := 0
	if true {
		for _, property := range c.propertiesContainer.GetProperties() {
			if !property.Visible() {
				continue
			}

			c.propsMap[property.Name] = property

			lblName := c.AddTextBlockOnGrid(0, index, property.Name)
			lblName.SetName("Prop " + property.Name)
			lblName.TextHAlign = canvas.HAlignLeft

			if property.Type == uiproperties.PropertyTypeBool {
				numEditor := c.AddCheckBoxOnGrid(1, index, "---")
				numEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				numEditor.OnCheckedChanged = c.CheckBoxChanged
				numEditor.SetUserData("propName", property.Name)
				numEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = numEditor
			}
			if property.Type == uiproperties.PropertyTypeInt {
				numEditor := c.AddSpinBoxOnGrid(1, index)
				numEditor.SetPrecision(0)
				numEditor.SetIncrement(1)
				numEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				numEditor.OnValueChanged = c.SpinBoxChanged
				numEditor.SetUserData("propName", property.Name)
				numEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = numEditor
			}
			if property.Type == uiproperties.PropertyTypeInt32 {
				numEditor := c.AddSpinBoxOnGrid(1, index)
				numEditor.SetPrecision(0)
				numEditor.SetIncrement(1)
				numEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				numEditor.OnValueChanged = c.SpinBoxChanged
				numEditor.SetUserData("propName", property.Name)
				numEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = numEditor
			}
			if property.Type == uiproperties.PropertyTypeString {
				txtEditor := c.AddTextBoxOnGrid(1, index)
				txtEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				txtEditor.OnTextChanged = c.TextBoxChanged
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = txtEditor
			}
			if property.Type == uiproperties.PropertyTypeMultiline {
				txtEditor := c.AddButtonOnGrid(1, index, "Edit ...", nil)
				txtEditor.SetOnPress(func(ev *uievents.Event) {
					dialog := NewMultilineEditor(c, txtEditor.TempData)
					dialog.OnAccept = func() {
						c.MultilineChanged(txtEditor, dialog.ResultText())
					}
					dialog.ShowDialog()
				})
				txtEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				c.propControls[property.Name] = txtEditor
			}
			if property.Type == uiproperties.PropertyTypeColor {
				txtEditor := uicontrols.NewColorPicker(c)
				txtEditor.SetPos(0, 0)
				txtEditor.SetSize(0, 20)
				txtEditor.SetGridX(1)
				txtEditor.SetGridY(index)
				txtEditor.OwnWindow = c.OwnWindow
				txtEditor.SetAnchors(uicontrols.ANCHOR_LEFT | uicontrols.ANCHOR_RIGHT | uicontrols.ANCHOR_TOP)
				txtEditor.OnColorChanged = c.ColorPickerChanged
				txtEditor.SetUserData("propName", property.Name)
				txtEditor.SetUserData("propType", property.Type)
				c.AddWidget(txtEditor)
				c.propControls[property.Name] = txtEditor
			}
			index++
		}
	}
	c.AddVSpacerOnGrid(0, index)
	c.EndUpdate()

	c.LoadPropertiesValues()
	c.loading = false
}

func (c *PropertiesEditor) LoadPropertiesValues() {
	c.BeginUpdate()
	c.loading = true
	for propName, widget := range c.propControls {
		value := c.propertiesContainer.PropertyValue(propName)

		if c.propsMap[propName].Type == uiproperties.PropertyTypeString {
			txtBox := widget.(*uicontrols.TextBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetText(fmt.Sprint(c.propertiesContainer.PropertyValue(propName)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeMultiline {
			txtBox := widget.(*uicontrols.Button)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.TempData = fmt.Sprint(c.propertiesContainer.PropertyValue(propName))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeInt {
			txtBox := widget.(*uicontrols.SpinBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetValue(float64(c.propertiesContainer.PropertyValue(propName).(int)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeInt32 {
			txtBox := widget.(*uicontrols.SpinBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetValue(float64(value.(int32)))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeBool {
			txtBox := widget.(*uicontrols.CheckBox)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetChecked(c.propertiesContainer.PropertyValue(propName).(bool))
			}
		}
		if c.propsMap[propName].Type == uiproperties.PropertyTypeColor {
			txtBox := widget.(*uicontrols.ColorPicker)
			txtBox.BeginUpdate()
			if value != nil {
				txtBox.SetColor(c.propertiesContainer.PropertyValue(propName).(color.Color))
			}
		}
	}

	for _, widget := range c.propControls {
		widget.EndUpdate()
	}
	c.loading = false
	c.EndUpdate()
}

func (c *PropertiesEditor) TextBoxChanged(txtBox *uicontrols.TextBox, oldValue string, newValue string) {
	if c.loading {
		return
	}
	propName := txtBox.UserData("propName").(string)
	if _, ok := c.propControls[propName]; ok {
		c.propertiesContainer.SetPropertyValue(propName, newValue)
	}
}

func (c *PropertiesEditor) MultilineChanged(txtBox *uicontrols.Button, newValue string) {
	if c.loading {
		return
	}
	propName := txtBox.UserData("propName").(string)
	if _, ok := c.propControls[propName]; ok {
		c.propertiesContainer.SetPropertyValue(propName, newValue)
	}
}

func (c *PropertiesEditor) SpinBoxChanged(spinBox *uicontrols.SpinBox, value float64) {
	if c.loading {
		return
	}
	propName := spinBox.UserData("propName").(string)
	propType := spinBox.UserData("propType").(uiproperties.PropertyType)
	if _, ok := c.propControls[propName]; ok {
		if propType == uiproperties.PropertyTypeInt {
			c.propertiesContainer.SetPropertyValue(propName, int(value))
		}
		if propType == uiproperties.PropertyTypeInt32 {
			c.propertiesContainer.SetPropertyValue(propName, int32(value))
		}
	}
}

func (c *PropertiesEditor) CheckBoxChanged(checkBox *uicontrols.CheckBox, checked bool) {
	if c.loading {
		return
	}
	propName := checkBox.UserData("propName").(string)
	if _, ok := c.propControls[propName]; ok {
		c.propertiesContainer.SetPropertyValue(propName, checked)
	}
}

func (c *PropertiesEditor) ColorPickerChanged(colorPicker *uicontrols.ColorPicker, color color.Color) {
	if c.loading {
		return
	}
	propName := colorPicker.UserData("propName").(string)
	if _, ok := c.propControls[propName]; ok {
		c.propertiesContainer.SetPropertyValue(propName, color)
	}
}

func (c *PropertiesEditor) OnPropertyChanged(prop *uiproperties.Property) {
	c.LoadPropertiesValues()
}
*/
