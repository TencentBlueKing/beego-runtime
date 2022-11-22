package migrations

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Migration_20221118_101625 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Migration_20221118_101625{}
	m.Created = "20221118_101625"

	migration.Register("Migration_20221118_101625", m)
}

// Run the migrations
func (m *Migration_20221118_101625) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	sql := "CREATE TABLE IF NOT EXISTS `schedule` (" +
		"`trace_i_d` varchar(255) NOT NULL PRIMARY KEY," +
		"`plugin_version` varchar(255) NOT NULL DEFAULT '' ," +
		"`state` tinyint NOT NULL DEFAULT 0 ," +
		"`invoke_count` integer NOT NULL DEFAULT 0 ," +
		"`inputs` varchar(255) NOT NULL DEFAULT '{}' ," +
		"`context_inputs` varchar(255) NOT NULL DEFAULT '{}' ," +
		"`context_store` varchar(255) NOT NULL DEFAULT '{}' ," +
		"`outputs` varchar(255) NOT NULL DEFAULT '{}' ," +
		"`create_at` date NOT NULL," +
		"`finished` bool NOT NULL DEFAULT FALSE ," +
		"`finish_at` datetime) ENGINE=InnoDB;"
	m.SQL(sql)
}

// Reverse the migrations
func (m *Migration_20221118_101625) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
