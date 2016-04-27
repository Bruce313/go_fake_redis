package main

import (
	"errors"
	"syscall"
	"io"
	"fmt"
	d "github.com/tj/go-debug"
	"os"
	"time"
)

var desch = d.Debug("go_redis:scheduler")

type scheduler struct {
	isEnd bool
	nwl       *networkListener
	fd2Client map[int]*client
	pcl       protocol
	dbm *dbManager
}

func newScheduler(nwl *networkListener, pcl protocol, dbm *dbManager) *scheduler {
	return &scheduler{
		nwl:       nwl,
		fd2Client: make(map[int]*client, 0),
		pcl:       pcl,
		dbm: dbm,
	}
}

func (shd *scheduler) end() {
	shd.isEnd = true
}

func (shd *scheduler) destroy() {
	for fd, c := range shd.fd2Client {
		syscall.Close(fd)
		c.destroy()
	}
}

func (shd *scheduler) loop() (err error) {
	for {
		//check ended by user or not 
		if shd.isEnd {
			desch("ended by user or signal, quit loop")
			return nil
		} 
		//check if there is fd event
		fes, errScan := shd.nwl.scan(time.Second * 1)
		if errScan != nil {
			return errScan
		}
		//handle file events
		for _, fe := range fes {
			errHandle := shd.handleFileEvent(fe)
			if errHandle != nil {
				return errHandle
			}
		}
		//
	}
}

func (shd *scheduler) handleFileEvent(fe *fileEvent) (err error) {
	desch("handle file event")
	if fe.eType == fileEventTypeNew {
		//new client
		desch("new client")
		nc := client{
			fd: fe.fd,
			db: shd.dbm.getDefaultDB(),
		}
		shd.fd2Client[fe.fd] = &nc
		connInfo := `connected`
		w := os.NewFile(uintptr(fe.fd), "write to conn fd")
		//defer w.Close()
		err = shd.pcl.Send(connInfo, w)
		return
	}
	//get client
	oc, ok := shd.fd2Client[fe.fd]
	if !ok {
		return fmt.Errorf("client with fd:%d not found in fd2Client", fe.fd)
	}
	//read from conn fd
	r := os.NewFile(uintptr(fe.fd), "read from conn fd")
	//defer r.Close()
	cmdStr, err := shd.pcl.Receive(r)
	if err != nil {
		return errors.New("receive from conn err:" + err.Error())
	}
	desch("get cmdStr:%s", cmdStr)
	cmd, err := lookupCommand(cmdStr)
	if err == io.EOF {
		desch("connection close, rm fd from epoll:%d", fe.fd)
		shd.nwl.rmFd(fe.fd)
		return nil 
	}
	if err != nil {
		return errors.New("look up command err:" + err.Error())
	}
	res, err := cmd.exec(oc)
	if err != nil {
		return errors.New("exec command:" + err.Error())
	}
	w := os.NewFile(uintptr(fe.fd), "write to conn fd")
	//defer w.Close()
	err = shd.pcl.Send(res, w)
	return
}
