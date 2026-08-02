package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
	JWT   JWTConfig
}

type AppConfig struct {
	Name                 string
	Port                 int
	Env                  string
	LogLevel             string `mapstructure:"log_level"`
	MarketUpdateInterval string `mapstructure:"market_update_interval"`
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int `mapstructure:"expiration_hours"`
}

var App *Config

func LoadConfig(pDir string) {
	viper.AddConfigPath(pDir)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	if lErr := viper.ReadInConfig(); lErr != nil {
		log.Fatalf("Error reading config file, %s", lErr)
	}

	lEnv := viper.GetString("app.env")
	if lEnv != "" && lEnv != "config" {
		viper.SetConfigName(lEnv)
		if lErr := viper.MergeInConfig(); lErr != nil {
			log.Printf("No environment specific config found for %s, or error: %s", lEnv, lErr)
		}
	}

	var lConfig Config
	if lErr := viper.Unmarshal(&lConfig); lErr != nil {
		log.Fatalf("Unable to decode into struct, %v", lErr)
	}

	App = &lConfig
	log.Println("Configuration loaded successfully. Environment:", App.App.Env)
}
