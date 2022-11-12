package runner

import (
	"fmt"
	"log"

	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/TencentBlueKing/beego-runtime/routers"
	runtimeUtils "github.com/TencentBlueKing/beego-runtime/utils"
	"github.com/beego/beego/v2/server/web"
)

func runServer() {
	var staticDir string
	var viewPath string
	var err error
	if conf.IsDevMode() {
		staticDir, err = runtimeUtils.GetStaticDirPath()
		if err != nil {
			log.Fatalf("get static files dir failed: %v\n", err)
		}
	} else {
		staticDir = "static"
	}
	log.Printf("serve /static at %v\n", staticDir)

	if conf.IsDevMode() {
		viewPath, err = runtimeUtils.GetViewPath()
		if err != nil {
			log.Fatalf("get view path failed: %v\n", err)
		}
	} else {
		viewPath = "views"
	}
	log.Printf("serve views at %v\n", viewPath)

	web.BConfig.CopyRequestBody = true
	web.BConfig.WebConfig.ViewsPath = viewPath
	web.SetStaticPath("/static", staticDir)
	web.AddNamespace(routers.PluginApiNamespace)
	web.Run(fmt.Sprintf(":%v", conf.Port()))
}
