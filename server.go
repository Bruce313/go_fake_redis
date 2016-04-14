package main

import (
	"log"
)

func main() {
	nwl, err := newNetworkListener("0.0.0.0", 3333)
	if err != nil {
		log.Fatal("new net worker listener:", err)
	}
	pp := newProtocolPlain()
	sc := newScheduler(nwl, pp)
	err = sc.loop()
	if err != nil {
		log.Fatal("loop end:", err)
	}
}
