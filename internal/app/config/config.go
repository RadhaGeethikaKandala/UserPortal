package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server       ServerConfig
	Database     DBConfig
	DatabaseTest DBConfig
}
type ServerConfig struct {
	Host string
	Port string
}
type DBConfig struct {
	Username string
	Password string
	Dbname   string
	Host     string
	Port     string
	Sslmode  string
}

const (
	configName = "config"
	configType = "yaml"
)

func ReadConfig(configPath string) *Config {
	var config = &Config{}
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Unable to read config file. Error : %s", err.Error())
	}

	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("Error in unmarshalling config : %s", err.Error())
	}
	return config
}
