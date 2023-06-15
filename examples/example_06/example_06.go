package example06

import (
	"strconv"

	"github.com/ipoluianov/goforms/ui"
)

func newMainForm() *ui.Form {
	form := ui.NewForm()
	column := form.Panel().AddVPanel()
	cmb := column.AddComboBox()
	for i := 0; i < 50; i++ {
		cmb.AddItem("Item_"+strconv.FormatInt(int64(i), 10), i)
	}
	column.AddVSpacer()
	return form
}

func ExecMainForm() {
	ui.InitUI()
	ui.StartMainForm(newMainForm())
}
