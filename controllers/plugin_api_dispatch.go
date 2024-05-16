package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	web "github.com/beego/beego/v2/server/web"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type PluginApiDispatchController struct {
	PluginApiController
}

type PluginApiDispatchParam struct {
	Url      string          `json:"url"`
	Method   string          `json:"method"`
	Username string          `json:"username"`
	Data     json.RawMessage `json:"data"`
}

type PluginApiDispatchResponse struct {
	*BaseResponse
	Data interface{} `json:"data"`
}

func (c *PluginApiDispatchController) FindController(path string, method string) (string, bool) {
	// method is GET or POST
	upperMethod := strings.ToUpper(method)
	methods := web.PrintTree()["Data"].(web.M)
	path = strings.TrimRight(path, "/")
	for m, v := range methods {
		upperM := strings.ToUpper(m)
		if upperMethod != upperM {
			continue
		}
		for _, vv := range *v.(*[][]string) {
			p, controllerType := vv[0], vv[2]
			if strings.TrimRight(p, "/") != path {
				continue
			}
			return controllerType, true
		}
	}
	return "", false
}

func (c *PluginApiDispatchController) Post() {
	var param PluginApiDispatchParam
	if err := c.BindJSON(&param); err != nil {
		log.Errorf("param bind error: %v\n", err)
		c.Data["json"] = &PluginApiDispatchResponse{
			BaseResponse: &BaseResponse{
				Result:  false,
				Message: fmt.Sprintf("request param bind fail, %v", err),
			},
			Data: nil,
		}
		c.ServeJSON()
		return
	}

	parsedURL, err := url.Parse(param.Url)
	if err != nil {
		log.Errorf("param.Url parse fail, %v\n", err)
		c.Data["json"] = &PluginApiDispatchResponse{
			BaseResponse: &BaseResponse{
				Result:  false,
				Message: fmt.Sprintf("Url parse fail, %v", err),
			},
			Data: nil,
		}
		c.ServeJSON()
		return
	}

	path := parsedURL.Path
	upperMethod := strings.ToUpper(param.Method)
	_, ok := c.FindController(path, upperMethod)
	if !ok {
		log.Errorf("controller not found, path: %s, method: %s", path, upperMethod)
		c.Data["json"] = &PluginApiDispatchResponse{
			BaseResponse: &BaseResponse{
				Result:  false,
				Message: fmt.Sprintf("controller not found, path: %s, method: %s", path, upperMethod),
			},
			Data: nil,
		}
		c.ServeJSON()
		return
	}

	newRequest := new(http.Request)
	*newRequest = *c.Ctx.Request
	newRequest.URL = &url.URL{
		Scheme: c.Ctx.Request.URL.Scheme,
		Host:   c.Ctx.Request.URL.Host,
		Path:   path,
	}
	newRequest.Header = make(http.Header, len(c.Ctx.Request.Header))
	for key, values := range c.Ctx.Request.Header {
		newRequest.Header[key] = append([]string(nil), values...)
	}

	newRequest.Method = upperMethod
	if upperMethod == http.MethodGet {
		newRequest.URL.RawQuery = parsedURL.RawQuery
	} else if upperMethod == http.MethodPost {
		newRequest.Header.Set("Content-Type", "application/json")
		newRequest.Body = io.NopCloser(bytes.NewReader(param.Data))
		newRequest.ContentLength = int64(len(param.Data))
	} else {
		log.Errorf("dispatch method not supported, method: %s\n", upperMethod)
		c.Data["json"] = &PluginApiDispatchResponse{
			BaseResponse: &BaseResponse{
				Result:  false,
				Message: fmt.Sprintf("dispatch method not supported, method: %s\n", upperMethod),
			},
			Data: nil,
		}
		c.ServeJSON()
		return
	}

	web.BeeApp.Handlers.ServeHTTP(c.Ctx.ResponseWriter, newRequest)
}
