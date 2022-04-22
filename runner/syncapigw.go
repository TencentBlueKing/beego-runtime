package runner

import (
	"fmt"
	"log"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/manager"
	"github.com/homholueng/beego-runtime/conf"
	"github.com/homholueng/beego-runtime/utils"
)

func runSyncApigw() {
	// load data path
	definitionPath, err := utils.GetApigwDefinitionPath()
	if err != nil {
		log.Fatalf("get apigw definition path error: %v\n", err)
	}
	resourcesPath, err := utils.GetApigwResourcesPath()
	if err != nil {
		log.Fatalf("get apigw resources path error: %v\n", err)
	}

	// create manager
	config := bkapi.ClientConfig{
		Endpoint:  conf.ApigwEndpoint(),
		AppCode:   conf.PluginName(),
		AppSecret: conf.PluginSecret(),
	}

	definitionManager, err := manager.NewManagerFrom(
		conf.ApigwApiName(),
		config,
		definitionPath,
		map[string]interface{}{
			"BK_PLUGIN_APIGW_STAGE_NAME":   conf.Environment(),
			"BK_PLUGIN_APIGW_BACKEND_HOST": conf.ApigwBackendHost(),
		},
	)
	if err != nil {
		log.Fatalf("create apigw definition manager error :%v\n", err)
	}

	resourcesManager, err := manager.NewManagerFrom(
		conf.ApigwApiName(),
		config,
		resourcesPath,
		nil,
	)
	if err != nil {
		log.Fatalf("create apigw resources manager error :%v\n", err)
	}

	// sync start
	syncStageRes, err := definitionManager.SyncStageConfig("stage")
	fmt.Printf("sync apigw stage return: %v\n", syncStageRes)
	if err != nil {
		log.Fatalf("sync apigw stage error :%v\n", err)
	}

	_, err = resourcesManager.SyncResourcesConfig("")
	if err != nil {
		log.Fatalf("sync apigw resource config error :%v\n", err)
	}

	_, err = definitionManager.CreateResourceVersion("resource_version")
	if err != nil {
		log.Fatalf("create resource version error :%v\n", err)
	}

	_, err = definitionManager.Release("release")
	if err != nil {
		log.Fatalf("release stage error :%v\n", err)
	}
}
