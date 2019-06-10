package main

import (
	_ "WebServer/models"
	_ "WebServer/routers"
	"github.com/astaxie/beego"
	"strconv"
)

func main() {
	err := beego.AddFuncMap("ShowPrePage", ShowPrePage)
	if err != nil {
		beego.Info(err)
		return
	}
	err = beego.AddFuncMap("ShowNextPage", ShowNextPage)
	if err != nil {
		beego.Info(err)
		return
	}
	beego.Run()
}

func ShowPrePage(data int) string {
	dataTmp := data
	pre := dataTmp - 1
	if pre == 0 {
		pre = 1
	}
	return strconv.Itoa(pre)
}

func ShowNextPage(data int) string {
	data += 1
	return strconv.Itoa(data)
}
