package worker

import (
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type AsynqPoller struct {
	Client *asynq.Client
}

func (p *AsynqPoller) Poll(traceID string, after time.Duration) error {
	task, err := NewPollTask(traceID)
	if err != nil {
		return err
	}
	info, err := p.Client.Enqueue(task, asynq.ProcessIn(after))
	log.Printf("poll %s info: %v", traceID, info)
	return err
}
