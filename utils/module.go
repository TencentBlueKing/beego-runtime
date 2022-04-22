package utils

import (
	"os"
	"path"

	"github.com/homholueng/beego-runtime/info"
	"golang.org/x/mod/module"
)

func GetApigwDefinitionPath() (string, error) {
	baseDir, err := GetModulePath("github.com/homholueng/beego-runtime", info.Version())
	if err != nil {
		return "", err
	}

	return path.Join(baseDir, "data/api-definition.yml"), nil
}

func GetApigwResourcesPath() (string, error) {
	baseDir, err := GetModulePath("github.com/homholueng/beego-runtime", info.Version())
	if err != nil {
		return "", err
	}

	return path.Join(baseDir, "data/api-resources.yml"), nil
}

func GetStaticDirPath() (string, error) {
	baseDir, err := GetModulePath("github.com/homholueng/beego-runtime", info.Version())
	if err != nil {
		return "", err
	}

	return path.Join(baseDir, "static"), nil
}

func GetViewPath() (string, error) {
	baseDir, err := GetModulePath("github.com/homholueng/beego-runtime", info.Version())
	if err != nil {
		return "", err
	}

	return path.Join(baseDir, "views"), nil
}

func GetModulePath(name, version string) (string, error) {
	// first we need GOMODCACHE
	cache, ok := os.LookupEnv("GOMODCACHE")
	if !ok {
		cache = path.Join(os.Getenv("GOPATH"), "pkg", "mod")
	}

	// then we need to escape path
	escapedPath, err := module.EscapePath(name)
	if err != nil {
		return "", err
	}

	// version also
	escapedVersion, err := module.EscapeVersion(version)
	if err != nil {
		return "", err
	}

	return path.Join(cache, escapedPath+"@"+escapedVersion), nil
}
