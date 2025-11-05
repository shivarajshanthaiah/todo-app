package configs

import (
	"log"
	"strings"

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
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile("../.env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found or could not read it: %v", err)
	} else {
		log.Println("Loaded .env file successfully")
	}

	// Explicitly bind expected environment variables
	keys := []string{
		"JWTSECRET", "HOST", "DBUSER", "PASSWORD", "DBNAME",
		"PORT", "SERVERPORT", "SSL", "REDISHOST",
	}
	for _, key := range keys {
		_ = viper.BindEnv(key)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into config struct: %v", err)
	}

	log.Printf("Loaded config: Host=%s, DBUser=%s, DBName=%s, Port=%s", config.Host, config.DBUser, config.Database, config.DBPORT)

	return &config
}
