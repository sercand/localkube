package localkube

import (
	scheduler "k8s.io/kubernetes/plugin/cmd/kube-scheduler/app"
	"k8s.io/kubernetes/plugin/cmd/kube-scheduler/app/options"
)

const (
	SchedulerName = "scheduler"
)

func NewSchedulerServer() Server {
	return SimpleServer{
		ComponentName: SchedulerName,
		StartupFn:     StartSchedulerServer,
	}.NoShutdown()
}

func StartSchedulerServer() {
	config := options.NewSchedulerServer()

	// master details
	config.Master = APIServerURL

	// defaults from command
	config.EnableProfiling = true

	go scheduler.Run(config)
}
