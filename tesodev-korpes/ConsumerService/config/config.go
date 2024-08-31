package config

import "fmt"

type DbConfig struct {
	DBName  string
	ColName string
}

type ConsumerConfig struct {
	Port     string
	DbConfig DbConfig
}

var cfgs = map[string]ConsumerConfig{
	"prod": {
		Port: ":1938",
		DbConfig: DbConfig{
			DBName:  "tesodev",
			ColName: "finance",
		},
	},
	"qa": {
		Port: ":1938",
		DbConfig: DbConfig{
			DBName:  "tesodev",
			ColName: "finance",
		},
	},
	"dev": {
		Port: ":1938",
		DbConfig: DbConfig{
			DBName:  "tesodev",
			ColName: "finance",
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
