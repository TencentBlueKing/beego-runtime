package runner

import (
	"log"
	"os/exec"

	runtimeUtils "github.com/homholueng/beego-runtime/utils"
)

func runCollectstatics() {
	staticDir, err := runtimeUtils.GetStaticDirPath()
	if err != nil {
		log.Fatalf("get static files dir failed: %v\n", err)
	}
	viewPath, err := runtimeUtils.GetViewPath()
	if err != nil {
		log.Fatalf("get view path failed: %v\n", err)
	}

	cpStaticCmd := exec.Command("cp", "-r", staticDir, ".")
	cpViewCmd := exec.Command("cp", "-r", viewPath, ".")

	log.Printf("run collect static command: %v\n", cpStaticCmd)
	if err := cpStaticCmd.Run(); err != nil {
		log.Fatalf("collect static failed: %v\n", err)
	}

	log.Printf("run collect view command: %v\n", cpViewCmd)
	if err := cpViewCmd.Run(); err != nil {
		log.Fatalf("collect view failed: %v\n", err)
	}
}
