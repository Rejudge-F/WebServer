package main

import (
	_ "WebServer/models"
	_ "WebServer/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
