package shard

import (
	"strings"

	"github.com/ameyw07/go-kv-store/pkg/commands"
	"github.com/ameyw07/go-kv-store/pkg/store"
)

type Shard struct {
	shardId int
	store   *store.Store
	// for sending commands to the correct shard
	cmdsShare chan string
}

func NewShard(id int) *Shard {

	return &Shard{
		shardId:   id,
		store:     store.NewStore(),
		cmdsShare: make(chan string, 100),
	}
}

func (sh *Shard) ExecuteCommand() int32 {

	for {
		select {
		case cmd := <-sh.cmdsShare:
			cmdArgs := strings.Split(cmd, " ")
			cmdName := cmdArgs[0]
			commands.CmdRegistry.Handlers[cmdName](cmdArgs[1:], sh.store)

		}

	}
}
