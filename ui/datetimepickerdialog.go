package ui

import (
	"time"
)

type DateTimePickerDialog struct {
	Dialog
	lblYear   *TextBlock
	lblMonth  *TextBlock
	lblDay    *TextBlock
	lblHour   *TextBlock
	lblMinute *TextBlock
	lblSecond *TextBlock

	numYear *SpinBox
	//numMonth  *SpinBox
	//numDay    *SpinBox
	numHour   *SpinBox
	numMinute *SpinBox
	numSecond *SpinBox

	dayOfMonthSelector *DayOfMonthSelector
	//monthSelector      *MonthSelector
	cmbMonth *ComboBox

	btnOK     *Button
	btnCancel *Button

	dateTime time.Time

	loading           bool
	onDateTimeChanged func(dateTime time.Time)
}

func NewDateTimePickerDialog(parent Widget) *DateTimePickerDialog {
	var c DateTimePickerDialog
	c.InitControl(parent, &c)
	c.loading = true

	pTime := c.ContentPanel().AddPanelOnGrid(0, 0)
	pYearMonth := c.ContentPanel().AddPanelOnGrid(0, 1)
	pDay := c.ContentPanel().AddPanelOnGrid(0, 2)
	pButtons := c.ContentPanel().AddPanelOnGrid(0, 4)
	pButtons.SetXExpandable(true)

	lblTime := pTime.AddTextBlockOnGrid(0, 1, "Time")
	lblTime.SetMinWidth(60)
	c.lblHour = pTime.AddTextBlockOnGrid(1, 0, "Hour")
	c.numHour = pTime.AddSpinBoxOnGrid(1, 1)
	c.numHour.SetPrecision(0)
	c.numHour.SetMinValue(0)
	c.numHour.SetMaxValue(23)
	c.numHour.OnValueChanged = c.numsChanged

	c.lblMinute = pTime.AddTextBlockOnGrid(2, 0, "Min")
	c.numMinute = pTime.AddSpinBoxOnGrid(2, 1)
	c.numMinute.SetPrecision(0)
	c.numMinute.SetMinValue(0)
	c.numMinute.SetMaxValue(59)
	c.numMinute.OnValueChanged = c.numsChanged

	c.lblSecond = pTime.AddTextBlockOnGrid(3, 0, "Sec")
	c.numSecond = pTime.AddSpinBoxOnGrid(3, 1)
	c.numSecond.SetPrecision(0)
	c.numSecond.SetMinValue(0)
	c.numSecond.SetMaxValue(59)
	c.numSecond.OnValueChanged = c.numsChanged

	pTime.AddHSpacerOnGrid(4, 0)

	lblDate := pYearMonth.AddTextBlockOnGrid(0, 1, "Date")
	lblDate.SetMinWidth(60)
	lblYear := pYearMonth.AddTextBlockOnGrid(1, 0, "Year")
	lblYear.SetMinWidth(60)
	c.numYear = pYearMonth.AddSpinBoxOnGrid(1, 1)
	c.numYear.SetPrecision(0)
	c.numYear.OnValueChanged = c.numsChanged

	c.lblMonth = pYearMonth.AddTextBlockOnGrid(2, 0, "Month")
	c.lblMonth.SetMinWidth(60)
	c.cmbMonth = pYearMonth.AddComboBoxOnGrid(2, 1)
	c.cmbMonth.AddItem("01", 1)
	c.cmbMonth.AddItem("02", 2)
	c.cmbMonth.AddItem("03", 3)
	c.cmbMonth.AddItem("04", 4)
	c.cmbMonth.AddItem("05", 5)
	c.cmbMonth.AddItem("06", 6)
	c.cmbMonth.AddItem("07", 7)
	c.cmbMonth.AddItem("08", 8)
	c.cmbMonth.AddItem("09", 9)
	c.cmbMonth.AddItem("10", 10)
	c.cmbMonth.AddItem("11", 11)
	c.cmbMonth.AddItem("12", 12)
	c.cmbMonth.OnCurrentIndexChanged = func(event *ComboBoxEvent) {
		c.dateTimeChanged()
	}
	/*c.monthSelector = NewMonthSelector(pDay)
	c.monthSelector.SetMonth(1)
	c.monthSelector.MonthChanged = c.dateTimeChanged
	pMonth.AddWidgetOnGrid(c.monthSelector, 1, 0)*/

	//c.lblDay = pDay.AddTextBlockOnGrid(0, 0, "Day")
	//c.lblDay.SetMinWidth(60)
	c.dayOfMonthSelector = NewDayOfMonthSelector(pDay)
	c.dayOfMonthSelector.SetYearAndMonth(1970, 1)
	c.dayOfMonthSelector.SetDay(1)
	c.dayOfMonthSelector.DayChanged = c.dateTimeChanged
	pDay.AddWidgetOnGrid(c.dayOfMonthSelector, 0, 1)

	pButtons.AddHSpacerOnGrid(0, 0)
	c.btnOK = pButtons.AddButtonOnGrid(1, 0, "OK", nil)
	c.btnOK.SetMinWidth(70)
	c.btnCancel = pButtons.AddButtonOnGrid(2, 0, "Cancel", nil)
	c.btnCancel.SetMinWidth(70)

	c.SetAcceptButton(c.btnOK)
	c.SetRejectButton(c.btnCancel)

	c.loading = false

	c.SetDateTime(c.dateTime)

	return &c
}

func (c *DateTimePickerDialog) OnInit() {
	c.Dialog.OnInit()
	c.Resize(400, 450)
	c.SetTitle("Select date/time ...")
}

func (c *DateTimePickerDialog) SetDateTime(dateTime time.Time) {
	c.dateTime = dateTime

	c.loading = true
	if c.numYear != nil {
		c.numYear.SetValue(float64(c.dateTime.Year()))
	}

	/*if c.monthSelector != nil {
		c.monthSelector.SetMonth(int(c.dateTime.Month()))
	}*/
	if c.cmbMonth != nil {
		c.cmbMonth.SetCurrentItemIndex(int(c.dateTime.Month()) - 1)
		c.dayOfMonthSelector.SetYearAndMonth(c.dateTime.Year(), int(c.dateTime.Month()))
	}

	if c.dayOfMonthSelector != nil {
		c.dayOfMonthSelector.SetYearAndMonth(c.dateTime.Year(), int(c.dateTime.Month()))
		c.dayOfMonthSelector.SetDay(c.dateTime.Day())
	}

	if c.numHour != nil {
		c.numHour.SetValue(float64(c.dateTime.Hour()))
	}
	if c.numMinute != nil {
		c.numMinute.SetValue(float64(c.dateTime.Minute()))
	}
	if c.numSecond != nil {
		c.numSecond.SetValue(float64(c.dateTime.Second()))
	}
	c.loading = false
}

func (c *DateTimePickerDialog) DateTime() time.Time {
	loc, _ := time.LoadLocation("")
	//dateTime := time.Date(int(c.numYear.Value()), time.Month(c.monthSelector.Month()), c.dayOfMonthSelector.Day(), int(c.numHour.Value()), int(c.numMinute.Value()), int(c.numSecond.Value()), 0, loc)
	dateTime := time.Date(int(c.numYear.Value()), time.Month(c.cmbMonth.CurrentItemIndex+1), c.dayOfMonthSelector.Day(), int(c.numHour.Value()), int(c.numMinute.Value()), int(c.numSecond.Value()), 0, loc)
	return dateTime
}

func (c *DateTimePickerDialog) numsChanged(spinBox *SpinBox, value float64) {
	c.dateTimeChanged()
}

func (c *DateTimePickerDialog) dateTimeChanged() {
	if c.loading {
		return
	}
	if c.onDateTimeChanged != nil {
		c.dateTime = c.DateTime()
		c.onDateTimeChanged(c.dateTime)
	}

	//c.dayOfMonthSelector.SetYearAndMonth(int(c.numYear.Value()), c.monthSelector.Month())
	c.dayOfMonthSelector.SetYearAndMonth(int(c.numYear.Value()), c.cmbMonth.CurrentItemIndex+1)
}
