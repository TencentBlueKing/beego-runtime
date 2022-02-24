package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type ScheduleTable_20220223_150313 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ScheduleTable_20220223_150313{}
	m.Created = "20220223_150313"

	migration.Register("ScheduleTable_20220223_150313", m)
}

// Run the migrations
func (m *ScheduleTable_20220223_150313) Up() {
	m.SQL("CREATE TABLE `schedule` (" +
		"`trace_id` VARCHAR(33) NOT NULL," +
		"`plugin_version` VARCHAR(128) NOT NULL," +
		"`state` TINYINT NOT NULL," +
		"`invoke_count` INT NOT NULL," +
		"`create_at` DATETIME NOT NULL," +
		"`finish_at` DATETIME," +
		"`data` LONGTEXT NOT NULL," +
		"PRIMARY KEY (`trace_id`)" +
		")")
}

// Reverse the migrations
func (m *ScheduleTable_20220223_150313) Down() {
	m.SQL("DROP TABLE `schedule`")
}
