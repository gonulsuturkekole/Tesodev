package config

import (
	"time"
)

// explain why we have the "shared" folder, why we have a config here and another config in seperate projects in the lecture?
type DbConfig struct {
	MongoDuration  time.Duration
	MongoClientURI string
}

var cfgs = map[string]DbConfig{
	"prod": {
		MongoDuration:  time.Second * 100,
		MongoClientURI: "mongodb+srv://melike:melike123@cluster0.yttfy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
		//MongoClientURI: "mongodb://root:root1234@localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.3.0",
	},
	"qa": {
		MongoDuration:  time.Second * 100,
		MongoClientURI: "mongodb+srv://melike:melike123@cluster0.yttfy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
		//MongoClientURI: "mongodb://root:root1234@localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.3.0",
	},
	"dev": {
		MongoDuration:  time.Second * 100,
		MongoClientURI: "mongodb+srv://melike:melike123@cluster0.yttfy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
		//MongoClientURI: "mongodb://root:root1234@localhost:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.3.0",
	},
}

func GetDBConfig(env string) *DbConfig {
	config, isExist := cfgs[env]
	if !isExist {
		panic("config does not exist")
	}
	return &config
}
