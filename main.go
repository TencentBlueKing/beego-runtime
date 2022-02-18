package main

import (
	_ "github.com/homholueng/beego-runtime/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run("localhost:8000")
}
