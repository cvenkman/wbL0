package config

import (
	"errors"
	"strings"
	"github.com/spf13/viper"
)

// struct for config file
type Config struct {
	Bind_addr string
	DB DBconfig
}

// struct for database config
type DBconfig struct {
	Name string
	Table string
	Username string
	Password string
	Host string
}

// add info from config to config and db struct's
func ReadConfig(configPath string) (Config, error) {
	slashIndex := strings.Index(configPath, "/")
	configName := configPath[slashIndex:strings.Index(configPath, ".")]
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath[:slashIndex])
	viper.SetConfigType("toml")

	var config Config
	err := viper.ReadInConfig()
	if err != nil {
		return config, errors.New("Can't read config: " + err.Error())
	}
	config.Bind_addr = viper.GetString("bind_addr")

	dbInfo := viper.GetStringMapString("database")
	config.DB.Name = dbInfo["name"]
	config.DB.Table = dbInfo["table"]
	config.DB.Password = dbInfo["password"]
	config.DB.Host = dbInfo["host"]
	config.DB.Username = dbInfo["username"]
	return config, nil
}