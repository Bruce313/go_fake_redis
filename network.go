package main

import (
	"fmt"
	d "github.com/tj/go-debug"
	"math"
	"net"
	"syscall"
	"time"
)

var debugNet = d.Debug("go_redis:network")

const maxPOLLEVENTS = 10
const listenBACKLOG = 10

type fileEventType int

const (
	fileEventTypeNew fileEventType = iota
	fileEventTypeReadable
)

type fileEvent struct {
	fd    int
	eType fileEventType
}

type networkListener struct {
	sockFd int
	epFd   int
}

func newNetworkListener(ip string, port int) (nl *networkListener, err error) {
	sockFd, err := syscall.Socket(syscall.AF_INET,
		syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return
	}
	//set no block
	if err = syscall.SetNonblock(sockFd, true); err != nil {
		return
	}
	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP(ip).To4())
	err = syscall.Bind(sockFd, &addr)
	if err != nil {
		return
	}
	//
	err = syscall.Listen(sockFd, listenBACKLOG)
	if err != nil {
		return
	}
	epFd, err := syscall.EpollCreate1(0)
	if err != nil {
		return
	}
	ee := &syscall.EpollEvent{
		Fd:     int32(sockFd),
		Events: syscall.EPOLLIN,
	}
	err = syscall.EpollCtl(epFd, syscall.EPOLL_CTL_ADD, sockFd, ee)
	if err != nil {
		return
	}
	nl = &networkListener{
		sockFd: sockFd,
		epFd:   epFd,
	}
	return
}

func (nl *networkListener) scan(t time.Duration) (firedFd []*fileEvent, err error) {
	firedFd = make([]*fileEvent, 0)
	msec := int(float64(t) / (math.Pow(float64(10), float64(6))))
	var events [maxPOLLEVENTS]syscall.EpollEvent
	debugNet("begin epoll_wait\n")
	count, err := syscall.EpollWait(nl.epFd, events[:], msec)
	debugNet("get count:%d for this epoll_wait", count)
	if err != nil {
		return nil, fmt.Errorf("epoll wait:%s, epFd:%d\n", err.Error(), nl.epFd)
	}
	debugNet("get count:%d for this epoll_wait", count)
	if count == 0 {
		debugNet("on file event happen after epoll wait for %d miliseconds", msec)
		return
	}
	for _, e := range events[:count] {
		debugNet("get epollEvent.events:%x, fd:%d, sockFd:%d", e.Events, e.Fd, nl.sockFd)
		var fe fileEvent
		if int(e.Fd) == nl.sockFd {
			//accept new conn
			nfd, _, errAcc := syscall.Accept(nl.sockFd)
			if errAcc != nil {
				return nil, fmt.Errorf("accept:%s, epFd:%d", errAcc.Error(), nl.sockFd)
			}
			fe.eType = fileEventTypeNew
			fe.fd = nfd
			//epoll add
			ee := &syscall.EpollEvent{
				Fd:     int32(nfd),
				Events: syscall.EPOLLIN,
			}
			errEctl := syscall.EpollCtl(nl.epFd, syscall.EPOLL_CTL_ADD, nfd, ee)
			if errEctl != nil {
				return nil, errEctl
			}
		} else {
			fe.fd = int(e.Fd)
			fe.eType = fileEventTypeReadable
		}
		firedFd = append(firedFd, &fe)
	}
	return
}
//rm fd when conn close
func (nl *networkListener) rmFd(fd int) {
	ee := &syscall.EpollEvent{
		Fd: int32(fd),
		Events: syscall.EPOLLIN,
	}
	//ignore errors 
	syscall.EpollCtl(nl.epFd, syscall.EPOLL_CTL_DEL, fd, ee)
	return
}
func (nl *networkListener) destroy() {
	syscall.Close(nl.sockFd)
	syscall.Close(nl.epFd)
}
