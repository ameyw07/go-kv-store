package commands

import (
	"errors"
	"strconv"
	"time"

	"github.com/ameyw07/go-kv-store/pkg/store"
)

var cSCMeta = &CommandMETA{
	commandName: "GET",
}

func init() {
	CmdRegistry.AddCommand(cSCMeta.commandName, SetCmdHandler)
}

func SetCmdHandler(args []string, st *store.Store) (string, error) {

	key := args[0]
	value := args[1]

	hasExpiry := false
	var futureTime *time.Time = nil
	if len(args) == 3 {

		expirationInt, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return "", err
		}

		if expirationInt <= 0 {
			return "", errors.New("Expiry duration should be positive")
		}
		hasExpiry = true
		durationToAdd := time.Duration(expirationInt) * time.Millisecond
		now := time.Now()
		ft := now.Add(durationToAdd)
		futureTime = &ft
	}
	st.VersionSeq++

	st.Data[key] = append(
		st.Data[key],
		store.GetNewEntry(
			value,
			st.VersionSeq,
			futureTime,
			hasExpiry,
		),
	)

	return "OK", nil
}
