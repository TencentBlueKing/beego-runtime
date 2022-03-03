package controllers

import (
	web "github.com/beego/beego/v2/server/web"
	"github.com/homholueng/beego-runtime/conf"
	runtimeInfo "github.com/homholueng/beego-runtime/info"
	"github.com/homholueng/bk-plugin-framework-go/hub"
	frameworkInfo "github.com/homholueng/bk-plugin-framework-go/info"
)

type MetaController struct {
	web.Controller
}

type MetaGetData struct {
	Code             string   `json:"code"`
	Description      string   `json:"description"`
	Versions         []string `json:"versions"`
	Language         string   `json:"language"`
	FrameworkVersion string   `json:"framework_version"`
	RuntimeVersion   string   `json:"runtime_version"`
}

type MetaGetResponse struct {
	*BaseResponse
	Data *MetaGetData `json:"data"`
}

func (c *MetaController) Get() {
	c.Data["json"] = &MetaGetResponse{
		BaseResponse: &BaseResponse{
			Result:  true,
			Message: "success",
		},
		Data: &MetaGetData{
			Code:             conf.PluginName(),
			Description:      "meta desciption",
			Versions:         hub.GetPluginVersions(),
			Language:         "go",
			FrameworkVersion: frameworkInfo.Version(),
			RuntimeVersion:   runtimeInfo.Version(),
		},
	}
	c.ServeJSON()
}
