{
  description = "Emrys nix-darwin system configuration";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    nix-darwin.url = "github:LnL7/nix-darwin";
    nix-darwin.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = inputs@{ self, nix-darwin, nixpkgs }:
  {
    # Single configuration that auto-detects the system (Apple Silicon or Intel)
    darwinConfigurations."emrys" = nix-darwin.lib.darwinSystem {
      system = builtins.currentSystem;
      modules = [ ./darwin-configuration.nix ];
    };
  };
}
