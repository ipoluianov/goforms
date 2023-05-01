package uiproperties

import "reflect"

type PropertiesChangesListItem struct {
	Name          string
	Value         interface{}
	PropContainer IPropertiesContainer
}

type PropertiesChangesList struct {
	items []*PropertiesChangesListItem
}

func (c *PropertiesChangesList) Items() []*PropertiesChangesListItem {
	return c.items
}

func NewPropertiesChangesList() *PropertiesChangesList {
	var c PropertiesChangesList
	c.items = make([]*PropertiesChangesListItem, 0)
	return &c
}

func (c *PropertiesChangesList) AddItem(propContainer IPropertiesContainer, propName string, value interface{}) {
	var item PropertiesChangesListItem
	item.PropContainer = propContainer
	item.Name = propName
	item.Value = value
	c.items = append(c.items, &item)
}

type PropertiesChangesStack struct {
	lists []*PropertiesChangesList
}

func NewPropertiesChangesStack() *PropertiesChangesStack {
	var c PropertiesChangesStack
	c.lists = make([]*PropertiesChangesList, 0)
	return &c
}

func (c *PropertiesChangesStack) AddList(list *PropertiesChangesList) {
	c.lists = append([]*PropertiesChangesList{list}, c.lists...)
}

func (c *PropertiesChangesStack) Undo() *PropertiesChangesList {
	list := NewPropertiesChangesList()

	if len(c.lists) < 2 {
		return list
	}

	topList := c.lists[0] // Names
	for _, nameItem := range topList.items {
		var valueSourceItem *PropertiesChangesListItem
		for listIndex := 1; listIndex < len(c.lists); listIndex++ {
			for _, item := range c.lists[listIndex].items {
				if item.Name == nameItem.Name {
					p1 := reflect.ValueOf(item.PropContainer).Pointer()
					p2 := reflect.ValueOf(nameItem.PropContainer).Pointer()
					if p1 == p2 {
						valueSourceItem = item
						break
					}
				}
			}

			if valueSourceItem != nil {
				break
			}
		}

		if valueSourceItem != nil {
			list.items = append(list.items, valueSourceItem)
		}
	}

	c.lists = append([]*PropertiesChangesList{}, c.lists[1:]...)

	return list
}
