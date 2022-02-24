package conf

import (
	"fmt"
	"os"

	"github.com/beego/beego/v2/core/config"
)

const configData string = `
pluginname = ${BKPAAS_APP_ID}
runmode = ${BKPAAS_ENVIRONMENT||dev}

dbname = ${GCS_MYSQL_NAME}
dbuser = ${GCS_MYSQL_USER}
dbpasswd = ${GCS_MYSQL_PASSWORD}
dbhost = ${GCS_MYSQL_HOST}
dbport = ${GCS_MYSQL_PORT}

dev_dbname = ${BKPAAS_APP_ID}
dev_dbuser = ${BK_PLUGIN_RUNTIME_DB_USER}
dev_dbpasswd = ${BK_PLUGIN_RUNTIME_DB_PWD}
dev_dbhost = ${BK_PLUGIN_RUNTIME_DB_HOST||127.0.0.1}
dev_dbport = ${BK_PLUGIN_RUNTIME_DB_PORT||3306}
`

var Settings config.Configer
var DataBase DataBaseSetting

const ()

type DataBaseSetting struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func onDevEnv() bool {
	return Settings.DefaultString("runmode", "dev") == "dev"
}

func initDataBase() {
	configs := make(map[string]string, 5)
	for _, key := range []string{"dbname", "dbuser", "dbpasswd", "dbhost", "dbport"} {
		var val string
		var err error
		if onDevEnv() {
			val, err = Settings.String(fmt.Sprintf("dev_%v", key))
		} else {
			val, err = Settings.String(key)
		}
		if err != nil {
			fmt.Printf("settings %v read error: %v\n", key, err)
			os.Exit(2)
		}
		configs[key] = val
	}
	DataBase = DataBaseSetting{
		User:     configs["dbuser"],
		Password: configs["dbpasswd"],
		Host:     configs["dbhost"],
		Port:     configs["dbport"],
		DBName:   configs["dbname"],
	}
}

func init() {
	var err error
	Settings, err = config.NewConfigData("ini", []byte(configData))
	if err != nil {
		fmt.Printf("runtime config load error: %v\n", err)
		os.Exit(2)
	}

	// init database info
	initDataBase()
}
