package worker

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"time"
)

type MachineryPoller struct {
}

func (p *MachineryPoller) Poll(traceID string, after time.Duration) error {
	task := NewPollTask(traceID, after)

	server, err := StartServer()
	if err != nil {
		return err
	}

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()

	server.SendTaskWithContext(ctx, task)

	return nil
}
