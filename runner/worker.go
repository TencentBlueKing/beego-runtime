package runner

import (
	"log"

	"github.com/TencentBlueKing/beego-runtime/worker"
)

func runWorker() {
	err := worker.MachineryWorkerRun()
	if err != nil {
		log.Fatalf("start worker error: %v\n", err)
	}
}
