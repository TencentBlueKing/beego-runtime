package worker

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/TencentBlueKing/beego-runtime/runtime"
	"github.com/TencentBlueKing/bk-plugin-framework-go/executor"
	"github.com/hibiken/asynq"
)

const TypePoll = "poll"

type PollPayload struct {
	TraceID string
}

func NewPollTask(traceID string) (*asynq.Task, error) {
	payload, err := json.Marshal(PollPayload{TraceID: traceID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePoll, payload), nil
}

func HandlePollTask(ctx context.Context, t *asynq.Task) error {
	var p PollPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	traceLogger := log.WithField("trace_id", p.TraceID)
	traceLogger.Info("prepare schedule")

	rss := runtime.MysqlScheduleStore{}

	schedule, err := rss.Get(p.TraceID)
	if err != nil {
		traceLogger.Errorf("schedule get error: %v\n", err)
		return err
	}

	reader := runtime.JSONContextReader{
		Inputs:        []byte(schedule.Inputs),
		ContextInputs: []byte(schedule.ContextInputs),
	}

	runtime := runtime.NewScheduleExecuteRuntime(schedule, &rss, &AsynqPoller{Client: conf.AsynqClient()})

	err = executor.Schedule(
		p.TraceID,
		schedule.PluginVersion,
		schedule.InvokeCount+1,
		&reader,
		runtime,
		traceLogger,
	)
	if err != nil {
		log.Error("schedule execute error: %v", err)
	}
	return err
}
