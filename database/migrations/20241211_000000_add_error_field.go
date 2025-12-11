package migrations

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// AddErrorField_20241211_000000 添加 error 字段用于存储插件执行错误信息
type AddErrorField_20241211_000000 struct {
	migration.Migration
}

func init() {
	m := &AddErrorField_20241211_000000{}
	m.Created = "20241211_000000"

	migration.Register("AddErrorField_20241211_000000", m)
}

// Up 执行迁移：添加 error 字段
func (m *AddErrorField_20241211_000000) Up() {
	// 添加 error 字段，用于存储插件执行失败时的错误信息
	// 供 SOPS 通过 err 字段展示给用户
	sql := "ALTER TABLE schedule ADD COLUMN error TEXT NULL AFTER outputs;"
	m.SQL(sql)
}

// Down 回滚迁移：删除 error 字段
func (m *AddErrorField_20241211_000000) Down() {
	sql := "ALTER TABLE schedule DROP COLUMN error;"
	m.SQL(sql)
}
