package commands

import "github.com/ameyw07/go-kv-store/pkg/store"

type commandHandler func(args string, store *store.Store) (string, error)

type CommandRegistry struct {
	handlers map[string]commandHandler
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers: make(map[string]commandHandler),
	}
}

func (cr *CommandRegistry) AddCommand(name string, cmdHandler commandHandler) {
	cr.handlers[name] = cmdHandler
}

var CmdRegistry CommandRegistry = CommandRegistry{
	handlers: make(map[string]commandHandler),
}
