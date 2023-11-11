package runner

import (
	"fmt"
	"log"

	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/TencentBlueKing/beego-runtime/routers"
	"github.com/beego/beego/v2/server/web"
)

func runServer() {
	var staticDir string
	var viewPath string
	staticDir = "static"
	log.Printf("serve /static at %v\n", staticDir)
	viewPath = "views"
	log.Printf("serve views at %v\n", viewPath)

	web.BConfig.CopyRequestBody = true
	web.BConfig.WebConfig.ViewsPath = viewPath
	web.SetStaticPath("/static", staticDir)
	web.AddNamespace(routers.PluginApiNamespace)
	web.Run(fmt.Sprintf(":%v", conf.Port()))
}
