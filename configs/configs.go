package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	SECRETKEY  string `mapstructure:"JWTSECRET"`
	Host       string `mapstructure:"HOST"`
	DBUser     string `mapstructure:"DBUSER"`
	Password   string `mapstructure:"PASSWORD"`
	Database   string `mapstructure:"DBNAME"`
	DBPORT     string `mapstructure:"PORT"`
	SERVERPORT string `mapstructure:"SERVERPORT"`
	Sslmode    string `mapstructure:"SSL"`
	REDISHOST  string `mapstructure:"REDISHOST"`
}

func LoadConfig() *Config {
	var config Config

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into config struct: %v", err)
	}

	return &config
}
