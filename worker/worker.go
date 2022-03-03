package worker

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/homholueng/beego-runtime/conf"
)

func Run() error {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: conf.RedisAddr()},
		asynq.Config{Concurrency: conf.WorkerConcurrency()},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypePoll, HandlePollTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
		return err
	}
	return nil
}
