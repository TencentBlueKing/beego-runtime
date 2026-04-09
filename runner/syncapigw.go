package runner

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TencentBlueKing/beego-runtime/conf"
	"github.com/TencentBlueKing/bk-apigateway-sdks/apigateway"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/manager"
	"github.com/flosch/pongo2/v5"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

func loadDefinitionWithVars(path string, vars map[string]string) (*manager.Definition, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	rendered, err := renderDefinition(string(content), vars)
	if err != nil {
		return nil, err
	}

	def, err := newDefinitionFromRendered(rendered, isLegacyDefinitionTemplate(string(content)))
	if err != nil {
		log.Printf("[DEBUG] definition file path: %s\n", path)
		log.Printf("[DEBUG] rendered yaml content:\n%s\n", rendered)
		return nil, err
	}
	return def, nil
}

func renderDefinition(content string, vars map[string]string) (string, error) {
	if isLegacyDefinitionTemplate(content) {
		template, err := pongo2.FromString(content)
		if err != nil {
			return "", fmt.Errorf("failed to parse legacy definition template: %w", err)
		}

		rendered, err := template.Execute(pongo2.Context{
			"data": legacyTemplateData(vars),
		})
		if err != nil {
			return "", fmt.Errorf("failed to render legacy definition template: %w", err)
		}
		return rendered, nil
	}

	return os.Expand(content, func(key string) string {
		if v, ok := vars[key]; ok {
			return v
		}
		return os.Getenv(key)
	}), nil
}

func isLegacyDefinitionTemplate(content string) bool {
	return strings.Contains(content, "{{") || strings.Contains(content, "{%")
}

func legacyTemplateData(vars map[string]string) map[string]interface{} {
	data := make(map[string]interface{}, len(vars))
	for key, value := range vars {
		data[key] = value
	}

	rawMaintainers, ok := vars["BK_APIGW_MANAGER_MAINTAINERS"]
	if !ok || rawMaintainers == "" {
		rawMaintainers = os.Getenv("BK_APIGW_MANAGER_MAINTAINERS")
	}
	if rawMaintainers != "" {
		data["BK_APIGW_MANAGER_MAINTAINERS"] = splitAndTrim(rawMaintainers)
	}

	return data
}

func splitAndTrim(raw string) []string {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func newDefinitionFromRendered(rendered string, legacy bool) (*manager.Definition, error) {
	if !legacy {
		return manager.NewDefinitionFromYaml([]byte(rendered))
	}

	var definition map[string]interface{}
	err := yaml.Unmarshal([]byte(rendered), &definition)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	normalizeLegacyDefinition(definition)
	return manager.NewDefinition(definition), nil
}

func normalizeLegacyDefinition(definition map[string]interface{}) {
	if _, ok := definition["stages"]; ok {
		return
	}

	stage, ok := definition["stage"]
	if !ok {
		return
	}

	definition["stages"] = []interface{}{stage}
}

func runSyncApigw() {
	logger := logrus.New()

	apiGwFilePath := conf.ApigwFilePath()
	definitionPath := fmt.Sprintf("%s/%s", apiGwFilePath, "api-definition.yml")
	resourcesPath := fmt.Sprintf("%s/%s", apiGwFilePath, "api-resources.yml")

	version := fmt.Sprintf("1.0.0+%v", time.Now().Unix())

	templateVars := map[string]string{
		"BK_APIGW_MANAGER_MAINTAINERS":     strings.Join(conf.ApigwManagerMaintainers(), ","),
		"BK_PLUGIN_APIGW_NAME":             conf.ApigwApiName(),
		"BK_PLUGIN_APIGW_STAGE_NAME":       conf.Environment(),
		"BK_PLUGIN_APIGW_BACKEND_HOST":     conf.ApigwBackendHost(),
		"BK_PLUGIN_APIGW_BACKEND_SCHEME":   conf.ApigwScheme(),
		"BK_PLUGIN_APIGW_BACKEND_SUB_PATH": conf.ApigwSubPath(),
		"BK_PLUGIN_APIGW_RESOURCE_VERSION": version,
		"RESOURCES_FILE_PATH":              resourcesPath,
	}

	logger.Printf("[DEBUG] apigw file path: %s\n", apiGwFilePath)
	for k, v := range templateVars {
		logger.Printf("[DEBUG] templateVar %s = %q\n", k, v)
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
		apigateway.New,
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
	syncResourcesRes, err := mgr.SyncResourcesConfig(map[string]interface{}{
		"content": string(resourcesContent),
	})
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
