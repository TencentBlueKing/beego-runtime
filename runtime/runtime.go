package runtime

import (
	"encoding/json"
	"time"

	"github.com/TencentBlueKing/bk-plugin-framework-go/constants"
	"github.com/TencentBlueKing/bk-plugin-framework-go/runtime"
	"github.com/beego/beego/v2/client/orm"
)

type Poller interface {
	Poll(traceID string, after time.Duration) error
}

type ExecuteRuntime struct {
	Inputs        []byte
	ContextInputs []byte
	OutputsStore  *SimpleObjectStore
	ContextStore  *SimpleObjectStore
	ScheduleStore ScheduleStore
	Poller        Poller
}

func (r *ExecuteRuntime) GetOutputsStore() runtime.ObjectStore {
	return r.OutputsStore
}

func (r *ExecuteRuntime) GetContextStore() runtime.ObjectStore {
	return r.ContextStore
}

func (r *ExecuteRuntime) SetPoll(traceID string, version string, invokeCount int, after time.Duration) error {

	storeData, err := json.Marshal(r.ContextStore.Data)
	if err != nil {
		return err
	}

	outputsData, err := json.Marshal(r.OutputsStore.Data)
	if err != nil {
		return err
	}

	if err = r.ScheduleStore.Set(&Schedule{
		TraceID:       traceID,
		PluginVersion: version,
		State:         constants.StatePoll,
		InvokeCount:   invokeCount,
		Inputs:        orm.JSONField(r.Inputs),
		ContextInputs: orm.JSONField(r.ContextInputs),
		ContextStore:  orm.JSONField(storeData),
		Outputs:       orm.JSONField(outputsData),
		CreateAt:      time.Now(),
		Finished:      false,
	}); err != nil {
		return err
	}

	if err = r.Poller.Poll(traceID, after); err != nil {
		return err
	}

	return nil
}

type ScheduleExecuteRuntime struct {
	Schedule      *Schedule
	ContextStore  *JSONObjectStore
	OutputsStore  *JSONObjectStore
	ScheduleStore ScheduleStore
	Poller        Poller
}

func (r *ScheduleExecuteRuntime) updateSchedule(state constants.State, invokeCount int, finished bool) error {

	storeData, err := json.Marshal(r.ContextStore.Data)

	if err != nil {
		return err
	}

	outputsData, err := json.Marshal(r.OutputsStore.Data)
	if err != nil {
		return err
	}

	r.Schedule.ContextStore = orm.JSONField(storeData)
	r.Schedule.Outputs = orm.JSONField(outputsData)
	r.Schedule.State = state
	r.Schedule.InvokeCount = invokeCount
	if finished {
		r.Schedule.Finished = true
		r.Schedule.FinishAt = time.Now()
	}

	return nil
}

func NewScheduleExecuteRuntime(schedule *Schedule, scheduleStore ScheduleStore, poller Poller) *ScheduleExecuteRuntime {
	return &ScheduleExecuteRuntime{
		Schedule:      schedule,
		ContextStore:  &JSONObjectStore{JSON: []byte(schedule.ContextStore)},
		OutputsStore:  &JSONObjectStore{JSON: []byte(schedule.Outputs)},
		ScheduleStore: scheduleStore,
		Poller:        poller,
	}
}

func (r *ScheduleExecuteRuntime) GetOutputsStore() runtime.ObjectStore {
	return r.OutputsStore
}

func (r *ScheduleExecuteRuntime) GetContextStore() runtime.ObjectStore {
	return r.ContextStore
}

func (r *ScheduleExecuteRuntime) SetPoll(traceID string, version string, invokeCount int, after time.Duration) error {

	// update schedule
	if err := r.updateSchedule(constants.StatePoll, invokeCount, false); err != nil {
		return err
	}

	if err := r.ScheduleStore.Set(r.Schedule); err != nil {
		return err
	}

	if err := r.Poller.Poll(traceID, after); err != nil {
		return err
	}

	return nil
}

func (r *ScheduleExecuteRuntime) SetFail(traceID string, err error) error {
	if err := r.updateSchedule(constants.StateFail, r.Schedule.InvokeCount, true); err != nil {
		return err
	}
	return r.ScheduleStore.Set(r.Schedule)
}

func (r *ScheduleExecuteRuntime) SetSuccess(traceID string) error {
	if err := r.updateSchedule(constants.StateSuccess, r.Schedule.InvokeCount, true); err != nil {
		return err
	}
	return r.ScheduleStore.Set(r.Schedule)
}
