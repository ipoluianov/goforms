package ui

type WidgetEvent struct {
	listeners []func()
}

func (c *WidgetEvent) AddListener(l func()) int {
	c.listeners = append(c.listeners, l)
	return len(c.listeners) - 1
}

func (c *WidgetEvent) RemoveListener(id int) {

}

func (c *WidgetEvent) RemoveAllListeners() {
}

func (c *WidgetEvent) Invoke() {
	for _, l := range c.listeners {
		l()
	}
}
