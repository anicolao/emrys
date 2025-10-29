# Emrys Bootstrap Process Design

## Overview

This document specifies the complete bootstrap process for Emrys, from the current state (nix-darwin installation complete) to a fully operational system where Emrys runs as a TUI application in tmux with voice output capabilities, auto-starts on boot, and can recover automatically from power outages.

## Current State

The emrys binary has successfully:
- Detected and installed Nix (using the Lix installer)
- Installed nix-darwin with a basic flake-based configuration
- Created initial configuration files at `~/.nixpkgs/darwin-configuration.nix` and `~/.nixpkgs/flake.nix`
- Exited with instructions for the user to restart their terminal

## Target State

The fully bootstrapped system will have:
- Ollama installed and running locally for AI inference
- Tmux installed and configured for session management
- Emrys TUI application running persistently in a tmux session
- Voice output capability using macOS `say` command with Jamie (Premium) voice configured via nix-darwin
- SSH server configured via nix-darwin with user's public key for remote access
- Auto-login configured via nix-darwin for the dedicated user account
- Emrys set to auto-start on login as a login item
- Automatic recovery to working state after power outages
- All necessary packages, dependencies, and system settings managed through nix-darwin

## Bootstrap Phases

### Phase 1: Package Installation via Nix-Darwin

#### Objectives
Extend the nix-darwin configuration to include all packages required for Emrys operation.

#### Required Packages
- **ollama**: Local AI inference engine for running language models
- **tmux**: Terminal multiplexer for persistent sessions
- **go**: For building the Emrys binary if needed during development
- **jq**: JSON processing utility for configuration management
- Any additional utility packages identified during implementation

#### Configuration Strategy
The emrys binary should be enhanced to:
- Detect whether the bootstrap packages are installed
- Update the nix-darwin configuration to include required packages
- Configure SSH server settings via nix-darwin (enable service, disable password auth)
- Configure auto-login via nix-darwin settings
- Configure Jamie (Premium) voice installation via nix-darwin
- Trigger `darwin-rebuild switch` to apply the configuration
- Verify successful installation of each required package
- Handle installation failures gracefully with clear error messages

#### User Experience
After running the enhanced emrys binary, the user should see:
- Clear progress indicators for each package being installed
- Estimated time remaining for the installation process
- Success confirmation for each package
- Instructions for what happens next

### Phase 2: Ollama Setup and Configuration

#### Objectives
Initialize Ollama with appropriate models and ensure it runs as a service.

#### Ollama Service Configuration
- Determine the appropriate method for running Ollama persistently
- Consider using launchd (macOS native service manager) for Ollama daemon
- Configure Ollama to start automatically on boot
- Set appropriate resource limits for Ollama process

#### Model Download and Installation
- Select an appropriate default model for Emrys (e.g., llama3.2, mistral, qwen2.5)
- Implement model download with progress indication
- Verify model integrity after download
- Configure model location and cache settings
- Plan for future model updates and management

#### API Accessibility
- Ensure Ollama API is accessible at `http://localhost:11434`
- Verify API responsiveness with test queries
- Handle cases where Ollama service fails to start
- Implement health check mechanism for Ollama service

### Phase 3: Voice Output Configuration

#### Objectives
Set up macOS text-to-speech capabilities with the Jamie (Premium) voice.

#### Voice Installation
- Configure Jamie (Premium) voice installation via nix-darwin
- Verify voice installation and availability
- Set Jamie (Premium) as the default voice for Emrys

#### Voice Output Integration
- Create a voice output module within Emrys
- Implement queuing system for voice messages to prevent overlap
- Add volume control and speech rate configuration
- Handle cases where voice output is not available (fallback to text-only)

#### User Controls
- Allow user to enable/disable voice output
- Provide controls for adjusting speech parameters
- Implement "quiet hours" functionality if desired
- Add voice output testing utility that speaks a confirmation phrase when voice is working

### Phase 4: TUI Application Development

#### Objectives
Build the Terminal User Interface application using Bubbletea framework.

#### Core TUI Components
- Main application view with status dashboard
- Command input interface for user interaction
- Log viewer for system and AI agent activities
- Task monitor showing active and completed operations
- Configuration interface for settings management

#### Status Display Elements
- Ollama service status and current model
- System resource usage (CPU, memory, network)
- Active tmux session information
- Voice output status
- Network connectivity status
- Last successful operation timestamp

#### Command Interface
- Command history and recall functionality
- Auto-completion for common commands
- Help system with available commands
- Error handling and user feedback

#### Visual Design
- Use Lipgloss for consistent styling and theming
- Implement responsive layout that adapts to terminal size
- Clear visual hierarchy for different information types
- Color-coded status indicators (green=good, yellow=warning, red=error)
- Support for both light and dark terminal themes

### Phase 5: Tmux Session Management

#### Objectives
Configure tmux for persistent Emrys session management.

#### Tmux Configuration
- Create custom tmux configuration optimized for Emrys
- Set appropriate session name (e.g., "emrys-main")
- Configure status bar with Emrys-relevant information
- Set up mouse support for easier interaction
- Configure scrollback buffer size
- Enable UTF-8 support for proper rendering

#### Session Persistence
- Implement automatic session creation on first run
- Handle session reattachment when Emrys is already running
- Prevent duplicate Emrys instances
- Graceful handling of session termination and restart

#### Remote Access Support
- Ensure SSH compatibility for remote access
- Support multiple simultaneous viewers
- Implement read-only attachment capability
- Document SSH setup and connection procedures

#### SSH Server Configuration
- Configure SSH server via nix-darwin to enable Remote Login
- Set SSH server for optimal security (disable password auth, use keys only)
- Test SSH connectivity on local network

#### SSH Key Setup
- Look for `id_rsa.pub` file in the same directory as the emrys binary
- If found, use it as the default and prompt user for confirmation
- If not found or user declines, prompt for file path to their public SSH key
- Validate SSH public key format before accepting
- Add the public key to `~/.ssh/authorized_keys` with correct permissions
- Set appropriate permissions on SSH directory and files (`chmod 700 ~/.ssh`, `chmod 600 ~/.ssh/authorized_keys`)
- Verify the public key was added correctly
- Provide instructions for testing SSH access from remote machine

### Phase 6: Auto-Start Configuration

#### Objectives
Configure the system to automatically start Emrys on boot and user login.

#### Login Items Setup
- Create a launchd plist for Emrys auto-start
- Configure the plist to run after user login
- Set appropriate working directory and environment variables
- Ensure tmux session is created automatically
- Handle cases where tmux or Emrys fail to start

#### Launch Agent Configuration
- Place launch agent plist in appropriate location (`~/Library/LaunchAgents/`)
- Set correct permissions on the plist file
- Configure stdout/stderr logging for debugging
- Implement retry logic for failed starts
- Set resource limits and watchdog timers

#### Startup Sequence
- Wait for network availability before starting Ollama
- Ensure Ollama is running before starting Emrys TUI
- Verify all dependencies are available
- Graceful degradation if optional components fail
- Log all startup events for troubleshooting

### Phase 7: Auto-Login Configuration

#### Objectives
Enable automatic login for the dedicated Mac Mini user account to ensure full system availability after power outages.

#### Security Considerations
- Document security implications of auto-login on dedicated hardware
- Recommend physical security measures for the Mac Mini
- Consider encrypted disk requirements
- Advise on network security (firewall rules, SSH key-only access)

#### Implementation Strategy
- Configure auto-login via nix-darwin settings
- Verify auto-login configuration is correctly applied
- Test auto-login with system restart
- Document how to disable auto-login if needed

#### FileVault Compatibility
- Address FileVault encryption compatibility with auto-login
- Document limitations if FileVault is enabled
- Provide alternative approaches for encrypted systems
- Consider security vs. convenience trade-offs

### Phase 8: Power Outage Recovery

#### Objectives
Ensure the system returns to operational state automatically after power loss.

#### System-Level Recovery
- Verify Mac Mini BIOS/NVRAM settings for auto-power-on after power loss
- Configure energy saver settings to restart automatically
- Test power loss scenario and recovery
- Document hardware-specific configuration steps

#### Service Recovery
- Ensure all launchd services are configured for automatic restart
- Implement health checks that trigger service restart if needed
- Add monitoring for Ollama service availability
- Configure tmux to auto-create sessions if missing

#### State Preservation
- Implement checkpoint mechanism for Emrys state
- Save critical data periodically to survive crashes
- Restore previous state on restart when possible
- Log power loss events and recovery actions

#### Monitoring and Alerting
- Implement system health monitoring
- Detect when services fail to recover automatically
- Consider external notification methods (email, push notifications)
- Log all recovery attempts for post-mortem analysis

### Phase 9: Configuration Management

#### Objectives
Provide a cohesive configuration system for all Emrys components.

#### Configuration File Structure
- Define YAML or JSON schema for Emrys configuration
- Include sections for: Ollama settings, voice output, TUI preferences, logging, auto-start behavior
- Support user overrides of default settings
- Validate configuration on load with helpful error messages

#### Configuration Location
- Store configuration in standard location (e.g., `~/.config/emrys/config.yaml`)
- Create configuration directory automatically if missing
- Support environment variable overrides for advanced users
- Document all configuration options

#### Configuration UI
- Provide TUI interface for common configuration changes
- Allow editing configuration file directly with validation
- Implement live reload of configuration changes where possible
- Show current configuration in status view

### Phase 10: Error Handling and Recovery

#### Objectives
Make the system resilient to common failure modes.

#### Failure Scenarios
- Ollama service fails to start or crashes
- Network connectivity lost
- Disk space exhausted
- Model files corrupted
- Tmux session unexpectedly terminated
- TUI application crashes
- Voice output unavailable

#### Recovery Strategies
- Automatic service restart with exponential backoff
- Graceful degradation (continue with reduced functionality)
- Clear error messages displayed in TUI
- Voice announcements of critical errors (when available)
- Detailed logging of all errors for debugging
- User notification when manual intervention required

#### Diagnostic Tools
- Built-in health check command
- System diagnostic report generation
- Log collection utility for troubleshooting
- Network connectivity testing
- Ollama API testing utility

### Phase 11: Testing and Validation

#### Objectives
Verify the complete bootstrap process works reliably.

#### Test Scenarios
- Fresh Mac Mini installation from scratch
- Bootstrap after system updates
- Recovery from power outage (simulated)
- Recovery from individual service failures
- SSH key configuration from file (id_rsa.pub in binary directory or user-specified path)
- Voice output functionality
- Configuration changes and reloads
- Disk space exhaustion recovery

#### Validation Checklist
- All packages installed correctly via nix-darwin
- SSH server enabled and configured via nix-darwin
- Ollama service running and responsive
- Default model downloaded and functional
- User's SSH public key correctly installed in authorized_keys
- SSH key-based authentication working from remote machine
- TUI application starts and displays correctly
- Tmux session persists across disconnections
- Voice output works with Jamie (Premium) voice
- Auto-login functions after restart
- Emrys auto-starts on login
- System recovers after simulated power loss
- SSH remote access to tmux session works correctly
- Error conditions handled gracefully

#### Documentation Requirements
- User guide for bootstrap process
- Troubleshooting guide for common issues
- Configuration reference documentation
- Development guide for extending Emrys
- Security best practices document

## Implementation Workflow

### Step 1: Enhance Emrys Binary
Modify the main emrys binary to include bootstrap functionality:
- Automatic bootstrap after initial install
- Implement SSH public key collection and configuration
- Implement package installation orchestration
- Add Ollama setup and model download
- Integrate voice setup with confirmation phrase
- Create launch agent configuration

### Step 2: Develop TUI Application
Build the Bubbletea-based TUI:
- Create basic application structure
- Implement status views
- Add command interface
- Integrate with Ollama
- Add voice output support

### Step 3: Configure Auto-Start
Set up persistent operation:
- Create launchd plists
- Configure tmux session management
- Implement startup scripts
- Test auto-start functionality

### Step 4: System Hardening
Ensure reliability:
- Implement error recovery
- Add health monitoring
- Create diagnostic tools
- Comprehensive logging

### Step 5: Documentation
Create user-facing documentation:
- Installation guide
- User manual
- Troubleshooting guide
- Configuration reference

## Bootstrap Command Flow

When the user runs the enhanced emrys binary after initial nix-darwin installation:

1. Welcome message explaining the bootstrap process
2. Pre-flight checks (system requirements, disk space, network)
3. Prompt user for confirmation to proceed
4. Collect user's SSH public key for remote access
5. Update nix-darwin configuration with required packages, SSH server, auto-login, and Jamie voice
6. Run darwin-rebuild to install packages and apply configuration (with progress indication)
7. Initialize Ollama and download default model
8. Verify Jamie (Premium) voice installation and speak a confirmation phrase
9. Install user's public key for SSH access
10. Build and install Emrys TUI binary
11. Create tmux configuration
12. Set up launch agent for auto-start
13. Offer to start Emrys immediately or on next login
14. Display success message with next steps and SSH access instructions
15. Launch Emrys in tmux (if user confirmed immediate start)

## Success Criteria

The bootstrap process is successful when:

1. User can run a single command (the emrys binary) to go from nix-darwin installation to fully operational Emrys
2. All required packages are installed automatically via nix-darwin
3. Ollama is running with an appropriate model downloaded
4. Emrys TUI application launches in tmux and displays status
5. Voice output is functional and uses Jamie (Premium) voice
6. SSH server is configured via nix-darwin and user's public key is correctly installed for remote access
7. System auto-starts Emrys on boot/login
8. System recovers automatically after power outages
9. Remote SSH access to tmux session works correctly using the configured key
10. Clear error messages and recovery procedures for any failures
11. Complete documentation is available for users

## Security Considerations

### Physical Security
- Mac Mini should be in a physically secure location
- Auto-login requires physical access control
- Consider hardware security features (T2 chip, etc.)

### Network Security
- SSH access should use key-based authentication only
- Consider firewall rules limiting SSH access
- Monitor for unauthorized access attempts
- Keep system updated with security patches

### Data Security
- Sensitive data should be encrypted at rest if possible
- Consider implications of auto-login for FileVault
- API keys and credentials securely stored
- Regular backups of important data

### Service Security
- Ollama API accessible only on localhost by default
- TUI session access controls via tmux permissions
- Log files protected with appropriate permissions
- Service accounts with minimal required privileges

## Future Enhancements

Capabilities that could be added after initial bootstrap:

- Multiple model support with easy switching
- Custom voice selection and configuration
- Web interface for remote management (alternative to SSH)
- Mobile app for status monitoring
- Integration with HomeKit for Mac Mini control
- Backup and restore functionality for Emrys state
- Multi-Mac deployment with shared configuration
- Performance monitoring and optimization tools
- Custom agent capabilities and plugins
- Integration with external services (email, calendar, etc.)

## Conclusion

This bootstrap design provides a comprehensive path from the current nix-darwin installation to a fully operational, resilient Emrys system. The phased approach allows for incremental implementation and testing, while the focus on automation and error recovery ensures a reliable system that can operate unattended and recover from power outages automatically.

The design prioritizes:
- **Simplicity**: Single command to bootstrap from current state
- **Reliability**: Automatic recovery from common failures
- **Usability**: Clear progress indication and error messages
- **Maintainability**: Declarative configuration via nix-darwin
- **Security**: Appropriate controls while enabling auto-operation

Implementation of this design will result in a Mac Mini that can truly operate as an autonomous AI assistant, automatically recovering from power outages and ready to serve the user at all times.
