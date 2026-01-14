package shard

import (
	"log"
	"strings"
	"unsafe"

	"github.com/ameyw07/go-kv-store/pkg/commands"
	"github.com/ameyw07/go-kv-store/pkg/store"
	"golang.org/x/sys/unix"
)

type Shard struct {
	shardId int
	store   *store.Store
	// for sending commands to the correct shard
	CmdsShare chan *CmdData
}

func NewShard(id int) *Shard {

	return &Shard{
		shardId:   id,
		store:     store.NewStore(),
		CmdsShare: make(chan *CmdData, CMDDATA_CHANNEL_BOUND),
	}
}

func (sh *Shard) ExecuteCommand() {

	for {

		cmdData := <-sh.CmdsShare
		cmdArgs := strings.Split(cmdData.Cmd, " ")
		cmdName := cmdArgs[0]
		resp, err := commands.CmdRegistry.Handlers[cmdName](cmdArgs[1:], sh.store)

		var outResp string
		if err != nil {
			outResp = err.Error()
		} else {
			outResp = resp
		}

		cmdData.OutResp = outResp

		var val uint64 = 1
		// eventfd write requires an 8-byte value
		buf := (*(*[8]byte)(unsafe.Pointer(&val)))[:]
		_, err = unix.Write(cmdData.WakeUpFd, buf)
		if err != nil {
			log.Printf("eventfd write error: %v\n", err)
		}
		log.Println("Wrote to eventfd to wake up epoll_wait")

	}

}
