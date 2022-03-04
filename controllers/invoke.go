package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	web "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"github.com/homholueng/beego-runtime/conf"
	"github.com/homholueng/beego-runtime/runtime"
	"github.com/homholueng/beego-runtime/worker"
	"github.com/homholueng/bk-plugin-framework-go/constants"
	"github.com/homholueng/bk-plugin-framework-go/executor"
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
	version := c.Ctx.Input.Param(":version")
	traceLogger := log.WithField("trace_id", traceID)

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
			ScheduleStore: &runtime.RedisScheduleStore{
				Client:             conf.RedisClient(),
				Expiration:         conf.ScheduleExpiration(),
				FinishedExpiration: conf.FinishedScheduleExpiration(),
			},
			Poller: &worker.AsynqPoller{
				Client: conf.AsynqClient(),
			},
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
