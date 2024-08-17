package config

type ConsumerConfig struct {
	Port     string
	DbConfig struct {
		DBName  string
		ColName string
	}
}

var cfgs = map[string]ConsumerConfig{
	"prod": {
		Port: ":8003",
		// This setup supports the independent operation of each service
		//whether they run on the same server or different servers
		// optimizing application performance and management.
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "finance",
		},
	},
	"qa": {
		Port: ":8003",
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "finance",
		},
	},
	"dev": {
		Port: ":8003",
		DbConfig: struct {
			DBName  string
			ColName string
		}{
			DBName:  "tesodev",
			ColName: "finance",
		},
	},
}

func GetConsumerConfig(env string) *ConsumerConfig {
	config, isExist := cfgs[env]
	if !isExist {
		panic("config does not exist")
	}
	return &config
}
