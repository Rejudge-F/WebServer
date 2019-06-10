package controllers

import (
	"WebServer/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) ShowIndex() {
	o := orm.NewOrm()
	var article []models.Article
	qs := o.QueryTable("article")

	count, err := qs.Count()
	if err != nil {
		beego.Info("select failed")
		return
	}

	pageIndex, err := strconv.Atoi(c.GetString("pageIndex"))
	if err != nil {
		pageIndex = 1
	}
	if pageIndex == 1 {
		FirPage := true
		c.Data["FirPage"] = FirPage
	}

	start := 10 * (pageIndex - 1)

	_, err = qs.Limit(10, start, start+10).All(&article)
	if err != nil {
		beego.Info(err)
		return
	}

	pageSize := count / 10
	if count%10 != 0 {
		pageSize += 1
	}

	if pageIndex == int(pageSize) {
		LastPage := true
		c.Data["LastPage"] = LastPage
	}

	c.Data["count"] = count
	c.Data["pageSize"] = pageSize
	c.Data["pageIndex"] = pageIndex
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
			fileName = h.Filename
			fileName = "./static/img/" + time.Now().Format("2006-01-02 15-04-05") + ":" + fileName
			defer f.Close()
			err = c.SaveToFile("uploadname", fileName)
			if err != nil {
				beego.Info("save img failed")
				return
			}

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
