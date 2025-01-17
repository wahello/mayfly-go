package routers

import (
	"mayfly-go/base/ctx"
	"mayfly-go/server/devops/apis"
	"mayfly-go/server/devops/application"

	"github.com/gin-gonic/gin"
)

func InitDbRouter(router *gin.RouterGroup) {
	db := router.Group("dbs")
	{
		d := &apis.Db{DbApp: application.DbApp}
		// 获取所有数据库列表
		db.GET("", func(c *gin.Context) {
			rc := ctx.NewReqCtxWithGin(c)
			rc.Handle(d.Dbs)
		})

		saveDb := ctx.NewLogInfo("保存数据库信息")
		db.POST("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).
				WithLog(saveDb).
				Handle(d.Save)
		})

		deleteDb := ctx.NewLogInfo("删除数据库信息")
		db.DELETE(":id", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).
				WithLog(deleteDb).
				Handle(d.DeleteDb)
		})

		db.GET(":dbId/t-infos", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(d.TableInfos)
		})

		db.GET(":dbId/t-index", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(d.TableIndex)
		})

		db.GET(":dbId/t-create-ddl", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).WithNeedToken(false).Handle(d.GetCreateTableDdl)
		})

		// db.GET(":dbId/exec-sql", controllers.SelectData)
		db.GET(":dbId/exec-sql", func(g *gin.Context) {
			rc := ctx.NewReqCtxWithGin(g).WithLog(ctx.NewLogInfo("执行Sql语句"))
			rc.Handle(d.ExecSql)
		})

		db.GET(":dbId/t-metadata", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(d.TableMA)
		})

		db.GET(":dbId/c-metadata", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(d.ColumnMA)
		})
		db.GET(":dbId/hint-tables", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(d.HintTables)
		})

		db.POST(":dbId/sql", func(c *gin.Context) {
			rc := ctx.NewReqCtxWithGin(c)
			rc.Handle(d.SaveSql)
		})

		db.GET(":dbId/sql", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(d.GetSql)
		})
	}
}
