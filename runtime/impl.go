package runtime

import (
	"encoding/json"

	"github.com/homholueng/beego-runtime/models"
	"github.com/homholueng/bk-plugin-framework-go/constants"
	"github.com/homholueng/bk-plugin-framework-go/runtime"
)

type RequestPhaseContextReader struct {
	Inputs        []byte
	ContextInputs []byte
}

func (r *RequestPhaseContextReader) ReadInputs(v interface{}) error {
	return json.Unmarshal(r.Inputs, v)
}

func (r *RequestPhaseContextReader) ReadContextInputs(v interface{}) error {
	return json.Unmarshal(r.ContextInputs, v)
}

type SimpleObjectStore struct {
	Data interface{}
}

func (s *SimpleObjectStore) Write(traceID string, v interface{}) error {
	s.Data = v
	return nil
}

func (s *SimpleObjectStore) Read(traceID string, v interface{}) error {
	v = s.Data
	return nil
}

type ExecuteRuntime struct {
	OutputsStore *SimpleObjectStore
	ContextStore *SimpleObjectStore
}

func (r *ExecuteRuntime) GetOutputsStore() runtime.ObjectStore {
	return r.OutputsStore
}

func (r *ExecuteRuntime) GetContextStore() runtime.ObjectStore {
	return r.ContextStore
}

func (r *ExecuteRuntime) SetPoll(traceID string, version string, invokeCount int, interval int) error {
	data, err := json.Marshal(map[string]interface{}{
		"store":   r.ContextStore.Data,
		"outputs": r.OutputsStore.Data,
	})
	if err != nil {
		return err
	}

	schedule := models.Schedule{
		TraceID:       traceID,
		PluginVersion: version,
		State:         int(constants.StatePoll),
		InvokeCount:   invokeCount,
		Data:          string(data),
	}
	_, err = models.Ormer().Insert(&schedule)
	return err
}
