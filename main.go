package main

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kanyways/wechat-robot/config"
	"github.com/kanyways/wechat-robot/model"
	"github.com/kanyways/wechat-robot/router"
	"io"
	"os"
)

func main() {
	fmt.Println("gin.Version: ", gin.Version)
	if config.ServerConfig.Env != model.DevelopmentMode {
		gin.SetMode(gin.ReleaseMode)
		// Disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()
		// Logging to a file.
		logFile, err := os.OpenFile(config.ServerConfig.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(-1)
		}
		gin.DefaultWriter = io.MultiWriter(logFile)
	}
	app := gin.Default()
	router.Route(app)
	app.Run(":" + fmt.Sprintf("%d", config.ServerConfig.Port))
}
