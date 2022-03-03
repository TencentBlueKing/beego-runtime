package runner

import (
	"os"

	"github.com/beego/beego/v2/core/logs"
	"github.com/homholueng/beego-runtime/worker"
)

func runWorker() {
	err := worker.Run()
	if err != nil {
		logs.Error("start worker error: %v", err)
		os.Exit(2)
	}
}
