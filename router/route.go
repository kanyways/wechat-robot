package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kanyways/wechat-robot/config"
	"github.com/kanyways/wechat-robot/controller/receive"
	"github.com/kanyways/wechat-robot/public"
	"html/template"
	"net/http"
)

// Route 路由
func Route(router *gin.Engine) {

	// 配置资源的文件目录
	router.StaticFS("/assets", public.Assets())
	// 配置网站的图标
	// router.StaticFile("/favicon.ico", "public/assets/favicon.ico")

	// 设置模板
	templates := template.Must(template.New("").ParseFS(public.Templates, "templates/*.html"))
	router.SetHTMLTemplate(templates)
	// 加载模板
	// 默认的欢迎页面
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   config.ServerConfig.SiteName,
			"message": "你来了呀？",
		})
	})

	PublicGroup := router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
		// 设置图标，采用 embed 将文件内嵌后，构建相对路径
		PublicGroup.GET("/favicon.ico", func(c *gin.Context) {
			c.FileFromFS("favicon.ico", public.Assets())
		})
	}

	apiPrefix := config.ServerConfig.ApiPrefix
	api := router.Group(apiPrefix)
	{
		// 收到Post请求
		api.POST("/receive/:id", receive.PostMessage)
		// 收到Get请求
		api.GET("/receive/:id", receive.GetMessage)
	}
}
