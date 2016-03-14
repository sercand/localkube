package localkube

import (
	"net"
	"os"

	apiserver "k8s.io/kubernetes/cmd/kube-apiserver/app"
	"k8s.io/kubernetes/cmd/kube-apiserver/app/options"
)

const (
	APIServerName = "apiserver"
)

var (
	ServiceIPRange = "10.1.30.0/24"
)

func init() {
	if ipRange := os.Getenv("SERVICE_IP_RANGE"); len(ipRange) != 0 {
		ServiceIPRange = ipRange
	}
}

func NewAPIServer() Server {
	return &SimpleServer{
		ComponentName: APIServerName,
		StartupFn:     StartAPIServer,
	}
}

func StartAPIServer() {
	config := options.NewAPIServer()

	// use localkube etcd
	config.EtcdServerList = EtcdClientURLs

	// set Service IP range
	_, ipnet, err := net.ParseCIDR(ServiceIPRange)
	if err != nil {
		panic(err)
	}
	config.ServiceClusterIPRange = *ipnet

	// defaults from apiserver command
	config.EnableProfiling = true
	config.EnableWatchCache = true
	config.MinRequestTimeout = 1800

	// start API server in it's own goroutine
	go apiserver.Run(config)
}
