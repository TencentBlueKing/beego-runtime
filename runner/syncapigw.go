package runner

import (
	"log"

	manager "github.com/TencentBlueKing/bk-apigateway-sdks/apigw-manager"
	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
)

func runSyncApigw() {
	manager, err := manager.NewManagerFrom(
		"my-api",
		bkapi.ClientConfig{
			Endpoint:  "http://bkapi.example.com",
			AppCode:   "my-app-code",
			AppSecret: "my-app-secret",
		},
		"/path/to/definition.yaml",
		map[string]interface{}{
			"key": "value",
		},
	)
	if err != nil {
		log.Fatalf("create apigw manager error :%v\n", err)
	}

	_, err = manager.SyncStageConfig("stage") // 同步环境信息
	if err != nil {
		log.Fatalf("sync apigw stage error :%v\n", err)
	}
	_, err = manager.SyncResourcesConfig("resources")
	if err != nil {
		log.Fatalf("sync apigw resource config error :%v\n", err)
	} // 同步资源配置
	_, err = manager.CreateResourceVersion("resource_version") // 创建资源版本
	if err != nil {
		log.Fatalf("create resource version error :%v\n", err)
	}
	_, err = manager.Release("release") // 发布资源
	if err != nil {
		log.Fatalf("release stage error :%v\n", err)
	}
}
