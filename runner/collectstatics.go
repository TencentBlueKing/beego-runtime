package runner

import (
	"os"
	"os/exec"

	"github.com/beego/beego/v2/core/logs"
	runtimeUtils "github.com/homholueng/beego-runtime/utils"
)

func runCollectstatics() {
	staticDir, err := runtimeUtils.GetStaticDirPath()
	if err != nil {
		logs.Error("get static files dir failed: %v", err)
		os.Exit(2)
	}
	viewPath, err := runtimeUtils.GetViewPath()
	if err != nil {
		logs.Error("get view path failed: %v", err)
		os.Exit(2)
	}

	cpStaticCmd := exec.Command("cp", "-r", staticDir, ".")
	cpViewCmd := exec.Command("cp", "-r", viewPath, ".")

	logs.Info("run collect static command: ", cpStaticCmd)
	if err := cpStaticCmd.Run(); err != nil {
		logs.Error("collect static failed: %v", err)
	}

	logs.Info("run collect view command: ", cpViewCmd)
	if err := cpViewCmd.Run(); err != nil {
		logs.Error("collect view failed: %v", err)
	}
}
