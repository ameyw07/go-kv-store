package shard

import (
	"strings"

	"github.com/ameyw07/go-kv-store/pkg/store"
)

type ShardController struct {
	shardId int
	store   *store.Store
	// for sending commands to the correct shard
	cmdsShare chan string
}

func NewShardController(id int) *ShardController {

	return &ShardController{
		shardId:   id,
		store:     store.NewStore(),
		cmdsShare: make(chan string, 100),
	}
}

func (sc *ShardController) chooseShard(key string) {

	var h maphash.Hasher
	h.Write([]byte(key))
	hashValue := h.Sum64()

}

func (sc *ShardController) ExecuteCommand(data []byte) {
	cmd := string(data)
	cmd_name := strings.Split(cmd, " ")[0]

}
