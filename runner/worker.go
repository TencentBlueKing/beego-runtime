package runner

import (
	"log"

	"github.com/homholueng/beego-runtime/worker"
)

func runWorker() {
	err := worker.Run()
	if err != nil {
		log.Fatalf("start worker error: %v\n", err)
	}
}
