package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

// Env holds the configuration from the environment
var Env struct {
	ListenAddress string `envconfig:"LISTEN_ADDRESS" default:":80"`
	SitePassword  string `envconfig:"SITE_PASSWORD" required:"true"`
	DatabasePath  string `envconfig:"DB_PATH" default:"parts.db"`
	MouserAPIKey  string `envconfig:"MOUSER_API_KEY"`
	DigikeyAPIKey string `envconfig:"DIGIKEY_API_KEY"`
}

// Process parses the environment
func Process() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		envconfig.Usage("", &Env)
		os.Exit(0)
	}
	envconfig.MustProcess("", &Env)
}
