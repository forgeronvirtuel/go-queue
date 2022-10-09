package command

import (
	"bytes"
	"strings"
)

type Manager struct {
	cmds    []Item
	actions []func([]string) func() error
}

func NewManager() *Manager {
	return &Manager{}
}

func (c *Manager) Add(item Item, action func([]string) error) {
	c.cmds = append(c.cmds, item)
	f := func(params []string) func() error {
		return func() error {
			return action(params)
		}
	}
	c.actions = append(c.actions, f)
}

func (c *Manager) Parse(cmd []byte) func() error {
	for idx, item := range c.cmds {
		if params, ok := item.parse(cmd); ok {
			return c.actions[idx](params)
		}
	}
	return nil
}

type Item struct {
	name []byte
}

func NewItem(name string) Item {
	nm := []byte(name)
	return Item{
		name: nm,
	}
}

func (i Item) parse(command []byte) ([]string, bool) {
	if len(command) < len(i.name) {
		return nil, false
	}
	sub := command[:len(i.name)]
	if bytes.Compare(sub, i.name) != 0 {
		return nil, false
	}
	remaining := string(command[len(i.name)+1:])
	return strings.Split(remaining, " "), true
}
