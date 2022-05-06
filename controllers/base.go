package controllers

import (
	"net/http"
	"time"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/manager"
	"github.com/homholueng/beego-runtime/conf"
)

type BaseResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

func parseApigwJWT(r *http.Request) (manager.ApigatewayJwtClaims, error) {
	jwt := r.Header["HTTP_X_BKAPI_JWT"][0]
	config := bkapi.ClientConfig{
		Endpoint:  conf.ApigwEndpoint(),
		AppCode:   conf.PluginName(),
		AppSecret: conf.PluginSecret(),
	}
	cache := manager.NewPublicKeyMemoryCache(config, 1*time.Hour, func(apiName string, config bkapi.ClientConfig) (*manager.Manager, error) {
		return manager.NewDefaultManager(apiName, config)
	})
	jwtParser := manager.NewRsaJwtTokenParser(cache)
	return jwtParser.Parse(jwt)
}
