package main

import (
	"math/rand"
	"os"
	"time"

	"k8s.io/component-base/logs"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"

	"github.com/SataQiu/k8s-scheduler-example/pkg/plugins/nodememoryusagelimit"

	// ensure scheme package is initialized.
	_ "github.com/SataQiu/k8s-scheduler-example/api/scheme"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logs.InitLogs()
	defer logs.FlushLogs()

	command := app.NewSchedulerCommand(
		app.WithPlugin(nodememoryusagelimit.Name, nodememoryusagelimit.New),
	)

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
