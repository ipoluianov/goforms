package forms

import (
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiforms"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"reflect"
)

type MainForm struct {
	uiforms.Form
	currentX int
	currentY int

	WindowToCreate uiinterfaces.Window
}

func (c *MainForm) OnInit() {

	c.currentX = 0
	c.currentY = 0

	c.SetTitle("UI test application")
	c.Resize(1200, 800)

	//c.Panel().AddButtonOnGrid(0, 0, "123", nil)
	//c.Panel().AddButtonOnGrid(1, 0, "123", nil)
	//c.Panel().AddButtonOnGrid(0, 1, "123", nil)
	//c.Panel().AddVSpacerOnGrid(0, 5)
	//c.Panel().AddHSpacerOnGrid(1, 0)

	/*t1 := c.Panel().AddTextBlockOnGrid(0, 0, "123")
	t1.SetMaxWidth(100)
	t2 := c.Panel().AddTextBlockOnGrid(0, 1, "456")
	t2.SetMaxWidth(150)
	c.Panel().AddImageBoxOnGrid(0, 2, uiresources.ResImageAdjusted("icons/material/image/drawable-hdpi/ic_blur_on_black_48dp.png", c.Panel().ForeColor()))
	*/
	c.addButton(FormTextBlock{})
	c.addButton(FormImageBox{})
	c.addButton(FormButton{})
	c.addButton(FormCheckBox{})
	c.addButton(FormComboBox{})
	c.addButton(FormProgressBar{})
	c.addButton(FormRadioButton{})
	c.addButton(FormSpinBox{})
	c.addButton(FormTextBox{})
	c.addButton(FormTable{})
	c.addButton(FormTreeView{})
	c.addButton(FormSplitPanel{})
	c.addButton(FormTabControl{})
	c.addButton(FormContextMenu{})
	c.addButton(FormDateTimePicker{})
	c.addButton(FormColorPicker{})
	c.addButton(FormListView{})
	c.addSeparator()

	c.addButton(FormStatusBar{})
	c.addButton(FormToolBar{})
	c.addButton(FormMenu{})
	c.addSeparator()

	c.addButton(FormDialogColor{})
	c.addButton(FormDialogDirectory{})
	c.addButton(FormDialogFile{})
	c.addButton(FormDialogFont{})
	c.addButton(FormMessageBox{})
	c.addButton(FormFont{})
	c.addButton(FormGridLayout{})
	c.addButton(FormAbsLayout{})
	c.addButton(FormCanvas{})
	c.addButton(FormDialog{})
	c.addSeparator()
}

func (c *MainForm) addButton(i interface{}) {
	if separator {
		if c.currentX == 0 {
			c.Panel().AddVSpacerOnGrid(0, 50)
		}
		c.currentY = 0
		c.currentX++
		separator = false
	}

	btnForm := c.Panel().AddButtonOnGrid(c.currentX, c.currentY, "Форма "+reflect.TypeOf(i).Name(), c.onBtnForm)
	btnForm.SetUserData("formType", i)
	c.currentY++
}

var separator = false

func (c *MainForm) addSeparator() {
	separator = true
}

func (c *MainForm) onBtnForm(ev *uievents.Event) {
	c.WindowToCreate = reflect.New(reflect.TypeOf(ev.Sender.(uiinterfaces.Widget).UserData("formType"))).Interface().(uiinterfaces.Window)
	c.Close()
}
