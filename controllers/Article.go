package controllers

import (
	"WebServer/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"path"
	"strconv"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) HandleIndex() {

}

func (c *MainController) ShowIndex() {
	name := c.GetSession("userName")
	if name == nil {
		c.Redirect("/login", 302)
		return
	}
	o := orm.NewOrm()
	var article []models.Article
	qs := o.QueryTable("article")

	count, err := qs.RelatedSel("AtypeName").Count()
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

	pageSize := count / 10
	if count%10 != 0 {
		pageSize += 1
	}

	if pageIndex >= int(pageSize) {
		LastPage := true
		c.Data["LastPage"] = LastPage
	}
	var artiType []models.ArticleType
	_, err = o.QueryTable("article_type").All(&artiType)
	if err != nil {
		beego.Info(err)
	}

	typeName := c.GetString("select")
	// 初始的时候没有类型
	if typeName == "" {
		_, err = qs.Limit(10, start, start+10).RelatedSel("AtypeName").All(&article)
		if err != nil {
			beego.Info(err)
			return
		}
	} else {
		_, err = qs.Limit(10, start, start+10).RelatedSel("AtypeName").Filter("AtypeName__TypeName", typeName).All(&article)
		if err != nil {
			beego.Info(err)
		}
	}

	c.Data["types"] = artiType
	c.Data["count"] = count
	c.Data["pageSize"] = pageSize
	c.Data["pageIndex"] = pageIndex
	c.Data["articles"] = article
	c.TplName = "index.html"
}

func (c *MainController) ShowAdd() {
	name := c.GetSession("userName")
	if name == nil {
		c.Redirect("/login", 302)
		return
	}

	var artiType []*models.ArticleType
	o := orm.NewOrm()
	_, err := o.QueryTable("article_type").All(&artiType)
	if err != nil {
		beego.Info(err)
		c.Redirect("/index", 302)
		return
	}
	c.Data["artiType"] = artiType
	c.TplName = "add.html"
}

func (c *MainController) HandleAdd() {
	name := c.GetSession("userName")
	if name == nil {
		c.Redirect("/login", 302)
		return
	}

	artiName := c.GetString("articleName")
	artiContent := c.GetString("content")
	id, err := c.GetInt("select")
	if err != nil {
		beego.Info("获取类型错误")
		return
	}
	f, h, err := c.GetFile("uploadname")
	defer f.Close()

	fileext := path.Ext(h.Filename)

	if fileext != ".jpg" && fileext != "png" {
		beego.Info("上传文件格式错误")
		return
	}

	if h.Size > 50000000 {
		beego.Info("上传文件过大")
		return
	}
	filename := time.Now().Format("2006-01-02 15:04:05") + fileext //6-1-2 3:4:5

	if err != nil {
		beego.Info("上传文件失败")
		return
	} else {
		c.SaveToFile("uploadname", "./static/img/"+filename)
	}

	if artiContent == "" || artiName == "" {
		beego.Info("添加文章数据错误")
		return
	}
	o := orm.NewOrm()
	arti := models.Article{}
	arti.Aname = artiName
	arti.Acontent = artiContent
	arti.Aimg = "./static/img/" + filename

	artiType := models.ArticleType{Id: id}
	o.Read(&artiType)
	if artiType.TypeName == "" {
		beego.Info("null type")
		c.Redirect("/addarticle", 302)
		return
	}
	arti.AtypeName = &artiType

	_, err = o.Insert(&arti)
	if err != nil {
		beego.Info("插入数据库错误")
		c.Redirect("/addarticle", 302)
		return
	}

	c.Redirect("/index", 302)
}

func (c *MainController) ShowContent() {
	name := c.GetSession("userName")
	if name == nil {
		c.Redirect("/login", 302)
		return
	}

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
	name := c.GetSession("userName")
	if name == nil {
		c.Redirect("/login", 302)
		return
	}

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
	name := c.GetSession("userName")
	if name == nil {
		c.Redirect("/login", 302)
		return
	}

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
	name := c.GetSession("userName")
	if name == nil {
		c.Redirect("/login", 302)
		return
	}

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

func (c *MainController) ShowType() {
	name := c.GetSession("userName")
	if name == nil {
		c.Redirect("/login", 302)
		return
	}

	var articleType []models.ArticleType
	o := orm.NewOrm()
	_, err := o.QueryTable("article_type").All(&articleType)
	if err != nil {
		beego.Info(err)
	}
	c.Data["articleType"] = articleType
	c.TplName = "addType.html"
}

func (c *MainController) HandleType() {
	name := c.GetSession("userName")
	if name == nil {
		c.Redirect("/login", 302)
		return
	}

	typeName := c.GetString("typeName")
	if typeName == "" {
		beego.Info("null name")
		c.Redirect("/addtype", 302)
		return
	}
	var Type models.ArticleType
	Type.TypeName = typeName
	o := orm.NewOrm()
	_, err := o.Insert(&Type)
	if err != nil {
		beego.Info(err)
	}
	c.Redirect("/addtype", 302)
}
