package main

import (
	_ "beego-runtime/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run("localhost:8000")
}
