package localkube

import (
	kubelet "k8s.io/kubernetes/cmd/kubelet/app"
	"k8s.io/kubernetes/cmd/kubelet/app/options"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

const (
	KubeletName = "kubelet"
)

var (
	DockerDaemonSock = "/var/run/weave.sock"
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

	// defaults from command
	config.ResolverConfig = kubetypes.ResolvConfDefault
	config.DockerEndpoint = DockerDaemonSock

	go kubelet.Run(config, nil)
}
