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

type PluginApiDispatchParam struct {
	Url        string                  `json:"url" form:"url"`
	Method     string                  `json:"method" form:"method"`
	Username   string                  `json:"username" form:"username"`
	Data       json.RawMessage         `json:"data" form:"data"`
	Files      []*multipart.FileHeader `json:"-" form:"-"`
	DumpedData string                  `json:"dumped_data" form:"dumped_data"`
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

func parseMixedRequest(c *PluginApiDispatchController) (*PluginApiDispatchParam, bool, error) {
	param := &PluginApiDispatchParam{}
	bodyBytes, _ := io.ReadAll(c.Ctx.Request.Body)
	c.Ctx.Request.Body.Close()
	c.Ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 根据内容特征判断真实类型
	isMultipart := c.Ctx.Request.MultipartForm != nil && len(c.Ctx.Request.MultipartForm.File) > 0

	log.Infof("dispatch requests: %v", c.Ctx.Request)

	if !isMultipart {
		// 普通情况直接按 json 解析
		err := c.BindJSON(param)
		return param, isMultipart, err
	} else {

		err := c.Ctx.Request.ParseMultipartForm(32 << 20) // 32 MB
		if err != nil {
			// 处理解析错误
			return nil, isMultipart, err
		}

		// 绑定表单字段
		if err := c.BindForm(param); err != nil {
			return nil, isMultipart, err
		}

		// 获取上传文件
		for _, headers := range c.Ctx.Request.MultipartForm.File {
			param.Files = append(param.Files, headers...)
		}

		return param, isMultipart, nil
	}
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
	var param *PluginApiDispatchParam
	var errMsg string

	// recover panic to ensure error logs are written and a proper response is returned
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("plugin api dispatch panic recovered: %v", r)
			handleErrResponse(c, fmt.Errorf("%v", r), "internal server error")
		}
	}()

	param, isMultipart, err := parseMixedRequest(c)
	if err != nil {
		errMsg = "parse request error"
		log.Errorf("%s, %v", errMsg, err)
		handleErrResponse(c, err, errMsg)
		return
	}
	log.Infof("dispatch params: %v", param)

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

		if isMultipart {
			var buffer bytes.Buffer
			writer := multipart.NewWriter(&buffer)
			defer writer.Close()

			if len(param.DumpedData) > 0 {
				// 适配标准运维的转发逻辑 form 格式下 使用 dumped_data 解析为 string 且不支持嵌套kv
				err = writer.WriteField("dumped_data", string(param.DumpedData))
				if err != nil {
					errMsg = "write request dumped_data failed"
					log.Errorf("%s, %v", errMsg, err)
					handleErrResponse(c, err, errMsg)
					return
				}
			}

			for _, fileHeader := range param.Files {
				file, err := fileHeader.Open()
				defer file.Close()
				if err != nil {
					errMsg = "open file fail"
					log.Errorf("%s, %v", errMsg, err)
					handleErrResponse(c, err, errMsg)
					return
				}
				part, err := writer.CreateFormFile(fileHeader.Filename, fileHeader.Filename)
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
