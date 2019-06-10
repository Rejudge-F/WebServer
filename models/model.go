package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id   int
	Name string `orm:"unique"`
	Pwd  string
}

type Article struct {
	Id       int       `orm:"pk;auto"`
	Aname    string    `orm:"size(200)"`
	Atime    time.Time `orm:"auto_now"`
	Acount   int       `orm:"default(0);null"`
	Acontent string
	Aimg     string
}

func init() {
	orm.RegisterDataBase("default", "mysql", "zhangfeng:980530@tcp(localhost:3306)/users?charset=utf8")

	orm.RegisterModel(new(User), new(Article))

	orm.RunSyncdb("default", false, true)
}
