// TencentBlueKing is pleased to support the open source community by making
// 蓝鲸智云-gopkg available.
// Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

// Package kit collect the basic tool for developer to
// develop a bk-plugin.
package controllers

import (
	"net/http"
	"strings"

	"github.com/TencentBlueKing/beego-runtime/conf"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"

	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type BKUser struct {
	Username string
	Token    string
}

type PluginApiController struct {
	beego.Controller
	User BKUser
}

type PluginApiBaseResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

func (p *PluginApiController) setUser(username string, token string) {
	p.User = BKUser{
		Username: username,
		Token:    token,
	}
}

func (p *PluginApiController) Prepare() {
	pluginApiTraceID := strings.Replace(uuid.NewString(), "-", "", -1)
	pluginApiLogger := log.WithField("plugin_api_trace_id", pluginApiTraceID)
	pluginApiLogger.Infof("request plugin api %v", p.Ctx.Request.URL.Path)

	token := ""
	username := ""
	if conf.IsDevMode() {
		bkUid, err := p.Ctx.Request.Cookie("bk_uid")
		if err == http.ErrNoCookie {
			username = conf.PluginApiDebugUsername()
			pluginApiLogger.Infof("Get bk_uid as username from cookie fail,use env PLUGIN_API_DEBUG_USERNAME[%v] as username", username)
		} else {
			username = bkUid.Value
		}

		tokenKey := conf.UserTokenKeyName()
		bkToken, err := p.Ctx.Request.Cookie(tokenKey)
		if err != nil {
			token = ""
			pluginApiLogger.Infof("Get %v from cookie fail,plase check env UserTokenKeyName", tokenKey)
		} else {
			token = bkToken.Value
		}
	} else {
		bkToken, ok := p.Ctx.Request.Header["X-Bkapi-Jwt"]
		if !ok {
			pluginApiLogger.Errorf("This API can only be accessed through API gateway")
			p.Data["json"] = &PluginApiBaseResponse{
				Result:  false,
				Message: "This API can only be accessed through API gateway",
			}
			p.ServeJSON()
			return
		} else {
			token = bkToken[0]
		}

		claims, err := parseApigwJWT(p.Ctx.Request)
		if err != nil {
			pluginApiLogger.Errorf("This API can only be accessed through API gateway")
			p.Data["json"] = &PluginApiBaseResponse{
				Result:  false,
				Message: "This API can only be accessed through API gateway",
			}
			p.ServeJSON()
			return
		} else {
			username = claims.User.Username
		}
	}
	p.setUser(username, token)
}

func (p *PluginApiController) GetBkapiAuthorizationInfo() string {
	authInfo := map[string]string{
		"bk_app_code":           conf.PluginName(),
		"bk_app_secret":         conf.PluginSecret(),
		conf.UserTokenKeyName(): p.User.Token,
	}

	if !conf.IsDevMode() {
		// TO DO: Get App access_token
		authInfo["access_token"] = "access_token"
	}
	b, _ := json.Marshal(authInfo)
	return string(b)
}
