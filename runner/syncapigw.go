package runner

import (
	"fmt"
	"log"
	"time"

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
		Stage:     conf.Environment(),
		AppCode:   conf.PluginName(),
		AppSecret: conf.PluginSecret(),
	}

	definitionManager, err := manager.NewManagerFrom(
		conf.ApigwApiName(),
		config,
		definitionPath,
		map[string]interface{}{
			"BK_PLUGIN_APIGW_STAGE_NAME":       conf.Environment(),
			"BK_PLUGIN_APIGW_BACKEND_HOST":     conf.ApigwBackendHost(),
			"BK_PLUGIN_APIGW_RESOURCE_VERSION": fmt.Sprintf("1.0.0+%v", time.Now().Unix()),
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

	syncResourcesRes, err := resourcesManager.SyncResourcesConfig("")
	fmt.Printf("sync apigw resources return: %v\n", syncResourcesRes)
	if err != nil {
		log.Fatalf("sync apigw resources error :%v\n", err)
	}

	createResourceRes, err := definitionManager.CreateResourceVersion("resource_version")
	fmt.Printf("create apigw resources version return: %v\n", createResourceRes)
	if err != nil {
		log.Fatalf("create resource version error :%v\n", err)
	}

	releaseRes, err := definitionManager.Release("release")
	fmt.Printf("release stage return: %v\n", releaseRes)
	if err != nil {
		log.Fatalf("release stage error :%v\n", err)
	}
}
