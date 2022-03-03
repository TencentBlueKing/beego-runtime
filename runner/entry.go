package runner

import (
	"flag"
	"log"
	"os"

	"github.com/beego/bee/v2/cmd"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/homholueng/beego-runtime/routers"
)

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

	case "server":
		runServer()
	case "collectstatics":
		runCollectstatics()
	case "worker":
		runWorker()
	default:
		logs.Error("Unknown subcommand: %v", args[0])
		os.Exit(2)
	}
}
