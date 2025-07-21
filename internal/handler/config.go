package handler

import (
	"fmt"
	"net/http"
	nginx2 "nginx_manager/internal/nginx"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"nginx_manager/internal/config"
)

type ConfigHandler struct {
	configManager *nginx2.ConfigManager
	nginxService  *nginx2.Service
}

type ConfigRequest struct {
	Content string `json:"content" binding:"required"`
}

func NewConfigHandler() *ConfigHandler {
	cfg := config.AppConfig
	configManager := nginx2.NewConfigManager(
		cfg.Nginx.ConfigPath,
		cfg.Backup.BackupDir,
		cfg.Backup.MaxBackups,
	)

	nginxService := nginx2.NewService(
		cfg.Nginx.ExecutablePath,
		cfg.Nginx.ConfigPath,
		cfg.Nginx.LogPath,
		cfg.Nginx.PidFile,
	)

	return &ConfigHandler{
		configManager: configManager,
		nginxService:  nginxService,
	}
}

// GetConfig 获取nginx配置文件内容
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	content, err := h.configManager.ReadConfig()
	if err != nil {
		logrus.Error("Failed to read config: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    content,
	})
}

// SaveConfig 保存nginx配置文件
func (h *ConfigHandler) SaveConfig(c *gin.Context) {
	var req ConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if err := h.configManager.WriteConfig(req.Content); err != nil {
		logrus.Error("Failed to save config: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Configuration saved successfully",
	})
}

// ValidateConfig 验证nginx配置文件语法
func (h *ConfigHandler) ValidateConfig(c *gin.Context) {
	var req ConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	// 先保存当前配置作为临时备份
	originalContent, err := h.configManager.ReadConfig()
	if err != nil {
		logrus.Error("Failed to read original config: ", err)
	}

	// 写入新配置进行测试
	if err := h.configManager.WriteConfig(req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to write test config",
		})
		return
	}

	// 测试配置
	err = h.nginxService.TestConfig()

	// 恢复原配置
	if originalContent != "" {
		if restoreErr := h.configManager.WriteConfig(originalContent); restoreErr != nil {
			logrus.Error("Failed to restore original config: ", restoreErr)
		}
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"valid":   false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"valid":   true,
		"message": "Configuration is valid",
	})
}

// GetBackups 获取备份列表
func (h *ConfigHandler) GetBackups(c *gin.Context) {
	backups, err := h.configManager.ListBackups()
	if err != nil {
		logrus.Error("Failed to list backups: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    backups,
	})
}

// RestoreBackup 恢复指定备份
func (h *ConfigHandler) RestoreBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Backup ID is required",
		})
		return
	}

	if err := h.configManager.RestoreBackup(backupID); err != nil {
		logrus.Error("Failed to restore backup: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Backup restored successfully",
	})
}

// DeleteBackup 删除指定备份
func (h *ConfigHandler) DeleteBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Backup ID is required",
		})
		return
	}

	if err := h.configManager.DeleteBackup(backupID); err != nil {
		logrus.Error("Failed to delete backup: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Backup deleted successfully",
	})
}

// DownloadBackup 下载指定备份文件
func (h *ConfigHandler) DownloadBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Backup ID is required",
		})
		return
	}

	// 获取备份文件路径
	backupPath, err := h.configManager.GetBackupPath(backupID)
	if err != nil {
		logrus.Error("Failed to get backup path: ", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", backupID))
	c.Header("Content-Transfer-Encoding", "binary")

	// 发送文件
	c.File(backupPath)
}

// GetTemplate 获取nginx配置模板
func (h *ConfigHandler) GetTemplate(c *gin.Context) {
	template := h.configManager.GetConfigTemplate()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    template,
	})
}
