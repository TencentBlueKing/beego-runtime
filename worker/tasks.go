package worker

import (
	"github.com/RichardKnop/machinery/v2/tasks"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/TencentBlueKing/beego-runtime/runtime"
	"github.com/TencentBlueKing/bk-plugin-framework-go/executor"
)

func NewPollTask(traceID string, after time.Duration) *tasks.Signature {
	eta := time.Now().UTC().Add(time.Second * after)

	Task := tasks.Signature{
		Name: "HandlePollTask",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: traceID,
			},
		},
		ETA: &eta,
	}

	return &Task
}

func HandlePollTask(traceID string) error {
	traceLogger := log.WithField("trace_id", traceID)
	traceLogger.Info("prepare schedule")

	rss := runtime.GetScheduleStore()

	schedule, err := rss.Get(traceID)
	if err != nil {
		traceLogger.Errorf("schedule get error: %v", err)
		return err
	}

	reader := runtime.JSONContextReader{
		Inputs:        []byte(schedule.Inputs),
		ContextInputs: []byte(schedule.ContextInputs),
	}

	runtime := runtime.NewScheduleExecuteRuntime(schedule, rss, &MachineryPoller{})

	err = executor.Schedule(
		traceID,
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
