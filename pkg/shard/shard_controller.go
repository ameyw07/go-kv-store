package shard

import (
	"hash/maphash"
	"log"
)

type ShardController struct {
	Shards    []*Shard
	numShards int
}

func NewShardController(numShards int) *ShardController {

	sc := &ShardController{
		Shards:    make([]*Shard, 0),
		numShards: numShards,
	}

	for i := range numShards {
		sc.Shards = append(sc.Shards, NewShard(i))
	}

	return sc
}

func (sc *ShardController) ChooseShard(key string) uint64 {

	var h maphash.Hash
	h.Write([]byte(key))
	hashValue := h.Sum64()
	shardIdx := hashValue % uint64(sc.numShards)
	log.Printf("%d shard chosen for key = %s", shardIdx, key)
	return shardIdx
}
