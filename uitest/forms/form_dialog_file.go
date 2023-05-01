package forms

import (
	"github.com/gazercloud/gazerui/uiforms"
)

type FormDialogFile struct {
	uiforms.Form
}

func (c *FormDialogFile) OnInit() {
	c.Resize(600, 400)
	c.SetTitle("FormDialogFile")

	/*c.Panel().AddButtonOnGrid(0, 0, "Clear", func(event *uievents.Event) {
		runtime.GC()
		debug.FreeOSMemory()
	})*/
}
