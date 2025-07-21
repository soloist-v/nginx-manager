package main

import (
	"fmt"
	"log"
	"nginx_manager/internal/config"
	"nginx_manager/internal/handler"
	"nginx_manager/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("./configs/config.yaml"); err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// 配置日志
	if config.AppConfig.Server.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORSMiddleware())

	// 创建处理器
	nginxHandler := handler.NewNginxHandler()
	configHandler := handler.NewConfigHandler()
	wsHandler := handler.NewWebSocketHandler()

	// API路由
	api := r.Group("/api")
	{
		// nginx服务管理
		nginx := api.Group("/nginx")
		{
			nginx.GET("/status", nginxHandler.GetStatus)
			nginx.POST("/start", nginxHandler.Start)
			nginx.POST("/stop", nginxHandler.Stop)
			nginx.POST("/restart", nginxHandler.Restart)
			nginx.POST("/reload", nginxHandler.Reload)
		}

		// 配置文件管理
		configRouter := api.Group("/config")
		{
			configRouter.GET("", configHandler.GetConfig)
			configRouter.PUT("", configHandler.SaveConfig)
			configRouter.POST("/validate", configHandler.ValidateConfig)
			configRouter.GET("/template", configHandler.GetTemplate)
		}

		// 备份管理
		backup := api.Group("/backup")
		{
			backup.GET("", configHandler.GetBackups)
			backup.GET("/download/:id", configHandler.DownloadBackup)
			backup.POST("/restore/:id", configHandler.RestoreBackup)
			backup.DELETE("/:id", configHandler.DeleteBackup)
		}
	}

	// WebSocket端点
	r.GET("/ws/status", wsHandler.HandleWebSocket)

	// 静态文件服务 (生产环境中用于服务前端文件)
	r.Static("/assets", "static/assets")
	r.StaticFile("/favicon.ico", "static/favicon.ico")
	r.LoadHTMLFiles("static/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	logrus.Infof("Starting server on %s://%s", "http", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
