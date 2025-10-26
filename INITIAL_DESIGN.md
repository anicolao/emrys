# Emrys Initial Design

## Overview

Emrys will be built as a Go-based console application with a Terminal User Interface (TUI) using the [Bubbletea](https://github.com/charmbracelet/bubbletea) framework. The TUI will be accessible both locally and via the web through a browser-based interface powered by [xterm.js](https://xtermjs.org/), providing both read-only monitoring and active interaction modes.

## Architecture

### 1. Core Application Layer (Go)

The core of Emrys is a Go application that serves dual purposes:
- Runs the TUI console interface locally
- Provides web server functionality for remote access

#### Components:
- **Main Application**: Entry point and lifecycle management
- **TUI Engine**: Bubbletea-based interactive terminal interface
- **Web Server**: HTTP(S) server for browser-based access
- **PTY Manager**: Pseudo-terminal allocation and management for web sessions
- **Agent Core**: AI agent orchestration and task execution
- **Tool Registry**: Available system tools and capabilities

### 2. Terminal User Interface (Bubbletea)

The TUI provides a rich console interface for interacting with Emrys:

#### Features:
- **Interactive Console**: Command input and output display
- **Status Dashboard**: System state, running tasks, resource usage
- **Task Monitor**: View active and completed agent tasks
- **Log Viewer**: Real-time display of agent activities and system logs
- **Configuration Interface**: Settings and preferences management

#### Design Principles:
- Clean, readable text-based interface
- Keyboard-driven navigation and commands
- Responsive layout that adapts to terminal size
- Clear visual hierarchy using colors and formatting

### 3. Web Interface Layer

The web interface makes the TUI accessible through a browser, enabling remote monitoring and control.

#### Architecture:
- **HTTP Server**: Serves static frontend assets and WebSocket connections
- **WebSocket Handler**: Bidirectional communication between browser and TUI
- **Session Manager**: Handles multiple concurrent connections
- **Authentication Middleware**: Secures access to the interface
- **Mode Controller**: Manages read-only vs. active access permissions

#### Frontend Stack:
- **xterm.js**: Terminal emulator rendering in the browser
- **WebSocket**: Real-time communication with the Go backend
- **Static HTML/CSS/JS**: Minimal frontend dependencies

### 4. Access Modes

The system supports two distinct modes of web access:

#### Read-Only Mode:
- View TUI output and system status
- Monitor agent activities and logs
- No command input or control capabilities
- Multiple simultaneous read-only viewers permitted
- Ideal for monitoring dashboards or status screens

#### Active Mode:
- Full interactive access to the TUI
- Execute commands and interact with the agent
- Exclusive session (single active user at a time)
- Complete control over agent operations
- Requires elevated authentication

### 5. Initial Capability Goals

The first iteration of Emrys will demonstrate core capabilities through two primary tasks:

#### Goal 1: Self-Bootstrap via Nix Installation
- **Objective**: Prove Emrys can modify its own system environment
- **Tasks**:
  - Detect current macOS environment
  - Download and install the Nix package manager
  - Verify Nix installation and basic functionality
  - Configure Nix for optimal macOS usage
  
- **Technical Requirements**:
  - Shell command execution capability
  - Installation script management
  - Environment variable handling
  - Post-installation verification

#### Goal 2: Package Management via Nix
- **Objective**: Demonstrate controlled software installation
- **Tasks**:
  - Use Nix to install basic utilities
  - Manage package profiles
  - Update and maintain installed software
  - Clean up unused packages

- **Technical Requirements**:
  - Nix command orchestration
  - Package search and selection
  - Installation monitoring and verification
  - Rollback capabilities

#### Goal 3: Web Browser Access
- **Objective**: Prove Emrys can interact with the web
- **Tasks**:
  - Install a headless browser (via Nix)
  - Launch and control browser instances
  - Navigate to websites and extract information
  - Demonstrate basic web automation

- **Technical Requirements**:
  - Browser process management
  - Playwright or Selenium integration
  - DOM interaction and data extraction
  - Screenshot and debugging capabilities

## Component Communication Flow

```
┌─────────────────────────────────────────────────────────────┐
│                        User Access                          │
├─────────────────┬───────────────────────────────────────────┤
│   Local TTY     │        Web Browser (xterm.js)             │
└────────┬────────┴──────────────────┬────────────────────────┘
         │                           │
         │                           │ WebSocket
         │                           │
         ▼                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   Bubbletea TUI Engine                      │
│  ┌──────────────┬──────────────┬──────────────────────┐    │
│  │   Console    │    Status    │    Task Monitor      │    │
│  │   Interface  │   Dashboard  │    & Logs            │    │
│  └──────────────┴──────────────┴──────────────────────┘    │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          │ Commands
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                   Agent Core (Go)                           │
│  ┌──────────────┬──────────────┬──────────────────────┐    │
│  │  Task        │   Tool       │    LLM               │    │
│  │  Planner     │   Executor   │    Integration       │    │
│  └──────────────┴──────────────┴──────────────────────┘    │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          │ System Calls
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    Tool Layer                               │
│  ┌──────────────┬──────────────┬──────────────────────┐    │
│  │   Shell      │  AppleScript │    Browser           │    │
│  │   Executor   │  Executor    │    Automation        │    │
│  └──────────────┴──────────────┴──────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## Technology Stack

### Backend (Go)
- **Framework**: Standard library + Bubbletea
- **TUI Library**: [Bubbletea](https://github.com/charmbracelet/bubbletea)
- **TUI Components**: [Bubbles](https://github.com/charmbracelet/bubbles) (text input, viewport, etc.)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) (terminal styling)
- **Web Server**: net/http (standard library)
- **WebSocket**: [gorilla/websocket](https://github.com/gorilla/websocket) or standard library websocket
- **PTY Management**: [creack/pty](https://github.com/creack/pty)

### Frontend (Web)
- **Terminal Emulator**: [xterm.js](https://xtermjs.org/)
- **Communication**: WebSocket API (browser native)
- **UI Framework**: Vanilla JavaScript (minimal dependencies)
- **Styling**: CSS (minimal framework)

### AI/Agent
- **LLM Runtime**: Ollama (local inference)
- **Models**: Open-source models (Llama 3, Mistral, Qwen, etc.)
- **Browser Automation**: Playwright (Go bindings)

### System Tools
- **Package Management**: Nix
- **Shell**: zsh/bash (standard macOS)
- **Scripting**: AppleScript (macOS automation)

## Data Flow

### Local TUI Interaction:
1. User inputs command in terminal
2. Bubbletea handles input event
3. Command dispatched to Agent Core
4. Agent processes and executes via Tool Layer
5. Results rendered back to TUI
6. Display updated in terminal

### Web Interface Interaction:
1. User types in xterm.js browser terminal
2. Input sent via WebSocket to Go server
3. Server writes input to PTY
4. PTY connected to Bubbletea TUI
5. TUI processes command (same as local flow)
6. TUI output written to PTY
7. Server reads PTY output and sends via WebSocket
8. xterm.js renders output in browser

### Read-Only Mode:
1. Browser connects via WebSocket
2. Server flags session as read-only
3. TUI output streamed to browser
4. Browser input disabled or ignored
5. Updates continue in real-time

## Security Considerations

### Web Access:
- **Authentication**: Token-based or session-based auth for web access
- **HTTPS**: TLS encryption for web interface
- **WebSocket Security**: WSS (WebSocket Secure) protocol
- **Access Control**: Separate tokens for read-only vs. active modes
- **Rate Limiting**: Prevent abuse of web endpoints

### Agent Operations:
- **Sandboxing**: Initial operations restricted to safe commands
- **Confirmation**: Destructive operations require explicit confirmation
- **Logging**: All agent actions logged for audit trail
- **Rollback**: Ability to undo changes where possible

### System Security:
- **Credential Management**: Secure storage for API keys and tokens
- **File Permissions**: Appropriate restrictions on config and data files
- **Network Security**: Firewall rules for web server port
- **Update Mechanism**: Secure method for updating Emrys itself

## Configuration

### Application Configuration:
```
# Server settings
server.host: localhost
server.port: 8080
server.tls: true
server.cert: /path/to/cert.pem
server.key: /path/to/key.pem

# Authentication
auth.enabled: true
auth.read_only_token: <generated-token>
auth.active_token: <generated-token>

# Agent settings
agent.model: llama3
agent.ollama_url: http://localhost:11434
agent.max_concurrent_tasks: 3

# Logging
log.level: info
log.file: /var/log/emrys/emrys.log
log.max_size_mb: 100
```

## Development Phases

### Phase 1: Foundation (Current Focus)
- Basic Go application structure
- Simple Bubbletea TUI with command input
- Static content display
- Logging system

### Phase 2: Web Interface
- HTTP server implementation
- WebSocket communication
- xterm.js integration
- PTY management
- Basic authentication

### Phase 3: Agent Integration
- Ollama LLM connection
- Basic command execution
- Shell tool implementation
- Task logging and display

### Phase 4: Initial Capabilities
- Nix installation workflow
- Package management via Nix
- Browser automation setup
- Web browsing capability demonstration

### Phase 5: Polish & Enhancement
- Multi-user session management
- Enhanced read-only mode features
- Improved error handling
- Comprehensive logging
- Configuration management

## File Structure

```
emrys/
├── cmd/
│   └── emrys/
│       └── main.go                 # Application entry point
├── internal/
│   ├── agent/
│   │   ├── core.go                 # Agent orchestration
│   │   ├── planner.go              # Task planning
│   │   └── executor.go             # Task execution
│   ├── tui/
│   │   ├── app.go                  # Bubbletea application
│   │   ├── commands.go             # Command handling
│   │   ├── views.go                # UI views/screens
│   │   └── components.go           # Reusable UI components
│   ├── server/
│   │   ├── http.go                 # HTTP server
│   │   ├── websocket.go            # WebSocket handler
│   │   ├── auth.go                 # Authentication
│   │   └── session.go              # Session management
│   ├── tools/
│   │   ├── shell.go                # Shell command execution
│   │   ├── nix.go                  # Nix package management
│   │   ├── browser.go              # Browser automation
│   │   └── applescript.go          # AppleScript execution
│   └── llm/
│       ├── client.go               # Ollama client
│       └── models.go               # Model management
├── web/
│   ├── static/
│   │   ├── index.html              # Main web page
│   │   ├── app.js                  # Frontend JavaScript
│   │   └── styles.css              # Styling
│   └── assets/
│       └── xterm/                  # xterm.js library
├── config/
│   └── config.yaml                 # Configuration file
├── go.mod                          # Go module definition
├── go.sum                          # Go dependencies
└── README.md                       # Project documentation
```

## Success Criteria

The initial design succeeds when:

1. **TUI Functionality**: Bubbletea interface runs smoothly with command input/output
2. **Web Access**: Browser can connect and view/interact with the TUI via xterm.js
3. **Mode Switching**: Read-only and active modes work as specified
4. **Nix Installation**: Emrys can successfully install Nix on a fresh macOS system
5. **Package Management**: Emrys can install and manage packages using Nix
6. **Web Browsing**: Emrys can launch a browser and navigate to websites
7. **Logging**: All operations are logged and visible in the TUI
8. **Security**: Web access is properly authenticated and secured

## Future Considerations

While not part of the initial design, these considerations inform architectural decisions:

### Scalability:
- Support for multiple agent instances
- Distributed task execution
- Cloud sync for configuration and logs (optional)

### Extensibility:
- Plugin system for custom tools
- Integration API for third-party services
- Custom model fine-tuning support

### User Experience:
- Mobile-responsive web interface
- Voice interaction capabilities
- Natural language command interface
- Customizable themes and layouts

### Advanced Features:
- Scheduled task execution
- Event-driven automation
- Multi-modal interactions (text, images, audio)
- Collaborative features (multiple users working with same agent)

## Conclusion

This initial design establishes a solid foundation for Emrys as a Go-based console application with both local TUI access and web-based remote access. The architecture supports the core vision of a local, privacy-first AI assistant while providing flexibility for future enhancements. The focus on Nix installation and web browsing as initial goals provides concrete, achievable milestones that demonstrate the system's capabilities.
