package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/TencentBlueKing/beego-runtime/runtime"
	"github.com/TencentBlueKing/beego-runtime/worker"
	"github.com/TencentBlueKing/bk-plugin-framework-go/constants"
	"github.com/TencentBlueKing/bk-plugin-framework-go/executor"
	web "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
	traceLogger := log.WithField("trace_id", traceID)

	// 只有非DEV环境才会去进行网关认证，方便本地调试
	if !conf.IsDevMode() {
		_, err := parseApigwJWT(c.Ctx.Request)
		if err != nil {
			c.Data["json"] = &InvokePostResponse{
				BaseResponse: &BaseResponse{
					Result:  false,
					Message: fmt.Sprintf("please make sure request is from apigw, %v", err),
				},
				TraceID: traceID,
				Data:    &InvokePostData{},
			}
			c.ServeJSON()
			return
		}
	}

	version := c.Ctx.Input.Param(":version")

	var param InvokePostParam
	if err := c.BindJSON(&param); err != nil {
		traceLogger.Errorf("param bind error: %v\n", err)
		c.Data["json"] = &InvokePostResponse{
			BaseResponse: &BaseResponse{
				Result:  false,
				Message: fmt.Sprintf("request param bind fail, %v", err),
			},
			TraceID: traceID,
			Data:    &InvokePostData{},
		}
		c.ServeJSON()
		return
	}

	reader := runtime.JSONContextReader{Inputs: param.Inputs, ContextInputs: param.Context}

	contextStore := runtime.SimpleObjectStore{}
	outputStore := runtime.SimpleObjectStore{}
	state, err := executor.Execute(
		traceID,
		version,
		&reader,
		&runtime.ExecuteRuntime{
			Inputs:        param.Inputs,
			ContextInputs: param.Context,
			OutputsStore:  &outputStore,
			ContextStore:  &contextStore,
			ScheduleStore: runtime.GetScheduleStore(),
			Poller:        &worker.MachineryPoller{},
		},
		traceLogger,
	)

	if err != nil {
		traceLogger.Errorf("param bind error: %v\n", err)
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
