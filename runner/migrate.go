package runner

import (
	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/TencentBlueKing/beego-runtime/info"
	"github.com/TencentBlueKing/beego-runtime/utils"
	"github.com/beego/bee/v2/cmd/commands/migrate"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

func runMigrate() {
	baseDir, err := utils.GetModulePath("github.com/TencentBlueKing/beego-runtime", info.Version())
	if err != nil {
		log.Fatalf("get baseDir path failed: %v\n", err)
	}
	migrate.MigrateUpdate(baseDir, "mysql", conf.MysqlConAddr(), filepath.Join(baseDir, "database/migrations"))
}
