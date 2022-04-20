package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["BK_STATIC_URL"] = "/static"
	c.Data["SITE_URL"] = "/bk_plugin"
	c.TplName = "debug_panel.tpl"
}
