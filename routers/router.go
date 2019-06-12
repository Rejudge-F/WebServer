package routers

import (
	"WebServer/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/register", &controllers.UserController{})
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/index", &controllers.MainController{}, "get:ShowIndex;post:HandleIndex")
	beego.Router("/addarticle", &controllers.MainController{}, "get:ShowAdd;post:HandleAdd")
	beego.Router("/content", &controllers.MainController{}, "get:ShowContent")
	beego.Router("/update", &controllers.MainController{}, "get:ShowUpdate;post:HandleUpdate")
	beego.Router("/delete", &controllers.MainController{}, "get:HandleDelete")
	beego.Router("/addtype", &controllers.MainController{}, "get:ShowType;post:HandleType")
	beego.Router("/logout", &controllers.UserController{}, "get:HandleLogout")
}
