# Phase 1 Bootstrap Implementation

This document describes the Phase 1 bootstrap implementation for Emrys.

## Overview

Phase 1 implements package installation via nix-darwin, as specified in BOOTSTRAP.md. The implementation automatically detects and installs required packages and provides auto-login configuration for dedicated hardware.

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

### Auto-Login Configuration

Configures auto-login for the dedicated Mac Mini (enabled by default):
- Auto-login is enabled for unattended operation and power outage recovery
- Replaces `__EMRYS_USERNAME__` with actual username from configuration
- Designed for dedicated, physically secure hardware

**Note on SSH:** SSH configuration on macOS should be managed through System Preferences > Sharing > Remote Login, or via `sudo systemsetup -setremotelogin on`. nix-darwin does not provide the same SSH service configuration options as NixOS.

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

## Implementation Details

### Package Detection

The `IsPhase1Complete()` function checks if all required packages are available in the system PATH using `exec.LookPath()`.

### Configuration Update

The `UpdateNixDarwinConfiguration()` function:
1. Reads the current nix-darwin configuration from `~/.nixpkgs/darwin-configuration.nix`
2. Checks if Phase 1 packages are already included (idempotent)
3. Adds Phase 1 packages to the `environment.systemPackages` section
4. Adds auto-login configuration (enabled by default)
5. Extracts username from existing configuration or environment
6. Replaces username placeholder in auto-login configuration
7. Writes the updated configuration back to disk

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

### Auto-Login

- Auto-login is enabled by default for dedicated, physically secure hardware
- Designed for unattended operation and automatic recovery from power outages
- Username is automatically extracted from the existing nix-darwin configuration
- Should only be used on physically secure Mac Mini systems

### SSH Access

- SSH configuration on macOS should be managed through System Preferences or systemsetup command
- For secure remote access, enable SSH in System Preferences > Sharing > Remote Login
- Configure SSH key-based authentication manually for security
- Disable password authentication in `/etc/ssh/sshd_config` if needed
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

### Packages not found after installation

If packages are not found in PATH after installation:
1. Restart your terminal to source the updated environment
2. Check if nix-daemon is running: `ps aux | grep nix-daemon`
3. Manually source Nix: `. /nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh`

### darwin-rebuild fails

If `darwin-rebuild` fails:
1. Check the error message for specific issues
2. Verify the configuration syntax: `nix flake check ~/.nixpkgs#emrys`
3. Try rebuilding with verbose output: `darwin-rebuild switch --flake ~/.nixpkgs#emrys --show-trace`

### Permission errors

If you encounter permission errors:
1. Ensure you have sudo access
2. Check that nix-darwin is properly installed: `which darwin-rebuild`
3. Verify Nix is properly installed: `which nix`

## Configuration File Locations

- nix-darwin configuration: `~/.nixpkgs/darwin-configuration.nix`
- Flake configuration: `~/.nixpkgs/flake.nix`
- System configuration: `/etc/nix/nix.conf`
- SSH configuration: Managed through macOS System Preferences or `/etc/ssh/sshd_config`
