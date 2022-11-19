package worker

import (
	"fmt"
	"github.com/RichardKnop/machinery/v2/example/tracers"
	"github.com/RichardKnop/machinery/v2/tasks"
	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strings"
)

func MachineryWorkerRun() error {
	// worker 需要生成随机的id用于多worker标示身份
	workerTag := strings.Replace(uuid.NewString(), "-", "", -1)[:5]

	consumerTag := fmt.Sprintf("worker-%s", workerTag)

	cleanup, err := tracers.SetupTracer(consumerTag)
	if err != nil {
		return err
	}
	defer cleanup()

	server, err := StartServer()
	if err != nil {
		return err
	}

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker(consumerTag, conf.WorkerConcurrency())

	errorHandler := func(err error) {
		log.Fatalf("it has some err of task, error:%s", err)
	}

	preTaskHandler := func(signature *tasks.Signature) {
		log.Infof("Received task , name: %s", signature.Name)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.Infof("the task is end : %s", signature.Name)
	}

	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetErrorHandler(errorHandler)
	worker.SetPreTaskHandler(preTaskHandler)

	return worker.Launch()

}
