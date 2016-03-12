package main

import (
	"fmt"
	"os"
	"os/signal"

	"rsprd.com/localkube"
)

var (
	Servers = localkube.Servers{}
)

func init() {
	// setup etc
	etcd := localkube.NewEtcd()
	Servers = append(Servers, etcd)

	// setup apiserver
	// setup controller-manager
	// setup scheduler
	// setup kubelet (configured for weave proxy)
}

func main() {
	Servers.StartAll()
	defer Servers.StopAll()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	<-interruptChan
	fmt.Printf("\nShutting down...\n")
}
