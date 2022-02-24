package main

import (
	"fmt"

	"github.com/homholueng/beego-runtime/runner"
	"github.com/homholueng/beego-runtime/utils"
)

func main() {
	fmt.Println(utils.GetMigrationDirPath())
	return
	runner.Run()
}
