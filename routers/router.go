package routers

import (
	"github.com/homholueng/beego-runtime/conf"
	"github.com/homholueng/beego-runtime/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

var PluginApiNamespace = beego.NewNamespace("/bk_plugin/plugin_api/")

func init() {
	beego.Router("/bk_plugin/meta", &controllers.MetaController{})
	beego.Router("/bk_plugin/detail/:version", &controllers.DetailController{})
	beego.Router("/bk_plugin/invoke/:version", &controllers.InvokeController{})
	beego.Router("/bk_plugin/schedule/:trace_id", &controllers.ScheduleController{})
	if conf.IsDevMode() {
		beego.Router("/", &controllers.DebugPannelController{})
		beego.Router("/bk_plugin/logs/:trace_id", &controllers.LogsController{})
	}
}
