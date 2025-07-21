package handler

import (
	"net/http"
	"nginx_manager/internal/nginx"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"nginx_manager/internal/config"
)

type NginxHandler struct {
	service *nginx.Service
}

func NewNginxHandler() *NginxHandler {
	cfg := config.AppConfig.Nginx
	service := nginx.NewService(
		cfg.ExecutablePath,
		cfg.ConfigPath,
		cfg.LogPath,
		cfg.PidFile,
	)

	return &NginxHandler{
		service: service,
	}
}

// GetStatus 获取nginx状态
func (h *NginxHandler) GetStatus(c *gin.Context) {
	status := h.service.GetStatus()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}

// Start 启动nginx服务
func (h *NginxHandler) Start(c *gin.Context) {
	if err := h.service.Start(); err != nil {
		logrus.Error("Failed to start nginx: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Nginx started successfully",
	})
}

// Stop 停止nginx服务
func (h *NginxHandler) Stop(c *gin.Context) {
	if err := h.service.Stop(); err != nil {
		logrus.Error("Failed to stop nginx: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Nginx stopped successfully",
	})
}

// Restart 重启nginx服务
func (h *NginxHandler) Restart(c *gin.Context) {
	if err := h.service.Restart(); err != nil {
		logrus.Error("Failed to restart nginx: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Nginx restarted successfully",
	})
}

// Reload 重新加载nginx配置
func (h *NginxHandler) Reload(c *gin.Context) {
	if err := h.service.Reload(); err != nil {
		logrus.Error("Failed to reload nginx: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Nginx configuration reloaded successfully",
	})
}
