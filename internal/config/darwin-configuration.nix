{ config, pkgs, lib, ... }:

{
  # Basic nix-darwin configuration for Emrys
  # This is a minimal configuration that will be used during initial setup
  
  # Set the host platform (will be auto-detected from the system)
  nixpkgs.hostPlatform = lib.mkDefault "aarch64-darwin";
  
  # Enable nix-darwin
  system.stateVersion = 5;

  # Enable nix flakes and new nix command
  nix.settings.experimental-features = [ "nix-command" "flakes" ];

  # Auto upgrade nix package and the daemon service
  services.nix-daemon.enable = true;

  # Enable Touch ID for sudo
  security.pam.enableSudoTouchIdAuth = true;

  # Basic system packages
  environment.systemPackages = with pkgs; [
    vim
    git
    curl
    wget
  ];

  # System defaults
  system.defaults = {
    dock.autohide = true;
    finder.AppleShowAllExtensions = true;
    NSGlobalDomain.AppleShowAllExtensions = true;
  };

  # Auto-optimize nix store
  nix.optimise.automatic = true;

  # Garbage collection
  nix.gc = {
    automatic = true;
    interval = { Weekday = 0; Hour = 0; Minute = 0; };
    options = "--delete-older-than 30d";
  };
}
