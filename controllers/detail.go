package controllers

import (
	"fmt"

	"github.com/TencentBlueKing/bk-plugin-framework-go/hub"
	web "github.com/beego/beego/v2/server/web"
)

type DetailController struct {
	web.Controller
}

type DetailGetData struct {
	Version       string                 `json:"version"`
	Desc          string                 `json:"desc"`
	Inputs        map[string]interface{} `json:"inputs"`
	ContextInputs map[string]interface{} `json:"context_inputs"`
	Outputs       map[string]interface{} `json:"outputs"`
	Forms         map[string]interface{} `json:"forms"`
}

type DetailGetResponse struct {
	*BaseResponse
	Data *DetailGetData `json:"data"`
}

func (c *DetailController) Get() {
	version := c.Ctx.Input.Param(":version")
	detail, err := hub.GetPluginDetail(version)
	if err != nil {
		c.Data["json"] = &DetailGetResponse{
			BaseResponse: &BaseResponse{
				Result:  false,
				Message: fmt.Sprintf("get plugin detail fail, %v", err),
			},
			Data: &DetailGetData{},
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = &DetailGetResponse{
		BaseResponse: &BaseResponse{
			Result:  true,
			Message: "success",
		},
		Data: &DetailGetData{
			Version:       version,
			Desc:          detail.Plugin().Desc(),
			Inputs:        detail.InputsSchemaJSON(),
			ContextInputs: detail.ContextInputsSchemaJSON(),
			Outputs:       detail.OutputsSchemaJSON(),
			Forms:         map[string]interface{}{"renderform": nil},
		},
	}
	c.ServeJSON()

}
