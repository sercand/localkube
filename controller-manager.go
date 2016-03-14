package localkube

import (
	controllerManager "k8s.io/kubernetes/cmd/kube-controller-manager/app"
	"k8s.io/kubernetes/cmd/kube-controller-manager/app/options"
)

const (
	ControllerManagerName = "controller-manager"
)

func NewControllerManagerServer() Server {
	return SimpleServer{
		ComponentName: ControllerManagerName,
		StartupFn:     StartControllerManagerServer,
	}.NoShutdown()
}

func StartControllerManagerServer() {
	config := options.NewCMServer()

	// defaults from command
	config.DeletingPodsQps = 0.1
	config.DeletingPodsBurst = 10
	config.EnableProfiling = true

	// start controller manager in it's own goroutine
	go controllerManager.Run(config)
}
