package runner

import (
	"embed"
	"github.com/TencentBlueKing/beego-runtime/conf"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed views
var views embed.FS

//go:embed static
var static embed.FS

//go:embed data
var apiGwData embed.FS

func copyFiles(sourceDir, targetDir string, dirEntries []fs.DirEntry, Fs embed.FS) {
	for _, entry := range dirEntries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		targetPath := filepath.Join(targetDir, entry.Name())

		if entry.IsDir() {
			// 如果是文件夹，则递归拷贝文件夹中的内容
			subDirEntries, err := Fs.ReadDir(sourcePath)
			if err != nil {
				panic(err)
			}
			err = os.MkdirAll(targetPath, 0755)
			if err != nil {
				panic(err)
			}
			copyFiles(sourcePath, targetPath, subDirEntries, Fs)
		} else {
			// 如果是文件，则拷贝文件内容
			fileData, err := Fs.ReadFile(sourcePath)
			if err != nil {
				panic(err)
			}
			err = os.WriteFile(targetPath, fileData, 0644)
			if err != nil {
				panic(err)
			}
		}
	}
}

func syncFile() {
	ViewDirEntries, err := views.ReadDir("views")
	if err != nil {
		panic(err)
	}
	StaticDirEntries, err := static.ReadDir("static")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("./views", 0755)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("./static", 0755)
	if err != nil {
		panic(err)
	}

	copyFiles("views", "./views", ViewDirEntries, views)
	copyFiles("static", "./static", StaticDirEntries, static)

	if !conf.IsDevMode() {

		err = os.MkdirAll("./data", 0755)
		if err != nil {
			panic(err)
		}

		ApiGwDirEntries, err := apiGwData.ReadDir("data")
		if err != nil {
			panic(err)
		}
		copyFiles("data", "./data", ApiGwDirEntries, apiGwData)
	}
}

func runCollectstatics() {
	syncFile()
}
