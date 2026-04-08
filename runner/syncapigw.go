package runner

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/manager"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

func loadDefinitionWithVars(path string, vars map[string]string) (*manager.Definition, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	rendered := os.Expand(string(content), func(key string) string {
		if v, ok := vars[key]; ok {
			return v
		}
		return os.Getenv(key)
	})

	return manager.NewDefinitionFromYaml([]byte(rendered))
}

func runSyncApigw() {
	logger := logrus.New()

	apiGwFilePath := conf.ApigwFilePath()
	definitionPath := fmt.Sprintf("%s/%s", apiGwFilePath, "api-definition.yml")
	resourcesPath := fmt.Sprintf("%s/%s", apiGwFilePath, "api-resources.yml")

	version := fmt.Sprintf("1.0.0+%v", time.Now().Unix())

	templateVars := map[string]string{
		"BK_APIGW_MANAGER_MAINTAINERS":     fmt.Sprintf("%v", conf.ApigwManagerMaintainers()),
		"BK_PLUGIN_APIGW_NAME":             conf.ApigwApiName(),
		"BK_PLUGIN_APIGW_STAGE_NAME":       conf.Environment(),
		"BK_PLUGIN_APIGW_BACKEND_HOST":     conf.ApigwBackendHost(),
		"BK_PLUGIN_APIGW_BACKEND_SCHEME":   conf.ApigwScheme(),
		"BK_PLUGIN_APIGW_BACKEND_SUB_PATH": conf.ApigwSubPath(),
		"BK_PLUGIN_APIGW_RESOURCE_VERSION": version,
		"RESOURCES_FILE_PATH":              resourcesPath,
	}

	definition, err := loadDefinitionWithVars(definitionPath, templateVars)
	if err != nil {
		log.Fatalf("load api definition error: %v\n", err)
	}

	config := bkapi.ClientConfig{
		AppCode:     conf.PluginName(),
		AppSecret:   conf.PluginSecret(),
		AppTenantID: conf.PluginTenantID(),
	}

	mgr, err := manager.NewManager(
		conf.ApigwApiName(),
		config,
		definition,
		nil,
	)
	if err != nil {
		log.Fatalf("create apigw manager error: %v\n", err)
	}

	syncBaseInfoRes, err := mgr.SyncBasicInfo()
	logger.Printf("sync apigw baseinfo return: %v\n", syncBaseInfoRes)
	if err != nil {
		log.Fatalf("sync apigw baseinfo error: %v\n", err)
	}

	syncStagesRes, err := mgr.SyncStagesConfig()
	logger.Printf("sync apigw stages return: %v\n", syncStagesRes)
	if err != nil {
		log.Fatalf("sync apigw stages error: %v\n", err)
	}

	resourcesContent, err := os.ReadFile(resourcesPath)
	if err != nil {
		log.Fatalf("read resources file error: %v\n", err)
	}
	var resourcesData map[string]interface{}
	if err := yaml.Unmarshal(resourcesContent, &resourcesData); err != nil {
		log.Fatalf("parse resources file error: %v\n", err)
	}
	syncResourcesRes, err := mgr.SyncResourcesConfig(resourcesData)
	logger.Printf("sync apigw resources return: %v\n", syncResourcesRes)
	if err != nil {
		log.Fatalf("sync apigw resources error: %v\n", err)
	}

	createResourceRes, err := mgr.CreateResourceVersion(version, "auto sync")
	logger.Printf("create apigw resources version return: %v\n", createResourceRes)
	if err != nil {
		log.Fatalf("create resource version error: %v\n", err)
	}

	releaseRes, err := mgr.Release(version)
	logger.Printf("release stage return: %v\n", releaseRes)
	if err != nil {
		log.Fatalf("release stage error: %v\n", err)
	}
}
