package runtime

import (
	"github.com/beego/beego/v2/adapter/orm"
)

var ormClient orm.Ormer

func initOrm() orm.Ormer {
	if ormClient != nil {
		return ormClient
	}
	ormClient = orm.NewOrm()
	return ormClient
}

type MysqlScheduleStore struct {
}

func (rss *MysqlScheduleStore) Set(s *Schedule) error {

	o := initOrm()
	schedule := Schedule{TraceID: s.TraceID}
	err := o.Read(&schedule)
	if err == nil {
		if s.ContextStore == "null" {
			// 说明有数据，更新数据
			s.ContextStore = schedule.ContextStore
		}
		if s.ContextStore == "null" {
			s.Outputs = schedule.Outputs
		}
		_, updateErr := o.Update(s)
		if updateErr != nil {
			return updateErr
		}
	} else {
		_, insertErr := o.Insert(s)
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}

func (rss *MysqlScheduleStore) Get(traceID string) (*Schedule, error) {
	o := initOrm()
	schedule := Schedule{TraceID: traceID}
	err := o.Read(&schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}