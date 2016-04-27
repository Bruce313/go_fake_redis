package main

import (
	"log"
	"os" 
	"os/signal"
	"syscall"
	d "github.com/tj/go-debug"
)

var deMain = d.Debug("go_redis:main")

func main() {
	nwl, err := newNetworkListener("0.0.0.0", 3333)
	if err != nil {
		log.Fatal("new net worker listener:", err)
	}
	pp := newProtocolPlain()
	dbm := newDBManager()
	sc := newScheduler(nwl, pp, dbm)
	go sc.loop()
	//wait ctrl + C
	chSig := make(chan os.Signal)
	signal.Notify((chan<-os.Signal)(chSig), syscall.SIGINT)
	<- chSig
	deMain("read signal SIGINT, quiting")
	sc.end()
	sc.destroy()
	nwl.destroy()
}
