package localkube

import (
	scheduler "k8s.io/kubernetes/plugin/cmd/kube-scheduler/app"
	"k8s.io/kubernetes/plugin/cmd/kube-scheduler/app/options"
	"os"
	"time"
)

const (
	SchedulerName = "scheduler"
)

var (
	SchedulerStop chan struct{}
)

func NewSchedulerServer() Server {
	return &SimpleServer{
		ComponentName: SchedulerName,
		StartupFn:     StartSchedulerServer,
		ShutdownFn: func() {
			close(SchedulerStop)
		},
	}
}

func StartSchedulerServer() {
	SchedulerStop = make(chan struct{})
	config := options.NewSchedulerServer()

	// master details
	config.Master = APIServerURL

	// defaults from command
	config.EnableProfiling = true

	schedFn := func() error {
		return scheduler.Run(config)
	}

	go until(schedFn, os.Stdout, SchedulerName, 200*time.Millisecond, SchedulerStop)
}
