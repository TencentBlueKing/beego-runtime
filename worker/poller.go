package worker

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
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

	server.SendTaskWithContext(ctx, task)

	return nil
}
