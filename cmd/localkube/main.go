package main

import (
	"fmt"
	"os"
	"os/signal"

	"rsprd.com/localkube"
)

var LK *localkube.LocalKube

func load() {
	LK = new(localkube.LocalKube)

	// setup etc
	etcd := localkube.NewEtcd()
	LK.Add(etcd)

	// setup apiserver
	apiserver := localkube.NewAPIServer()
	LK.Add(apiserver)

	// setup controller-manager
	controllerManager := localkube.NewControllerManagerServer()
	LK.Add(controllerManager)

	// setup scheduler
	// setup kubelet (configured for weave proxy)
}

func main() {
	// check for network

	// if first
	load()
	err := LK.Run(os.Args, os.Stderr)
	if err != nil {
		fmt.Printf("localkube errored: %v\n", err)
		os.Exit(1)
	}
	defer LK.StopAll()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	<-interruptChan
	fmt.Printf("\nShutting down...\n")
}
