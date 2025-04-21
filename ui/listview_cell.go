package ui

type listViewCell struct {
	text string
}

func newListViewCell(text string) *listViewCell {
	var c listViewCell
	c.text = text
	return &c
}
