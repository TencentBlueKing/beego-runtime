package runner

import (
	"log"
	"os/exec"

	"github.com/TencentBlueKing/beego-runtime/utils"
)

func runCollectstatics() {
	staticDir, err := utils.GetStaticDirPath()
	if err != nil {
		log.Fatalf("get static files dir failed: %v\n", err)
	}
	viewPath, err := utils.GetViewPath()
	if err != nil {
		log.Fatalf("get view path failed: %v\n", err)
	}
	definitionPath, err := utils.GetApigwDefinitionPath()
	if err != nil {
		log.Fatalf("get apigw definition path failed: %v\n", err)
	}
	resourcesPath, err := utils.GetApigwResourcesPath()
	if err != nil {
		log.Fatalf("get apigw resources path failed: %v\n", err)
	}

	cmds := []*exec.Cmd{
		exec.Command("cp", "-r", staticDir, "."),
		exec.Command("cp", "-r", viewPath, "."),
		exec.Command("cp", "-r", definitionPath, "."),
		exec.Command("cp", "-r", resourcesPath, "."),
	}

	for _, c := range cmds {
		log.Printf("run collect static command: %v\n", c)
		if err := c.Run(); err != nil {
			log.Fatalf("collect static failed: %v\n", err)
		}
	}
}
