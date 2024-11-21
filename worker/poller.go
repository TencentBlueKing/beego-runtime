package worker

import (
	"context"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type MachineryPoller struct {
}

func (p *MachineryPoller) Poll(traceID string, after time.Duration) error {
	task := NewPollTask(traceID, after)

	server, err := NewServer()
	if err != nil {
		return err
	}

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()

	_, err = server.SendTaskWithContext(ctx, task)

	if err != nil {
		log.Errorf("SendTaskWithContext error %s", err)
	}

	return nil
}
