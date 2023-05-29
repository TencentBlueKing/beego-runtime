package migrations

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Schedule_20230529_155339 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Schedule_20230529_155339{}
	m.Created = "20230529_155339"

	migration.Register("Schedule_20230529_155339", m)
}

// Run the migrations
func (m *Schedule_20230529_155339) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	sql := "ALTER TABLE schedule MODIFY inputs LONGTEXT, MODIFY context_inputs LONGTEXT, MODIFY context_store LONGTEXT, MODIFY outputs LONGTEXT; "
	m.SQL(sql)

}

// Reverse the migrations
func (m *Schedule_20230529_155339) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
