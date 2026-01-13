package server

import (
	"log"
	"net"

	"github.com/ameyw07/go-kv-store/pkg/shard"
	"golang.org/x/sys/unix"
)

type IOMultiplexer struct {
	wakeUpFdToConnFd map[int32]int32
}

func NewIOMultiplexer() *IOMultiplexer {

	return &IOMultiplexer{
		wakeUpFdToConnFd: make(map[int32]int32),
	}
}

func (iomx *IOMultiplexer) StartPollWorker(sc *shard.ShardController) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Get file descriptor and set non-blocking
	file, _ := ln.(*net.TCPListener).File()
	fd := int(file.Fd())
	unix.SetNonblock(fd, true)

	epfd, err := unix.EpollCreate1(0)
	defer unix.Close(epfd)

	event := &unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd:     int32(fd),
	}
	unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, fd, event)
	events := make([]unix.EpollEvent, 10)

	for {
		n_events, err := unix.EpollWait(epfd, events, 1)
		if err != nil {
			log.Println(err)
			continue
		}

		for i := 0; i < n_events; i++ {
			if events[i].Fd == int32(fd) {
				// accept conn
				conn, err := ln.Accept()
				if err != nil {
					log.Println(err)
					continue
				}
				file, _ := conn.(*net.TCPConn).File()
				connFd := int32(file.Fd())

				connEvent := &unix.EpollEvent{
					Events: unix.EPOLLIN,
					Fd:     int32(connFd),
				}
				unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(connFd), connEvent)

			} else {
				// handle client conn

				clientFd := int(events[i].Fd)
				buf := make([]byte, 1024)
				n, err := unix.Read(clientFd, buf)
				if n <= 0 || err != nil {
					log.Printf("Closing the connection for descriptor %d", clientFd)
					unix.EpollCtl(epfd, unix.EPOLL_CTL_DEL, clientFd, nil)
					unix.Close(clientFd)
					continue
				}
			}

		}

	}

}
