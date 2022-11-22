package migrations

import (
	"database/sql"
	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/beego/bee/v2/cmd/commands/migrate"
	beeLogger "github.com/beego/bee/v2/logger"
	"github.com/beego/beego/v2/client/orm/migration"
	"os"
	"strings"
)

func checkMigrationsTable() {
	// 调用beego 框架内的方法创建并检查migrations表
	db, err := sql.Open("mysql", conf.MysqlConAddr())
	if err != nil {
		beeLogger.Log.Fatalf("Could not connect to database using '%s': %s", conf.MysqlConAddr(), err)
	}
	defer db.Close()

	showTableSQL := "SHOW TABLES LIKE 'migrations'"
	if rows, err := db.Query(showTableSQL); err != nil {
		beeLogger.Log.Fatalf("Could not show migrations table: %s", err)
	} else if !rows.Next() {
		// No migrations table, create new ones
		createTableSQL := migrate.MYSQLMigrationDDL
		beeLogger.Log.Infof("Creating 'migrations' table...")
		if _, err := db.Query(createTableSQL); err != nil {
			beeLogger.Log.Fatalf("Could not create migrations table: %s", err)
		}
	}

	// Checking that migrations table schema are expected
	selectTableSQL := "DESC migrations"
	if rows, err := db.Query(selectTableSQL); err != nil {
		beeLogger.Log.Fatalf("Could not show columns of migrations table: %s", err)
	} else {
		for rows.Next() {
			var fieldBytes, typeBytes, nullBytes, keyBytes, defaultBytes, extraBytes []byte
			if err := rows.Scan(&fieldBytes, &typeBytes, &nullBytes, &keyBytes, &defaultBytes, &extraBytes); err != nil {
				beeLogger.Log.Fatalf("Could not read column information: %s", err)
			}
			fieldStr, typeStr, nullStr, keyStr, defaultStr, extraStr :=
				string(fieldBytes), string(typeBytes), string(nullBytes), string(keyBytes), string(defaultBytes), string(extraBytes)
			if fieldStr == "id_migration" {
				if keyStr != "PRI" || extraStr != "auto_increment" {
					beeLogger.Log.Hint("Expecting KEY: PRI, EXTRA: auto_increment")
					beeLogger.Log.Fatalf("Column migration.id_migration type mismatch: KEY: %s, EXTRA: %s", keyStr, extraStr)
				}
			} else if fieldStr == "name" {
				if !strings.HasPrefix(typeStr, "varchar") || nullStr != "YES" {
					beeLogger.Log.Hint("Expecting TYPE: varchar, NULL: YES")
					beeLogger.Log.Fatalf("Column migration.name type mismatch: TYPE: %s, NULL: %s", typeStr, nullStr)
				}
			} else if fieldStr == "created_at" {
				if typeStr != "timestamp" || (!strings.EqualFold(defaultStr, "CURRENT_TIMESTAMP") && !strings.EqualFold(defaultStr, "CURRENT_TIMESTAMP()")) {
					beeLogger.Log.Hint("Expecting TYPE: timestamp, DEFAULT: CURRENT_TIMESTAMP || CURRENT_TIMESTAMP()")
					beeLogger.Log.Fatalf("Column migration.timestamp type mismatch: TYPE: %s, DEFAULT: %s", typeStr, defaultStr)
				}
			}
		}
	}

}

func Migrate() {
	checkMigrationsTable()

	task := "upgrade"
	switch task {
	case "upgrade":
		if err := migration.Upgrade(0); err != nil {
			os.Exit(2)
		}
	case "rollback":
		if err := migration.Rollback(""); err != nil {
			os.Exit(2)
		}
	case "reset":
		if err := migration.Reset(); err != nil {
			os.Exit(2)
		}
	case "refresh":
		if err := migration.Refresh(); err != nil {
			os.Exit(2)
		}
	}
}
