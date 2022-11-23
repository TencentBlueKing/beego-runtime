package runtime

import (
	"encoding/json"
	"time"

	"github.com/TencentBlueKing/bk-plugin-framework-go/constants"
	"github.com/beego/beego/v2/client/orm"

	"github.com/TencentBlueKing/beego-runtime/conf"
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
	TraceID       string `orm:"PK"`
	PluginVersion string
	State         constants.State
	InvokeCount   int
	Inputs        orm.JSONField
	ContextInputs orm.JSONField
	ContextStore  orm.JSONField
	Outputs       orm.JSONField
	CreateAt      time.Time `orm:"auto_now_add;type(date)"`
	Finished      bool
	FinishAt      time.Time `orm:"null"`
}
type ScheduleStore interface {
	Set(s *Schedule) error
	Get(traceID string) (*Schedule, error)
}

func GetScheduleStore() ScheduleStore {
	if conf.ScheduleStoreMode() == "redis" {
		return &RedisScheduleStore{
			Client:             conf.RedisClient(),
			Expiration:         conf.ScheduleExpiration(),
			FinishedExpiration: conf.FinishedScheduleExpiration()}
	} else {
		return &MysqlScheduleStore{}
	}
}

func init() {
	orm.RegisterModel(new(Schedule))
}
