package runtime

import (
	"github.com/beego/beego/v2/adapter/orm"
)

var ormClient orm.Ormer

func GetOrmClient() orm.Ormer {
	if ormClient != nil {
		return ormClient
	}
	ormClient = orm.NewOrm()
	return ormClient
}

type MysqlScheduleStore struct {
}

func (rss *MysqlScheduleStore) Set(s *Schedule) error {

	o := GetOrmClient()
	schedule := Schedule{TraceID: s.TraceID}
	err := o.Read(&schedule)

	if err != nil {
		_, err := o.Insert(s)
		if err != nil {
			return err
		}
	} else {
		if s.ContextStore == "null" {
			// 说明有数据，更新数据
			s.ContextStore = schedule.ContextStore
		}
		if s.Outputs == "null" {
			s.Outputs = schedule.Outputs
		}
		_, err := o.Update(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rss *MysqlScheduleStore) Get(traceID string) (*Schedule, error) {
	o := GetOrmClient()
	schedule := Schedule{TraceID: traceID}
	err := o.Read(&schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}
