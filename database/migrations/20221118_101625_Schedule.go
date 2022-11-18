package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Schedule_20221118_101625 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Schedule_20221118_101625{}
	m.Created = "20221118_101625"

	migration.Register("Schedule_20221118_101625", m)
}

// Run the migrations
func (m *Schedule_20221118_101625) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	sql := "\tCREATE TABLE IF NOT EXISTS `schedule` (\n\t\t`trace_i_d` varchar(255) NOT NULL PRIMARY KEY,\n\t\t`plugin_version` varchar(255) NOT NULL DEFAULT '' ,\n\t\t`state` tinyint NOT NULL DEFAULT 0 ,\n\t\t`invoke_count` integer NOT NULL DEFAULT 0 ,\n\t\t`inputs` varchar(255) NOT NULL DEFAULT '{}' ,\n\t\t`context_inputs` varchar(255) NOT NULL DEFAULT '{}' ,\n\t\t`context_store` varchar(255) NOT NULL DEFAULT '{}' ,\n\t\t`outputs` varchar(255) NOT NULL DEFAULT '{}' ,\n\t\t`create_at` date NOT NULL,\n\t\t`finished` bool NOT NULL DEFAULT FALSE ,\n\t\t`finish_at` datetime\n\t) ENGINE=InnoDB;"
	m.SQL(sql)
}

// Reverse the migrations
func (m *Schedule_20221118_101625) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
