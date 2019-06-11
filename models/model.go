package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id       int
	Name     string `orm:"unique"`
	Pwd      string
	Articles []*Article `orm:"rel(m2m)"`
}

type Article struct {
	Id        int       `orm:"pk;auto"`
	Aname     string    `orm:"size(100)"`
	Atime     time.Time `orm:"auto_now"`
	Acount    int       `orm:"default(0);null"`
	Acontent  string
	Aimg      string
	AtypeName *ArticleType `orm:"rel(fk)"`
	Users     []*User      `orm:"reverse(many)"`
}

type ArticleType struct {
	Id       int
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"`
}

func init() {
	orm.RegisterDataBase("default", "mysql", "zhangfeng:980530@tcp(localhost:3306)/users?charset=utf8")

	orm.RegisterModel(new(User), new(Article), new(ArticleType))

	orm.RunSyncdb("default", false, true)
}
