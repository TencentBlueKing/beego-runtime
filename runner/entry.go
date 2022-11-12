package runner

import (
	"flag"
	"log"
	"os"

	_ "github.com/TencentBlueKing/beego-runtime/routers"
	"github.com/beego/bee/v2/cmd"
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
	case "syncapigw":
		runSyncApigw()
	default:
		log.Fatalf("Unknown subcommand: %v\n", args[0])
	}
}
