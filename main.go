package main

import (
	"fmt"
	"os"
)

func main() {
	// setup etc

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting etcd with '%s' as data directory.", wd)
	etcd, err := Etcd(wd)
	if err != nil {
		panic(err)
	}
	etcd.Start()

	println("Stats: " + string(etcd.SelfStats()) + " ENDSTATS")

	// setup apiserver
	// setup controller-manager
	// setup scheduler
	// setup kubelet (configured for weave proxy)
}
