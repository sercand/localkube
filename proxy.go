package localkube

import (
	kubeproxy "k8s.io/kubernetes/cmd/kube-proxy/app"
	"k8s.io/kubernetes/cmd/kube-proxy/app/options"
)

const (
	ProxyName = "proxy"
)

func NewProxyServer() Server {
	return SimpleServer{
		ComponentName: ProxyName,
		StartupFn:     StartProxyServer,
	}.NoShutdown()
}

func StartProxyServer() {
	config := options.NewProxyConfig()

	// master details
	config.Master = APIServerURL

	server, err := kubeproxy.NewProxyServerDefault(config)
	if err != nil {
		panic(err)
	}

	go server.Run()
}
