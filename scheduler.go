package main

import (
	"errors"
	"fmt"
	. "github.com/tj/go-debug"
	"os"
	"time"
)

var desch = Debug("go_redis:scheduler")

type scheduler struct {
	nwl       *networkListener
	fd2Client map[int]*client
	pcl       protocol
}

func newScheduler(nwl *networkListener, pcl protocol) *scheduler {
	return &scheduler{
		nwl:       nwl,
		fd2Client: make(map[int]*client, 0),
		pcl:       pcl,
	}
}

func (self *scheduler) loop() (err error) {
	for {
		//see if there is fd event
		fes, errScan := self.nwl.scan(time.Second * 1)
		if errScan != nil {
			return errScan
		}
		//handle file events
		for _, fe := range fes {
			errHandle := self.handleFileEvent(fe)
			if errHandle != nil {
				return errHandle
			}
		}
		//
	}
	return
}

func (self *scheduler) handleFileEvent(fe *fileEvent) (err error) {
	desch("handle file event")
	if fe.eType == fileEventTypeNew {
		//new client
		desch("new client")
		nc := client{
			fd: fe.fd,
		}
		self.fd2Client[fe.fd] = &nc
		connInfo := `connected`
		w := os.NewFile(uintptr(fe.fd), "write to conn fd")
		//defer w.Close()
		err = self.pcl.Send(connInfo, w)
		return
	}
	//get client
	oc, ok := self.fd2Client[fe.fd]
	if !ok {
		return errors.New(fmt.Sprintf("client with fd:%d not found in fd2Client", fe.fd))
	}
	//read from conn fd
	r := os.NewFile(uintptr(fe.fd), "read from conn fd")
	//defer r.Close()
	cmdStr, err := self.pcl.Receive(r)
	if err != nil {
		return errors.New("receive from conn err:" + err.Error())
	}
	desch("get cmdStr:%s", cmdStr)
	cmd, err := lookupCommand(cmdStr)
	if err != nil {
		return errors.New("look up command err:" + err.Error())
	}
	res, err := cmd.exec(oc)
	if err != nil {
		return errors.New("exec command:" + err.Error())
	}
	w := os.NewFile(uintptr(fe.fd), "write to conn fd")
	//defer w.Close()
	err = self.pcl.Send(res, w)
	return
}
