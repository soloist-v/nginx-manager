# Nginx Manager

A modern web-based Nginx service management tool built with Go backend and Vue.js frontend, designed for efficient Nginx configuration and monitoring on Windows systems.

## 🚀 Features

### Core Functionality
- **Service Management**: Start, stop, restart, and reload Nginx service with real-time status monitoring
- **Configuration Editor**: Advanced online editor with Monaco Editor providing syntax highlighting, auto-completion, and real-time validation
- **Backup System**: Automatic configuration backups with version control, restore capabilities, and download functionality
- **Real-time Monitoring**: WebSocket-powered live status updates with automatic reconnection
- **Log Management**: Real-time viewing of Nginx access and error logs with filtering capabilities
- **Security**: Optional basic authentication for web interface protection
- **System Monitoring**: Real-time system performance metrics using gopsutil

### User Interface
- **Modern UI**: Built with Vuetify 3 for a professional Material Design experience
- **Responsive Design**: Works seamlessly on desktop and mobile devices
- **Real-time Updates**: Live status indicators and automatic data refresh
- **Intuitive Navigation**: Clean, organized interface with easy access to all features

## 🛠 Technology Stack

### Backend (Go)
- **Framework**: Gin Web Framework for high-performance HTTP routing
- **WebSocket**: Gorilla WebSocket for real-time communication
- **Configuration**: Viper for flexible configuration management
- **Logging**: Logrus for structured logging
- **System Monitoring**: gopsutil for system metrics
- **CORS**: Built-in CORS middleware for cross-origin requests

### Frontend (Vue.js)
- **Framework**: Vue 3 with Composition API
- **UI Library**: Vuetify 3 for Material Design components
- **State Management**: Pinia for reactive state management
- **Routing**: Vue Router for SPA navigation
- **Code Editor**: Monaco Editor (VS Code's editor) with Nginx syntax support
- **HTTP Client**: Axios for API communication
- **Build Tool**: Vite for fast development and optimized builds

## 📁 Project Structure

```
nginx_manager/
├── main.go                    # Application entry point
├── go.mod                     # Go module dependencies
├── go.sum                     # Go module checksums
├── configs/
│   ├── config.yaml           # Application configuration
│   └── nginx_template.conf   # Nginx configuration template
├── internal/                  # Internal Go packages
│   ├── config/               # Configuration management
│   │   └── config.go
│   ├── handler/              # HTTP request handlers
│   │   ├── nginx.go          # Nginx service operations
│   │   ├── config.go         # Configuration file management
│   │   └── websocket.go      # WebSocket connections
│   ├── middleware/           # HTTP middleware
│   │   └── cors.go           # CORS handling
│   └── nginx/                # Nginx-specific utilities
│       ├── config.go         # Nginx configuration operations
│       └── service.go        # Nginx service management
├── frontend/                  # Vue.js frontend application
│   ├── src/
│   │   ├── api/              # API client functions
│   │   │   └── nginx.js
│   │   ├── stores/           # Pinia state stores
│   │   │   └── nginx.js
│   │   ├── views/            # Page components
│   │   │   ├── Dashboard.vue      # Main dashboard
│   │   │   ├── ConfigEditor.vue   # Configuration editor
│   │   │   ├── BackupManager.vue  # Backup management
│   │   │   └── LogViewer.vue      # Log viewing
│   │   ├── router/           # Vue Router configuration
│   │   │   └── index.js
│   │   ├── plugins/          # Vue plugins
│   │   │   └── vuetify.js
│   │   ├── css/              # Custom stylesheets
│   │   │   └── css2.css
│   │   ├── App.vue           # Root component
│   │   └── main.js           # Application entry point
│   ├── index.html            # HTML template
│   ├── package.json          # Node.js dependencies
│   └── vite.config.js        # Vite configuration
├── backups/                   # Automatic backup storage
├── logs/                      # Application logs
├── static/                    # Built frontend assets (production)
├── web/                       # Alternative static assets
└── README.md                  # This file
```

## 🚀 Quick Start

### Prerequisites
- **Go**: 1.23.0 or higher
- **Node.js**: 18.0 or higher
- **Nginx**: Installed on Windows system
- **Git**: For cloning the repository

### 1. Clone and Setup

```bash
git clone <repository-url>
cd nginx_manager
```

### 2. Configure the Application

Edit `configs/config.yaml` to match your Nginx installation:

```yaml
server:
  host: "127.0.0.1"
  port: 8080
  debug: true

nginx:
  executable_path: "C:/nginx/nginx.exe"        # Your nginx.exe path
  config_path: "C:/nginx/conf/nginx.conf"     # Your nginx.conf path
  log_path: "C:/nginx/logs"                   # Your logs directory
  pid_file: "C:/nginx/logs/nginx.pid"         # Your PID file path

security:
  enable_auth: false                           # Set to true for authentication
  username: "admin"                           # Authentication username
  password: "password"                        # Authentication password

backup:
  enable: true
  backup_dir: "./backups"
  max_backups: 10
```

### 3. Start the Backend

```bash
# Install Go dependencies
go mod tidy

# Run the server
go run main.go
```

The server will start at `http://localhost:8080`

### 4. Frontend Development (Optional)

For development with hot reload:

```bash
cd frontend
npm install
npm run dev
```

Frontend dev server will start at `http://localhost:5173`

### 5. Production Build

```bash
# Build frontend
cd frontend
npm install
npm run build

# Build backend (includes frontend assets)
cd ..
go build -o nginx-manager.exe main.go

# Run production binary
./nginx-manager.exe
```

## 📡 API Reference

### Nginx Service Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/nginx/status` | Get current Nginx service status |
| `POST` | `/api/nginx/start` | Start Nginx service |
| `POST` | `/api/nginx/stop` | Stop Nginx service |
| `POST` | `/api/nginx/restart` | Restart Nginx service |
| `POST` | `/api/nginx/reload` | Reload Nginx configuration |

### Configuration Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/config` | Get current configuration content |
| `PUT` | `/api/config` | Save configuration file |
| `POST` | `/api/config/validate` | Validate configuration syntax |
| `GET` | `/api/config/template` | Get configuration template |

### Backup Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/backup` | Get list of available backups |
| `GET` | `/api/backup/download/:id` | Download specific backup file |
| `POST` | `/api/backup/restore/:id` | Restore configuration from backup |
| `DELETE` | `/api/backup/:id` | Delete specific backup |

### WebSocket
| Endpoint | Description |
|----------|-------------|
| `WS /ws/status` | Real-time status updates with auto-reconnection |

## 🎯 Feature Details

### Dashboard
- **Real-time Status**: Live display of Nginx service status
- **Process Information**: PID, version, and uptime details
- **Quick Actions**: One-click service control buttons
- **System Metrics**: CPU, memory, and disk usage
- **Connection Status**: WebSocket connection indicator

### Configuration Editor
- **Advanced Editor**: Monaco Editor with Nginx syntax highlighting
- **Real-time Validation**: Instant syntax checking
- **Auto-completion**: Intelligent code suggestions
- **Auto-backup**: Automatic backup before saving changes
- **Template Support**: Pre-built configuration templates

### Backup Manager
- **Automatic Backups**: Timestamped backups on configuration changes
- **Version Control**: Browse and compare backup versions
- **One-click Restore**: Instant rollback to any previous version
- **Download Support**: Export backup files
- **Cleanup**: Automatic removal of old backups

### Log Viewer
- **Real-time Logs**: Live access and error log viewing
- **Filtering**: Search and filter log entries
- **Auto-refresh**: Automatic log updates
- **Export**: Download log files

## 🔧 Configuration Options

### Server Configuration
- `host`: Server binding address (default: "127.0.0.1")
- `port`: Server port (default: 8080)
- `debug`: Enable debug mode (default: true)

### Nginx Configuration
- `executable_path`: Path to nginx.exe
- `config_path`: Path to nginx.conf
- `log_path`: Directory containing Nginx logs
- `pid_file`: Path to Nginx PID file

### Security Configuration
- `enable_auth`: Enable basic authentication
- `username`: Authentication username
- `password`: Authentication password

### Backup Configuration
- `enable`: Enable automatic backups
- `backup_dir`: Backup storage directory
- `max_backups`: Maximum number of backups to keep

## 🛡️ Security Features

- **Basic Authentication**: Optional username/password protection
- **CORS Protection**: Configurable cross-origin request handling
- **Input Validation**: Comprehensive validation of all inputs
- **Error Handling**: Secure error responses without sensitive information
- **File Permissions**: Proper file access controls

## 📊 Monitoring & Logging

### Real-time Monitoring
- **Service Status**: Continuous monitoring of Nginx process
- **System Metrics**: CPU, memory, and disk usage tracking
- **WebSocket Updates**: Real-time status broadcasting
- **Connection Health**: Automatic reconnection on failures

### Logging
- **Structured Logging**: JSON-formatted logs with Logrus
- **Multiple Levels**: Debug, Info, Warning, Error levels
- **File Rotation**: Automatic log file management
- **Performance Tracking**: Request timing and performance metrics

## 🚨 Troubleshooting

### Common Issues

1. **Permission Denied**
   - Ensure the application has permissions to access Nginx files
   - Run as administrator if required

2. **Nginx Not Found**
   - Verify the `executable_path` in config.yaml
   - Ensure Nginx is properly installed

3. **Port Already in Use**
   - Change the port in config.yaml
   - Check for other services using port 8080

4. **WebSocket Connection Failed**
   - Check firewall settings
   - Verify CORS configuration

### Debug Mode
Enable debug mode in `config.yaml`:
```yaml
server:
  debug: true
```

This will provide detailed logging and error information.

## 📦 Dependencies

### Backend Dependencies
- `github.com/gin-gonic/gin` - Web framework
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/spf13/viper` - Configuration management
- `github.com/sirupsen/logrus` - Structured logging
- `github.com/gin-contrib/cors` - CORS middleware
- `github.com/shirou/gopsutil/v3` - System monitoring

### Frontend Dependencies
- `vue@^3.4.0` - Vue.js framework
- `vuetify@^3.4.0` - Material Design UI
- `pinia@^2.1.7` - State management
- `vue-router@^4.2.5` - Routing
- `axios@^1.6.2` - HTTP client
- `monaco-editor@^0.44.0` - Code editor

## 🤝 Contributing

We welcome contributions! Please feel free to submit issues and pull requests.

### Development Setup
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

### Code Style
- Follow Go formatting standards (`gofmt`)
- Use meaningful commit messages
- Add comments for complex logic
- Follow Vue.js style guide for frontend code

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - Fast HTTP web framework
- [Vuetify](https://vuetifyjs.com/) - Material Design component framework
- [Monaco Editor](https://microsoft.github.io/monaco-editor/) - Code editor
- [Vue.js](https://vuejs.org/) - Progressive JavaScript framework

---

**Nginx Manager** - Simplifying Nginx management on Windows systems. 