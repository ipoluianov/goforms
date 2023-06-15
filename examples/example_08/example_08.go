package example08

import (
	"strconv"

	"github.com/ipoluianov/goforms/ui"
)

func newMainForm() *ui.Form {
	form := ui.NewForm()
	column := form.Panel().AddVPanel()
	listView := column.AddListView()
	listView.AddColumn("Column1", 200)
	listView.AddColumn("Column2", 100)
	listView.AddColumn("Column2", 100)
	for r := 0; r < 100; r++ {
		listView.AddItem3("row_"+strconv.FormatInt(int64(r), 10), "222", "333")
	}
	return form
}

func ExecMainForm() {
	ui.InitUI()
	ui.StartMainForm(newMainForm())
}
