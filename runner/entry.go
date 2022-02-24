package runner

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/beego/bee/v2/cmd"
	"github.com/beego/bee/v2/cmd/commands"
	"github.com/beego/bee/v2/config"
	"github.com/beego/bee/v2/utils"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
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

func Run() {
	flag.Usage = cmd.Usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()

	if len(args) < 1 {
		cmd.Usage()
		os.Exit(2)
		return
	}

	if args[0] == "help" {
		cmd.Help(args[1:])
		return
	}

	switch args[0] {

	case "migrate":
		migDir, err := runtimeUtils.GetMigrationDirPath()
		if err != nil {
			fmt.Printf("get migration files dir failed: %v\n", err)
			os.Exit(2)
		}

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

	case "server":
		orm.RegisterDataBase(
			"default",
			"mysql",
			fmt.Sprintf(
				"%v:%v@tcp(%v:%v)/%v",
				conf.DataBase.User,
				conf.DataBase.Password,
				conf.DataBase.Host,
				conf.DataBase.Port,
				conf.DataBase.DBName,
			),
		)
		beego.BConfig.CopyRequestBody = true
		beego.SetStaticPath("/static", "static")
		beego.Run(":5000")

	default:
		fmt.Printf("Unknown subcommand: %v\n", args[0])
		os.Exit(2)
	}
}
