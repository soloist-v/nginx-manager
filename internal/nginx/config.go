package nginx

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type ConfigManager struct {
	ConfigPath string
	BackupDir  string
	MaxBackups int
}

type BackupInfo struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

func NewConfigManager(configPath, backupDir string, maxBackups int) *ConfigManager {
	return &ConfigManager{
		ConfigPath: configPath,
		BackupDir:  backupDir,
		MaxBackups: maxBackups,
	}
}

// ReadConfig 读取nginx配置文件
func (cm *ConfigManager) ReadConfig() (string, error) {
	if _, err := os.Stat(cm.ConfigPath); os.IsNotExist(err) {
		return "", fmt.Errorf("config file does not exist: %s", cm.ConfigPath)
	}

	content, err := os.ReadFile(cm.ConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	return string(content), nil
}

// WriteConfig 保存nginx配置文件
func (cm *ConfigManager) WriteConfig(content string) error {
	// 首先创建备份
	if err := cm.CreateBackup(); err != nil {
		logrus.Warn("Failed to create backup before saving config: ", err)
	}

	// 写入新配置
	if err := os.WriteFile(cm.ConfigPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	logrus.Info("Config file saved successfully")
	return nil
}

// CreateBackup 创建配置文件备份
func (cm *ConfigManager) CreateBackup() error {
	// 确保备份目录存在
	if err := os.MkdirAll(cm.BackupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// 读取当前配置
	content, err := cm.ReadConfig()
	if err != nil {
		return fmt.Errorf("failed to read config for backup: %w", err)
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupFilename := fmt.Sprintf("nginx_conf_%s.backup", timestamp)
	backupPath := filepath.Join(cm.BackupDir, backupFilename)

	// 保存备份
	if err := os.WriteFile(backupPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	// 清理旧备份
	cm.cleanOldBackups()

	logrus.Infof("Config backup created: %s", backupFilename)
	return nil
}

// RestoreBackup 恢复指定的备份
func (cm *ConfigManager) RestoreBackup(backupID string) error {
	backupPath := filepath.Join(cm.BackupDir, backupID)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupID)
	}

	// 读取备份内容
	content, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// 在恢复前创建当前配置的备份
	if err := cm.CreateBackup(); err != nil {
		logrus.Warn("Failed to backup current config before restore: ", err)
	}

	// 恢复配置
	if err := os.WriteFile(cm.ConfigPath, content, 0644); err != nil {
		return fmt.Errorf("failed to restore config: %w", err)
	}

	logrus.Infof("Config restored from backup: %s", backupID)
	return nil
}

// ListBackups 列出所有备份
func (cm *ConfigManager) ListBackups() ([]BackupInfo, error) {
	if _, err := os.Stat(cm.BackupDir); os.IsNotExist(err) {
		return []BackupInfo{}, nil
	}

	files, err := os.ReadDir(cm.BackupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []BackupInfo
	for _, file := range files {
		if file.IsDir() || !isBackupFile(file.Name()) {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		backup := BackupInfo{
			ID:        file.Name(),
			Filename:  file.Name(),
			CreatedAt: info.ModTime(),
			Size:      info.Size(),
		}
		backups = append(backups, backup)
	}

	return backups, nil
}

// DeleteBackup 删除指定备份
func (cm *ConfigManager) DeleteBackup(backupID string) error {
	backupPath := filepath.Join(cm.BackupDir, backupID)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupID)
	}

	if err := os.Remove(backupPath); err != nil {
		return fmt.Errorf("failed to delete backup: %w", err)
	}

	logrus.Infof("Backup deleted: %s", backupID)
	return nil
}

// GetConfigTemplate 获取基础nginx配置模板
func (cm *ConfigManager) GetConfigTemplate() string {
	return `# 基本的nginx配置模板 - 适用于Windows
worker_processes auto;

events {
    worker_connections 1024;
}

http {
    include       mime.types;
    default_type  application/octet-stream;

    # 日志格式
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    # 日志文件
    access_log logs/access.log main;
    error_log  logs/error.log;

    # 基本设置
    sendfile        on;
    keepalive_timeout  65;

    # 默认服务器配置
    server {
        listen       80;  # 监听80端口，适用于所有IP地址
        server_name  localhost;

        # 根目录
        location / {
            root   html;
            index  index.html index.htm;
        }

        # 错误页面
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
    }
}`
}

// cleanOldBackups 清理超过最大数量的旧备份
func (cm *ConfigManager) cleanOldBackups() {
	if cm.MaxBackups <= 0 {
		return
	}

	backups, err := cm.ListBackups()
	if err != nil {
		logrus.Warn("Failed to list backups for cleanup: ", err)
		return
	}

	if len(backups) <= cm.MaxBackups {
		return
	}

	// 按时间排序，删除最旧的备份
	// 这里简化处理，实际应该按创建时间排序
	for i := 0; i < len(backups)-cm.MaxBackups; i++ {
		if err := cm.DeleteBackup(backups[i].ID); err != nil {
			logrus.Warn("Failed to delete old backup: ", err)
		}
	}
}

// isBackupFile 检查是否为备份文件
func isBackupFile(filename string) bool {
	return filepath.Ext(filename) == ".backup"
}

// GetBackupPath 获取备份文件的完整路径
func (cm *ConfigManager) GetBackupPath(backupID string) (string, error) {
	if !isBackupFile(backupID) {
		return "", fmt.Errorf("invalid backup file: %s", backupID)
	}

	backupPath := filepath.Join(cm.BackupDir, backupID)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return "", fmt.Errorf("backup file not found: %s", backupID)
	}

	return backupPath, nil
}
