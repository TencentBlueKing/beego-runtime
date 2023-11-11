package runner

import (
	"fmt"
	"log"
	"time"

	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/manager"
	"github.com/sirupsen/logrus"
)

func runSyncApigw() {
	logger := logrus.New()
	// load data path

	apiGwFilePath := conf.ApigwFilePath()

	definitionPath := fmt.Sprintf("%s/%s", apiGwFilePath, "api-definition.yml")
	resourcesPath := fmt.Sprintf("%s/%s", apiGwFilePath, "api-resources.yml")

	// create manager
	config := bkapi.ClientConfig{
		AppCode:   conf.PluginName(),
		AppSecret: conf.PluginSecret(),
	}

	manager, err := manager.NewManagerFrom(
		conf.ApigwApiName(),
		config,
		definitionPath,
		map[string]interface{}{
			"BK_APIGW_MANAGER_MAINTAINERS":     conf.ApigwManagerMaintainers(),
			"BK_PLUGIN_APIGW_NAME":             conf.ApigwApiName(),
			"BK_PLUGIN_APIGW_STAGE_NAME":       conf.Environment(),
			"BK_PLUGIN_APIGW_BACKEND_HOST":     conf.ApigwBackendHost(),
			"BK_PLUGIN_APIGW_RESOURCE_VERSION": fmt.Sprintf("1.0.0+%v", time.Now().Unix()),
			"RESOURCES_FILE_PATH":              resourcesPath,
		},
	)
	if err != nil {
		log.Fatalf("create apigw  manager error :%v\n", err)
	}

	// create gateway
	syncBaseInfoRes, err := manager.SyncBasicInfo("apigateway")
	logger.Printf("sync apigw baseinfo return: %v\n", syncBaseInfoRes)
	if err != nil {
		log.Fatalf("sync apigw baseinfo error :%v\n", err)
	}

	// sync start
	syncStageRes, err := manager.SyncStageConfig("stage")
	logger.Printf("sync apigw stage return: %v\n", syncStageRes)
	if err != nil {
		log.Fatalf("sync apigw stage error :%v\n", err)
	}

	syncResourcesRes, err := manager.SyncResourcesConfig("resources")
	logger.Printf("sync apigw resources return: %v\n", syncResourcesRes)
	if err != nil {
		log.Fatalf("sync apigw resources error :%v\n", err)
	}

	createResourceRes, err := manager.CreateResourceVersion("resource_version")
	logger.Printf("create apigw resources version return: %v\n", createResourceRes)
	if err != nil {
		log.Fatalf("create resource version error :%v\n", err)
	}

	releaseRes, err := manager.Release("release")
	logger.Printf("release stage return: %v\n", releaseRes)
	if err != nil {
		log.Fatalf("release stage error :%v\n", err)
	}
}
