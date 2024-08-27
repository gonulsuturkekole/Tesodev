package config

import "fmt"

type AppConfig struct {
	SecretKey string
}

var cfg = map[string]AppConfig{
	"prod": {
		SecretKey: "079c9b74-24a7-4341-ae15-5b7a42f8bfb7",
	},
	"qa": {
		SecretKey: "079c9b74-24a7-4341-ae15-5b7a42f8bfb7",
	},
	"dev": {
		SecretKey: "079c9b74-24a7-4341-ae15-5b7a42f8bfb7",
	},
}

func GetAppConfig(env string) *AppConfig {
	config, exists := cfg[env]
	if !exists {
		panic(fmt.Sprintf("Config for environment '%s' does not exist", env))
	}
	return &config
}
