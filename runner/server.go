package runner

import (
	"fmt"
	"os"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/homholueng/beego-runtime/conf"
	runtimeUtils "github.com/homholueng/beego-runtime/utils"
)

func runServer() {
	var staticDir string
	var viewPath string
	var err error
	if conf.IsDevMode() {
		staticDir, err = runtimeUtils.GetStaticDirPath()
		if err != nil {
			logs.Error("get static files dir failed: %v", err)
			os.Exit(2)
		}
	} else {
		staticDir = "static"
	}
	logs.Info("serve /static at %v", staticDir)

	if conf.IsDevMode() {
		viewPath, err = runtimeUtils.GetViewPath()
		if err != nil {
			logs.Error("get view path failed: %v", err)
			os.Exit(2)
		}
	} else {
		viewPath = "views"
	}
	logs.Info("serve views at %v", viewPath)

	web.BConfig.CopyRequestBody = true
	web.BConfig.WebConfig.ViewsPath = viewPath
	web.SetStaticPath("/static", staticDir)
	web.Run(fmt.Sprintf(":%v", conf.Port()))
}
