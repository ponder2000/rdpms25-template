package config

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	Port int

	App  *AppConf
	Log  *logConf
	Db   *Db
	Grpc *GrpcConf
}

type AppConf struct{}

type GrpcConf struct {
	Host string
	Port int
	Size struct {
		Send    int
		Receive int
	}
}

type logConf struct {
	Level     string
	Grpc      bool
	OnConsole bool
	IsJson    bool
	Trace     bool
	MaxAge    int
	MaxSize   int
}

type Db struct {
	Name           string
	Host           string
	Port           int
	User           string
	Password       string
	IdleConnection int
	OpenConnection int
}

func loadConf() *Configuration {
	c := Configuration{}
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	if e := viper.ReadInConfig(); e != nil {
		log.Fatal(e)
	}
	if e := viper.Unmarshal(&c); e != nil {
		log.Fatal(e)
	}
	return &c
}
