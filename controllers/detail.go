package controllers

import (
	"fmt"

	web "github.com/beego/beego/v2/server/web"
	"github.com/homholueng/bk-plugin-framework-go/hub"
)

type DetailController struct {
	web.Controller
}

type DetailGetData struct {
	Version       string                 `json:"version"`
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
			Inputs:        detail.InputsSchemaJSON(),
			ContextInputs: detail.ContextInputsSchemaJSON(),
			Outputs:       detail.OutputsSchemaJSON(),
			Forms:         map[string]interface{}{"renderform": nil},
		},
	}
	c.ServeJSON()

}
