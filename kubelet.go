package localkube

import (
	kubelet "k8s.io/kubernetes/cmd/kubelet/app"
	"k8s.io/kubernetes/cmd/kubelet/app/options"
)

const (
	KubeletName = "kubelet"
)

var (
	WeaveProxySock = "unix:///var/run/weave/weave.sock"
)

func NewKubeletServer() Server {
	return SimpleServer{
		ComponentName: KubeletName,
		StartupFn:     StartKubeletServer,
	}.NoShutdown()
}

func StartKubeletServer() {
	config := options.NewKubeletServer()

	// master details
	config.APIServerList = []string{APIServerURL}

	// Docker
	config.Containerized = true
	config.DockerEndpoint = WeaveProxySock

	// Networking
	config.ResolverConfig = "/dev/null"

	go kubelet.Run(config, nil)
}
