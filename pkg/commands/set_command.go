package commands

import "github.com/ameyw07/go-kv-store/pkg/store"

const CMD_NAME = "ADD"

func init() {
	CmdRegistry.AddCommand(CMD_NAME, SetCmdHandler)
}

func SetCmdHandler(args string, store *store.Store) (string, error) {
	return "OK", nil
}
