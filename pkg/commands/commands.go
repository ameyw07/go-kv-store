package commands

import "github.com/ameyw07/go-kv-store/pkg/store"

type commandHandler func(args []string, store *store.Store) (string, error)

type CommandRegistry struct {
	Handlers map[string]commandHandler
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		Handlers: make(map[string]commandHandler),
	}
}

func (cr *CommandRegistry) AddCommand(name string, cmdHandler commandHandler) {
	cr.Handlers[name] = cmdHandler
}

var CmdRegistry CommandRegistry = CommandRegistry{
	Handlers: make(map[string]commandHandler),
}

type CommandMETA struct {
	commandName string
}

const INF = "INF"
