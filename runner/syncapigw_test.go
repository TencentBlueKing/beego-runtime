package runner

import (
	"os"
	"path/filepath"
	"testing"
)

const legacyDefinitionTemplate = `apigateway:
  description: {{ data.BK_PLUGIN_APIGW_NAME }} apigw
  is_public: true
  maintainers:
    {% for member in data.BK_APIGW_MANAGER_MAINTAINERS %}
    - {{ member }}
    {% endfor %}

stage:
  name: {{ data.BK_PLUGIN_APIGW_STAGE_NAME }}
  vars:
    api_sub_path: "{{ data.BK_PLUGIN_APIGW_BACKEND_SUB_PATH }}"
  proxy_http:
    timeout: 60
    upstreams:
      loadbalance: "roundrobin"
      hosts:
        - host: {{ data.BK_PLUGIN_APIGW_BACKEND_SCHEME }}://{{ data.BK_PLUGIN_APIGW_BACKEND_HOST }}/
          weight: 100

resource_version:
  title: "auto-release-version"
  version: {{ data.BK_PLUGIN_APIGW_RESOURCE_VERSION }}

resources:
  include_file: {{ data.RESOURCES_FILE_PATH }}

release:
  version: {{ data.BK_PLUGIN_APIGW_RESOURCE_VERSION }}
  resource_version_name: "auto-release-version"
  stage_names:
    - {{ data.BK_PLUGIN_APIGW_STAGE_NAME }}
  comment: "auto release by bk-plugin-runtime"
`

const modernDefinitionTemplate = `apigateway:
  description: ${BK_PLUGIN_APIGW_NAME} apigw
  is_public: true
  maintainers:
    - admin

stages:
  - name: ${BK_PLUGIN_APIGW_STAGE_NAME}
    vars:
      api_sub_path: "${BK_PLUGIN_APIGW_BACKEND_SUB_PATH}"
    proxy_http:
      timeout: 60
      upstreams:
        loadbalance: "roundrobin"
        hosts:
          - host: ${BK_PLUGIN_APIGW_BACKEND_SCHEME}://${BK_PLUGIN_APIGW_BACKEND_HOST}/
            weight: 100
`

func TestLoadDefinitionWithVars_LegacyTemplate(t *testing.T) {
	t.Parallel()

	path := writeTempDefinition(t, legacyDefinitionTemplate)
	definition, err := loadDefinitionWithVars(path, map[string]string{
		"BK_PLUGIN_APIGW_NAME":             "bp-plugin-test01",
		"BK_PLUGIN_APIGW_STAGE_NAME":       "prod",
		"BK_PLUGIN_APIGW_BACKEND_HOST":     "apps.sg.bk2game.com",
		"BK_PLUGIN_APIGW_BACKEND_SCHEME":   "http",
		"BK_PLUGIN_APIGW_BACKEND_SUB_PATH": "prod--default--plugin-test01/",
		"BK_PLUGIN_APIGW_RESOURCE_VERSION": "1.0.0+123",
		"RESOURCES_FILE_PATH":              "./data/api-resources.yml",
		"BK_APIGW_MANAGER_MAINTAINERS":     "admin,operator",
	})
	if err != nil {
		t.Fatalf("loadDefinitionWithVars returned error: %v", err)
	}

	gateway, err := definition.Get("apigateway")
	if err != nil {
		t.Fatalf("definition.Get(apigateway) returned error: %v", err)
	}
	maintainers, ok := gateway["maintainers"].([]interface{})
	if !ok {
		t.Fatalf("maintainers has unexpected type: %T", gateway["maintainers"])
	}
	if len(maintainers) != 2 {
		t.Fatalf("expected 2 maintainers, got %d", len(maintainers))
	}

	stages, err := definition.GetArray("stages")
	if err != nil {
		t.Fatalf("definition.GetArray(stages) returned error: %v", err)
	}
	if len(stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(stages))
	}
	if stages[0]["name"] != "prod" {
		t.Fatalf("expected normalized stage name to be prod, got %v", stages[0]["name"])
	}
}

func TestLoadDefinitionWithVars_ModernTemplate(t *testing.T) {
	t.Parallel()

	path := writeTempDefinition(t, modernDefinitionTemplate)
	definition, err := loadDefinitionWithVars(path, map[string]string{
		"BK_PLUGIN_APIGW_NAME":             "bp-plugin-test01",
		"BK_PLUGIN_APIGW_STAGE_NAME":       "prod",
		"BK_PLUGIN_APIGW_BACKEND_HOST":     "apps.sg.bk2game.com",
		"BK_PLUGIN_APIGW_BACKEND_SCHEME":   "http",
		"BK_PLUGIN_APIGW_BACKEND_SUB_PATH": "prod--default--plugin-test01/",
	})
	if err != nil {
		t.Fatalf("loadDefinitionWithVars returned error: %v", err)
	}

	stages, err := definition.GetArray("stages")
	if err != nil {
		t.Fatalf("definition.GetArray(stages) returned error: %v", err)
	}
	if len(stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(stages))
	}
	if stages[0]["name"] != "prod" {
		t.Fatalf("expected stage name to be prod, got %v", stages[0]["name"])
	}
}

func writeTempDefinition(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "api-definition.yml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp definition: %v", err)
	}
	return path
}
