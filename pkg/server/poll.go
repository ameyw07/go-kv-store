package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"syscall"

	"github.com/ameyw07/go-kv-store/pkg/shard"
	"golang.org/x/sys/unix"
)

type IOMultiplexer struct {
	wakeUpFdToConnFd map[int32]*shard.CmdData
}

func NewIOMultiplexer() *IOMultiplexer {

	return &IOMultiplexer{
		wakeUpFdToConnFd: make(map[int32]*shard.CmdData),
	}
}

func (iomx *IOMultiplexer) StartPollWorker(sc *shard.ShardController, wg *sync.WaitGroup) {
	defer wg.Done()
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// Allow address reuse
				_ = unix.SetsockoptInt(
					int(fd),
					unix.SOL_SOCKET,
					unix.SO_REUSEADDR,
					1,
				)

				// Allow port reuse (Linux)
				_ = unix.SetsockoptInt(
					int(fd),
					unix.SOL_SOCKET,
					unix.SO_REUSEPORT,
					1,
				)
			})
		},
	}

	ln, err := lc.Listen(context.Background(), "tcp", ":9999")
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

			if _, ok := err.(net.Error); ok {
				continue
			}

			// Handle EINTR explicitly
			if errors.Is(err, syscall.EINTR) {
				continue
			}

			// log.Println("Hi")
			log.Println(err)
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

			} else if cmdData, ok := iomx.wakeUpFdToConnFd[events[i].Fd]; ok {

				ob := []byte(cmdData.OutResp)
				unix.Write(cmdData.ClientFd, ob)

			} else {
				// handle client conn

				clientFd := int(events[i].Fd)
				buf := make([]byte, 1024)
				var allData []byte

				for {
					n, err := unix.Read(clientFd, buf)
					if n > 0 {
						// Append the read data to the total slice
						allData = append(allData, buf[:n]...)
					}

					if err != nil {
						if err != io.EOF {
							log.Printf("Read error: %v\n", err)
						}
						break
					}

					if n <= 0 || err != nil {
						log.Printf("Closing the connection for descriptor %d", clientFd)
						unix.EpollCtl(epfd, unix.EPOLL_CTL_DEL, clientFd, nil)
						unix.Close(clientFd)
						continue
					}
				}

				command := string(allData)

				shardIdx := sc.ChooseShard(strings.Split(command, " ")[1])

				wakeUpFd, err := unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC)
				if err != nil {
					fmt.Printf("eventfd create error: %v\n", err)
					return
				}
				defer unix.Close(wakeUpFd)

				wakeUpEvent := unix.EpollEvent{
					Events: unix.EPOLLIN | unix.EPOLLET, // Use EPOLLET for edge-triggered behavior
					Fd:     int32(wakeUpFd),
				}

				if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, wakeUpFd, &wakeUpEvent); err != nil {
					fmt.Printf("epoll ctl add error: %v\n", err)
					return
				}

				cmdData := &shard.CmdData{
					ClientFd: clientFd,
					Cmd:      command,
					WakeUpFd: wakeUpFd,
				}

				select {
				case sc.Shards[shardIdx].CmdsShare <- cmdData:
					fmt.Println("Sent the command to the hard")
				default:

					fmt.Println(CMDOVERLOAD_RESP)

					unix.Write(clientFd, []byte(CMDOVERLOAD_RESP))
				}

			}

		}

	}

}
