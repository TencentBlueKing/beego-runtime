package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/TencentBlueKing/beego-runtime/controllers"
)

var PluginApiNamespace = beego.NewNamespace("/bk_plugin/plugin_api/")

func init() {
	beego.Router("/bk_plugin/meta", &controllers.MetaController{})
	beego.Router("/bk_plugin/detail/:version", &controllers.DetailController{})
	beego.Router("/bk_plugin/invoke/:version", &controllers.InvokeController{})
	beego.Router("/bk_plugin/schedule/:trace_id", &controllers.ScheduleController{})
	beego.Router("/bk_plugin/plugin_api_dispatch", &controllers.PluginApiDispatchController{})
	if conf.IsDevMode() {
		beego.Router("/", &controllers.DebugPannelController{})
		beego.Router("/bk_plugin/logs/:trace_id", &controllers.LogsController{})
	}
}
