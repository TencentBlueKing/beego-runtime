package routers

import (
	"github.com/homholueng/beego-runtime/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/bk_plugin/meta", &controllers.MetaController{})
	beego.Router("/bk_plugin/detail/:version", &controllers.DetailController{})
	beego.Router("/bk_plugin/invoke/:version", &controllers.InvokeController{})
	beego.Router("/bk_plugin/schedule/:trace_id", &controllers.ScheduleController{})
}
