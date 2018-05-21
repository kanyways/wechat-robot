package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/kanyways/wechat-robot/controller/receive"
	"github.com/kanyways/wechat-robot/config"
	"github.com/kanyways/wechat-robot/utils"
)

// Route 路由
func Route(router *gin.Engine) {

	//配置今天资源的文件目录
	router.Static("/static", utils.GetFilePath("static"))
	// 配置网站的图标
	router.StaticFile("/favicon.ico", utils.GetFilePath("static/favicon.ico"))
	//加载模板
	router.LoadHTMLGlob(utils.GetFilePath("views/*"))

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   config.ServerConfig.SiteName,
			"message": "你来了呀？",
		})
	})

	apiPrefix := config.ServerConfig.APIPrefix
	api := router.Group(apiPrefix)
	{
		//接收到消息
		api.POST("/receive/:id", receive.PostMessage)
		api.GET("/receive/:id", receive.GetMessage)
	}
}
