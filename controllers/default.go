package controllers

import (
	"WebServer/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplName = "register.html"
}

func (c *MainController) Post() {
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

func (c *MainController) ShowLogin() {
	c.TplName = "login.html"
}

func (c *MainController) HandleLogin() {
	userName := c.GetString("username")
	userPwd := c.GetString("password")
	if userName == "" || userPwd == "" {
		c.Redirect("/", 302)
	}
	o := orm.NewOrm()

	err := o.Read(&models.User{rand.Int(), userName, userPwd}, "name", "pwd")
	if err != nil {
		c.Ctx.WriteString("login failed")
		beego.Info("login failed")
	} else {
		//c.Ctx.WriteString("login successful")
		c.Redirect("/index", 302)
	}
}

func (c *MainController) ShowIndex() {
	o := orm.NewOrm()
	var article []models.Article
	_, err := o.QueryTable("article").All(&article)
	if err != nil {
		beego.Info("query failed")
		return
	}
	c.Data["articles"] = article
	c.TplName = "index.html"
}

func (c *MainController) ShowAdd() {
	c.TplName = "add.html"
}

func (c *MainController) HandleAdd() {
	artiName := c.GetString("articleName")
	artiContent := c.GetString("content")
	f, h, err := c.GetFile("uploadname")

	beego.Info(artiName, artiContent)
	if artiName == "" || artiContent == "" {
		beego.Info("content invaild")
		return
	}

	if err != nil {
		beego.Info("upload image failed")
	} else {
		defer f.Close()
		err = c.SaveToFile("uploadname", "./static/img/"+time.Now().Format("2006-01-02 15-04-05")+":"+h.Filename)
		if err != nil {
			beego.Info("save img failed")
			return
		}
	}

	o := orm.NewOrm()
	arti := models.Article{
		Aname:    artiName,
		Acontent: artiContent,
		Aimg:     "./static/img/" + time.Now().Format("2006-01-02 15-04-05") + ":" + h.Filename,
		Acount:   0,
		Atime:    time.Now().Local(),
	}
	_, err = o.Insert(&arti)
	if err != nil {
		beego.Info("insert failed", err)
		return
	}
	c.Redirect("/index", 302)
}

func (c *MainController) ShowContent() {
	id, err := c.GetInt("id")
	if err != nil {
		beego.Info("get id failed")
		return
	}
	arti := models.Article{Id: id}
	o := orm.NewOrm()
	err = o.Read(&arti)
	if err != nil {
		beego.Info("read article faield")
	}
	arti.Acount = arti.Acount + 1
	_, err = o.Update(&arti)
	c.Data["article"] = arti
	c.TplName = "content.html"
}

func (c *MainController) ShowUpdate() {
	id, err := c.GetInt("id")
	if err != nil {
		beego.Info("get id failed", err)
		return
	}
	arti := models.Article{Id: id}
	o := orm.NewOrm()
	err = o.Read(&arti)
	if err != nil {
		beego.Info("read failed")
		return
	}
	c.Data["article"] = arti
	c.TplName = "update.html"
}

func (c *MainController) HandleUpdate() {
	id, _ := c.GetInt("id")
	artiName := c.GetString("articleName")
	artiContent := c.GetString("content")
	f, h, err := c.GetFile("uploadname")

	beego.Info(artiName, artiContent)
	if artiName == "" || artiContent == "" {
		beego.Info("content invaild")
		return
	}
	var fileName string
	if err != nil {
		beego.Info("upload image failed")
		fileName = ""
	} else {
		if f != nil {
			defer f.Close()
			err = c.SaveToFile("uploadname", "./static/img/"+time.Now().Format("2006-01-02 15-04-05")+":")
			if err != nil {
				beego.Info("save img failed")
				return
			}
			fileName = h.Filename
			fileName = "./static/img/" + time.Now().Format("2006-01-02 15-04-05") + ":" + fileName
		}
	}

	o := orm.NewOrm()
	tmpArti := models.Article{Id: id}
	o.Read(&tmpArti)
	if fileName == "" {
		fileName = tmpArti.Aimg
	}
	arti := models.Article{
		Id:       id,
		Aname:    artiName,
		Acontent: artiContent,
		Aimg:     fileName,
		Acount:   0,
		Atime:    time.Now().Local(),
	}
	_, err = o.Update(&arti)
	if err != nil {
		beego.Info("insert failed", err)
		return
	}
	c.Redirect("/index", 302)
}

func (c *MainController) HandleDelete() {
	deleteId, err := c.GetInt("id")
	if err != nil {
		beego.Info(err)
		return
	}
	o := orm.NewOrm()
	arti := models.Article{Id: deleteId}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("dont has artical", err)
		return
	}
	_, err = o.Delete(&arti)
	c.Redirect("/index", 302)
}
