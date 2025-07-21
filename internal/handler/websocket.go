package handler

import (
	"encoding/json"
	"net/http"
	"nginx_manager/internal/nginx"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"nginx_manager/internal/config"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 在生产环境中应该检查Origin
		return true
	},
}

type WebSocketHandler struct {
	nginxService *nginx.Service
	clients      map[*websocket.Conn]bool
	broadcast    chan []byte
}

type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	Time time.Time   `json:"time"`
}

func NewWebSocketHandler() *WebSocketHandler {
	cfg := config.AppConfig.Nginx
	service := nginx.NewService(
		cfg.ExecutablePath,
		cfg.ConfigPath,
		cfg.LogPath,
		cfg.PidFile,
	)

	handler := &WebSocketHandler{
		nginxService: service,
		clients:      make(map[*websocket.Conn]bool),
		broadcast:    make(chan []byte),
	}

	// 启动广播协程
	go handler.handleBroadcast()

	// 启动状态监控协程
	go handler.monitorStatus()

	return handler
}

// HandleWebSocket 处理WebSocket连接
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Error("Failed to upgrade websocket: ", err)
		return
	}
	defer conn.Close()

	// 注册客户端
	h.clients[conn] = true
	logrus.Info("New WebSocket client connected")

	// 发送当前状态
	h.sendCurrentStatus(conn)

	// 监听客户端消息
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			logrus.Debug("WebSocket client disconnected: ", err)
			delete(h.clients, conn)
			break
		}
	}
}

// handleBroadcast 处理广播消息
func (h *WebSocketHandler) handleBroadcast() {
	for {
		message := <-h.broadcast

		// 向所有连接的客户端发送消息
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logrus.Debug("Error writing to websocket: ", err)
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}

// monitorStatus 监控nginx状态变化
func (h *WebSocketHandler) monitorStatus() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var lastStatus *nginx.Status

	for {
		select {
		case <-ticker.C:
			currentStatus := h.nginxService.GetStatus()

			// 如果状态发生变化，广播给所有客户端
			if lastStatus == nil || statusChanged(lastStatus, currentStatus) {
				h.broadcastStatus(currentStatus)
				lastStatus = currentStatus
			}
		}
	}
}

// sendCurrentStatus 向特定客户端发送当前状态
func (h *WebSocketHandler) sendCurrentStatus(conn *websocket.Conn) {
	status := h.nginxService.GetStatus()
	message := WSMessage{
		Type: "status",
		Data: status,
		Time: time.Now(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		logrus.Error("Failed to marshal status message: ", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logrus.Error("Failed to send current status: ", err)
	}
}

// broadcastStatus 广播状态更新
func (h *WebSocketHandler) broadcastStatus(status *nginx.Status) {
	message := WSMessage{
		Type: "status",
		Data: status,
		Time: time.Now(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		logrus.Error("Failed to marshal status message: ", err)
		return
	}

	select {
	case h.broadcast <- data:
	default:
		logrus.Warn("Broadcast channel full, dropping message")
	}
}

// BroadcastEvent 广播事件消息
func (h *WebSocketHandler) BroadcastEvent(eventType, message string) {
	wsMessage := WSMessage{
		Type: "event",
		Data: map[string]string{
			"type":    eventType,
			"message": message,
		},
		Time: time.Now(),
	}

	data, err := json.Marshal(wsMessage)
	if err != nil {
		logrus.Error("Failed to marshal event message: ", err)
		return
	}

	select {
	case h.broadcast <- data:
	default:
		logrus.Warn("Broadcast channel full, dropping event message")
	}
}

// statusChanged 检查状态是否发生变化
func statusChanged(old, new *nginx.Status) bool {
	if old.IsRunning != new.IsRunning {
		return true
	}
	if old.PID != new.PID {
		return true
	}
	if old.ConfigValid != new.ConfigValid {
		return true
	}
	return false
}
