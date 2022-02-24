package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

var ormer orm.Ormer

func Ormer() orm.Ormer {
	if ormer == nil {
		ormer = orm.NewOrm()
	}
	return ormer
}

type Schedule struct {
	TraceID       string    `orm:"column(trace_id);pk;size(33)"`
	PluginVersion string    `orm:"column(plugin_version);size(128)"`
	State         int       `orm:"column(state)"`
	InvokeCount   int       `orm:"column(invoke_count);"`
	CreateAt      time.Time `orm:"column(create_at);auto_now_add;type(datetime)"`
	FinishAt      time.Time `orm:"column(finish_at);type(datetime);null"`
	Data          string    `orm:"column(data);type(text)"`
}

func init() {
	orm.RegisterModel(new(Schedule))
}
