package routers

import (
	"strings"

	"github.com/homholueng/beego-runtime/conf"
	"github.com/homholueng/beego-runtime/controllers"

	"github.com/beego/beego/v2/server/web/context"

	beego "github.com/beego/beego/v2/server/web"
)

var FilterPluginApi = func(ctx *context.Context) {
	ctx.Request.RequestURI = strings.Replace(ctx.Request.RequestURI, "/bk_plugin/plugin_api", "", 1)
	ctx.Request.URL.Path = strings.Replace(ctx.Request.URL.Path, "/bk_plugin/plugin_api", "", 1)
}

func init() {
	beego.Router("/bk_plugin/meta", &controllers.MetaController{})
	beego.Router("/bk_plugin/detail/:version", &controllers.DetailController{})
	beego.Router("/bk_plugin/invoke/:version", &controllers.InvokeController{})
	beego.Router("/bk_plugin/schedule/:trace_id", &controllers.ScheduleController{})
	beego.InsertFilter("/bk_plugin/plugin_api/*", beego.BeforeRouter, FilterPluginApi)
	if conf.IsDevMode() {
		beego.Router("/", &controllers.DebugPannelController{})
		beego.Router("/bk_plugin/logs/:trace_id", &controllers.LogsController{})
	}
}
