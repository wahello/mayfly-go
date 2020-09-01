package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/go-sql-driver/mysql"
	_ "mayfly-go/routers"
	scheduler "mayfly-go/scheudler"
	"net/http"
	"strings"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:111049@tcp(localhost:3306)/mayfly-job?charset=utf8&loc=Local")
}

func main() {
	orm.Debug = true
	// 跨域配置
	beego.InsertFilter("/**", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	scheduler.Start()
	defer scheduler.Stop()
	beego.Run()
}

// 解决beego无法访问根目录静态文件
func TransparentStatic(ctx *context.Context) {
	if strings.Index(ctx.Request.URL.Path, "api/") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/"+ctx.Request.URL.Path)
}