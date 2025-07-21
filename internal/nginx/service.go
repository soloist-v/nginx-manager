package nginx

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	ExecutablePath string
	ConfigPath     string
	LogPath        string
	PidFile        string
}

type Status struct {
	IsRunning   bool      `json:"is_running"`
	PID         int       `json:"pid"`
	Uptime      string    `json:"uptime"`
	Version     string    `json:"version"`
	ConfigValid bool      `json:"config_valid"`
	LastError   string    `json:"last_error,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewService(execPath, configPath, logPath, pidFile string) *Service {
	return &Service{
		ExecutablePath: execPath,
		ConfigPath:     configPath,
		LogPath:        logPath,
		PidFile:        pidFile,
	}
}

// Start 启动nginx服务
func (s *Service) Start() error {
	// 检查是否已经运行
	if s.IsRunning() {
		return fmt.Errorf("nginx is already running")
	}

	// 检查可执行文件是否存在
	if _, err := os.Stat(s.ExecutablePath); os.IsNotExist(err) {
		return fmt.Errorf("nginx executable not found at: %s", s.ExecutablePath)
	}

	// 验证配置文件
	if err := s.TestConfig(); err != nil {
		return fmt.Errorf("config test failed: %w", err)
	}

	// 启动nginx
	cmd := exec.Command(s.ExecutablePath, "-c", s.ConfigPath)
	cmd.Dir = filepath.Dir(s.ExecutablePath)

	logrus.Infof("starting nginx executable: %s", s.ExecutablePath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start nginx: %w", err)
	}

	// 等待进程完全启动，最多等待5秒
	logrus.Info("waiting for nginx process to initialize...")
	for i := 0; i < 1000; i++ { // 最多等待5秒 (1000 * 5ms)
		if cmd.Process != nil {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}

	// 验证是否启动成功
	if cmd.Process == nil {
		return fmt.Errorf("nginx process failed to start")
	}

	// 进一步验证nginx是否真正运行
	//if !s.IsRunning() {
	//	return fmt.Errorf("nginx failed to start - check error logs")
	//}

	logrus.Infof("nginx started successfully (PID: %d)", cmd.Process.Pid)

	logrus.Info("Nginx started successfully")
	return nil
}

// Stop 停止nginx服务
func (s *Service) Stop() error {
	if !s.IsRunning() {
		return fmt.Errorf("nginx is not running")
	}

	// 使用nginx -s quit优雅停止
	cmd := exec.Command(s.ExecutablePath, "-c", s.ConfigPath, "-s", "quit")
	cmd.Dir = filepath.Dir(s.ExecutablePath) // 设置工作目录为nginx安装目录
	if err := cmd.Run(); err != nil {
		// 如果优雅停止失败，尝试强制杀死进程
		logrus.Warn("Graceful shutdown failed, trying force kill")
		return s.forceKill()
	}
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	// 等待进程完全停止
	deadline := time.NewTimer(5 * time.Second)
	select {
	case err := <-done:
		if err != nil {
			return err
		}
	case <-deadline.C:
		logrus.Warn("Graceful shutdown timed out")
	}
	return fmt.Errorf("nginx did not stop within timeout")
}

// Restart 重启nginx服务
func (s *Service) Restart() error {
	if s.IsRunning() {
		if err := s.Stop(); err != nil {
			return fmt.Errorf("failed to stop nginx: %w", err)
		}
	}

	// 等待一下确保完全停止
	time.Sleep(1 * time.Second)

	return s.Start()
}

// Reload 重新加载配置文件
func (s *Service) Reload() error {
	if !s.IsRunning() {
		return fmt.Errorf("nginx is not running")
	}

	// 先测试配置文件
	if err := s.TestConfig(); err != nil {
		return fmt.Errorf("config test failed: %w", err)
	}

	// 发送reload信号
	cmd := exec.Command(s.ExecutablePath, "-c", s.ConfigPath, "-s", "reload")
	cmd.Dir = filepath.Dir(s.ExecutablePath) // 设置工作目录为nginx安装目录
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload nginx: %w", err)
	}

	logrus.Info("Nginx reloaded successfully")
	return nil
}

// IsRunning 检查nginx是否正在运行
func (s *Service) IsRunning() bool {
	// 方法1：检查PID文件
	if pid := s.getPIDFromFile(); pid > 0 {
		logrus.Infof("nginx pid: %d", pid)
		if s.isPIDRunning(pid) {
			return true
		}
		// PID文件存在但进程不存在，清理PID文件
		//_ = os.Remove(s.PidFile)
	}

	// 方法2：检查是否有使用我们配置的nginx进程
	return s.checkManagedNginxProcess()
}

// GetStatus 获取nginx详细状态
func (s *Service) GetStatus() *Status {
	status := &Status{
		UpdatedAt: time.Now(),
	}

	status.IsRunning = s.IsRunning()
	if status.IsRunning {
		status.PID = s.getPIDFromFile()
		status.Uptime = s.getUptime()
	}

	status.Version = s.getVersion()
	status.ConfigValid = s.TestConfig() == nil

	return status
}

// TestConfig 测试nginx配置文件语法
func (s *Service) TestConfig() error {
	cmd := exec.Command(s.ExecutablePath, "-t", "-c", s.ConfigPath)
	cmd.Dir = filepath.Dir(s.ExecutablePath) // 设置工作目录为nginx安装目录
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("config test failed: %s", string(output))
	}
	return nil
}

// getPIDFromFile 从PID文件读取进程ID
func (s *Service) getPIDFromFile() int {
	if _, err := os.Stat(s.PidFile); os.IsNotExist(err) {
		return 0
	}

	content, err := os.ReadFile(s.PidFile)
	if err != nil {
		return 0
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return 0
	}

	return pid
}

// isPIDRunning 检查指定PID的进程是否存在
func (s *Service) isPIDRunning(pid int) bool {
	exist, err := process.PidExists(int32(pid))
	if err != nil {
		logrus.Infof("检查进程出错: %v\n", err)
		return false
	}
	if exist {
		logrus.Infof("进程 PID %d 存在\n", pid)
		return true
	}
	logrus.Infof("进程 PID %d 不存在\n", pid)
	return false
}

// checkProcessList 检查进程列表中是否有nginx
func (s *Service) checkProcessList() bool {
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq nginx.exe", "/FO", "CSV")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// 检查是否有nginx进程
	lines := strings.Split(string(output), "\n")
	return len(lines) > 1 // 第一行是标题，如果有nginx进程会有更多行
}

// checkManagedNginxProcess 检查是否有我们管理的nginx进程在运行
func (s *Service) checkManagedNginxProcess() bool {
	// 更保守的检查：只依赖PID文件存在和进程存在
	// 如果没有PID文件，则认为nginx没有运行（至少不是我们管理的）
	return false // 让主要检查逻辑依赖PID文件
}

// forceKill 强制杀死nginx进程
func (s *Service) forceKill() error {
	cmd := exec.Command("taskkill", "/F", "/IM", "nginx.exe")
	return cmd.Run()
}

// getVersion 获取nginx版本
func (s *Service) getVersion() string {
	cmd := exec.Command(s.ExecutablePath, "-v")
	cmd.Dir = filepath.Dir(s.ExecutablePath) // 设置工作目录为nginx安装目录
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	// 解析版本信息
	version := strings.TrimSpace(string(output))
	if strings.Contains(version, "nginx version:") {
		parts := strings.Split(version, "nginx version: ")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1])
		}
	}

	return "unknown"
}

// getUptime 获取运行时间
func (s *Service) getUptime() string {
	pid := s.getPIDFromFile()
	if pid <= 0 {
		return "unknown"
	}

	// 在Windows上，这个功能比较复杂，暂时返回简单信息
	return "running"
}
