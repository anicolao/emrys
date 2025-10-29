{
  description = "Emrys nix-darwin system configuration";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    nix-darwin.url = "github:LnL7/nix-darwin";
    nix-darwin.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = inputs@{ self, nix-darwin, nixpkgs }:
  let
    # Support both Apple Silicon and Intel Macs
    systems = [ "aarch64-darwin" "x86_64-darwin" ];
    
    # Helper to create configurations for each system
    forAllSystems = nixpkgs.lib.genAttrs systems;
  in
  {
    darwinConfigurations = forAllSystems (system:
      nix-darwin.lib.darwinSystem {
        inherit system;
        modules = [ ./darwin-configuration.nix ];
      }
    );
    
    # Default configuration (will be used when system is auto-detected)
    darwinConfigurations."emrys" = nix-darwin.lib.darwinSystem {
      system = builtins.currentSystem;
      modules = [ ./darwin-configuration.nix ];
    };
  };
}
