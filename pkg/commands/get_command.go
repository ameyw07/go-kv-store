package commands

import (
	"errors"
	"slices"
	"time"

	"github.com/ameyw07/go-kv-store/pkg/store"
)

var cGCMeta = &CommandMETA{
	commandName: "GET",
}

func init() {
	CmdRegistry.AddCommand(cGCMeta.commandName, GetCmdHandler)
}

func GetCmdHandler(args []string, st *store.Store) (string, error) {

	key := args[0]

	currTime := time.Now().UTC()

	if existingEntries, ok := st.Data[key]; ok {
		slices.Reverse(existingEntries)
		for _, entry := range existingEntries {
			if entry.ExpirationDate.After(currTime) {
				return entry.Value, nil
			}

		}
	}

	return "", errors.New("Value for the key not found")
}
