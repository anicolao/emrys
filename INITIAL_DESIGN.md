# Emrys Initial Design

## Overview

Emrys will be built as a Go-based console application with a Terminal User Interface (TUI) using the [Bubbletea](https://github.com/charmbracelet/bubbletea) framework. The TUI will run in a screen or tmux session on the dedicated Mac Mini, accessible remotely via SSH. This provides both local console access and remote monitoring/control through standard terminal multiplexing tools.

## Architecture

### 1. Core Application Layer (Go)

The core of Emrys is a Go application that provides:
- TUI console interface using Bubbletea
- Runs persistently in a screen/tmux session
- Agent orchestration and task execution

#### Components:
- **Main Application**: Entry point and lifecycle management
- **TUI Engine**: Bubbletea-based interactive terminal interface
- **Session Manager**: Integration with screen/tmux for persistence
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

### 3. Remote Access via SSH

Remote access to Emrys is accomplished through standard SSH and terminal multiplexing:

#### Architecture:
- **SSH Server**: macOS built-in SSH server for secure remote access
- **Terminal Multiplexer**: screen or tmux session hosting the Emrys TUI
- **Session Persistence**: TUI remains running even when disconnected
- **Multi-User Support**: Multiple users can attach to view (read-only via screen/tmux options)

#### Access Patterns:
- **Local Access**: Direct execution on the Mac Mini console
- **Remote Interactive**: SSH into Mac Mini and attach to the session
- **Remote Monitoring**: SSH with read-only attachment to observe activity
- **Detached Operation**: TUI continues running when no users are attached

#### Benefits:
- Standard, well-understood SSH security model
- No custom web infrastructure to maintain
- Native terminal experience with full color and formatting support
- Leverages existing screen/tmux capabilities for session management
- Simple firewall configuration (SSH port only)

### 4. Initial Capability Goals

The first iteration of Emrys will demonstrate core capabilities through two primary tasks:

#### Goal 1: Self-Bootstrap via Nix-Darwin Installation
- **Objective**: Prove Emrys can modify its own system environment
- **Tasks**:
  - Detect current macOS environment
  - Download and install Nix and nix-darwin
  - Verify Nix installation and basic functionality
  - Configure nix-darwin for system control
  
- **Technical Requirements**:
  - Shell command execution capability
  - Installation script management
  - Environment variable handling
  - Post-installation verification

#### Goal 2: System Configuration via Nix-Darwin
- **Objective**: Demonstrate controlled system configuration and software management
- **Tasks**:
  - Write nix-darwin configuration files to define system state
  - Use configuration to install and manage system packages
  - Apply declarative system settings and preferences
  - Update and maintain system configuration
  - Clean up unused packages

- **Technical Requirements**:
  - Configuration file generation and management
  - nix-darwin command orchestration (darwin-rebuild)
  - Package search and selection
  - System state monitoring and verification
  - Configuration rollback capabilities

#### Goal 3: Web Browser Access via ChatGPT Atlas
- **Objective**: Prove Emrys can interact with the web using an AI-enhanced browser
- **Tasks**:
  - Install ChatGPT Atlas (Chromium fork with ChatGPT integration) via Nix
  - Launch and control ChatGPT Atlas instances
  - Leverage Atlas's built-in AI capabilities for web navigation
  - Demonstrate AI-assisted browsing and task completion

- **Technical Requirements**:
  - Browser process management for ChatGPT Atlas
  - Investigation of headless mode support (may or may not be available)
  - Integration with Atlas's command/control interface
  - Utilization of Atlas's ChatGPT-powered context understanding
  - Screenshot and debugging capabilities
  
- **Notes**:
  - ChatGPT Atlas provides AI-assisted browsing capabilities, reducing the complexity of web automation
  - Atlas's Chromium foundation ensures compatibility with modern web standards
  - The AI integration in Atlas should make Emrys more capable at complex web tasks

## Component Communication Flow

```
┌─────────────────────────────────────────────────────────────┐
│                        User Access                          │
├─────────────────┬───────────────────────────────────────────┤
│   Local TTY     │     SSH + screen/tmux (Remote)            │
└────────┬────────┴──────────────────┬────────────────────────┘
         │                           │
         │                           │ Terminal
         │                           │ Session
         ▼                           ▼
┌─────────────────────────────────────────────────────────────┐
│           screen/tmux Session (Persistent)                  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │             Bubbletea TUI Engine                      │  │
│  │  ┌──────────────┬──────────────┬──────────────────┐  │  │
│  │  │   Console    │    Status    │   Task Monitor   │  │  │
│  │  │   Interface  │   Dashboard  │   & Logs         │  │  │
│  │  └──────────────┴──────────────┴──────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
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
│  │   Shell      │  AppleScript │  ChatGPT Atlas       │    │
│  │   Executor   │  Executor    │  (AI Browser)        │    │
│  └──────────────┴──────────────┴──────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## Technology Stack

### Backend (Go)
- **Framework**: Standard library + Bubbletea
- **TUI Library**: [Bubbletea](https://github.com/charmbracelet/bubbletea)
- **TUI Components**: [Bubbles](https://github.com/charmbracelet/bubbles) (text input, viewport, etc.)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) (terminal styling)

### Remote Access
- **SSH**: macOS built-in SSH server (OpenSSH)
- **Terminal Multiplexer**: screen or tmux for session persistence
- **Authentication**: SSH keys and standard SSH authentication

### AI/Agent
- **LLM Runtime**: Ollama (local inference)
- **Models**: Open-source models (Llama 3, Mistral, Qwen, etc.)
- **Browser**: ChatGPT Atlas (Chromium fork with ChatGPT integration)
- **Browser Automation**: Integration with ChatGPT Atlas's control interface

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

### Remote SSH Interaction:
1. User connects via SSH to Mac Mini
2. User attaches to screen/tmux session running Emrys
3. User inputs command in terminal (same as local)
4. Bubbletea handles input event
5. Command dispatched to Agent Core
6. Agent processes and executes via Tool Layer
7. Results rendered back to TUI
8. Display updated in terminal (visible over SSH)

### Read-Only Monitoring:
1. User connects via SSH to Mac Mini
2. User attaches to screen/tmux session in read-only mode
   - screen: `screen -x` allows shared viewing
   - tmux: `tmux attach -r` for read-only access
3. TUI output visible but input is restricted
4. Multiple users can monitor simultaneously

## Security Considerations

### SSH Access:
- **Authentication**: SSH key-based authentication (no password access)
- **Encryption**: All traffic encrypted via SSH protocol
- **Firewall**: SSH port (22 or custom) properly configured
- **Access Control**: User accounts and SSH authorized_keys management
- **Audit Logging**: SSH connection logs and session activity tracking

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
# Session settings
session.multiplexer: tmux  # or "screen"
session.name: emrys
session.auto_create: true

# Agent settings
agent.model: llama3
agent.ollama_url: http://localhost:11434
agent.max_concurrent_tasks: 3

# Browser settings
browser.type: chatgpt-atlas
browser.headless: auto  # Use headless if available, otherwise GUI
browser.executable: /path/to/chatgpt-atlas

# Logging
log.level: info
log.file: /var/log/emrys/emrys.log
log.max_size_mb: 100

# SSH/Access (managed by macOS)
# SSH configuration in /etc/ssh/sshd_config
# Authorized keys in ~/.ssh/authorized_keys
```

## Development Phases

### Phase 1: Foundation (Current Focus)
- Basic Go application structure
- Simple Bubbletea TUI with command input
- Static content display
- Logging system
- Integration with screen/tmux for session persistence

### Phase 2: Remote Access Setup
- SSH server configuration on Mac Mini
- screen/tmux session management
- Auto-start on boot configuration
- Session recovery and reconnection handling

### Phase 3: Agent Integration
- Ollama LLM connection
- Basic command execution
- Shell tool implementation
- Task logging and display

### Phase 4: Initial Capabilities
- Nix installation workflow
- Package management via Nix
- ChatGPT Atlas installation and integration
- AI-assisted web browsing demonstration

### Phase 5: Polish & Enhancement
- Session sharing and read-only access refinement
- Improved error handling
- Comprehensive logging
- Configuration management
- Auto-recovery from crashes

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
│   ├── session/
│   │   ├── manager.go              # screen/tmux session management
│   │   └── recovery.go             # Session recovery and persistence
│   ├── tools/
│   │   ├── shell.go                # Shell command execution
│   │   ├── nix.go                  # Nix package management
│   │   ├── atlas.go                # ChatGPT Atlas browser control
│   │   └── applescript.go          # AppleScript execution
│   └── llm/
│       ├── client.go               # Ollama client
│       └── models.go               # Model management
├── scripts/
│   ├── setup-ssh.sh                # SSH configuration helper
│   ├── start-session.sh            # Launch Emrys in screen/tmux
│   └── install-nix.sh              # Nix installation script
├── config/
│   └── config.yaml                 # Configuration file
├── go.mod                          # Go module definition
├── go.sum                          # Go dependencies
└── README.md                       # Project documentation
```

## Success Criteria

The initial design succeeds when:

1. **TUI Functionality**: Bubbletea interface runs smoothly with command input/output
2. **Session Persistence**: TUI runs reliably in screen/tmux and survives disconnections
3. **SSH Access**: Remote users can connect and interact with the TUI via SSH
4. **Read-Only Monitoring**: Multiple users can observe TUI activity without interfering
5. **Nix Installation**: Emrys can successfully install Nix on a fresh macOS system
6. **Package Management**: Emrys can install and manage packages using Nix
7. **ChatGPT Atlas Integration**: Emrys can launch and control ChatGPT Atlas for web tasks
8. **AI-Assisted Browsing**: Atlas's ChatGPT integration enhances web interaction capabilities
9. **Logging**: All operations are logged and visible in the TUI
10. **Security**: SSH access is properly secured with key-based authentication

## Conclusion

This initial design establishes a solid foundation for Emrys as a Go-based console application with persistent TUI access via screen/tmux and secure SSH remote access. The architecture supports the core vision of a local, privacy-first AI assistant while providing flexibility for future enhancements. The integration of ChatGPT Atlas as an AI-enhanced browser provides powerful web interaction capabilities that leverage existing AI technology. The focus on Nix installation and AI-assisted web browsing as initial goals provides concrete, achievable milestones that demonstrate the system's capabilities.
