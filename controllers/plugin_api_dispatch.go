package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	web "github.com/beego/beego/v2/server/web"
	log "github.com/sirupsen/logrus"
)

type PluginApiDispatchController struct {
	PluginApiController
}

type PluginApiDispatchBase struct {
	Url      string `json:"url" form:"url"`
	Method   string `json:"method" form:"method"`
	Username string `json:"username" form:"username"`
}

type PluginApiDispatchParam struct {
	PluginApiDispatchBase
	Data json.RawMessage         `json:"data" form:"data"`
	File []*multipart.FileHeader `form:"file"`
}

type PluginApiDispatchResponse struct {
	*BaseResponse
	Data interface{} `json:"data"`
}

func handleErrResponse(c *PluginApiDispatchController, err error, msg string) {
	c.Data["json"] = &PluginApiDispatchResponse{
		BaseResponse: &BaseResponse{
			Result:  false,
			Message: fmt.Sprintf("%s, %v", msg, err),
		},
		Data: nil,
	}
	c.ServeJSON()
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
	var errMsg string

	contentType := c.Ctx.Request.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// 处理multipart/form-data请求
		if err := c.BindForm(&param); err != nil {
			errMsg = "request param bind fail"
			log.Errorf("%s %v\n", errMsg, err)
			handleErrResponse(c, err, errMsg)
			return
		}

		// 文件无法直接绑定 需要手动解析
		if err := c.Ctx.Request.ParseMultipartForm(32 << 20); err != nil {
			errMsg = "failed to parse multipart form"
			log.Errorf("%s %v\n", errMsg, err)
			handleErrResponse(c, err, errMsg)
			return
		}

		fileHeaders := c.Ctx.Request.MultipartForm.File["file"]
		if len(fileHeaders) == 0 {
			log.Warn("no files uploaded")
		}

		for _, fileHeader := range fileHeaders {
			param.File = append(param.File, fileHeader)
		}

	} else {
		if err := c.BindJSON(&param); err != nil {
			errMsg = "param bind error"
			log.Errorf("%s %v\n", errMsg, err)
			handleErrResponse(c, err, errMsg)
			return
		}
	}
	parsedURL, err := url.Parse(param.Url)
	if err != nil {
		errMsg = "param.Url parse fail"
		log.Errorf("%s, %v\n", errMsg, err)
		handleErrResponse(c, err, errMsg)
		return
	}

	path := parsedURL.Path
	upperMethod := strings.ToUpper(param.Method)
	_, ok := c.FindController(path, upperMethod)
	if !ok {
		errMsg = fmt.Sprintf("controller not found, path: %s, method: %s", path, upperMethod)
		log.Errorf(errMsg)
		handleErrResponse(c, nil, errMsg)
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

		if strings.HasPrefix(contentType, "multipart/form-data") {
			var buffer bytes.Buffer
			writer := multipart.NewWriter(&buffer)

			if len(param.Data) > 0 {
				_ = writer.WriteField("data", string(param.Data))
			}

			for _, fileHeader := range param.File {
				file, err := fileHeader.Open()
				if err != nil {
					errMsg = "open file fail"
					log.Errorf("%s, %v", errMsg, err)
					handleErrResponse(c, err, errMsg)
					return
				}
				part, err := writer.CreateFormFile("file", fileHeader.Filename)
				if err != nil {
					errMsg = "create form file fail"
					log.Errorf("%s, %v", errMsg, err)
					handleErrResponse(c, err, errMsg)
					return
				}

				if _, err := io.Copy(part, file); err != nil {
					errMsg = "copy file fail"
					log.Errorf("%s, %v", errMsg, err)
					handleErrResponse(c, err, errMsg)
					return
				}

				if err := writer.Close(); err != nil {
					errMsg = "close form file fail"
					log.Errorf("%s, %v", errMsg, err)
					handleErrResponse(c, err, errMsg)
					return
				}
				newRequest.Body = io.NopCloser(&buffer)
			}
		} else {
			newRequest.Body = io.NopCloser(bytes.NewReader(param.Data))
			newRequest.ContentLength = int64(len(param.Data))
		}

	} else {
		errMsg = fmt.Sprintf("dispatch method not supported, method: %s\n", upperMethod)
		log.Errorf(errMsg)
		handleErrResponse(c, err, errMsg)
		return
	}
	web.BeeApp.Handlers.ServeHTTP(c.Ctx.ResponseWriter, newRequest)
}
