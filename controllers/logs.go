package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type LogsController struct {
	beego.Controller
}

type LogGetData struct {
	Log string `json:"log"`
}

type LogsGetResponse struct {
	*BaseResponse
	Data *LogGetData `json:"data"`
}

func (c *LogsController) Get() {
	c.Data["json"] = &LogsGetResponse{
		BaseResponse: &BaseResponse{
			Result:  true,
			Message: "success",
		},
		Data: &LogGetData{
			Log: "bk-plugin-framework-go does not support log api, check your stdout window please.",
		},
	}
	c.ServeJSON()
}
