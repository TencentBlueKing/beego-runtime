package runner

import (
	"fmt"
	"os"

	"github.com/beego/bee/v2/cmd/commands"
	"github.com/beego/bee/v2/config"
	"github.com/beego/bee/v2/utils"
	"github.com/beego/beego/v2/core/logs"
	"github.com/homholueng/beego-runtime/conf"
	_ "github.com/homholueng/beego-runtime/models"
	_ "github.com/homholueng/beego-runtime/routers"
	runtimeUtils "github.com/homholueng/beego-runtime/utils"
)

var migrateCommand *commands.Command

func init() {
	for _, c := range commands.AvailableCommands {
		if c.Name() == "migrate" {
			migrateCommand = c
		}
	}
	if migrateCommand == nil {
		utils.PrintErrorAndExit("can not load bee migrate command", "")
	}
}

func runBeeCommand(c *commands.Command, args []string) {
	if c.Run == nil {
		return
	}
	c.Flag.Usage = func() { c.Usage() }
	if c.CustomFlags {
		args = args[1:]
	} else {
		c.Flag.Parse(args[1:])
		args = c.Flag.Args()
	}

	if c.PreRun != nil {
		c.PreRun(c, args)
	}

	config.LoadConfig()
	os.Exit(c.Run(c, args))
	return
}

func runMigrate() {
	migDir, err := runtimeUtils.GetMigrationDirPath()
	if err != nil {
		logs.Error("get migration files dir failed: %v\n", err)
		os.Exit(2)
	}
	logs.Info("use migration dir: %v", migDir)

	migrateArgs := []string{
		"migrate",
		"-driver=mysql",
		fmt.Sprintf(
			"-conn=%v:%v@tcp(%v:%v)/%v",
			conf.DataBase.User,
			conf.DataBase.Password,
			conf.DataBase.Host,
			conf.DataBase.Port,
			conf.DataBase.DBName,
		),
		fmt.Sprintf("-dir=%v", migDir),
	}
	runBeeCommand(migrateCommand, migrateArgs)
}
