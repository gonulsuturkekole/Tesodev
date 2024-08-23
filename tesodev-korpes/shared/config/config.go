package config

import (
	"time"
)

// explain why we have the "shared" folder, why we have a config here and another config in seperate projects in the lecture?
type DbConfig struct {
	MongoDuration  time.Duration
	MongoClientURI string
	SecretKey      string
}

var cfgs = map[string]DbConfig{
	"prod": {
		MongoDuration:  time.Second * 10,
		MongoClientURI: "mongodb://root:root1234@mongodb_docker:27017",
		SecretKey:      "079c9b74-24a7-4341-ae15-5b7a42f8bfb7",
	},
	"qa": {
		MongoDuration:  time.Second * 10,
		MongoClientURI: "mongodb://root:root1234@mongodb_docker:27017",
		SecretKey:      "079c9b74-24a7-4341-ae15-5b7a42f8bfb7",
	},
	"dev": {
		MongoDuration:  time.Second * 10,
		MongoClientURI: "mongodb://localhost:27017/",
		SecretKey:      "079c9b74-24a7-4341-ae15-5b7a42f8bfb7",
	},
}

func GetDBConfig(env string) *DbConfig {
	config, isExist := cfgs[env]
	if !isExist {
		panic("config does not exist")
	}
	return &config
}
