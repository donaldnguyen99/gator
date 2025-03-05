package cli

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]handler
}

func (c *commands) register(name string, f handler) {
	if _, ok := c.handlers[name]; ok {
		panic(fmt.Sprintf("command %s already registered", name))
	}
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	h, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("command %s not registered", cmd.name)
	}
	
	if err := h(s, cmd); err != nil {
		return fmt.Errorf("failed to execute %s: %w", cmd.name, err)
	}
	return nil
}


