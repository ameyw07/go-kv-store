package shard

type CmdData struct {
	ClientFd int
	WakeUpFd int
	Cmd      string
	OutResp  string
}

const (
	CMDDATA_CHANNEL_BOUND = 1000
)
