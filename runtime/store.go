package runtime

import (
	"encoding/json"
	"time"

	"github.com/homholueng/bk-plugin-framework-go/constants"
)

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

type JSONObjectStore struct {
	Data interface{}
	JSON []byte
}

func (s *JSONObjectStore) Write(traceID string, v interface{}) error {
	s.Data = v
	return nil
}

func (s *JSONObjectStore) Read(traceID string, v interface{}) error {
	return json.Unmarshal(s.JSON, v)
}

type Schedule struct {
	TraceID       string
	PluginVersion string
	State         constants.State
	InvokeCount   int
	Inputs        []byte
	ContextInputs []byte
	ContextStore  []byte
	Outputs       []byte
	CreateAt      time.Time
	Finished      bool
	FinishAt      time.Time
}

type ScheduleStore interface {
	Set(s *Schedule) error
	Get(traceID string) (*Schedule, error)
}
