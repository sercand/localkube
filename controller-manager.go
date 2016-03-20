package localkube

import (
	"os"
	"time"

	controllerManager "k8s.io/kubernetes/cmd/kube-controller-manager/app"
	"k8s.io/kubernetes/cmd/kube-controller-manager/app/options"
)

const (
	ControllerManagerName = "controller-manager"
)

var (
	CMStop chan struct{}
)

func NewControllerManagerServer() Server {
	return &SimpleServer{
		ComponentName: ControllerManagerName,
		StartupFn:     StartControllerManagerServer,
		ShutdownFn: func() {
			close(CMStop)
		},
	}
}

func StartControllerManagerServer() {
	config := options.NewCMServer()

	// defaults from command
	config.DeletingPodsQps = 0.1
	config.DeletingPodsBurst = 10
	config.EnableProfiling = true

	fn := func() error {
		return controllerManager.Run(config)
	}

	// start controller manager in it's own goroutine
	go until(fn, os.Stdout, ControllerManagerName, 200*time.Millisecond, SchedulerStop)
}
