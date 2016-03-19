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

func NewKubeletServer(clusterDomain, clusterDNS string) Server {
	return SimpleServer{
		ComponentName: KubeletName,
		StartupFn:     StartKubeletServer(clusterDomain, clusterDNS),
	}.NoShutdown()
}

func StartKubeletServer(clusterDomain, clusterDNS string) func() {
	config := options.NewKubeletServer()

	// master details
	config.APIServerList = []string{APIServerURL}

	// Docker
	config.Containerized = true
	config.DockerEndpoint = WeaveProxySock

	// Networking
	config.ClusterDomain = clusterDomain
	config.ClusterDNS = clusterDNS
	config.ResolverConfig = "/dev/null"

	return func() {
		go kubelet.Run(config, nil)
	}
}
