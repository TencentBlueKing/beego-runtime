package runner

import (
	"flag"
	"log"
	"os"

	"github.com/beego/bee/v2/cmd"
	"github.com/beego/bee/v2/cmd/commands"
	"github.com/beego/bee/v2/config"
	"github.com/beego/bee/v2/utils"
	beego "github.com/beego/beego/v2/server/web"
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
		runBeeCommand(migrateCommand, args)
	case "server":
		beego.Run(":8000")
	default:
		utils.PrintErrorAndExit("Unknown subcommand", cmd.ErrorTemplate)
	}
}
