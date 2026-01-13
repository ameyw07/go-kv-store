package shard

import (
	"hash/maphash"
	"log"
)

type ShardController struct {
	shards    []*Shard
	numShards int
}

func NewShardController(numShards int) *ShardController {

	sc := &ShardController{
		shards:    make([]*Shard, 0),
		numShards: numShards,
	}

	for i := 0; i < numShards; i++ {
		sc.shards = append(sc.shards, NewShard(i))
	}

	return sc
}

func (sc *ShardController) chooseShard(key string) uint64 {

	var h maphash.Hash
	h.Write([]byte(key))
	hashValue := h.Sum64()
	shardIdx := hashValue % uint64(sc.numShards)
	log.Printf("%d shard chosen for key = %s", shardIdx, key)
	return shardIdx
}
