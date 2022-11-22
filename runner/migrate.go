package runner

import (
	"github.com/TencentBlueKing/beego-runtime/database/migrations"
)

func runMigrate() {
	migrations.Migrate()
}
