package runner

import (
	"github.com/TencentBlueKing/beego-runtime/worker"
	"log"
)

func runWorker() {
	err := worker.MachineryWorkerRun()
	if err != nil {
		log.Fatalf("start worker error: %v\n", err)
	}
}
