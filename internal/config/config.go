package config

import _ "embed"

// DefaultNixDarwinConfig contains the embedded default nix-darwin configuration
//
//go:embed darwin-configuration.nix
var DefaultNixDarwinConfig string
