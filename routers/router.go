package routers

import (
	"github.com/homholueng/beego-runtime/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
