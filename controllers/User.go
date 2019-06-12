package controllers

import (
	"WebServer/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) ShowLogin() {
	userName := c.Ctx.GetCookie("userName")
	if userName != "" {
		c.Data["check"] = "checked"
		c.Data["userName"] = userName
	} else {
		c.Data["userName"] = ""
	}

	c.TplName = "login.html"
}

func (c *UserController) HandleLogin() {
	userName := c.GetString("username")
	userPwd := c.GetString("password")
	if userName == "" || userPwd == "" {
		c.Redirect("/", 302)
	}
	o := orm.NewOrm()

	err := o.Read(&models.User{Name: userName, Pwd: userPwd}, "name", "pwd")
	if err != nil {
		c.Ctx.WriteString("login failed")
		beego.Info("login failed")
	} else {
		remember := c.GetString("remember")
		if remember == "on" {
			c.Ctx.SetCookie("userName", userName, time.Second*3600)
		} else {
			c.Ctx.SetCookie("userName", "null", -1)
		}
		c.SetSession("userName", userName)
		//c.Ctx.WriteString("login successful")
		c.Redirect("/index", 302)
	}
}
func (c *UserController) Get() {
	c.TplName = "register.html"
}

func (c *UserController) Post() {
	userName := c.GetString("username")
	userPwd := c.GetString("password")
	if userName == "" || userPwd == "" {
		c.Redirect("/", 302)
	}
	o := orm.NewOrm()
	_, err := o.Insert(&models.User{Name: userName, Pwd: userPwd})
	if err != nil {
		beego.Info("register failed", err)
		c.TplName = "register.html"
	} else {
		c.Data["username"] = userName
		c.TplName = "login.html"
	}
}

func (c *UserController) HandleLogout() {
	c.DelSession("userName")
	c.Redirect("/login", 302)
}
