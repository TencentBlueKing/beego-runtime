package controllers

import (
	"net/http"
	"time"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/manager"
	"github.com/homholueng/beego-runtime/conf"
	"github.com/pkg/errors"
)

type BaseResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

func parseApigwJWT(r *http.Request) (manager.ApigatewayJwtClaims, error) {
	jwt, ok := r.Header["HTTP_X_BKAPI_JWT"]
	if !ok {
		return manager.ApigatewayJwtClaims{}, errors.Errorf("can not find HTTP_X_BKAPI_JWT header in request")
	}
	config := bkapi.ClientConfig{
		Endpoint:  conf.ApigwEndpoint(),
		AppCode:   conf.PluginName(),
		AppSecret: conf.PluginSecret(),
	}
	cache := manager.NewPublicKeyMemoryCache(config, 1*time.Hour, func(apiName string, config bkapi.ClientConfig) (*manager.Manager, error) {
		return manager.NewDefaultManager(apiName, config)
	})
	jwtParser := manager.NewRsaJwtTokenParser(cache)
	return jwtParser.Parse(jwt[0])
}
