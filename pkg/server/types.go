package server

type CmdData struct {
	wakeUpFd int
	cmd      byte
	resp     string
	err      error
}
