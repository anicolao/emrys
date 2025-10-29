# Phase 1 Bootstrap Implementation

This document describes the Phase 1 bootstrap implementation for Emrys.

## Overview

Phase 1 implements package installation via nix-darwin, as specified in BOOTSTRAP.md. The implementation automatically detects and installs required packages, enables SSH server, and provides auto-login configuration for dedicated hardware.

## Architecture

The Phase 1 implementation is organized in the `internal/bootstrap` package:

- `phase1.go`: Core Phase 1 bootstrap functionality
- `phase1_test.go`: Comprehensive tests for bootstrap functionality

## Features

### Package Installation

Phase 1 installs the following packages via nix-darwin:

1. **ollama** - Local AI inference engine for running language models (latest stable version from nixpkgs)
2. **tmux** - Terminal multiplexer for persistent sessions (latest stable version from nixpkgs)
3. **go** - Go programming language for building Emrys (latest stable version from nixpkgs)
4. **jq** - JSON processing utility for configuration management (latest stable version from nixpkgs)

**Note:** All packages are installed from nixpkgs-unstable and are automatically kept up-to-date through nix-darwin. Ollama models are managed separately and can be downloaded after installation.

### SSH Server Configuration

Automatically enables the SSH server (Remote Login) via nix-darwin:
- SSH service enabled through `services.openssh.enable = true`
- Managed declaratively through nix-darwin configuration
- Users should configure SSH keys manually in `~/.ssh/authorized_keys`

### Auto-Login Configuration

Configures auto-login for the dedicated Mac Mini (enabled by default):
- Auto-login is enabled for unattended operation and power outage recovery
- Replaces `__EMRYS_USERNAME__` with actual username from configuration
- Designed for dedicated, physically secure hardware

## Usage

After nix-darwin is installed, running the `emrys` binary will:

1. Detect if Phase 1 packages are installed
2. Prompt the user to proceed with Phase 1 bootstrap if needed
3. Update nix-darwin configuration automatically
4. Run `darwin-rebuild switch` to apply changes
5. Verify all packages are installed correctly
6. Display next steps to the user

Example:

```bash
./emrys
```

Output when Phase 1 is needed:

```
╔════════════════════════════════════════╗
║           Emrys Setup                  ║
║  Your Personal AI Assistant on macOS  ║
╚════════════════════════════════════════╝

✓ nix-darwin is already installed!

⚠ Phase 1 bootstrap is not yet complete.

Would you like to run Phase 1 bootstrap now? (y/n): y

═══════════════════════════════════════
  Phase 1: Package Installation
═══════════════════════════════════════

Missing packages:
  - ollama
  - tmux
  - go
  - jq

Step 1: Updating nix-darwin configuration...
✓ Updated configuration at /Users/username/.nixpkgs/darwin-configuration.nix

Step 2: Applying configuration...
Applying nix-darwin configuration...
Note: This may take several minutes and will require sudo access

[darwin-rebuild output...]

✓ Configuration applied successfully

Step 3: Verifying installation...
Verifying package installation...
✓ All Phase 1 packages verified:
  - ollama     /run/current-system/sw/bin/ollama
  - tmux       /run/current-system/sw/bin/tmux
  - go         /run/current-system/sw/bin/go
  - jq         /run/current-system/sw/bin/jq

═══════════════════════════════════════
✓ Phase 1 Bootstrap Complete!
═══════════════════════════════════════
```

Phase 1 is now complete. Phase 2 will automatically configure Ollama service and download models.

## Phase 2: Ollama Setup and Configuration

Phase 2 implements Ollama service configuration and model management, as specified in BOOTSTRAP.md.

### Features

#### Ollama Service Configuration

Phase 2 configures Ollama to run as a persistent service:

1. **Launch Agent Creation** - Creates a launchd plist for macOS service management
2. **Automatic Startup** - Configures Ollama to start automatically on boot
3. **Keep Alive** - Ensures the service restarts if it crashes
4. **Logging** - Redirects stdout/stderr to log files for debugging

#### Model Management

Phase 2 downloads and configures the default model:

1. **Default Model** - llama3.2 (configurable via `DefaultModel` constant)
2. **Progress Indication** - Shows download progress in real-time
3. **Model Verification** - Tests the model with a simple inference query
4. **Integrity Check** - Ensures the model was downloaded correctly

#### API Health Checks

Phase 2 verifies Ollama API accessibility:

1. **Service Detection** - Checks if Ollama is running at http://localhost:11434
2. **API Testing** - Verifies API endpoints respond correctly
3. **Model Listing** - Confirms models can be queried via API

### Usage

After Phase 1 is complete, running the `emrys` binary will:

1. Detect if Phase 2 is complete
2. Prompt the user to proceed with Phase 2 bootstrap if needed
3. Create and load the Ollama launch agent
4. Start the Ollama service
5. Test API connectivity
6. Download the default model (llama3.2)
7. Verify model integrity
8. Display success confirmation

Example:

```bash
./emrys
```

Output when Phase 2 is needed:

```
╔════════════════════════════════════════╗
║           Emrys Setup                  ║
║  Your Personal AI Assistant on macOS  ║
╚════════════════════════════════════════╝

✓ nix-darwin is already installed!

✓ Phase 1 bootstrap is complete!

⚠ Phase 2 bootstrap is not yet complete.

Would you like to run Phase 2 bootstrap now? (y/n): y

═══════════════════════════════════════
  Phase 2: Ollama Setup
═══════════════════════════════════════

Step 1: Starting Ollama service...
✓ Created launch agent at /Users/username/Library/LaunchAgents/com.ollama.service.plist
Starting Ollama service...
✓ Ollama service started successfully

Step 2: Testing Ollama API...
Testing Ollama API connectivity...
✓ Ollama API is accessible and responding

Step 3: Downloading default model...
Downloading model 'llama3.2'...
Note: This may take several minutes depending on your internet connection

[download progress output...]

✓ Model 'llama3.2' downloaded successfully

Step 4: Verifying model...
Verifying model 'llama3.2'...
✓ Model 'llama3.2' verified successfully

═══════════════════════════════════════
✓ Phase 2 Bootstrap Complete!
═══════════════════════════════════════

Ollama is running at http://localhost:11434
Default model: llama3.2

Next steps:
  - Phase 3 will configure voice output
  - Phase 4 will set up the TUI application
```

### Testing

Phase 2 includes comprehensive tests in `phase2_test.go`:

- `TestIsOllamaRunning`: Tests service detection
- `TestIsModelInstalled`: Tests model detection
- `TestGetInstalledModels`: Tests model listing
- `TestIsPhase2Complete`: Tests Phase 2 completion detection
- `TestCreateOllamaLaunchAgent`: Tests launch agent creation with idempotency
- `TestTestOllamaAPI`: Tests API connectivity checking
- `TestVerifyModelIntegrity`: Tests model verification
- `TestDownloadModel`: Tests model download error handling
- `TestDefaultModelConstant`: Verifies default model is set
- `TestOllamaAPIURLConstant`: Verifies API URL is correct

Run tests with:

```bash
go test ./internal/bootstrap/... -v
```

### Ollama Service Management

The Ollama service is managed through macOS launchd:

**Start service:**
```bash
launchctl load ~/Library/LaunchAgents/com.ollama.service.plist
```

**Stop service:**
```bash
launchctl unload ~/Library/LaunchAgents/com.ollama.service.plist
```

**Check service status:**
```bash
launchctl list | grep ollama
```

**View logs:**
```bash
tail -f ~/Library/Logs/ollama.log
tail -f ~/Library/Logs/ollama-error.log
```

### Model Management

List installed models:
```bash
ollama list
```

Download additional models:
```bash
ollama pull mistral
ollama pull qwen2.5
```

Remove a model:
```bash
ollama rm model-name
```

## Phase 1 Implementation Details

### Package Detection

The `IsPhase1Complete()` function checks if all required packages are available in the system PATH using `exec.LookPath()`.

### Configuration Update

The `UpdateNixDarwinConfiguration()` function:
1. Reads the current nix-darwin configuration from `~/.nixpkgs/darwin-configuration.nix`
2. Checks if Phase 1 packages are already included (idempotent)
3. Adds Phase 1 packages to the `environment.systemPackages` section
4. Enables SSH server via `services.openssh.enable = true`
5. Adds auto-login configuration (enabled by default)
6. Extracts username from existing configuration or environment
7. Replaces username placeholder in auto-login configuration
8. Writes the updated configuration back to disk

### Configuration Application

The `ApplyConfiguration()` function:
1. Runs `darwin-rebuild switch --flake ~/.nixpkgs#emrys`
2. Handles sudo password prompts
3. Displays command output to the user

### Package Verification

The `VerifyPackageInstallation()` function:
1. Checks each package using `exec.LookPath()`
2. Displays the full path of each installed package
3. Returns an error if any packages are missing

## Testing

The implementation includes comprehensive tests:

- `TestIsPackageInstalled`: Tests package detection logic
- `TestGetMissingPackages`: Tests missing package identification
- `TestIsPhase1Complete`: Tests Phase 1 completion detection
- `TestUpdateNixDarwinConfiguration`: Tests configuration update logic with full idempotency testing

Run tests with:

```bash
go test ./internal/bootstrap/... -v
```

## Security Considerations

### SSH Server

- SSH server is enabled via nix-darwin's `services.openssh.enable = true`
- SSH access is managed declaratively through the nix-darwin configuration
- Configure SSH key-based authentication manually in `~/.ssh/authorized_keys`
- For additional security, password authentication can be disabled in `/etc/ssh/sshd_config`
- Remote Login will be enabled on system activation

### Auto-Login

- Auto-login is enabled by default for dedicated, physically secure hardware
- Designed for unattended operation and automatic recovery from power outages
- Username is automatically extracted from the existing nix-darwin configuration
- Should only be used on physically secure Mac Mini systems
- May have implications for FileVault encryption (see BOOTSTRAP.md Phase 7)

## Next Steps

After Phase 1 is complete, the next phases are:

- **Phase 2**: Ollama setup and configuration (model download, service configuration)
- **Phase 3**: Voice output configuration (Jamie voice installation and testing)
- **Phase 4**: TUI application development using Bubbletea
- **Phase 5**: Tmux session management
- **Phase 6**: Auto-start configuration
- **Phase 7**: Auto-login testing and FileVault compatibility
- **Phase 8**: Power outage recovery testing

## Troubleshooting

### Phase 1 Issues

#### Packages not found after installation

If packages are not found in PATH after installation:
1. Restart your terminal to source the updated environment
2. Check if nix-daemon is running: `ps aux | grep nix-daemon`
3. Manually source Nix: `. /nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh`

#### darwin-rebuild fails

If `darwin-rebuild` fails:
1. Check the error message for specific issues
2. Verify the configuration syntax: `nix flake check ~/.nixpkgs#emrys`
3. Try rebuilding with verbose output: `darwin-rebuild switch --flake ~/.nixpkgs#emrys --show-trace`

#### Permission errors

If you encounter permission errors:
1. Ensure you have sudo access
2. Check that nix-darwin is properly installed: `which darwin-rebuild`
3. Verify Nix is properly installed: `which nix`

### Phase 2 Issues

#### Ollama service won't start

If the Ollama service fails to start:
1. Check if ollama binary is in PATH: `which ollama`
2. Try starting manually: `ollama serve`
3. Check the logs: `cat ~/Library/Logs/ollama-error.log`
4. Verify the launch agent exists: `ls ~/Library/LaunchAgents/com.ollama.service.plist`
5. Try reloading the launch agent: `launchctl unload ~/Library/LaunchAgents/com.ollama.service.plist && launchctl load ~/Library/LaunchAgents/com.ollama.service.plist`

#### Model download fails

If model download fails:
1. Check your internet connection
2. Verify Ollama is running: `curl http://localhost:11434`
3. Check available disk space: `df -h`
4. Try downloading manually: `ollama pull llama3.2`
5. Check for rate limiting or network issues

#### API connectivity issues

If Ollama API is not accessible:
1. Verify the service is running: `ps aux | grep ollama`
2. Check if the port is in use: `lsof -i :11434`
3. Test with curl: `curl http://localhost:11434`
4. Check firewall settings
5. Review error logs: `cat ~/Library/Logs/ollama-error.log`

#### Model verification fails

If model verification fails after download:
1. List installed models: `ollama list`
2. Try running the model manually: `ollama run llama3.2`
3. Check for corrupted downloads: remove and re-download the model
4. Ensure sufficient RAM is available (models require significant memory)
5. Check system logs for GPU/Metal-related errors on Apple Silicon

## Configuration File Locations

### Phase 1
- nix-darwin configuration: `~/.nixpkgs/darwin-configuration.nix`
- Flake configuration: `~/.nixpkgs/flake.nix`
- System configuration: `/etc/nix/nix.conf`
- SSH configuration: Enabled via nix-darwin, keys in `~/.ssh/authorized_keys`

### Phase 2
- Ollama launch agent: `~/Library/LaunchAgents/com.ollama.service.plist`
- Ollama logs: `~/Library/Logs/ollama.log`
- Ollama error logs: `~/Library/Logs/ollama-error.log`
- Ollama models: `~/.ollama/models/`
