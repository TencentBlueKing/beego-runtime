package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/TencentBlueKing/beego-runtime/runtime"
	"github.com/beego/beego/v2/server/web"
)

const timeFormat = "2006-01-02 15:04:05"

type ScheduleController struct {
	web.Controller
}

type ScheduleGetData struct {
	TraceID       string                 `json:"trace_id"`
	PluginVersion string                 `json:"plugin_version"`
	State         int                    `json:"state"`
	Outputs       map[string]interface{} `json:"outputs"`
	CreateAt      string                 `json:"create_at"`
	FinishAt      string                 `json:"finish_at"`
}

type ScheduleGetResponse struct {
	*BaseResponse
	Data *ScheduleGetData `json:"data"`
}

func (c *ScheduleController) Get() {
	traceID := c.Ctx.Input.Param(":trace_id")
	store := runtime.RedisScheduleStore{
		Client:             conf.RedisClient(),
		Expiration:         conf.ScheduleExpiration(),
		FinishedExpiration: conf.FinishedScheduleExpiration(),
	}
	schedule, err := store.Get(traceID)
	if err != nil {
		c.Data["json"] = &ScheduleGetResponse{
			BaseResponse: &BaseResponse{
				Result:  false,
				Message: fmt.Sprintf("get schedule fail, %v", err),
			},
			Data: nil,
		}
		c.ServeJSON()
		return
	}

	var outputs map[string]interface{}
	if err := json.Unmarshal(schedule.Outputs, &outputs); err != nil {
		c.Data["json"] = &ScheduleGetResponse{
			BaseResponse: &BaseResponse{
				Result:  false,
				Message: fmt.Sprintf("outputs data unmarshal fail, %v", err),
			},
			Data: nil,
		}
		c.ServeJSON()
		return
	}

	finishAt := ""
	if schedule.Finished {
		finishAt = schedule.FinishAt.Format(timeFormat)
	}

	c.Data["json"] = &ScheduleGetResponse{
		BaseResponse: &BaseResponse{
			Result:  true,
			Message: "success",
		},
		Data: &ScheduleGetData{
			TraceID:       schedule.TraceID,
			PluginVersion: schedule.PluginVersion,
			State:         int(schedule.State),
			Outputs:       outputs,
			CreateAt:      schedule.CreateAt.Format(timeFormat),
			FinishAt:      finishAt,
		},
	}
	c.ServeJSON()
}
