package config

import "fmt"

type DbConfig struct {
	DBName    string
	ColName   string
	SecretKey string
}

type ConsumerConfig struct {
	Port     string
	DbConfig DbConfig
}

var cfgs = map[string]ConsumerConfig{
	"prod": {
		Port: ":8003",
		DbConfig: DbConfig{
			DBName:    "tesodev",
			ColName:   "finance",
			SecretKey: "079c9b74-24a7-4341-ae15-5b7a42f8bfb7",
		},
	},
	"qa": {
		Port: ":8003",
		DbConfig: DbConfig{
			DBName:    "tesodev",
			ColName:   "finance",
			SecretKey: "079c9b74-24a7-4341-ae15-5b7a42f8bfb7",
		},
	},
	"dev": {
		Port: ":8003",
		DbConfig: DbConfig{
			DBName:    "tesodev",
			ColName:   "finance",
			SecretKey: "079c9b74-24a7-4341-ae15-5b7a42f8bfb7",
		},
	},
}

func GetConsumerConfig(env string) *ConsumerConfig {
	config, isExist := cfgs[env]
	if !isExist {
		panic(fmt.Sprintf("Config for environment '%s' does not exist", env))
	}
	return &config
}
