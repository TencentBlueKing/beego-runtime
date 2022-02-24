package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	web "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"github.com/homholueng/beego-runtime/runtime"
	"github.com/homholueng/bk-plugin-framework-go/constants"
	"github.com/homholueng/bk-plugin-framework-go/executor"
)

type InvokeController struct {
	web.Controller
}

type InvokePostParam struct {
	Inputs  json.RawMessage `json:"inputs"`
	Context json.RawMessage `json:"context"`
}

type InvokePostData struct {
	Outputs interface{}     `json:"outputs"`
	State   constants.State `json:"state"`
	Err     string          `json:"err"`
}

type InvokePostResponse struct {
	*BaseResponse
	TraceID string          `json:"trace_id"`
	Data    *InvokePostData `json:"data"`
}

func (c *InvokeController) Post() {
	traceID := strings.Replace(uuid.NewString(), "-", "", -1)
	version := c.Ctx.Input.Param(":version")

	var param InvokePostParam
	if err := c.BindJSON(&param); err != nil {
		c.Data["json"] = &InvokePostResponse{
			BaseResponse: &BaseResponse{
				Result:  false,
				Message: fmt.Sprintf("request param bind fail, %v", err),
			},
			TraceID: traceID,
			Data:    &InvokePostData{},
		}
		c.ServeJSON()
	}

	reader := runtime.RequestPhaseContextReader{Inputs: param.Inputs, ContextInputs: param.Context}
	contextStore := runtime.SimpleObjectStore{}
	outputStore := runtime.SimpleObjectStore{}
	state, err := executor.Execute(
		traceID,
		version,
		&reader,
		&runtime.ExecuteRuntime{
			OutputsStore: &outputStore,
			ContextStore: &contextStore,
		},
	)
	if err != nil {
		c.Data["json"] = &InvokePostResponse{
			BaseResponse: &BaseResponse{
				Result:  false,
				Message: fmt.Sprintf("plugin execute fail, %v", err),
			},
			TraceID: traceID,
			Data: &InvokePostData{
				Outputs: nil,
				State:   state,
				Err:     err.Error(),
			},
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = &InvokePostResponse{
		BaseResponse: &BaseResponse{
			Result:  true,
			Message: "success",
		},
		TraceID: traceID,
		Data: &InvokePostData{
			Outputs: outputStore.Data,
			State:   state,
			Err:     "",
		},
	}
	c.ServeJSON()
}
