package main

import (
	"fmt"
	"os"
	"os/signal"

	"rsprd.com/localkube"
)

var (
	Servers = localkube.Servers{}

	WorkingDirectory string
)

func init() {
	if wd, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		WorkingDirectory = wd
	}

	// setup etc
	etcd := localkube.NewEtcd(WorkingDirectory)
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
