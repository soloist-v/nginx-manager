package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Nginx    NginxConfig    `mapstructure:"nginx"`
	Security SecurityConfig `mapstructure:"security"`
	Backup   BackupConfig   `mapstructure:"backup"`
}

type ServerConfig struct {
	Host  string `mapstructure:"host"`
	Port  int    `mapstructure:"port"`
	Debug bool   `mapstructure:"debug"`
}

type NginxConfig struct {
	ExecutablePath string `mapstructure:"executable_path"`
	ConfigPath     string `mapstructure:"config_path"`
	LogPath        string `mapstructure:"log_path"`
	PidFile        string `mapstructure:"pid_file"`
}

type SecurityConfig struct {
	EnableAuth bool   `mapstructure:"enable_auth"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
}

type BackupConfig struct {
	Enable     bool   `mapstructure:"enable"`
	BackupDir  string `mapstructure:"backup_dir"`
	MaxBackups int    `mapstructure:"max_backups"`
}

var AppConfig *Config

func LoadConfig(configPath string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("./configs")
	}
	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 确保路径格式正确
	AppConfig.Nginx.ExecutablePath = filepath.Clean(AppConfig.Nginx.ExecutablePath)
	AppConfig.Nginx.ConfigPath = filepath.Clean(AppConfig.Nginx.ConfigPath)
	AppConfig.Nginx.LogPath = filepath.Clean(AppConfig.Nginx.LogPath)
	AppConfig.Nginx.PidFile = filepath.Clean(AppConfig.Nginx.PidFile)

	return nil
}

func setDefaults() {
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.debug", true)
	viper.SetDefault("nginx.executable_path", "C:/nginx/nginx.exe")
	viper.SetDefault("nginx.config_path", "C:/nginx/conf/nginx.conf")
	viper.SetDefault("nginx.log_path", "C:/nginx/logs")
	viper.SetDefault("nginx.pid_file", "C:/nginx/logs/nginx.pid")
	viper.SetDefault("security.enable_auth", false)
	viper.SetDefault("backup.enable", true)
	viper.SetDefault("backup.backup_dir", "./backups")
	viper.SetDefault("backup.max_backups", 10)
}
